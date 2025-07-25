package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/fourslash"
	"github.com/microsoft/typescript-go/testutil"
)

func TestExplainFilesNodeNextWithTypesReference(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /node_modules/react-hook-form/package.json
{
  "name": "react-hook-form",
  "main": "dist/index.cjs.js",
  "module": "dist/index.esm.js",
  "types": "dist/index.d.ts",
  "exports": {
    "./package.json": "./package.json",
    ".": {
      "import": "./dist/index.esm.js",
      "require": "./dist/index.cjs.js",
      "types": "./dist/index.d.ts"
    }
  }
}
// @Filename: /node_modules/react-hook-form/dist/index.cjs.js
module.exports = {};
// @Filename: /node_modules/react-hook-form/dist/index.esm.js
export function useForm() {}
// @Filename: /node_modules/react-hook-form/dist/index.d.ts
/// <reference types="react/**/" />
export type Foo = React.Whatever;
export function useForm(): any;
// @Filename: /node_modules/react/index.d.ts
declare namespace JSX {}
declare namespace React { export interface Whatever {} }
// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "module": "nodenext",
        "explainFiles": true
    }
    "files": ["./index.ts"]
}
// @Filename: /index.ts
import { useForm } from "react-hook-form";`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "")
}
