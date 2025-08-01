package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestPathCompletionsPackageJsonImportsSrcNoDistWildcard4(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /home/src/workspaces/project/tsconfig.json
{
  "compilerOptions": {
    "module": "nodenext",
    "rootDir": "src",
    "outDir": "dist"
  }
}
// @Filename: /home/src/workspaces/project/package.json
{
  "types": "index.d.ts",
  "imports": {
    "#*": "dist/*",
    "#foo/*": "dist/*",
    "#bar/*": "dist/*",
    "#exact-match": "dist/index.d.ts"
  }
}
// @Filename: /home/src/workspaces/project/nope.ts
export const nope = 0;
// @Filename: /home/src/workspaces/project/src/index.ts
export const index = 0;
// @Filename: /home/src/workspaces/project/src/blah.ts
export const blah = 0;
// @Filename: /home/src/workspaces/project/src/foo/onlyInFooFolder.ts
export const foo = 0;
// @Filename: /home/src/workspaces/project/src/subfolder/one.ts
export const one = 0;
// @Filename: /home/src/workspaces/project/src/a.mts
import { } from "/**/";`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "#a.mjs",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "#blah.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "#index.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "#foo",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
				&lsproto.CompletionItem{
					Label: "#subfolder",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
				&lsproto.CompletionItem{
					Label: "#bar",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
				&lsproto.CompletionItem{
					Label: "#exact-match",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
			},
		},
	})
	f.Insert(t, "#foo/")
	f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "a.mjs",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "blah.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "index.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "foo",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
				&lsproto.CompletionItem{
					Label: "subfolder",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
			},
		},
	})
	f.Insert(t, "foo/")
	f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "onlyInFooFolder.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
			},
		},
	})
}
