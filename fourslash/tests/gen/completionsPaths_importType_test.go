package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionsPaths_importType(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @moduleResolution: node
// @Filename: /ns.ts
file content not read
// @Filename: /node_modules/package/index.ts
file content not read
// @Filename: /usage.ts
type A = typeof import("p/*1*/");
type B = import(".//*2*/");
// @Filename: /user.js
/** @type {import("/*3*/")} */`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, []string{"1", "3"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "package",
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
					Label: "lib",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "lib.decorators",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "lib.decorators.legacy",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "ns",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "user",
					Kind:  ptrTo(lsproto.CompletionItemKindFile),
				},
				&lsproto.CompletionItem{
					Label: "node_modules",
					Kind:  ptrTo(lsproto.CompletionItemKindFolder),
				},
			},
		},
	})
}
