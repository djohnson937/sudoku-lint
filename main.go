package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: sudoku-lint <board-file>")
		os.Exit(2)
	}

	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sudoku-lint: %v\n", err)
		os.Exit(2)
	}
	defer f.Close()

	findings, err := Lint(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sudoku-lint: %v\n", err)
		os.Exit(2)
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Col < findings[j].Col
	})

	for _, fnd := range findings {
		fmt.Printf("%s:%s\n", path, fnd.String())
	}

	if len(findings) > 0 {
		os.Exit(1)
	}
}
