package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListForExportEquals(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /node_modules/foo/index.d.ts
export = Foo;
declare var Foo: Foo.Static;
declare namespace Foo {
    interface Static {
        foo(): void;
    }
}
// @Filename: /a.ts
import { /**/ } from "foo";`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"foo",
				"Static",
				&lsproto.CompletionItem{
					Label:    "type",
					SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
		},
	})
}
