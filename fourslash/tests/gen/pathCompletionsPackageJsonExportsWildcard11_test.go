package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestPathCompletionsPackageJsonExportsWildcard11(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: preserve
// @moduleResolution: bundler
// @jsx: react
// @Filename: /node_modules/repo/package.json
{
  "name": "repo",
  "exports": {
    "./*": "./src/*"
  }
}
// @Filename: /node_modules/repo/src/card.tsx
export {};
// @Filename: /main.ts
import { } from "repo//**/";`
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
					Label: "card.js",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
			},
		},
	})
}
