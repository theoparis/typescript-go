package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListObjectMembers(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = ` var object: {
     (bar: any): any;
     new (bar: any): any;
     [bar: any]: any;
     bar: any;
     foo(bar: any): any;
 };
object./**/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "bar",
					Detail: ptrTo("(property) bar: any"),
				},
				&lsproto.CompletionItem{
					Label:  "foo",
					Detail: ptrTo("(method) foo(bar: any): any"),
				},
			},
		},
	})
}
