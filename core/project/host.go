package project

import "github.com/microsoft/typescript-go/core/vfs"

type ServiceHost interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	NewLine() string
}
