package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestPathCompletionsPackageJsonExportsWildcard12(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: node18
// @Filename: /node_modules/foo/package.json
 {
   "name": "foo",
   "exports": {
     "./bar/_*/suffix": "./dist/*.js"
   }
 }
// @Filename: /node_modules/foo/dist/b.d.ts
export const x = 0;
// @Filename: /node_modules/foo/dist/dir/x.d.ts
/export const x = 0;
// @Filename: /a.mts
import {} from "foo/bar//*0*/";
import {} from "foo/bar/dir//*1*/"; // invalid
import {} from "foo/bar/_/*2*/";
import {} from "foo/bar/_dir//*3*/";`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "0", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "bar/_b/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "bar/_dir/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
			},
		},
	})
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "bar/_b/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "bar/_dir/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
			},
		},
	})
	f.VerifyCompletions(t, "2", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "b",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "dir",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
			},
		},
	})
	f.VerifyCompletions(t, "3", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "x",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
			},
		},
	})
}
