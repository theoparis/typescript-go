//// [tests/cases/compiler/jsFileCompilationWithMapFileAsJsWithOutDir.ts] ////

//// [a.ts]
class c {
}

//// [b.js.map]
function foo() {
}

//// [b.js]
function bar() {
}

//// [a.js]
class c {
}
//# sourceMappingURL=a.js.map
//// [b.js]
function bar() {
}
//# sourceMappingURL=b.js.map