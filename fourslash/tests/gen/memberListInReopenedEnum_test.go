package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestMemberListInReopenedEnum(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module M {
    enum E {
        A, B
    }
    enum E {
        C = 0, D
    }
    var x = E./*1*/
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "A",
					Detail: ptrTo("(enum member) E.A = 0"),
				},
				&lsproto.CompletionItem{
					Label:  "B",
					Detail: ptrTo("(enum member) E.B = 1"),
				},
				&lsproto.CompletionItem{
					Label:  "C",
					Detail: ptrTo("(enum member) E.C = 0"),
				},
				&lsproto.CompletionItem{
					Label:  "D",
					Detail: ptrTo("(enum member) E.D = 1"),
				},
			},
		},
	})
}
