package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionNoAutoInsertQuestionDotForTypeParameter(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strict: true
interface Address {
    city: string = "";
    "postal code": string = "";
}
function f<T extends Address>(x: T) {
    x[|./**/|]
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "city",
					Detail: ptrTo("(property) Address.city: string"),
				},
				&lsproto.CompletionItem{
					Label:      "postal code",
					InsertText: ptrTo("[\"postal code\"]"),
					Detail:     ptrTo("(property) Address[\"postal code\"]: string"),
				},
			},
		},
	})
}
