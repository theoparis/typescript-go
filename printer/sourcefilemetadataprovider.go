package printer

import (
	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}
