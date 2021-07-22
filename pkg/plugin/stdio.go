package plugin

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"golang.org/x/sync/errgroup"
)

const (
	ErrNoStreamingConnection PluginError = "no streaming connection established"

	stdout muxStream = iota
	stderr

	streamListenAddress = "127.0.0.1:"

	bufSize = 32 * 1024

	headerSize = 8
	fdIndex    = 0
	sizeIndex  = 4
)

type muxStream byte

type StdioServer struct {
	listener   net.Listener
	connection net.Conn
}

func NewStdioServer() (*StdioServer, error) {
	listener, err := net.Listen("tcp", streamListenAddress)
	if err != nil {
		return nil, err
	}

	s := &StdioServer{listener: listener}

	// TODO: wait (somewhere) that this is done before proceeding
	go func() {
		if conn, err := s.listener.Accept(); err != nil {
			log.Default().Error("could not accept client connection", "error", err)
		} else {
			s.connection = conn
		}
	}()

	return s, nil
}

func (s *StdioServer) Endpoint() string {
	return s.listener.Addr().String()
}

func (s *StdioServer) Copy(inStream io.Writer, outStream io.Reader, errStream io.Reader) error {
	if s.connection == nil {
		return ErrNoStreamingConnection
	}

	eg, _ := errgroup.WithContext(context.Background())
	inCopier := NewCancelCopier(s.connection, inStream)
	outCopier := NewMuxCopier(outStream, errStream, s.connection)

	defer func() {
		if err := s.connection.Close(); err != nil {
			log.Default().Error("could not close client connection", "error", err)
		}

		if err := s.listener.Close(); err != nil {
			log.Default().Error("could not close listener", "error", err)
		}
	}()

	eg.Go(func() error {
		defer inCopier.Close()

		if err := outCopier.Copy(); err != nil {
			log.L().Error("could not copy stdout", "error", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := inCopier.Copy(); err != nil {
			log.L().Error("could not copy stdin", "error", err)
		}

		return nil
	})

	return eg.Wait()
}

type StdioClient struct {
	connection net.Conn
}

func NewStdioClient(url string) (*StdioClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", url)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &StdioClient{connection: conn}, nil
}

func (c *StdioClient) Copy(inStream io.Reader, outStream io.Writer, errStream io.Writer) error {
	eg, _ := errgroup.WithContext(context.Background())
	inCopier := NewCancelCopier(inStream, c.connection)
	outCopier := NewDemuxCopier(c.connection, outStream, errStream)

	defer func() {
		if err := c.connection.Close(); err != nil {
			log.L().Warn("could not close streaming connection", "error", err)
		}
	}()

	eg.Go(func() error {
		defer inCopier.Close()

		if err := outCopier.Copy(); err != nil {
			log.L().Error("could not copy stdout", "error", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := inCopier.Copy(); err != nil {
			log.L().Error("could not copy stdin", "error", err)
		}

		return nil
	})

	return eg.Wait()
}

type Copier interface {
	Copy() error
}

type MuxCopier struct {
	SrcOut io.Reader
	SrcErr io.Reader
	Dst    io.Writer

	bufPool *sync.Pool
}

func NewMuxCopier(srcOut io.Reader, srcErr io.Reader, dst io.Writer) *MuxCopier {
	return &MuxCopier{
		SrcOut: srcOut,
		SrcErr: srcErr,
		Dst:    dst,
		bufPool: &sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(nil) },
		},
	}
}

func (c *MuxCopier) Copy() error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return c.copyFrom(c.SrcOut, stdout)
	})

	eg.Go(func() error {
		return c.copyFrom(c.SrcErr, stderr)
	})

	return eg.Wait()
}

func (c *MuxCopier) copyFrom(src io.Reader, stream muxStream) error {
	readBuf := make([]byte, bufSize)

	for {
		read, readErr := src.Read(readBuf)
		if read > 0 {
			writeBuf := c.bufPool.Get().(*bytes.Buffer)
			c.writeHeader(writeBuf, stream, read)
			writeBuf.Write(readBuf[0:read])

			written, writeErr := c.Dst.Write(writeBuf.Bytes())

			written -= headerSize
			if written < 0 {
				written = 0
			}

			writeBuf.Reset()
			c.bufPool.Put(writeBuf)

			if writeErr != nil {
				return writeErr
			}

			if read != written {
				return io.ErrShortWrite
			}
		}

		if readErr == io.EOF {
			return nil
		}

		if readErr != nil {
			return readErr
		}
	}
}

func (c *MuxCopier) writeHeader(w io.Writer, stream muxStream, size int) (int, error) {
	header := [headerSize]byte{fdIndex: byte(stream)}
	binary.BigEndian.PutUint32(header[sizeIndex:], uint32(size))
	return w.Write(header[:])
}

type DemuxCopier struct {
	Src    io.Reader
	DstOut io.Writer
	DstErr io.Writer

	buf   []byte
	index int
}

func NewDemuxCopier(src io.Reader, dstOut io.Writer, dstErr io.Writer) *DemuxCopier {
	return &DemuxCopier{
		Src:    src,
		DstOut: dstOut,
		DstErr: dstErr,
		buf:    make([]byte, bufSize+headerSize+1),
		index:  0,
	}
}

func (c *DemuxCopier) Copy() error {
	for {
		out, size, err := c.readHeader()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		c.growTo(size)

		for c.index < size {
			if read, err := c.Src.Read(c.buf[c.index:]); err == io.EOF {
				return nil
			} else if err != nil {
				return err
			} else {
				c.index += read
			}
		}

		if written, err := out.Write(c.buf[headerSize:size]); err != nil {
			return err
		} else if written != size-headerSize {
			return io.ErrShortWrite
		}

		copy(c.buf, c.buf[size:])
		c.index -= size
	}
}

func (c *DemuxCopier) growTo(size int) {
	if size > len(c.buf) {
		c.buf = append(c.buf, make([]byte, size-len(c.buf)+1)...)
	}
}

func (c *DemuxCopier) readHeader() (io.Writer, int, error) {
	for c.index < headerSize {
		read, err := c.Src.Read(c.buf[c.index:])
		c.index += read

		if err != nil {
			return nil, 0, err
		}
	}

	var out io.Writer

	switch muxStream(c.buf[fdIndex]) {
	case stdout:
		out = c.DstOut
	case stderr:
		out = c.DstErr
	default:
		return out, 0, fmt.Errorf("Unrecognized input header: %d", c.buf[fdIndex])
	}

	size := headerSize + int(binary.BigEndian.Uint32(c.buf[sizeIndex:sizeIndex+4]))

	return out, size, nil
}

type CancelCopier struct {
	Src io.Reader
	Dst io.Writer

	buf    []byte
	read   chan int
	closed bool
	lock   sync.Mutex
}

func NewCancelCopier(src io.Reader, dst io.Writer) *CancelCopier {
	return &CancelCopier{
		Src:  src,
		Dst:  dst,
		buf:  make([]byte, bufSize),
		read: make(chan int),
	}
}

func (c *CancelCopier) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.closed {
		close(c.read)
		c.closed = true
	}
}

func (c *CancelCopier) Copy() error {
	for {
		go func() {
			if read, err := c.Src.Read(c.buf); err != nil {
				c.Close()
			} else {
				c.lock.Lock()
				defer c.lock.Unlock()
				if !c.closed {
					c.read <- read
				}
			}
		}()

		select {
		case read, ok := <-c.read:
			if read > 0 {
				if written, err := c.Dst.Write(c.buf[0:read]); err != nil {
					return err
				} else if read != written {
					return io.ErrShortWrite
				}
			}

			if !ok {
				return nil
			}
		}
	}
}
