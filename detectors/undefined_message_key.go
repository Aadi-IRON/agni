package detectors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var targetMaps = []string{
	"Messages",
	"SuccessMessages",
	"FailMessages",
	"Message",
	"BugsnagMessages",
}

func DetectUnDefinedMessageKeys(filePath string) {
	if filePath == "" {
		fmt.Println("❌ Please enter a valid project folder name.")
		return
	}
	usedKeys := make(map[string]struct{})
	definedKeys := make(map[string]struct{})

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Extract defined keys only from messages.go
		if path == filePath+"config/message.go" || path == filePath+"config/Message.go" || path == filePath+"config/messages.go" || path == filePath+"config/Messages.go" {
			collectDefinedKeys(path, definedKeys)
		} else {
			collectUsedKeys(path, usedKeys)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking files: %v\n", err)
		return
	}
	fmt.Println()
	// Compare and report
	fmt.Println("🔍 Missing Keys (Used but not defined):")
	missing := 0
	for key := range usedKeys {
		if _, found := definedKeys[key]; !found {
			fmt.Println("-", key)
			missing++
		}
	}

	if missing == 0 {
		fmt.Println("🎉 No missing keys! Everything is defined.")
	}
}

func collectUsedKeys(filePath string, usedKeys map[string]struct{}) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", filePath, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		indexExpr, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		selector, ok := indexExpr.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := selector.X.(*ast.Ident)
		if !ok || ident.Name != "config" {
			return true
		}

		for _, mapName := range targetMaps {
			if selector.Sel.Name == mapName {
				keyLit, ok := indexExpr.Index.(*ast.BasicLit)
				if ok && keyLit.Kind == token.STRING {
					key := strings.Trim(keyLit.Value, `"`)
					usedKeys[key] = struct{}{}
				}
			}
		}

		return true
	})
}

func collectDefinedKeys(filePath string, definedKeys map[string]struct{}) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", filePath, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		vspec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}

		for _, value := range vspec.Values {
			compLit, ok := value.(*ast.CompositeLit)
			if !ok {
				continue
			}

			_, ok = compLit.Type.(*ast.MapType)
			if !ok {
				continue
			}

			for _, elt := range compLit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				keyLit, ok := kv.Key.(*ast.BasicLit)
				if ok && keyLit.Kind == token.STRING {
					key := strings.Trim(keyLit.Value, `"`)
					definedKeys[key] = struct{}{}
				}
			}
		}

		return true
	})
}
