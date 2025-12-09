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

var targetMaps = []string{
	"Messages",
	"SuccessMessages",
	"FailMessages",
	"Message",
	"BugsnagMessages",
}

func isTargetMap(name string) bool {
	for _, m := range targetMaps {
		if name == m {
			return true
		}
	}
	return false
}

func DetectUnDefinedMessageKeys(filePath string) {
	fmt.Println(config.CreateCompactBoxHeader("UNDEFINED MESSAGE KEYS", config.BoldPurple))
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
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		base := strings.ToLower(filepath.Base(path))
		if base == "message.go" || base == "messages.go" {
			CollectDefinedKeys(path, definedKeys)
		} else {
			CollectUsedKeys(path, usedKeys)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking files: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println(config.BoldYellow + "🔍 Missing Keys (Used but not defined):")
	fmt.Println()
	missing := 0
	for key := range usedKeys {
		if _, found := definedKeys[key]; !found {
			fmt.Println("-", key)
			missing++
		}
	}
	if missing == 0 {
		fmt.Println(config.Cyan + "🎉 No missing keys! Everything is defined.")
		return
	}
	fmt.Println()
}

func CollectUsedKeys(filePath string, usedKeys map[string]struct{}) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", filePath, err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		idx, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		switch x := idx.X.(type) {
		case *ast.SelectorExpr:
			if isTargetMap(x.Sel.Name) {
				if lit, ok := idx.Index.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					usedKeys[strings.Trim(lit.Value, `"`)] = struct{}{}
				}
			}
		case *ast.Ident:
			// In-package use: Messages["key"]
			if isTargetMap(x.Name) {
				if lit, ok := idx.Index.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					usedKeys[strings.Trim(lit.Value, `"`)] = struct{}{}
				}
			}
		}
		return true
	})
}

func CollectDefinedKeys(filePath string, definedKeys map[string]struct{}) {
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

		// Only consider target map vars
		target := false
		for _, name := range vspec.Names {
			if isTargetMap(name.Name) {
				target = true
				break
			}
		}
		if !target {
			return true
		}

		for _, value := range vspec.Values {
			cl, ok := value.(*ast.CompositeLit)
			if !ok {
				continue
			}

			// Accept map literals, named map types, selector types, or elided types
			switch cl.Type.(type) {
			case *ast.MapType, *ast.Ident, *ast.SelectorExpr, nil:
				// ok
			default:
				// still walk keys; being permissive avoids missing aliased types
			}

			for _, elt := range cl.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if lit, ok := kv.Key.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						definedKeys[strings.Trim(lit.Value, `"`)] = struct{}{}
					}
				}
			}
		}
		return true
	})
}
