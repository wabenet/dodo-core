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
	switch m.WhichType() {
	case api.Mount_Bind_case:
		return BindMountFromProto(m.GetBind())
	case api.Mount_Volume_case:
		return VolumeMountFromProto(m.GetVolume())
	case api.Mount_Tmpfs_case:
		return TmpfsMountFromProto(m.GetTmpfs())
	case api.Mount_Image_case:
		return ImageMountFromProto(m.GetImage())
	case api.Mount_Device_case:
		return DeviceMountFromProto(m.GetDevice())
	case api.Mount_Type_not_set_case:
		return nil
	default:
		return nil
	}
}

const TypeBind MountType = "bind"

var _ Mount = BindMount{}

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
	mnt := &api.BindMount{}

	mnt.SetHostPath(m.HostPath)
	mnt.SetContainerPath(m.ContainerPath)
	mnt.SetReadonly(m.Readonly)

	out := &api.Mount{}

	out.SetBind(mnt)

	return out
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

var _ Mount = VolumeMount{}

type VolumeMount struct {
	VolumeName    string
	ContainerPath string
	Subpath       string
	Readonly      bool
}

func (VolumeMount) Type() MountType {
	return TypeVolume
}

func VolumeMountFromProto(m *api.VolumeMount) VolumeMount {
	return VolumeMount{
		VolumeName:    m.GetVolumeName(),
		ContainerPath: m.GetContainerPath(),
		Subpath:       m.GetSubpath(),
		Readonly:      m.GetReadonly(),
	}
}

func (m VolumeMount) ToProto() *api.Mount {
	mnt := &api.VolumeMount{}

	mnt.SetVolumeName(m.VolumeName)
	mnt.SetContainerPath(m.ContainerPath)
	mnt.SetSubpath(m.Subpath)
	mnt.SetReadonly(m.Readonly)

	out := &api.Mount{}

	out.SetVolume(mnt)

	return out
}

const TypeTmpfs MountType = "tmpfs"

var _ Mount = TmpfsMount{}

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
	mnt := &api.TmpfsMount{}

	mnt.SetContainerPath(m.ContainerPath)
	mnt.SetSize(int64(m.Size))
	mnt.SetMode(m.Mode.String())

	out := &api.Mount{}

	out.SetTmpfs(mnt)

	return out
}

const TypeImage MountType = "image"

var _ Mount = ImageMount{}

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
	mnt := &api.ImageMount{}

	mnt.SetImage(m.Image)
	mnt.SetContainerPath(m.ContainerPath)
	mnt.SetSubpath(m.Subpath)
	mnt.SetReadonly(m.Readonly)

	out := &api.Mount{}

	out.SetImage(mnt)

	return out
}

const TypeDevice MountType = "device"

var _ Mount = DeviceMount{}

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
	mnt := &api.DeviceMount{}

	mnt.SetCgroupRule(m.CGroupRule)
	mnt.SetHostPath(m.HostPath)
	mnt.SetContainerPath(m.ContainerPath)
	mnt.SetPermissions(m.Permissions)

	out := &api.Mount{}

	out.SetDevice(mnt)

	return out
}

func DeviceMountFromSpec(spec string) (Mount, error) {
	out := DeviceMount{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return out, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	case 1:
		out.HostPath = values[0]
	case 2:
		out.HostPath = values[0]
		out.ContainerPath = values[1]
	case 3:
		out.HostPath = values[0]
		out.ContainerPath = values[1]
		out.Permissions = values[2]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	}

	return out, nil
}
