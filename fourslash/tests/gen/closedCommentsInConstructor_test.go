package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestClosedCommentsInConstructor(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class Foo {
    constructor(/* /**/ */) { }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", nil)
}
