package tstransforms

import (
	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/jsnum"
	"github.com/microsoft/typescript-go/printer"
)

func convertEntityNameToExpression(emitContext *printer.EmitContext, name *ast.EntityName) *ast.Expression {
	if ast.IsQualifiedName(name) {
		left := convertEntityNameToExpression(emitContext, name.AsQualifiedName().Left)
		right := name.AsQualifiedName().Right
		prop := emitContext.Factory.NewPropertyAccessExpression(left, nil /*questionDotToken*/, right, ast.NodeFlagsNone)
		emitContext.SetOriginal(prop, name)
		emitContext.AssignCommentAndSourceMapRanges(prop, name)
		return prop
	}
	return name.Clone(emitContext.Factory)
}

func constantExpression(value any, factory *printer.NodeFactory) *ast.Expression {
	switch value := value.(type) {
	case string:
		return factory.NewStringLiteral(value)
	case jsnum.Number:
		if value.IsInf() || value.IsNaN() {
			return nil
		}
		if value < 0 {
			return factory.NewPrefixUnaryExpression(ast.KindMinusToken, constantExpression(-value, factory))
		}
		return factory.NewNumericLiteral(value.String())
	}
	return nil
}

func isInstantiatedModule(node *ast.ModuleDeclarationNode, preserveConstEnums bool) bool {
	moduleState := ast.GetModuleInstanceState(node)
	return moduleState == ast.ModuleInstanceStateInstantiated ||
		(preserveConstEnums && moduleState == ast.ModuleInstanceStateConstEnumOnly)
}
