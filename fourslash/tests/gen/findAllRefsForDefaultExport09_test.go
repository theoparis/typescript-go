package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestFindAllRefsForDefaultExport09(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "esnext",
        "strict": true,
        "outDir": "./out",
        "allowSyntheticDefaultImports": true
    }
}
// @filename: /a.js
module.exports = [];
// @filename: /b.js
module.exports = 1;
// @filename: /c.ts
export = [];
// @filename: /d.ts
export = 1;
// @filename: /foo.ts
import * as /*0*/a from "./a.js"
import /*1*/aDefault from "./a.js"
import * as /*2*/b from "./b.js"
import /*3*/bDefault from "./b.js"

import * as /*4*/c from "./c"
import /*5*/cDefault from "./c"
import * as /*6*/d from "./d"
import /*7*/dDefault from "./d"`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "0", "1", "2", "3", "4", "5", "6", "7")
}
