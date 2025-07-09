package detectors

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Detects unused messages throughout the directory.
func DetectUnusedMessages(filePath string) {
	if filePath == "" {
		fmt.Println("Please pass a valid directory path. ")
		return
	}
	// Extract all keys from the Messages map in message.go
	keys, err := ExtractKeysFromMessages(filePath + "/config/message.go")
	var oldVersion bool
	if err != nil {
		keys, err = ExtractKeysFromMessages(filePath + "/config/Message.go")
		oldVersion = true
		if err != nil {
			fmt.Println("Error extracting keys from message.go:", err)
			return
		}
	}
	// Check if each key is used in the project
	var unusedKeys []string
	for _, key := range keys {
		found, err := SearchKeyInProject(filePath, key, oldVersion)
		if err != nil {
			fmt.Println("Error searching for key in the project:", err)
			return
		}
		if !found {
			unusedKeys = append(unusedKeys, key)
		}
	}
	// Print results
	if len(unusedKeys) == 0 {
		fmt.Println("✅  All keys in messages.go file are used in the project.")
	} else {
		fmt.Println()
		fmt.Println("Unused keys in messages.go file:-> ")
		for _, key := range unusedKeys {
			fmt.Println("- ", key)
		}
	}
	fmt.Println("----------------------------")
}

// SearchKeyInProject searches for a specific key across all .go files in the project
func SearchKeyInProject(rootDir, searchKey string, oldVersion bool) (bool, error) {
	found := false

	// Walk through all files in the project
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		condition := !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(path, "message.go")
		if oldVersion {
			condition = !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(path, "Message.go")
		}
		// Process only .go files, excluding the message.go file itself
		if condition {
			fileFound, err := SearchKeyInFile(path, searchKey)
			if err != nil {
				return err
			}
			if fileFound {
				found = true
			}
		}
		return nil
	})

	return found, err
}

// searchKeyInFile searches for a key in a specific .go file
func SearchKeyInFile(filePath, searchKey string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// Check each line for the search key
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchKey) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// ExtractKeysFromMessages extracts all keys from the Messages map in message.go
func ExtractKeysFromMessages(messageFilePath string) ([]string, error) {
	var keys []string

	// Read and parse the message.go file
	functionSet := token.NewFileSet()
	node, err := parser.ParseFile(functionSet, messageFilePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Println("1")
		return nil, err
	}

	// Traverse the AST to extract keys from the Messages map
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, value := range valueSpec.Values {
				compLit, ok := value.(*ast.CompositeLit)
				if !ok {
					continue
				}

				// Check if the map is named "Messages"
				for _, elt := range compLit.Elts {
					kvExpr, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					key := strings.Trim(kvExpr.Key.(*ast.BasicLit).Value, `"`)
					keys = append(keys, key)
				}
			}
		}
	}
	return keys, nil
}
