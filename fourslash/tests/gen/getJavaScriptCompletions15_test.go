package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestGetJavaScriptCompletions15(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowNonTsExtensions: true
// @Filename: refFile1.ts
export var V = 1;
// @Filename: refFile2.ts
export var V = "123"
// @Filename: refFile3.ts
export var V = "123"
// @Filename: main.js
import ref1 = require("./refFile1");
var ref2 = require("./refFile2");
ref1.V./*1*/;
ref2.V./*2*/;
var v = { x: require("./refFile3") };
v.x./*3*/;
v.x.V./*4*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"toExponential",
			},
		},
	})
	f.VerifyCompletions(t, "2", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"toLowerCase",
			},
		},
	})
	f.VerifyCompletions(t, "3", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"V",
				&lsproto.CompletionItem{
					Label:    "ref1",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "ref2",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "require",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "v",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "x",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
			},
		},
	})
	f.VerifyCompletions(t, "4", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"toLowerCase",
			},
		},
	})
}
