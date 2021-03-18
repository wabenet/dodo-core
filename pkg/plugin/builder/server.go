package builder

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type server struct {
	impl ImageBuilder
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo()
}

func (s *server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	id, err := s.impl.CreateImage(request.Config)
	if err != nil {
		return nil, err
	}

	return &api.CreateImageResponse{ImageId: id}, nil
}
