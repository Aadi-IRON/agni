package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Aadi-IRON/agni/detectors"
)

func main() {
	// Optional: Allow custom directory via flag
	dirPtr := flag.String("dir", ".", "Directory to run Agni checks in")
	flag.Parse()

	absPath, err := filepath.Abs(*dirPtr)
	if err != nil {
		fmt.Println("❌ Error getting absolute path:", err)
		os.Exit(1)
	}

	fmt.Println("🔥 Running Agni checks in:", absPath)
	detectors.RunAll(absPath)
}
