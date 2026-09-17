package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

// jsonFinding is the shape printed with -json. It carries the file path
// alongside a Finding since Finding itself doesn't know what file it came
// from.
type jsonFinding struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col,omitempty"`
	Message string `json:"message"`
}

func main() {
	jsonOutput := flag.Bool("json", false, "print findings as a JSON array instead of text")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: sudoku-lint [-json] <board-file>")
		fmt.Fprintln(os.Stderr, "       sudoku-lint [-json] -    (read board from stdin)")
	}
	flag.Parse()

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := "<stdin>"
	in := os.Stdin
	if flag.NArg() == 1 && flag.Arg(0) != "-" {
		path = flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sudoku-lint: %v\n", err)
			os.Exit(2)
		}
		defer f.Close()
		in = f
	}

	findings, err := Lint(in)
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

	if *jsonOutput {
		out := make([]jsonFinding, len(findings))
		for i, fnd := range findings {
			out[i] = jsonFinding{File: path, Line: fnd.Line, Col: fnd.Col, Message: fnd.Message}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintf(os.Stderr, "sudoku-lint: %v\n", err)
			os.Exit(2)
		}
	} else {
		for _, fnd := range findings {
			fmt.Printf("%s:%s\n", path, fnd.String())
		}
	}

	if len(findings) > 0 {
		os.Exit(1)
	}
}
