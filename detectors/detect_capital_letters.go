package detectors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Aadi-IRON/agni/config"
)

func DetectCapitalVars(path string) {
	fmt.Println(config.CreateCompactBoxHeader("CAPITAL LETTERS", config.BoldBlue))
	if path == "" {
		fmt.Println("Please pass a valid directory name.", path)
		return
	}
	fmt.Println()
	fmt.Printf(config.BoldBlue + "🔍 Detecting variables and function parameters with capital letters :---")
	fmt.Println()
	functionSet := token.NewFileSet()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() ||
			strings.HasSuffix(path, "_test.go") ||
			!strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "const.go") ||
			strings.HasSuffix(path, "Const.go") ||
			strings.HasSuffix(path, "message.go") ||
			strings.HasSuffix(path, "Message.go") ||
			strings.HasSuffix(path, "Messages.go") ||
			strings.HasSuffix(path, "messages.go") {
			return nil
		}
		CheckFile(path, functionSet)
		return nil
	})
	fmt.Println()
	fmt.Println(config.Green + "FINISHED")
	if err != nil {
		log.Fatal(err)
	}
}

// Checks a single Go file for capitalized variable/parameter names
func CheckFile(path string, fset *token.FileSet) {
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println("Error parsing:", path, err)
		return
	}

	// Track if we're inside a function
	var insideFunction bool

	ast.Inspect(node, func(astNode ast.Node) bool {
		switch stmt := astNode.(type) {

		case *ast.FuncDecl:
			// We're entering a function
			insideFunction = true

			// Function parameters
			if stmt.Type.Params != nil {
				for _, param := range stmt.Type.Params.List {
					for _, name := range param.Names {
						if IsCapitalized(name.Name) {
							position := fset.Position(name.Pos())
							log.Printf(config.Yellow+"Capitalized function parameter"+config.BoldRed+" '%s'"+config.Yellow+" at %s\n", name.Name, position)
						}
					}
				}
			}
			// Named return parameters
			if stmt.Type.Results != nil {
				for _, result := range stmt.Type.Results.List {
					for _, name := range result.Names {
						if IsCapitalized(name.Name) {
							position := fset.Position(name.Pos())
							log.Printf(config.Yellow+"Capitalized named return variable"+config.BoldRed+" '%s'"+config.Yellow+" at %s\n", name.Name, position)
						}
					}
				}
			}

		case *ast.FuncLit:
			// We're entering an anonymous function
			insideFunction = true

		case *ast.GenDecl:
			if stmt.Tok == token.CONST {
				// Skip all constants
				return true
			}
			// Handle var declarations based on context
			if stmt.Tok == token.VAR {
				if insideFunction {
					// Local var declarations inside functions (should be checked)
					for _, spec := range stmt.Specs {
						if valSpec, ok := spec.(*ast.ValueSpec); ok {
							for _, name := range valSpec.Names {
								if IsCapitalized(name.Name) {
									position := fset.Position(name.Pos())
									log.Printf(config.Yellow+"Capitalized local variable"+config.BoldRed+" '%s'"+config.Yellow+" at %s\n", name.Name, position)
								}
							}
						}
					}
				} else {
					// Global variables (should be skipped - they should be capitalized)
					return true
				}
			}

		// MyVar := ... (local variables)
		case *ast.AssignStmt:
			if stmt.Tok.String() == ":=" {
				for _, lhs := range stmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && IsCapitalized(ident.Name) {
						position := fset.Position(ident.Pos())
						log.Printf(config.Yellow+"Capitalized short variable"+config.BoldRed+" '%s'"+config.Yellow+" at %s\n", ident.Name, position)
					}
				}
			}
		}
		return true
	})
}

// Checks if a name starts with a capital letter
func IsCapitalized(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}
