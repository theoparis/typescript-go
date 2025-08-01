package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestTsxCompletion7(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: file.tsx
 declare module JSX {
     interface Element { }
     interface IntrinsicElements {
         div: { ONE: string; TWO: number; }
     }
 }
 let y = { ONE: '' };
 var x = <div {...y} /**/ />;`
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
					Label:    "TWO",
					Kind:     ptrTo(lsproto.CompletionItemKindField),
					SortText: ptrTo(string(ls.SortTextLocationPriority)),
				},
				&lsproto.CompletionItem{
					Label:    "ONE",
					Kind:     ptrTo(lsproto.CompletionItemKindField),
					SortText: ptrTo(string(ls.SortTextMemberDeclaredBySpreadAssignment)),
				},
			},
		},
	})
}
