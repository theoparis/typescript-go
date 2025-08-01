package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionListInObjectBindingPattern13(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface I {
    x: number;
    y: string;
    z: boolean;
}

interface J {
    x: string;
    y: string;
}

let { /**/ }: I | J = { x: 10 };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"x",
				"y",
			},
		},
	})
}
