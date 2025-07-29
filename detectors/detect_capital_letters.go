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
	if path == "" {
		fmt.Println("Please pass a valid directory name.", path)
		return
	}

	fmt.Printf(config.BoldBlue + "Detecting variables and function parameters with capital letters :---")

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
	fmt.Println(config.Green + "FINISHED")
	fmt.Println(config.Reset + "----------------------------")
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

	ast.Inspect(node, func(astNode ast.Node) bool {
		switch stmt := astNode.(type) {

		case *ast.GenDecl:
			if stmt.Tok == token.CONST {
				// Skip all constants
				return true
			}
			if stmt.Tok == token.VAR {
				// Handle var declarations
				for _, specific := range stmt.Specs {
					if valSpecific, ok := specific.(*ast.ValueSpec); ok {
						for _, name := range valSpecific.Names {
							if IsCapitalized(name.Name) {
								position := fset.Position(name.Pos())
								log.Printf("Capitalized variable '%s' at %s\n", name.Name, position)
							}
						}
					}
				}
			}
		// MyVar := ...
		case *ast.AssignStmt:
			if stmt.Tok.String() == ":=" {
				for _, lhs := range stmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && IsCapitalized(ident.Name) {
						position := fset.Position(ident.Pos())
						log.Printf("Capitalized short variable '%s' at %s\n", ident.Name, position)
					}
				}
			}
		// func FunctionName(Parameter Type) (ReturnType)
		case *ast.FuncDecl:
			// Function parameters
			if stmt.Type.Params != nil {
				for _, param := range stmt.Type.Params.List {
					for _, name := range param.Names {
						if IsCapitalized(name.Name) {
							position := fset.Position(name.Pos())
							log.Printf("Capitalized function parameter '%s' at %s\n", name.Name, position)
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
							log.Printf("Capitalized named return variable '%s' at %s\n", name.Name, position)
						}
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
