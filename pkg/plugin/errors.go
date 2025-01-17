package plugin

import (
	"fmt"

	api "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
)

const (
	FailedPlugin = "error"
)

func NewFailedPluginInfo(t Type, err error) *api.PluginInfo {
	return &api.PluginInfo{
		Name:   MkName(t, FailedPlugin),
		Fields: map[string]string{"error": err.Error()},
	}
}

type NotFoundError struct {
	Plugin *api.PluginName
}

func NewNotFoundError(t Type, name string) NotFoundError {
	return NotFoundError{Plugin: MkName(t, name)}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(
		"could not find plugin '%s' of type %s",
		e.Plugin.GetName(),
		e.Plugin.GetType(),
	)
}

type InvalidError struct {
	Plugin  *api.PluginName
	Message string
}

func NewInvalidError(t Type, name, msg string) InvalidError {
	return InvalidError{Plugin: MkName(t, name), Message: msg}
}

func (e InvalidError) Error() string {
	if e.Plugin == nil {
		return "invalid unknown plugin encountered: " + e.Message
	}

	return fmt.Sprintf(
		"invalid plugin '%s' of type %s: %s",
		e.Plugin.GetName(),
		e.Plugin.GetType(),
		e.Message,
	)
}
