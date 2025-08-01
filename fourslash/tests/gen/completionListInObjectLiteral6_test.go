package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListInObjectLiteral6(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const foo = {
    a: "a",
    b: "b"
};
function fn<T extends { [key: string]: any }>(obj: T, events: { [Key in ` + "`" + `on_${string & keyof T}` + "`" + `]?: Key }) {}

fn(foo, {
    /*1*/
})
fn({ a: "a", b: "b" }, {
    /*2*/
})`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, []string{"1", "2"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:      "on_a?",
					InsertText: ptrTo("on_a"),
					FilterText: ptrTo("on_a"),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
				&lsproto.CompletionItem{
					Label:      "on_b?",
					InsertText: ptrTo("on_b"),
					FilterText: ptrTo("on_b"),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
			},
		},
	})
}
