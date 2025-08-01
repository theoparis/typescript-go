package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestTsxCompletionNonTagLessThan(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.tsx
var x: Array<numb/*a*/;
[].map<numb/*b*/;
1 < Infini/*c*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, []string{"a", "b"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:    "number",
					SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
			Excludes: []string{
				"SVGNumber",
			},
		},
	})
	f.VerifyCompletions(t, "c", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:    "Infinity",
					SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
		},
	})
}
