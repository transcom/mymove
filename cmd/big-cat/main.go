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
		f.Close()
		count++
		if limit >= 0 && count == limit {
			break
		}
	}
}
