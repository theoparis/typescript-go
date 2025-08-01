package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionsPaths_pathMapping_nonTrailingWildcard1(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /src/b.ts
export const x = 0;
// @Filename: /src/dir/x.ts
/export const x = 0;
// @Filename: /src/a.ts
import {} from "foo//*0*/";
import {} from "foo/dir//*1*/"; // invalid
import {} from "foo/_/*2*/";
import {} from "foo/_dir//*3*/";
// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "baseUrl": ".",
        "paths": {
            "foo/_*/suffix": ["src/*.ts"]
        }
    }
}`
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
					Label: "foo/_a/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "foo/_b/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "foo/_dir/suffix",
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
					Label: "foo/_a/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "foo/_b/suffix",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "foo/_dir/suffix",
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
					Label: "a",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
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
