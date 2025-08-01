package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionsRedeclareModuleAsGlobal(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @esModuleInterop: true,
// @target: esnext
// @Filename: /myAssert.d.ts
declare function assert(value:any, message?:string):void;
export = assert;
export as namespace assert;
// @Filename: /ambient.d.ts
import assert from './myAssert';

type Assert = typeof assert;

declare global {
  const assert: Assert;
}
// @Filename: /index.ts
/// <reference path="./ambient.d.ts" />
asser/**/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:    "assert",
					SortText: ptrTo(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
		},
	})
}
