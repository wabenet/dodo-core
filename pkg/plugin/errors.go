package plugin

import (
	"fmt"
)

const (
	FailedPlugin = "error"
)

func NewFailedPluginInfo(t Type, err error) Metadata {
	return Metadata{
		ID:     ID{Type: t.String(), Name: FailedPlugin},
		Labels: Labels{"error": err.Error()},
	}
}

type NotFoundError struct {
	PluginID ID
}

func NewNotFoundError(t Type, name string) NotFoundError {
	return NotFoundError{PluginID: ID{Type: t.String(), Name: name}}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(
		"could not find plugin '%s' of type %s",
		e.PluginID.Name,
		e.PluginID.Type,
	)
}

type InvalidError struct {
	PluginID ID
	Message  string
}

func NewInvalidError(t Type, name, msg string) InvalidError {
	return InvalidError{PluginID: ID{Type: t.String(), Name: name}, Message: msg}
}

func (e InvalidError) Error() string {
	return fmt.Sprintf(
		"invalid plugin '%s' of type %s: %s",
		e.PluginID.Name,
		e.PluginID.Type,
		e.Message,
	)
}
