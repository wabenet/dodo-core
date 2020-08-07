package decoder

import (
	"fmt"
	"reflect"
	"strconv"
)

type ConfigError string

const (
	ErrUnexpectedType ConfigError = "unexpected type"
	ErrUnexpectedKey  ConfigError = "unexpected key"
)

func (e ConfigError) Error() string {
	return string(e)
}

func Kinds(lookup map[reflect.Kind]Decoding) Decoding {
	return func(d *Decoder, config interface{}) {
		kind := reflect.ValueOf(config).Kind()
		if decode, ok := lookup[kind]; ok {
			decode(d, config)
		} else {
			d.Error(fmt.Errorf("could not decode type %v: %w", kind, ErrUnexpectedType))
		}
	}
}

func Keys(lookup map[string]Decoding) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().(map[interface{}]interface{})
		if !ok {
			d.Error(fmt.Errorf("could not decode map: %w", ErrUnexpectedType))
			return
		}

		for k, v := range decoded {
			key := k.(string)
			if decode, ok := lookup[key]; ok {
				d.Run(key, decode, v)
			} else {
				d.Error(fmt.Errorf("could not decode map: %w", ErrUnexpectedKey))
			}
		}
	}
}

func Slice(produce Producer, target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().([]interface{})
		if !ok {
			d.Error(fmt.Errorf("could not decode list: %w", ErrUnexpectedType))
			return
		}

		items := reflect.ValueOf(target).Elem()

		for i, item := range decoded {
			ptr, decode := produce()
			d.Run(strconv.Itoa(i), decode, item)
			items.Set(reflect.Append(items, reflect.ValueOf(ptr).Elem()))
		}
	}
}

func Singleton(produce Producer, target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		items := reflect.ValueOf(target).Elem()
		ptr, decode := produce()
		d.Run("", decode, config)
		items.Set(reflect.Append(items, reflect.ValueOf(ptr).Elem()))
	}
}

func Map(produce Producer, target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().(map[interface{}]interface{})
		if !ok {
			d.Error(fmt.Errorf("could not decode map: %w", ErrUnexpectedType))
			return
		}

		items := reflect.ValueOf(target).Elem()

		for key, value := range decoded {
			ptr, decode := produce()
			d.Run(key.(string), decode, value)
			items.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(ptr).Elem())
		}
	}
}

func String(target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := config.(string)
		if !ok {
			d.Error(fmt.Errorf("could not decode string: %w", ErrUnexpectedType))
			return
		}

		templated, err := ApplyTemplate(d, decoded)
		if err != nil {
			d.Error(err)
			return
		}

		reflect.ValueOf(target).Elem().SetString(templated)
	}
}

func NewString() Producer {
	return func() (interface{}, Decoding) {
		var target string
		return &target, String(&target)
	}
}

func Bool(target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := config.(bool)
		if !ok {
			d.Error(fmt.Errorf("could not decode boolean: %w", ErrUnexpectedType))
			return
		}

		reflect.ValueOf(target).Elem().SetBool(decoded)
	}
}

func NewBool() Producer {
	return func() (interface{}, Decoding) {
		var target bool
		return &target, Bool(&target)
	}
}

func Int(target interface{}) Decoding {
	return func(d *Decoder, config interface{}) {
		decoded, ok := config.(int64)
		if !ok {
			d.Error(fmt.Errorf("could not decode number: %w", ErrUnexpectedType))
			return
		}

		reflect.ValueOf(target).Elem().SetInt(decoded)
	}
}

func NewInt() Producer {
	return func() (interface{}, Decoding) {
		var target int64
		return &target, Int(&target)
	}
}
