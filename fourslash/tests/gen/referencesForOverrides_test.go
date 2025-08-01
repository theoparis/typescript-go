package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestReferencesForOverrides(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module FindRef3 {
	module SimpleClassTest {
		export class Foo {
			public /*foo*/foo(): void {
			}
		}
		export class Bar extends Foo {
			public foo(): void {
			}
		}
	}

	module SimpleInterfaceTest {
		export interface IFoo {
			/*ifoo*/ifoo(): void;
		}
		export interface IBar extends IFoo {
			ifoo(): void;
		}
	}

	module SimpleClassInterfaceTest {
		export interface IFoo {
			/*icfoo*/icfoo(): void;
		}
		export class Bar implements IFoo {
			public icfoo(): void {
			}
		}
	}

	module Test {
		export interface IBase {
			/*field*/field: string;
			/*method*/method(): void;
		}

		export interface IBlah extends IBase {
			field: string;
		}

		export interface IBlah2 extends IBlah {
			field: string;
		}

		export interface IDerived extends IBlah2 {
			method(): void;
		}

		export class Bar implements IDerived {
			public field: string;
			public method(): void { }
		}

		export class BarBlah extends Bar {
			public field: string;
		}
	}

	function test() {
		var x = new SimpleClassTest.Bar();
		x.foo();

		var y: SimpleInterfaceTest.IBar = null;
		y.ifoo();

        var w: SimpleClassInterfaceTest.Bar = null;
        w.icfoo();

		var z = new Test.BarBlah();
		z.field = "";
        z.method();
	}
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "foo", "ifoo", "icfoo", "field", "method")
}
