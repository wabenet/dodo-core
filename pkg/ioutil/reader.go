package ioutil

import (
	"context"
	"fmt"
	"io"
)

const bufSize = 1024

//nolint:containedctx
type CancelableReader struct {
	wrapped io.Reader
	ctx     context.Context
	data    chan []byte
	err     error
}

func NewCancelableReader(ctx context.Context, r io.Reader) *CancelableReader {
	c := &CancelableReader{
		wrapped: r,
		ctx:     ctx,
		data:    make(chan []byte),
	}

	go c.begin()

	return c
}

func (c *CancelableReader) begin() {
	buf := make([]byte, bufSize)

	for {
		n, err := c.wrapped.Read(buf)

		if n > 0 {
			tmp := make([]byte, n)

			copy(tmp, buf[:n])

			c.data <- tmp
		}

		if err != nil {
			c.err = err
			close(c.data)

			return
		}
	}
}

func (c *CancelableReader) Read(p []byte) (int, error) {
	select {
	case <-c.ctx.Done():
		if err := c.ctx.Err(); err != nil {
			return 0, fmt.Errorf("error during async read: %w", err)
		}

		return 0, nil

	case d, ok := <-c.data:
		if ok {
			copy(p, d)

			return len(d), nil
		}

		if c.err != nil {
			return 0, fmt.Errorf("error during async read: %w", c.err)
		}

		return 0, nil
	}
}
