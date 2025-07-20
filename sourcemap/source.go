package sourcemap

import "github.com/microsoft/typescript-go/core"

type Source interface {
	Text() string
	FileName() string
	LineMap() []core.TextPos
}
