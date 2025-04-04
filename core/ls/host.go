package ls

import (
	"github.com/microsoft/typescript-go/core/ast"
	"github.com/microsoft/typescript-go/core/compiler"
	"github.com/microsoft/typescript-go/core/core"
	"github.com/microsoft/typescript-go/core/tspath"
	"github.com/microsoft/typescript-go/core/vfs"
)

type Host interface {
	FS() vfs.FS
	DefaultLibraryPath() string
	GetCurrentDirectory() string
	NewLine() string
	Trace(msg string)
	GetProjectVersion() int
	// GetRootFileNames was called GetScriptFileNames in the original code.
	GetRootFileNames() []string
	// GetCompilerOptions was called GetCompilationSettings in the original code.
	GetCompilerOptions() *core.CompilerOptions
	GetSourceFile(fileName string, path tspath.Path, languageVersion core.ScriptTarget) *ast.SourceFile
	// This responsibility was moved from the language service to the project,
	// because they were bidirectionally interdependent.
	GetProgram() *compiler.Program
	GetDefaultLibraryPath() string
}
