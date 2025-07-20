package ls

import (
	"github.com/microsoft/typescript-go/compiler"
	"github.com/microsoft/typescript-go/lsp/lsproto"
)

type Host interface {
	GetProgram() *compiler.Program
	GetPositionEncoding() lsproto.PositionEncodingKind
	GetLineMap(fileName string) *LineMap
}
