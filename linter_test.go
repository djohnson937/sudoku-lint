package main

import (
	"strings"
	"testing"
)

const validBoard = `534678912
672195348
198342567
859761423
426853791
713924856
961537284
287419635
345286179
`

func TestLintValidBoard(t *testing.T) {
	findings, err := Lint(strings.NewReader(validBoard))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestLintSkipsCommentsAndBlankLines(t *testing.T) {
	board := "# header\n\n" + validBoard
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestLintShortRow(t *testing.T) {
	board := strings.Replace(validBoard, "534678912\n", "53467891\n", 1)
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "row has 8 characters, want 9") {
		t.Fatalf("expected short row finding, got %v", findings)
	}
}

func TestLintInvalidCharacter(t *testing.T) {
	board := strings.Replace(validBoard, "534678912\n", "53467891X\n", 1)
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	want := Finding{Line: 1, Col: 9, Message: `invalid character 'X', want a digit 1-9 or '.'`}
	if !containsFinding(findings, want) {
		t.Fatalf("expected %v, got %v", want, findings)
	}
}

func TestLintDuplicateInRow(t *testing.T) {
	board := strings.Replace(validBoard, "534678912\n", "534678911\n", 1)
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "duplicate value 1 in row 1") {
		t.Fatalf("expected duplicate row finding, got %v", findings)
	}
}

func TestLintDuplicateInColumn(t *testing.T) {
	// Column 1 already has a 5 on line 1; put another 5 in column 1 of line 2.
	board := strings.Replace(validBoard, "672195348\n", "572195348\n", 1)
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "duplicate value 5 in column 1") {
		t.Fatalf("expected duplicate column finding, got %v", findings)
	}
}

func TestLintDuplicateInBox(t *testing.T) {
	// Top-left box holds 5,3,4,6,7,2,1,9,8. Turn the 8 at (3,3) into a 5,
	// which collides with (1,1) in the box (this also collides within row 3,
	// since a solved row already contains every digit once).
	board := strings.Replace(validBoard, "198342567\n", "195342567\n", 1)
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "duplicate value 5 in box 1") {
		t.Fatalf("expected duplicate box finding, got %v", findings)
	}
}

func TestLintTooFewRows(t *testing.T) {
	lines := strings.Split(strings.TrimRight(validBoard, "\n"), "\n")
	board := strings.Join(lines[:8], "\n") + "\n"
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "board has 8 rows, want 9") {
		t.Fatalf("expected row count finding, got %v", findings)
	}
}

func TestLintTooManyRows(t *testing.T) {
	board := validBoard + "123456789\n"
	findings, err := Lint(strings.NewReader(board))
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !containsMessage(findings, "unexpected extra row, board must have exactly 9 rows") {
		t.Fatalf("expected extra row finding, got %v", findings)
	}
}

func TestFindingString(t *testing.T) {
	rowFinding := Finding{Line: 4, Message: "row has 8 characters, want 9"}
	if got, want := rowFinding.String(), "line 4: row has 8 characters, want 9"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	cellFinding := Finding{Line: 6, Col: 1, Message: "invalid character 'X', want a digit 1-9 or '.'"}
	if got, want := cellFinding.String(), "line 6:1: invalid character 'X', want a digit 1-9 or '.'"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func containsMessage(findings []Finding, substr string) bool {
	for _, f := range findings {
		if strings.Contains(f.Message, substr) {
			return true
		}
	}
	return false
}

func containsFinding(findings []Finding, want Finding) bool {
	for _, f := range findings {
		if f == want {
			return true
		}
	}
	return false
}
