package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/ls"
	"github.com/microsoft/typescript-go/lsp/lsproto"
	"github.com/microsoft/typescript-go/testutil"
)

func TestTsxCompletion13(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: file.tsx
// @jsx: preserve
// @skipLibCheck: true
 declare module JSX {
     interface Element { }
     interface IntrinsicElements {
     }
     interface ElementAttributesProperty { props; }
 }
 interface ClickableProps {
     children?: string;
     className?: string;
 }
 interface ButtonProps extends ClickableProps {
     onClick(event?: React.MouseEvent<HTMLButtonElement>): void;
 }
 interface LinkProps extends ClickableProps {
     goTo: string;
 }
 declare function MainButton(buttonProps: ButtonProps): JSX.Element;
 declare function MainButton(linkProps: LinkProps): JSX.Element;
 declare function MainButton(props: ButtonProps | LinkProps): JSX.Element;
 let opt = <MainButton /*1*/ />;
 let opt = <MainButton children="chidlren" /*2*/ />;
 let opt = <MainButton onClick={()=>{}} /*3*/ />;
 let opt = <MainButton onClick={()=>{}} ignore-prop /*4*/ />;
 let opt = <MainButton goTo="goTo" /*5*/ />;
 let opt = <MainButton wrong /*6*/ />;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, []string{"1", "6"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"goTo",
				"onClick",
				&lsproto.CompletionItem{
					Label:      "children?",
					InsertText: ptrTo("children"),
					FilterText: ptrTo("children"),
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
				&lsproto.CompletionItem{
					Label:      "className?",
					InsertText: ptrTo("className"),
					FilterText: ptrTo("className"),
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
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
			Exact: []fourslash.CompletionsExpectedItem{
				"goTo",
				"onClick",
				&lsproto.CompletionItem{
					Label:      "className?",
					InsertText: ptrTo("className"),
					FilterText: ptrTo("className"),
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
			},
		},
	})
	f.VerifyCompletions(t, []string{"3", "4", "5"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:      "children?",
					InsertText: ptrTo("children"),
					FilterText: ptrTo("children"),
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
				&lsproto.CompletionItem{
					Label:      "className?",
					InsertText: ptrTo("className"),
					FilterText: ptrTo("className"),
					Kind:       ptrTo(lsproto.CompletionItemKindField),
					SortText:   ptrTo(string(ls.SortTextOptionalMember)),
				},
			},
		},
	})
}
