package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var copyFlag bool
	flag.BoolVar(&copyFlag, "c", false, "Copy source file to destination file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c] <source> <destination>\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := args[0]
	destFile := args[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Source file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}

	if copyFlag {
		if err := copyFile(sourceFile, destFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully copied '%s' to '%s'\n", sourceFile, destFile)
	} else {
		if err := copyFile(sourceFile, destFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully copied '%s' to '%s'\n", sourceFile, destFile)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}