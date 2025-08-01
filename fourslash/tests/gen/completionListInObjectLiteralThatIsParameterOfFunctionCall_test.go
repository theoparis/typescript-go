package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListInObjectLiteralThatIsParameterOfFunctionCall(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function f(a: { xa: number; xb: number; }) { }
var xc;
f({
    /**/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"xa",
				"xb",
			},
		},
	})
}
