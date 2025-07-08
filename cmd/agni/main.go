// cmd/agni/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Aadi-IRON/agni/detectors"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "agni",
		Short: "Agni 🔥 - Static analyzer for Go code",
	}

	var checkCmd = &cobra.Command{
		Use:   "check [flags]",
		Short: "Run all enabled Agni checks",
		Run: func(cmd *cobra.Command, args []string) {
			// ✅ STEP: Automatically use current working directory
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting working directory:", err)
				os.Exit(1)
			}
			absPath, _ := filepath.Abs(dir)
			fmt.Println("Running Agni checks in:", absPath)

			detectors.RunAll(absPath)
		},
	}

	rootCmd.AddCommand(checkCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
