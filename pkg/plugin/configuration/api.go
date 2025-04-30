package configuration

import (
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
)

type Configuration interface {
	dodo.Plugin

	ListBackdrops() ([]Backdrop, error)
	GetBackdrop(name string) (Backdrop, error)
}
