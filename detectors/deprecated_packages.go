package detectors

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aadi-IRON/agni/config"
)

// DeprecatedPackage represents a deprecated package with its details
type DeprecatedPackage struct {
	Name        string
	Description string
	Alternative string
	Since       string
}

// List of deprecated packages to check
var deprecatedPackages = []DeprecatedPackage{
	{
		Name:        "golang.org/x/crypto/ssh/terminal",
		Description: "Deprecated: use golang.org/x/term instead",
		Alternative: "golang.org/x/term",
		Since:       "Go 1.19",
	},
	{
		Name:        "io/ioutil",
		Description: "Deprecated: use io and os packages instead",
		Alternative: "io, os",
		Since:       "Go 1.16",
	},
	{
		Name:        "golang.org/x/net/context",
		Description: "Deprecated: use context package instead",
		Alternative: "context",
		Since:       "Go 1.7",
	},
}

// DetectDeprecatedPackages scans for deprecated package imports
func DetectDeprecatedPackages(path string) {
	fmt.Println(config.CreateCompactBoxHeader("DEPRECATED PACKAGES", config.BoldRed))
	fmt.Println()
	fmt.Println(config.BoldYellow + "🔍 Scanning for deprecated package imports:")
	fmt.Println()

	if path == "" {
		fmt.Println("❌ Please enter a valid project folder name.")
		return
	}

	var foundDeprecated []struct {
		filePath string
		line     int
		pkg      DeprecatedPackage
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process Go files
		if info.IsDir() || !strings.HasSuffix(filePath, ".go") {
			return nil
		}

		// Skip test files if needed
		if strings.HasSuffix(filePath, "_test.go") {
			return nil
		}

		deprecated := checkFileForDeprecatedPackages(filePath)
		foundDeprecated = append(foundDeprecated, deprecated...)

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking files: %v\n", err)
		return
	}

	// Report results
	if len(foundDeprecated) == 0 {
		fmt.Println(config.BoldGreen + "✅ No deprecated packages found!")
		return
	}

	fmt.Println(config.BoldYellow + "⚠️  Found deprecated packages:")
	fmt.Println()

	for _, item := range foundDeprecated {
		fmt.Printf(config.Red+"📁 %s:%d"+config.Reset+" - "+config.BoldYellow+"%s"+config.Reset+"\n",
			item.filePath, item.line, item.pkg.Name)
		fmt.Printf("   "+config.Cyan+"Description: %s"+config.Reset+"\n", item.pkg.Description)
		fmt.Printf("   "+config.Green+"Alternative: %s"+config.Reset+"\n", item.pkg.Alternative)
		fmt.Printf("   "+config.Purple+"Since: %s"+config.Reset+"\n", item.pkg.Since)
		fmt.Println()
	}
}

// checkFileForDeprecatedPackages checks a single file for deprecated imports
func checkFileForDeprecatedPackages(filePath string) []struct {
	filePath string
	line     int
	pkg      DeprecatedPackage
} {
	var found []struct {
		filePath string
		line     int
		pkg      DeprecatedPackage
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return found
	}

	// Check each import
	for _, importSpec := range node.Imports {
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		// Check if this import is deprecated
		for _, deprecated := range deprecatedPackages {
			if importPath == deprecated.Name {
				position := fset.Position(importSpec.Pos())
				found = append(found, struct {
					filePath string
					line     int
					pkg      DeprecatedPackage
				}{
					filePath: filePath,
					line:     position.Line,
					pkg:      deprecated,
				})
			}
		}
	}

	return found
}
