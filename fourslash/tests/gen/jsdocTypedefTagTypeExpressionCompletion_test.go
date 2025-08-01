package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestJsdocTypedefTagTypeExpressionCompletion(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface I {
    age: number;
}
 class Foo {
     property1: string;
     constructor(value: number) { this.property1 = "hello"; }
     static method1() {}
     method3(): number { return 3; }
     /**
      * @param {string} foo A value.
      * @returns {number} Another value
      * @mytag
      */
     method4(foo: string) { return 3; }
 }
 namespace Foo.Namespace { export interface SomeType { age2: number } }
 /**
  * @type { /*type1*/Foo./*typeFooMember*/Namespace./*NamespaceMember*/SomeType }
  */
var x;
/*globalValue*/
x./*valueMemberOfSomeType*/
var x1: Foo;
x1./*valueMemberOfFooInstance*/;
Foo./*valueMemberOfFoo*/;
 /**
  * @type { {/*propertyName*/ageX: number} }
  */
var y;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "type1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "Foo",
					Kind:  ptrTo(lsproto.CompletionItemKindClass),
				},
				&lsproto.CompletionItem{
					Label: "I",
					Kind:  ptrTo(lsproto.CompletionItemKindInterface),
				},
			},
			Excludes: []string{
				"Namespace",
				"SomeType",
				"x",
				"x1",
				"y",
				"method1",
				"property1",
				"method3",
				"method4",
				"foo",
			},
		},
	})
	f.VerifyCompletions(t, "typeFooMember", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "Namespace",
					Kind:  ptrTo(lsproto.CompletionItemKindModule),
				},
			},
		},
	})
	f.VerifyCompletions(t, "NamespaceMember", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "SomeType",
					Kind:  ptrTo(lsproto.CompletionItemKindInterface),
				},
			},
		},
	})
	f.VerifyCompletions(t, "globalValue", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "Foo",
					Kind:  ptrTo(lsproto.CompletionItemKindClass),
				},
				&lsproto.CompletionItem{
					Label: "x",
					Kind:  ptrTo(lsproto.CompletionItemKindVariable),
				},
				&lsproto.CompletionItem{
					Label: "x1",
					Kind:  ptrTo(lsproto.CompletionItemKindVariable),
				},
				&lsproto.CompletionItem{
					Label: "y",
					Kind:  ptrTo(lsproto.CompletionItemKindVariable),
				},
			},
			Excludes: []string{
				"I",
				"Namespace",
				"SomeType",
				"method1",
				"property1",
				"method3",
				"method4",
				"foo",
			},
		},
	})
	f.VerifyCompletions(t, "valueMemberOfSomeType", nil)
	f.VerifyCompletions(t, "valueMemberOfFooInstance", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "method3",
					Kind:  ptrTo(lsproto.CompletionItemKindMethod),
				},
				&lsproto.CompletionItem{
					Label: "method4",
					Kind:  ptrTo(lsproto.CompletionItemKindMethod),
				},
				&lsproto.CompletionItem{
					Label: "property1",
					Kind:  ptrTo(lsproto.CompletionItemKindField),
				},
			},
		},
	})
	f.VerifyCompletions(t, "valueMemberOfFoo", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: completionFunctionMembersPlus(
				[]fourslash.CompletionsExpectedItem{
					&lsproto.CompletionItem{
						Label:    "method1",
						Kind:     ptrTo(lsproto.CompletionItemKindMethod),
						SortText: ptrTo(string(ls.SortTextLocalDeclarationPriority)),
					},
					&lsproto.CompletionItem{
						Label:    "prototype",
						SortText: ptrTo(string(ls.SortTextLocationPriority)),
					},
				}),
		},
	})
	f.VerifyCompletions(t, "propertyName", nil)
}
