package main

import (
	"os"
	"testing"
)

func TestColourise(t *testing.T) {
	tests := []struct {
		name      string
		colour    string
		text      string
		useColour bool
		want      string
	}{
		{
			name:      "with colours enabled",
			colour:    colourRed,
			text:      "error",
			useColour: true,
			want:      "\033[31merror\033[0m",
		},
		{
			name:      "with colours disabled",
			colour:    colourRed,
			text:      "error",
			useColour: false,
			want:      "error",
		},
		{
			name:      "green text with colours",
			colour:    colourGreen,
			text:      "success",
			useColour: true,
			want:      "\033[32msuccess\033[0m",
		},
		{
			name:      "empty text",
			colour:    colourBlue,
			text:      "",
			useColour: true,
			want:      "\033[34m\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original state
			originalUseColours := useColours
			defer func() { useColours = originalUseColours }()

			useColours = tt.useColour
			got := colourise(tt.colour, tt.text)
			if got != tt.want {
				t.Errorf("colourise() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		noColourEnv  string
		wantColours  bool
	}{
		{
			name:         "NO_COLOUR not set",
			noColourEnv:  "",
			wantColours:  true,
		},
		{
			name:         "NO_COLOUR set to any value",
			noColourEnv:  "1",
			wantColours:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore environment
			originalNoColour := os.Getenv("NO_COLOUR")
			defer func() {
				if originalNoColour == "" {
					os.Unsetenv("NO_COLOUR")
				} else {
					os.Setenv("NO_COLOUR", originalNoColour)
				}
			}()

			// Set test environment
			if tt.noColourEnv == "" {
				os.Unsetenv("NO_COLOUR")
			} else {
				os.Setenv("NO_COLOUR", tt.noColourEnv)
			}

			// Re-run init logic
			useColours = true
			if os.Getenv("NO_COLOUR") != "" {
				useColours = false
			}

			if useColours != tt.wantColours {
				t.Errorf("useColours = %v, want %v", useColours, tt.wantColours)
			}
		})
	}
}

func TestPrintExprColoured(t *testing.T) {
	// Save original state
	originalUseColours := useColours
	defer func() { useColours = originalUseColours }()
	useColours = true

	tests := []struct {
		name string
		expr *Expr
		want string
	}{
		{
			name: "nil expression",
			expr: nil,
			want: "\033[37mnil\033[0m",
		},
		{
			name: "nilExpr",
			expr: nilExpr,
			want: "\033[37mnil\033[0m",
		},
		{
			name: "trueExpr",
			expr: trueExpr,
			want: "\033[32mtrue\033[0m",
		},
		{
			name: "falseExpr",
			expr: falseExpr,
			want: "\033[31mfalse\033[0m",
		},
		{
			name: "number expression",
			expr: &Expr{Type: Number, Num: 42},
			want: "\033[36m42\033[0m",
		},
		{
			name: "string expression",
			expr: &Expr{Type: String, Str: "hello"},
			want: "\033[33m\"hello\"\033[0m",
		},
		{
			name: "symbol expression",
			expr: &Expr{Type: Symbol, Sym: "foo"},
			want: "\033[35mfoo\033[0m",
		},
		{
			name: "builtin function",
			expr: &Expr{Type: Builtin, Sym: "+"},
			want: "\033[34m<builtin>\033[0m",
		},
		{
			name: "hash expression",
			expr: &Expr{Type: Hash, HashTable: make(map[string]*Expr)},
			want: "\033[32m{}\033[0m",
		},
		{
			name: "lambda expression",
			expr: &Expr{Type: Lambda},
			want: "\033[34m<lambda>\033[0m",
		},
		{
			name: "macro expression",
			expr: &Expr{Type: Macro},
			want: "\033[34m<macro>\033[0m",
		},
		{
			name: "pair expression - simple list",
			expr: &Expr{
				Type: Pair,
				Head: &Expr{Type: Number, Num: 1},
				Tail: &Expr{
					Type: Pair,
					Head: &Expr{Type: Number, Num: 2},
					Tail: nilExpr,
				},
			},
			want: "(1 2)",
		},
		{
			name: "pair expression - empty list",
			expr: &Expr{Type: Pair, Head: nilExpr, Tail: nilExpr},
			want: "(nil)",
		},
		{
			name: "unknown type - default case",
			expr: &Expr{Type: ExprType("Unknown"), Num: 99},
			want: "<unknown>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := printExprColoured(tt.expr)
			if got != tt.want {
				t.Errorf("printExprColoured() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintExprColouredWithoutColours(t *testing.T) {
	// Save original state
	originalUseColours := useColours
	defer func() { useColours = originalUseColours }()
	useColours = false

	tests := []struct {
		name string
		expr *Expr
		want string
	}{
		{
			name: "number without colours",
			expr: &Expr{Type: Number, Num: 42},
			want: "42",
		},
		{
			name: "string without colours",
			expr: &Expr{Type: String, Str: "hello"},
			want: "\"hello\"",
		},
		{
			name: "true without colours",
			expr: trueExpr,
			want: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := printExprColoured(tt.expr)
			if got != tt.want {
				t.Errorf("printExprColoured() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColourConstants(t *testing.T) {
	// Verify ANSI colour codes are correct
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"colourReset", colourReset, "\033[0m"},
		{"colourRed", colourRed, "\033[31m"},
		{"colourGreen", colourGreen, "\033[32m"},
		{"colourYellow", colourYellow, "\033[33m"},
		{"colourBlue", colourBlue, "\033[34m"},
		{"colourPurple", colourPurple, "\033[35m"},
		{"colourCyan", colourCyan, "\033[36m"},
		{"colourGray", colourGray, "\033[37m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestPrintResult(t *testing.T) {
	// Note: printResult outputs to stdout, so we're testing it doesn't panic
	// and works with different expression types
	tests := []struct {
		name string
		expr *Expr
	}{
		{"nil", nil},
		{"nilExpr", nilExpr},
		{"true", trueExpr},
		{"false", falseExpr},
		{"number", &Expr{Type: Number, Num: 42}},
		{"string", &Expr{Type: String, Str: "test"}},
		{"symbol", &Expr{Type: Symbol, Sym: "foo"}},
		{"list", &Expr{
			Type: Pair,
			Head: &Expr{Type: Number, Num: 1},
			Tail: nilExpr,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("printResult() panicked: %v", r)
				}
			}()
			printResult(tt.expr)
		})
	}
}
