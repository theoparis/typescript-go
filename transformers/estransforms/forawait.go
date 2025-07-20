package estransforms

import (
	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/printer"
	"github.com/microsoft/typescript-go/transformers"
)

type forawaitTransformer struct {
	transformers.Transformer
}

func (ch *forawaitTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newforawaitTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &forawaitTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
