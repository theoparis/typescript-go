package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestTsxFindAllReferences8(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: file.tsx
// @jsx: preserve
// @noLib: true
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
 /*1*/declare function /*2*/MainButton(buttonProps: ButtonProps): JSX.Element;
 /*3*/declare function /*4*/MainButton(linkProps: LinkProps): JSX.Element;
 /*5*/declare function /*6*/MainButton(props: ButtonProps | LinkProps): JSX.Element;
 let opt = /*7*/</*8*/MainButton />;
 let opt = /*9*/</*10*/MainButton children="chidlren" />;
 let opt = /*11*/</*12*/MainButton onClick={()=>{}} />;
 let opt = /*13*/</*14*/MainButton onClick={()=>{}} ignore-prop />;
 let opt = /*15*/</*16*/MainButton goTo="goTo" />;
 let opt = /*17*/</*18*/MainButton wrong />;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18")
}
