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

	"github.com/Aadi-IRON/agni/config"
)

// Detects unused messages throughout the directory.
func DetectUnusedMessages(filePath string) {
	fmt.Println(config.CreateCompactBoxHeader("UNUSED MESSAGES", config.BoldCyan))
	fmt.Println()
	if filePath == "" {
		fmt.Println("Please pass a valid directory path. ")
		return
	}
	fmt.Println(config.BoldYellow + "🔍 Detecting unused messages (Declared but not used):")
	fmt.Println()
	// Extract all keys from the Messages map in message.go
	keys, err := ExtractKeysFromMessages(filePath + "/config/message.go")
	messageFileName := "message.go"
	if err != nil {
		keys, err = ExtractKeysFromMessages(filePath + "/config/Message.go")
		messageFileName = "Message.go"
		if err != nil {
			keys, err = ExtractKeysFromMessages(filePath + "/config/messages.go")
			messageFileName = "messages.go"
			if err != nil {
				keys, err = ExtractKeysFromMessages(filePath + "/config/Messages.go")
				messageFileName = "Messages.go"
				if err != nil {
					fmt.Println("Error occurred while extracting keys", err)
					return
				}
			}
		}
	}
	// Check if each key is used in the project
	var unusedKeys []string
	for _, key := range keys {
		found, err := SearchKeyInProject(filePath, key, messageFileName)
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
		fmt.Println(config.BoldGreen + "✅  All keys in messages.go file are used in the project.")
	} else {
		fmt.Println()
		fmt.Println(config.BoldYellow + "Unused keys in messages.go file:-> ")
		for _, key := range unusedKeys {
			fmt.Println(config.Red+"- ", key)
		}
	}
	fmt.Println()
}

// SearchKeyInProject searches for a specific key across all .go files in the project
func SearchKeyInProject(rootDir, searchKey string, messageFileName string) (bool, error) {
	found := false

	// Walk through all files in the project
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		condition := !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(path, messageFileName)
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
