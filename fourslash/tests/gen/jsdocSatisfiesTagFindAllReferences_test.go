package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestJsdocSatisfiesTagFindAllReferences(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noEmit: true
// @allowJS: true
// @checkJs: true
// @filename: /a.js
/**
 * @typedef {Object} T
 * @property {number} a
 */

/** @satisfies {/**/T} comment */
const foo = { a: 1 };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "")
}
