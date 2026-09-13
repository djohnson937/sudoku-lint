# sudoku-lint

Sudoku boards that come from hand editing, a generator script, or a copy-paste
out of some other tool break in small, dumb ways: a row that's one character
short, a stray letter where a digit belongs, a duplicate value that makes the
puzzle unsolvable before you've even started. Those mistakes are easy to miss
by eye in a 9x9 grid of text. sudoku-lint checks a board file for exactly
those structural problems and prints the line (and column, when it's a single
cell at fault) so you can go straight to the mistake.

It does not solve puzzles or check that a puzzle has a unique solution. It
only checks that the board is well-formed: right shape, valid characters, no
value repeated in a row, column, or 3x3 box.

## Board format

A board is 9 lines of 9 characters each. Use a digit `1`-`9` for a filled
cell and `.` for an empty one. Blank lines and lines starting with `#` are
ignored, so you can add a header or comment above the grid.

```
# example puzzle
53..7....
6..195...
.98....6.
8...6...3
4..8.3..1
7...2...6
.6....28.
...419..5
....8..79
```

## Usage

```
go run . examples/valid.sudoku
```

A well-formed board prints nothing and exits 0.

```
go run . examples/invalid.sudoku
```

```
examples/invalid.sudoku:line 2:9: duplicate value 7 in row 1 (first seen at line 2, column 5)
examples/invalid.sudoku:line 4: row has 8 characters, want 9
examples/invalid.sudoku:line 6:1: invalid character 'X', want a digit 1-9 or '.'
```

Exit code is 1 when there are findings, 2 on a usage or I/O error, 0 when the
board is clean.

Pass `-json` to get findings as a JSON array instead, one object per
finding with `file`, `line`, `col` (omitted when the finding applies to a
whole row), and `message`:

```
go run . -json examples/invalid.sudoku
```

## Building

```
go build -o sudoku-lint .
./sudoku-lint path/to/board.sudoku
```

Standard library only, no dependencies to fetch.
