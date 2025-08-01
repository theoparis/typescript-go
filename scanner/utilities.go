package scanner

import (
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/ast"
	"github.com/microsoft/typescript-go/core"
)

func tokenIsIdentifierOrKeyword(token ast.Kind) bool {
	return token >= ast.KindIdentifier
}

func IdentifierToKeywordKind(node *ast.Identifier) ast.Kind {
	return textToKeyword[node.Text]
}

func GetSourceTextOfNodeFromSourceFile(sourceFile *ast.SourceFile, node *ast.Node, includeTrivia bool) string {
	return GetTextOfNodeFromSourceText(sourceFile.Text(), node, includeTrivia)
}

func GetTextOfNodeFromSourceText(sourceText string, node *ast.Node, includeTrivia bool) string {
	if ast.NodeIsMissing(node) {
		return ""
	}
	pos := node.Pos()
	if !includeTrivia {
		pos = SkipTrivia(sourceText, pos)
	}
	text := sourceText[pos:node.End()]
	// if (isJSDocTypeExpressionOrChild(node)) {
	//     // strip space + asterisk at line start
	//     text = text.split(/\r\n|\n|\r/).map(line => line.replace(/^\s*\*/, "").trimStart()).join("\n");
	// }
	return text
}

func GetTextOfNode(node *ast.Node) string {
	return GetSourceTextOfNodeFromSourceFile(ast.GetSourceFileOfNode(node), node, false /*includeTrivia*/)
}

func DeclarationNameToString(name *ast.Node) string {
	if name == nil || name.Pos() == name.End() {
		return "(Missing)"
	}
	return GetTextOfNode(name)
}

func IsIdentifierText(name string, languageVariant core.LanguageVariant) bool {
	ch, size := utf8.DecodeRuneInString(name)
	if !IsIdentifierStart(ch) {
		return false
	}
	for i := size; i < len(name); {
		ch, size = utf8.DecodeRuneInString(name[i:])
		if !IsIdentifierPartEx(ch, languageVariant) {
			return false
		}
		i += size
	}
	return true
}

func IsIntrinsicJsxName(name string) bool {
	return len(name) != 0 && (name[0] >= 'a' && name[0] <= 'z' || strings.ContainsRune(name, '-'))
}
