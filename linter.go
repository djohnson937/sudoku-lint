package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Finding is a single problem found in a board. Col is 0 when the problem
// applies to a whole row rather than one cell.
type Finding struct {
	Line    int
	Col     int
	Message string
}

func (f Finding) String() string {
	if f.Col == 0 {
		return fmt.Sprintf("line %d: %s", f.Line, f.Message)
	}
	return fmt.Sprintf("line %d:%d: %s", f.Line, f.Col, f.Message)
}

type cell struct {
	value int // 0 means empty
	line  int
	col   int
}

// Lint reads a sudoku board and returns every structural problem it finds:
// wrong row length or row count, invalid characters, and duplicate values
// within a row, column, or 3x3 box. Blank lines and lines starting with '#'
// are treated as comments and skipped, so line numbers refer to the file,
// not the board. When strict is true, those comment and blank lines are
// reported as findings instead of being skipped silently, for callers (like
// a pre-commit hook) that want board files free of stray formatting.
func Lint(r io.Reader, strict bool) ([]Finding, error) {
	var findings []Finding
	var grid [9][9]cell
	rowIdx := 0
	lineNum := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			if strict {
				findings = append(findings, Finding{
					Line:    lineNum,
					Message: "blank line not allowed in strict mode",
				})
			}
			continue
		}
		if strings.HasPrefix(text, "#") {
			if strict {
				findings = append(findings, Finding{
					Line:    lineNum,
					Message: "comment line not allowed in strict mode",
				})
			}
			continue
		}

		if rowIdx >= 9 {
			findings = append(findings, Finding{
				Line:    lineNum,
				Message: "unexpected extra row, board must have exactly 9 rows",
			})
			continue
		}

		if len(text) != 9 {
			findings = append(findings, Finding{
				Line:    lineNum,
				Message: fmt.Sprintf("row has %d characters, want 9", len(text)),
			})
		}

		for i := 0; i < len(text) && i < 9; i++ {
			ch := text[i]
			switch {
			case ch == '.':
				grid[rowIdx][i] = cell{value: 0, line: lineNum, col: i + 1}
			case ch >= '1' && ch <= '9':
				grid[rowIdx][i] = cell{value: int(ch - '0'), line: lineNum, col: i + 1}
			default:
				findings = append(findings, Finding{
					Line:    lineNum,
					Col:     i + 1,
					Message: fmt.Sprintf("invalid character %q, want a digit 1-9 or '.'", ch),
				})
			}
		}
		rowIdx++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if rowIdx < 9 {
		findings = append(findings, Finding{
			Line:    lineNum,
			Message: fmt.Sprintf("board has %d rows, want 9", rowIdx),
		})
	}

	findings = append(findings, checkGroups(grid)...)
	return findings, nil
}

// checkGroups runs the duplicate check across every row, column, and
// 3x3 box in the grid.
func checkGroups(grid [9][9]cell) []Finding {
	var findings []Finding

	for r := 0; r < 9; r++ {
		findings = append(findings, findDuplicates(grid[r][:], "row", r+1)...)
	}

	for c := 0; c < 9; c++ {
		var col [9]cell
		for r := 0; r < 9; r++ {
			col[r] = grid[r][c]
		}
		findings = append(findings, findDuplicates(col[:], "column", c+1)...)
	}

	box := 0
	for br := 0; br < 9; br += 3 {
		for bc := 0; bc < 9; bc += 3 {
			box++
			var cells [9]cell
			i := 0
			for r := br; r < br+3; r++ {
				for c := bc; c < bc+3; c++ {
					cells[i] = grid[r][c]
					i++
				}
			}
			findings = append(findings, findDuplicates(cells[:], "box", box)...)
		}
	}

	return findings
}

// findDuplicates reports every cell in cells whose value repeats one seen
// earlier in the same group, pointing back at the first occurrence.
func findDuplicates(cells []cell, kind string, index int) []Finding {
	seen := make(map[int]cell)
	var findings []Finding
	for _, c := range cells {
		if c.value == 0 {
			continue
		}
		if first, ok := seen[c.value]; ok {
			findings = append(findings, Finding{
				Line: c.line,
				Col:  c.col,
				Message: fmt.Sprintf(
					"duplicate value %d in %s %d (first seen at line %d, column %d)",
					c.value, kind, index, first.line, first.col,
				),
			})
			continue
		}
		seen[c.value] = c
	}
	return findings
}
