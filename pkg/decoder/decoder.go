package decoder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type (
	Decoding  func(*Decoder, interface{})
	Producer func() (interface{}, Decoding)

	Decoder struct {
		file   string
		path   []string
		errors []error
	}
)

func New(filename string) *Decoder {
	return &Decoder{
		file:   filename,
		path:   []string{},
		errors: []error{},
	}
}

func (d *Decoder) Error(msg string) {
	// TODO add file and path to error messages
	d.errors = append(d.errors, errors.New(msg))
}

func (d *Decoder) Errors() []error {
	return d.errors
}

func (d *Decoder) Run(name string, decode Decoding, value interface{}) {
	d.path = append(d.path, name)
	decode(d, value)
	d.path = d.path[:len(d.path)-1]
}

func (d *Decoder) DecodeYaml(content []byte, target interface{}, lookup map[string]Decoding) {
	var mapType map[interface{}]interface{}
	if err := yaml.Unmarshal(content, &mapType); err != nil {
		d.Error("invalid yaml")
		return
	}

	var dummySlice []struct{}
	var dummyItem struct{}

	includeHelper := func() (interface{}, Decoding) {
		return &dummyItem, Keys(map[string]Decoding{
			"text": func(d *Decoder, config interface{}) {
				var decoded string
				String(&decoded)(d, config)
				d.DecodeYaml([]byte(decoded), target, lookup)
			},
			"file": func(d *Decoder, config interface{}) {
				var decoded string
				String(&decoded)(d, config)

				bytes, err := readFile(decoded)
				if err != nil {
					d.Error("could not read file")
					return
				}

				sub := New(decoded)
				sub.DecodeYaml(bytes, target, lookup)
				d.errors = append(d.errors, sub.errors...)
			},
		})
	}

	lookup["include"] = Kinds(map[reflect.Kind]Decoding{
		reflect.Map:   Singleton(includeHelper, &dummySlice),
		reflect.Slice: Slice(includeHelper, &dummySlice),
	})

	d.Run("", Keys(lookup), mapType)
}

func readFile(filename string) ([]byte, error) {
	if !filepath.IsAbs(filename) {
		directory, err := os.Getwd()
		if err != nil {
			return []byte{}, err
		}
		filename, err = filepath.Abs(filepath.Join(directory, filename))
		if err != nil {
			return []byte{}, err
		}
	}
	return ioutil.ReadFile(filename)
}
