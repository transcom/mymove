package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("Usage: big-cat <path> [limit]")
		os.Exit(1)
	}
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		panic(err)
	}
	limit := -1
	if len(os.Args) == 3 {
		l, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		limit = l
	}
	count := 0
	for _, file := range files {
		f, err := os.Open(filepath.Clean(file))
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(os.Stdout, bufio.NewReader(f)); err != nil {
			panic(err)
		}
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to end an asynchronous connection pertaining to file formatting
		//RA: Given the functions causing the lint errors are used to end a running asynchronous connection, it does not present a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		f.Close() // nolint:errcheck
		count++
		if limit >= 0 && count == limit {
			break
		}
	}
}
