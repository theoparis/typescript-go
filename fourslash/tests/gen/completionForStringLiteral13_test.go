package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionForStringLiteral13(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface SymbolConstructor {
    readonly species: symbol;
}
var Symbol: SymbolConstructor;
interface PromiseConstructor {
  [Symbol.species]: PromiseConstructor;
}
var Promise: PromiseConstructor;
Promise["/*1*/"];`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{},
		},
	})
}
