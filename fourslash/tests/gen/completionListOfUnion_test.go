package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListOfUnion(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strictNullChecks: true
const x: { a: number, b: number } | { a: string, c: string } | { b: boolean } | number | null | undefined = { /*x*/ };
interface I { a: number; }
function f(...args: Array<I | I[]>) {}
f({ /*f*/ });`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "x", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "a",
					Detail: ptrTo("(property) a: string | number"),
				},
				&lsproto.CompletionItem{
					Label:  "b",
					Detail: ptrTo("(property) b: number | boolean"),
				},
				&lsproto.CompletionItem{
					Label:  "c",
					Detail: ptrTo("(property) c: string"),
				},
			},
		},
	})
	f.VerifyCompletions(t, "f", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "a",
					Detail: ptrTo("(property) I.a: number"),
				},
			},
		},
	})
}
