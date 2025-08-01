package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestReferencesForAmbients(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/*1*/declare module "/*2*/foo" {
    /*3*/var /*4*/f: number;
}

/*5*/declare module "/*6*/bar" {
    /*7*/export import /*8*/foo = require("/*9*/foo");
    var f2: typeof /*10*/foo./*11*/f;
}

declare module "baz" {
    /*12*/import bar = require("/*13*/bar");
    var f2: typeof bar./*14*/foo;
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14")
}
