package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestJsDocGenerics1(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowNonTsExtensions: true
// @Filename: ref.d.ts
 namespace Thing {
     export interface Thung {
         a: number;
     ]
 ]
// @Filename: Foo.js

 /** @type {Array<number>} */
 var v;
 v[0]./*1*/

 /** @type {{x: Array<Array<number>>}} */
 var w;
 w.x[0][0]./*2*/

 /** @type {Array<Thing.Thung>} */
 var x;
 x[0].a./*3*/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, f.Markers(), &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "toFixed",
					Kind:  ptrTo(lsproto.CompletionItemKindMethod),
				},
			},
		},
	})
}
