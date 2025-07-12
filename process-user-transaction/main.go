package main

import (
	"os"
	"path/filepath"
	"process-user-transaction/internal/factory"
)

func main() {
	filename := os.Args[1]
	inputDir := filepath.Join("files", "input", filename)
	outputDir := filepath.Join("files", "output", filename)

	f := factory.NewFactory()

	err := f.Start(inputDir, outputDir)
	if err != nil {
		panic(err)
	}
}
