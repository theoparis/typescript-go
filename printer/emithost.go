package printer

import (
	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/core"
	"github.com/microsoft/typescript-go/tsoptions"
	"github.com/microsoft/typescript-go/tspath"
)

type WriteFileData struct {
	SourceMapUrlPos int
	// BuildInfo BuildInfo
	Diagnostics      []*ast.Diagnostic
	DiffersOnlyInMap bool
	SkippedDtsWrite  bool
}

// NOTE: EmitHost operations must be thread-safe
type EmitHost interface {
	Options() *core.CompilerOptions
	SourceFiles() []*ast.SourceFile
	UseCaseSensitiveFileNames() bool
	GetCurrentDirectory() string
	CommonSourceDirectory() string
	IsEmitBlocked(file string) bool
	WriteFile(fileName string, text string, writeByteOrderMark bool, relatedSourceFiles []*ast.SourceFile, data *WriteFileData) error
	GetEmitModuleFormatOfFile(file ast.HasFileName) core.ModuleKind
	GetEmitResolver() EmitResolver
	GetOutputAndProjectReference(path tspath.Path) *tsoptions.OutputDtsAndProjectReference
	IsSourceFileFromExternalLibrary(file *ast.SourceFile) bool
}
