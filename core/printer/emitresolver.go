package printer

import (
	"github.com/microsoft/typescript-go/core/ast"
	"github.com/microsoft/typescript-go/core/binder"
)

type EmitResolver interface {
	binder.ReferenceResolver
	IsReferencedAliasDeclaration(node *ast.Node) bool
	IsValueAliasDeclaration(node *ast.Node) bool
	IsTopLevelValueImportEqualsWithEntityName(node *ast.Node) bool
	MarkLinkedReferencesRecursively(file *ast.SourceFile)
	GetExternalModuleFileFromDeclaration(node *ast.Node) *ast.SourceFile
}
