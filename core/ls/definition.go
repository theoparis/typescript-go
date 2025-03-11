package ls

import (
	"github.com/microsoft/typescript-go/core/ast"
	"github.com/microsoft/typescript-go/core/astnav"
	"github.com/microsoft/typescript-go/core/core"
	"github.com/microsoft/typescript-go/core/scanner"
)

func (l *LanguageService) ProvideDefinitions(fileName string, position int) []Location {
	program, file := l.getProgramAndFile(fileName)
	node := astnav.GetTouchingPropertyName(file, position)
	if node.Kind == ast.KindSourceFile {
		return nil
	}

	checker := program.GetTypeChecker()
	if symbol := checker.GetSymbolAtLocation(node); symbol != nil {
		if symbol.Flags&ast.SymbolFlagsAlias != 0 {
			if resolved, ok := checker.ResolveAlias(symbol); ok {
				symbol = resolved
			}
		}

		locations := make([]Location, 0, len(symbol.Declarations))
		for _, decl := range symbol.Declarations {
			file := ast.GetSourceFileOfNode(decl)
			loc := decl.Loc
			pos := scanner.GetTokenPosOfNode(decl, file, false /*includeJsDoc*/)

			locations = append(locations, Location{
				FileName: file.FileName(),
				Range:    core.NewTextRange(pos, loc.End()),
			})
		}
		return locations
	}
	return nil
}
