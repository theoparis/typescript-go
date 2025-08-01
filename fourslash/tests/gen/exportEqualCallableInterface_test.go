package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestExportEqualCallableInterface(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: exportEqualCallableInterface_file0.ts
interface x {
    (): Date;
    foo: string;
}
export = x;
// @Filename: exportEqualCallableInterface_file1.ts
///<reference path='exportEqualCallableInterface_file0.ts'/>
import test = require('./exportEqualCallableInterface_file0');
var t2: test;
t2./**/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: completionFunctionMembersWithPrototypePlus(
				[]fourslash.CompletionsExpectedItem{
					"foo",
				}),
		},
	})
}
