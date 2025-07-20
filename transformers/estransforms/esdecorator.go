package estransforms

import (
	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/printer"
	"github.com/microsoft/typescript-go/transformers"
)

type esDecoratorTransformer struct {
	transformers.Transformer
}

func (ch *esDecoratorTransformer) visit(node *ast.Node) *ast.Node {
	return node // !!!
}

func newESDecoratorTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &esDecoratorTransformer{}
	return tx.NewTransformer(tx.visit, emitContext)
}
