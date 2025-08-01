//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/core"
	"github.com/microsoft/typescript-go/parser"
	"github.com/microsoft/typescript-go/repo"
	"github.com/microsoft/typescript-go/tspath"
)

var (
	libInputDir     = filepath.Join(repo.TypeScriptSubmodulePath, "src", "lib")
	copyrightNotice = filepath.Join(repo.TypeScriptSubmodulePath, "scripts", "CopyrightNotice.txt")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	libs := readLibs()
	generateLibs(libs)
	generateLibList(libs)
	generateEmbedded(libs)
}

type lib struct {
	target  string   // target relative to libs dir
	sources []string // sources relative to src/lib dir
}

func generateLibs(libs []lib) {
	const outputDir = "libs"

	copyright := readCopyright()

	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatalf("failed to remove libs directory: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("failed to create libs directory: %v", err)
	}

	for _, lib := range libs {
		var output bytes.Buffer
		output.Write(copyright)

		for _, source := range lib.sources {
			sourcePath := filepath.Join(libInputDir, source)
			b, err := os.ReadFile(sourcePath)
			if err != nil {
				log.Fatalf("failed to read %s: %v", sourcePath, err)
			}

			output.WriteByte('\n')
			output.Write(removeCRLF(b))
		}

		outputPath := filepath.Join(outputDir, lib.target)
		if err := os.WriteFile(outputPath, output.Bytes(), 0o644); err != nil {
			log.Fatalf("failed to write %s: %v", outputPath, err)
		}
	}
}

func generateLibList(libs []lib) {
	var code bytes.Buffer
	code.WriteString("// Code generated by generate.go; DO NOT EDIT.\n\n")
	code.WriteString("package bundled\n\n")

	code.WriteString("// LibNames is the list of all bundled lib files, sorted by name.\n")
	code.WriteString("// For the list of libs sorted by load order, use [tsoptions.Libs].\n")
	code.WriteString("var LibNames = []string{\n")
	for _, lib := range libs {
		code.WriteString("\t\"" + lib.target + "\",\n")
	}
	code.WriteString("}\n")

	writeCode("libs_generated.go", code.Bytes())
}

func generateEmbedded(libs []lib) {
	libVarNames := make([]string, len(libs))
	for i, lib := range libs {
		libVarNames[i] = "libs_" + strings.ReplaceAll(lib.target, ".", "_")
	}

	var code bytes.Buffer
	code.WriteString("//go:build !noembed\n\n")
	code.WriteString("// Code generated by generate.go; DO NOT EDIT.\n\n")
	code.WriteString("package bundled\n\n")
	code.WriteString("import (\n")
	code.WriteString("\"io/fs\"\n\n")
	code.WriteString("_ \"embed\"\n")
	code.WriteString(")\n\n")

	code.WriteString("var (\n")
	for i, lib := range libs {
		varName := libVarNames[i]
		code.WriteString("//go:embed libs/" + lib.target + "\n")
		code.WriteString("" + varName + " string\n")
	}
	code.WriteString(")\n\n")

	code.WriteString("var embeddedContents = map[string]string{\n")
	for i, lib := range libs {
		varName := libVarNames[i]
		code.WriteString("\t\"libs/" + lib.target + "\": " + varName + ",\n")
	}
	code.WriteString("}\n\n")

	code.WriteString("var libsEntries = []fs.DirEntry{\n")
	for i, lib := range libs {
		varName := libVarNames[i]
		fmt.Fprintf(&code, "\t&fileInfo{name: %q, size: int64(len(%s))},\n", lib.target, varName)
	}
	code.WriteString("}\n")

	writeCode("embed_generated.go", code.Bytes())
}

func readLibs() []lib {
	type libsMeta struct {
		libs  []string
		paths map[string]string
	}

	libsFile := tspath.NormalizeSlashes(filepath.Join(libInputDir, "libs.json"))

	b, err := os.ReadFile(libsFile)
	if err != nil {
		log.Fatalf("failed to open libs.json: %v", err)
	}

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: libsFile,
		Path:     tspath.Path(libsFile),
	}, string(b), core.ScriptKindJSON)
	diags := sourceFile.Diagnostics()

	if len(diags) > 0 {
		for _, diag := range diags {
			log.Printf("%s", diag.Message())
		}
		log.Fatalf("failed to parse libs.json")
	}

	paths := make(map[string]string)
	var libNames []string

	props := sourceFile.Statements.Nodes[0].
		AsExpressionStatement().
		Expression.
		AsObjectLiteralExpression().
		Properties.
		Nodes

	for _, prop := range props {
		assign := prop.AsPropertyAssignment()
		name := assign.Name().Text()
		switch name {
		case "libs":
			for _, lib := range assign.Initializer.AsArrayLiteralExpression().Elements.Nodes {
				libNames = append(libNames, lib.AsStringLiteral().Text)
			}
		case "paths":
			for _, path := range assign.Initializer.AsObjectLiteralExpression().Properties.Nodes {
				prop := path.AsPropertyAssignment()
				key := prop.Name().Text()
				value := prop.Initializer.AsStringLiteral().Text
				paths[key] = value
			}
		default:
			log.Fatalf("unexpected property: %s", name)
		}
	}

	var libs []lib
	for _, libName := range libNames {
		sources := []string{"header.d.ts", libName + ".d.ts"}
		var target string
		if path, ok := paths[libName]; ok {
			target = path
		} else {
			target = "lib." + libName + ".d.ts"
		}
		libs = append(libs, lib{target: target, sources: sources})
	}

	slices.SortFunc(libs, func(a lib, b lib) int {
		return strings.Compare(a.target, b.target)
	})

	return libs
}

func readCopyright() []byte {
	b, err := os.ReadFile(copyrightNotice)
	if err != nil {
		log.Fatalf("failed to read copyright notice: %v", err)
	}
	return removeCRLF(b)
}

func removeCRLF(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
}

func writeCode(filename string, code []byte) {
	formatted, err := format.Source(code)
	if err != nil {
		log.Fatalf("failed to format source: %v", err)
	}

	if err := os.WriteFile(filename, formatted, 0o644); err != nil {
		log.Fatalf("failed to write %s: %v", filename, err)
	}
}
