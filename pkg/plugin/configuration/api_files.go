package configuration

import (
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/configuration/v1alpha2"
)

type FilesConfig []File

type File struct {
	FilePath string
	Contents []byte
}

func MergeFilesConfig(first, second FilesConfig) FilesConfig {
	return append(first, second...)
}

func FilesConfigFromProto(f []*api.File) FilesConfig {
	out := FilesConfig{}

	for _, file := range f {
		out = append(out, FileFromProto(file))
	}

	return out
}

func (f FilesConfig) ToProto() []*api.File {
	out := []*api.File{}

	for _, file := range f {
		out = append(out, file.ToProto())
	}

	return out
}

func FileFromProto(f *api.File) File {
	return File{
		FilePath: f.GetFilePath(),
		Contents: []byte(f.GetContents()),
	}
}

func (f File) ToProto() *api.File {
	out := &api.File{}

	out.SetFilePath(f.FilePath)
	out.SetContents(string(f.Contents))

	return out
}
