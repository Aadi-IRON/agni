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

// Detects unused params throughout the project.
func DetectUnusedParams(filePath string) {
	if filePath == "" {
		fmt.Println("Please pass a valid directory path. ")
		return
	}
	if err := ProcessDirectory(filePath); err != nil {
		fmt.Println("Error occurred :", err)
	}
	fmt.Println(config.Reset + "----------------------------")
}

// ProcessDirectory processes all .go files in the specified directory.
func ProcessDirectory(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %v", err)
		}
		// Process only Go files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			if err := CheckUnusedVars(path); err != nil {
				fmt.Printf("Error processing file '%s': %v\n", path, err)
			}
		}
		return nil
	})
	return err
}
func GetNode(fileName string) (node *ast.File, err error) {
	// Open the Go file
	file, err := os.Open(fileName)
	if err != nil {
		return node, err
	}
	defer file.Close()
	// Create a new scanner/tokenizer for the file
	fileSet := token.NewFileSet()
	// Parse the file into an AST (Abstract Syntax Tree)
	node, err = parser.ParseFile(fileSet, fileName, file, parser.AllErrors)
	if err != nil {
		return node, err
	}
	return node, err
}

// Checks for unused variables in a Go file, considering arguments in function calls.
func CheckUnusedVars(fileName string) error {
	node, err := GetNode(fileName)
	if err != nil {
		return fmt.Errorf("error while opening/parsing file: %v", err)
	}
	// Analyze function declarations in the file
	ast.Inspect(node, func(astNode ast.Node) bool {
		if function, ok := astNode.(*ast.FuncDecl); ok {
			AnalyzeFunc(function, fileName)
		}
		return true
	})
	return nil
}

// Inspects a function's parameters and body for unused variables.
func AnalyzeFunc(function *ast.FuncDecl, fileName string) {
	varUsed := InitializeVarUsage(function.Type.Params)
	// Mark variables as used if they appear in the function body
	MarkUsedVars(function.Body, varUsed)
	// Collect and print unused variables
	if unusedVars := GetUnusedVars(varUsed); len(unusedVars) > 0 {
		fmt.Printf(config.Yellow+"File '%s': Function '%s' has unused variables: %s\n", fileName, function.Name.Name, strings.Join(unusedVars, ", "))
	}
}

// Initializes a map to track variable usage for function parameters.
func InitializeVarUsage(params *ast.FieldList) map[string]bool {
	varUsed := make(map[string]bool)
	if params == nil {
		return varUsed
	}
	for _, param := range params.List {
		for _, paramName := range param.Names {
			varUsed[paramName.Name] = false
		}
	}
	return varUsed
}

// Inspects a function body to mark variables as used.
func MarkUsedVars(body *ast.BlockStmt, varUsed map[string]bool) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(astNode ast.Node) bool {
		switch node := astNode.(type) {
		case *ast.Ident:
			// Mark identifiers as used
			if _, exists := varUsed[node.Name]; exists {
				varUsed[node.Name] = true
			}
		case *ast.CallExpr:
			// Check arguments in function calls
			MarkArgsAsUsed(node.Args, varUsed)
		}
		return true
	})
}

// Marks variables used as arguments in function calls.
func MarkArgsAsUsed(args []ast.Expr, varUsed map[string]bool) {
	for _, arg := range args {
		if ident, ok := arg.(*ast.Ident); ok {
			if _, exists := varUsed[ident.Name]; exists {
				varUsed[ident.Name] = true
			}
		}
	}
}

// Collects variable names that were never marked as used.
func GetUnusedVars(varUsed map[string]bool) []string {
	unusedVars := []string{}
	for varName, used := range varUsed {
		if !used {
			unusedVars = append(unusedVars, varName)
		}
	}
	return unusedVars
}
