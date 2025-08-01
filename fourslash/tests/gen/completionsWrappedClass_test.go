package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestCompletionsWrappedClass(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class Client {
    private close() { }
    public open() { }
}
type Wrap<T> = T &
{
    [K in Extract<keyof T, string> as ` + "`" + `${K}Wrapped` + "`" + `]: T[K];
};
class Service {
    method() {
        let service = undefined as unknown as Wrap<Client>;
        const { /*a*/ } = service;
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "a", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"open",
				"openWrapped",
			},
		},
	})
}
