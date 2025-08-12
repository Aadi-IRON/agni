package detectors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aadi-IRON/agni/config"
)

type FuncInfo struct {
	Name          string
	FilePath      string
	Line          int
	Column        int
	Package       string
	UsedElsewhere bool
}

func DetectExportedButInternalFuncs(filePath string) {
	fmt.Println(config.CreateCompactBoxHeader("EXPORTED FUNCTIONS THAT SHOULD BE UNEXPORTED", config.BoldPurple))
	fmt.Println("")
	if filePath == "" {
		fmt.Println("❌ Please enter a valid project folder name.")
		return
	}

	var exportedFuncs []FuncInfo
	fset := token.NewFileSet()

	// Pass 1: Collect all exported functions
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Failed to parse %s: %v\n", path, err)
			return nil
		}

		for _, decl := range node.Decls {
			if function, ok := decl.(*ast.FuncDecl); ok && function.Name.IsExported() {
				position := fset.Position(function.Pos())
				exportedFuncs = append(exportedFuncs, FuncInfo{
					Name:     function.Name.Name,
					FilePath: path,
					Line:     position.Line,
					Column:   position.Column,
					Package:  node.Name.Name,
				})
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking files: %v\n", err)
		return
	}

	// Pass 2: Check if exported functions are used in other packages
	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Printf("Failed to parse %s: %v\n", path, err)
			return nil
		}

		ast.Inspect(node, func(astNode ast.Node) bool {
			if ident, ok := astNode.(*ast.Ident); ok {
				for idx := range exportedFuncs {
					if ident.Name == exportedFuncs[idx].Name &&
						node.Name.Name != exportedFuncs[idx].Package {
						exportedFuncs[idx].UsedElsewhere = true
					}
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking files: %v\n", err)
		return
	}

	// Output results
	missing := 0
	for _, function := range exportedFuncs {
		if !function.UsedElsewhere {
			fmt.Printf("%s:%d:%d - %s%s%s should be unexported (used only inside package '%s')\n",
				function.FilePath, function.Line, function.Column,
				config.BoldYellow, function.Name, config.Reset, function.Package)
			missing++
		}
	}

	if missing == 0 {
		fmt.Println(config.Cyan + "🎉 No incorrectly exported functions found.")
	}
	fmt.Println("")
	fmt.Println("")
}
