package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestGetJavaScriptCompletions14(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowNonTsExtensions: true
// @Filename: file1.js
/// <reference no-default-lib="true"/>
interface Number {
    toExponential(fractionDigits?: number): string;
}
var x = 1;
x./*1*/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "toExponential",
					Kind:  ptrTo(lsproto.CompletionItemKindMethod),
				},
			},
		},
	})
}
