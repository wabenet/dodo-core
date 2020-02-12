package decoder

import (
	"fmt"
	"reflect"
)

type ConfigError struct {
	Name            string
	UnsupportedType reflect.Kind
}

func (e *ConfigError) Error() string {
	if e.UnsupportedType != reflect.Invalid {
		return fmt.Sprintf("unsupported type of '%s': '%s'", e.Name, e.UnsupportedType.String())
	}
	return fmt.Sprintf("configuration error in '%s'", e.Name)
}
