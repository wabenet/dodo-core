package runtime

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	api "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
)

type MountConfig []Mount

func MountConfigFromProto(m []*api.Mount) MountConfig {
	out := MountConfig{}

	for _, mnt := range m {
		out = append(out, MountFromProto(mnt))
	}

	return out
}

func (m MountConfig) ToProto() []*api.Mount {
	out := []*api.Mount{}

	for _, mnt := range m {
		out = append(out, mnt.ToProto())
	}

	return out
}

type Mount interface {
	Type() MountType
	ToProto() *api.Mount
}

type MountType string

func MergeMountConfig(first, second MountConfig) MountConfig {
	return append(first, second...)
}

func MountFromProto(m *api.Mount) Mount {
	switch m := m.GetType().(type) {
	case *api.Mount_Bind:
		return BindMountFromProto(m.Bind)
	case *api.Mount_Volume:
		return VolumeMountFromProto(m.Volume)
	case *api.Mount_Tmpfs:
		return TmpfsMountFromProto(m.Tmpfs)
	case *api.Mount_Image:
		return ImageMountFromProto(m.Image)
	case *api.Mount_Device:
		return DeviceMountFromProto(m.Device)
	default:
		return nil
	}
}

const TypeBind MountType = "bind"

var ErrVolumeFormat = errors.New("invalid volume format")

type BindMount struct {
	HostPath      string
	ContainerPath string
	Readonly      bool
}

func (BindMount) Type() MountType {
	return TypeBind
}

func BindMountFromProto(m *api.BindMount) BindMount {
	return BindMount{
		HostPath:      m.GetHostPath(),
		ContainerPath: m.GetContainerPath(),
		Readonly:      m.GetReadonly(),
	}
}

func (m BindMount) ToProto() *api.Mount {
	return &api.Mount{
		Type: &api.Mount_Bind{
			Bind: &api.BindMount{
				HostPath:      m.HostPath,
				ContainerPath: m.ContainerPath,
				Readonly:      m.Readonly,
			},
		},
	}
}

func BindMountFromSpec(spec string) (Mount, error) {
	var source, target string

	var readonly bool

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return BindMount{}, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	case 1:
		source = values[0]
	case 2:
		source = values[0]
		target = values[1]
	case 3:
		source = values[0]
		target = values[1]
		readonly = (values[2] == "ro")
	default:
		return BindMount{}, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	}

	// If source is an absolute path, assume bind mount, otherwise assume a volume
	if source[0] == '/' {
		return BindMount{
			HostPath:      source,
			ContainerPath: target,
			Readonly:      readonly,
		}, nil
	}

	return VolumeMount{
		VolumeName:    source,
		ContainerPath: target,
		Readonly:      readonly,
	}, nil
}

const TypeVolume MountType = "volume"

type VolumeMount struct {
	VolumeName    string
	ContainerPath string
	Subpath       string
	Readonly      bool
}

func (VolumeMount) Type() MountType {
	return TypeVolume
}

func VolumeMountFromProto(m *api.VolumeMount) Mount {
	return VolumeMount{
		VolumeName:    m.GetVolumeName(),
		ContainerPath: m.GetContainerPath(),
		Subpath:       m.GetSubpath(),
		Readonly:      m.GetReadonly(),
	}
}

func (m VolumeMount) ToProto() *api.Mount {
	return &api.Mount{
		Type: &api.Mount_Volume{
			Volume: &api.VolumeMount{
				VolumeName:    m.VolumeName,
				ContainerPath: m.ContainerPath,
				Subpath:       m.Subpath,
				Readonly:      m.Readonly,
			},
		},
	}
}

const TypeTmpfs MountType = "tmpfs"

type TmpfsMount struct {
	ContainerPath string
	Size          int
	Mode          os.FileMode
}

func (TmpfsMount) Type() MountType {
	return TypeTmpfs
}

func TmpfsMountFromProto(m *api.TmpfsMount) TmpfsMount {
	// TODO error handling
	mode, _ := strconv.ParseUint(m.GetMode(), 8, 32)

	return TmpfsMount{
		ContainerPath: m.GetContainerPath(),
		Size:          int(m.GetSize()),
		Mode:          os.FileMode(mode),
	}
}

func (m TmpfsMount) ToProto() *api.Mount {
	return &api.Mount{
		Type: &api.Mount_Tmpfs{
			Tmpfs: &api.TmpfsMount{
				ContainerPath: m.ContainerPath,
				Size:          int64(m.Size),
				Mode:          m.Mode.String(),
			},
		},
	}
}

const TypeImage MountType = "image"

type ImageMount struct {
	Image         string
	ContainerPath string
	Subpath       string
	Readonly      bool
}

func (ImageMount) Type() MountType {
	return TypeImage
}

func ImageMountFromProto(m *api.ImageMount) ImageMount {
	return ImageMount{
		Image:         m.GetImage(),
		ContainerPath: m.GetContainerPath(),
		Subpath:       m.GetSubpath(),
		Readonly:      m.GetReadonly(),
	}
}

func (m ImageMount) ToProto() *api.Mount {
	return &api.Mount{
		Type: &api.Mount_Image{
			Image: &api.ImageMount{
				Image:         m.Image,
				ContainerPath: m.ContainerPath,
				Subpath:       m.Subpath,
				Readonly:      m.Readonly,
			},
		},
	}
}

const TypeDevice MountType = "device"

var ErrDeviceFormat = errors.New("invalid device format")

type DeviceMount struct {
	CGroupRule    string
	HostPath      string
	ContainerPath string
	Permissions   string
}

func (DeviceMount) Type() MountType {
	return TypeDevice
}

func DeviceMountFromProto(m *api.DeviceMount) DeviceMount {
	return DeviceMount{
		CGroupRule:    m.GetCgroupRule(),
		HostPath:      m.GetHostPath(),
		ContainerPath: m.GetContainerPath(),
		Permissions:   m.GetPermissions(),
	}
}

func (m DeviceMount) ToProto() *api.Mount {
	return &api.Mount{
		Type: &api.Mount_Device{
			Device: &api.DeviceMount{
				CgroupRule:    m.CGroupRule,
				HostPath:      m.HostPath,
				ContainerPath: m.ContainerPath,
				Permissions:   m.Permissions,
			},
		},
	}
}

func DeviceMountFromSpec(spec string) (Mount, error) {
	dev := DeviceMount{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return dev, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	case 1:
		dev.HostPath = values[0]
	case 2:
		dev.HostPath = values[0]
		dev.ContainerPath = values[1]
	case 3:
		dev.HostPath = values[0]
		dev.ContainerPath = values[1]
		dev.Permissions = values[2]
	default:
		return dev, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	}

	return dev, nil
}
