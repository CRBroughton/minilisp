package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

func handleREPLCommand(cmd string, env *Env, rl *readline.Instance) {
	switch cmd {
	case ":help":
		printHelp()
	case ":env", ":e":
		showEnvironment(env)
	case ":history", ":h":
		showHistory()
	case ":clear", ":c":
		readline.ClearScreen(rl)
	case ":quit", ":q":
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s (try :help)\n", cmd)
	}
}

func printHelp() {
	fmt.Print(`
Available commands:
  :help, :h       - Show this help
  :env,  :e       - Show environment bindings
  :history        - Show command history
  :clear, :c      - Clear screen
  :quit, :q       - Exit REPL (or press Ctrl+D)

Keyboard shortcuts:
  Ctrl+A    - Move to beginning of line
  Ctrl+E    - Move to end of line
  Ctrl+U    - Delete from cursor to beginning
  Ctrl+K    - Delete from cursor to end
  Ctrl+R    - Search command history
  ↑/↓       - Navigate command history
  Tab       - Autocomplete function names

`)
}

func showEnvironment(env *Env) {
	fmt.Println("\nCurrent environment:")

	// Collect and sort bindings
	names := make([]string, 0, len(env.bindings))
	for name := range env.bindings {
		names = append(names, name)
	}
	sort.Strings(names)

	// Print in columns
	for _, name := range names {
		val := env.bindings[name]
		typeStr := string(val.Type)
		fmt.Printf("  %-20s = <%s>\n", name, strings.ToLower(typeStr))
	}
	fmt.Println()
}

func showHistory() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not find home directory")
		return
	}
	historyFile := filepath.Join(homeDir, "minilisp_history")

	content, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Println("No history available")
		return
	}

	lines := strings.Split(string(content), "\n")

	// Show last 20 commands
	start := 0
	if len(lines) > 20 {
		start = len(lines) - 20
	}

	fmt.Println("\nRecent history:")
	count := 1
	for _, line := range lines[start:] {
		if line != "" {
			fmt.Printf("%3d  %s\n", count, line)
			count++
		}
	}
	fmt.Println()
}

// Check if expression has balanced parentheses
func isCompleteExpr(input string) bool {
	depth := 0
	inString := false
	inComment := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		// Handle comments
		if ch == ';' {
			inComment = true
		}
		if ch == '\n' {
			inComment = false
		}
		if inComment {
			continue
		}

		// Handle strings
		if ch == '"' {
			inString = !inString
		}
		if inString {
			continue
		}

		// Count parentheses
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		}
	}

	return depth == 0
}

func createCompleter(env *Env) *readline.PrefixCompleter {
	// Get all defined symbols from environment
	var items []readline.PrefixCompleterInterface

	// Add REPL commands
	items = append(items,
		readline.PcItem(":help"),
		readline.PcItem(":env"),
		readline.PcItem(":history"),
		readline.PcItem(":clear"),
		readline.PcItem(":quit"),
	)

	// Add builtin functions and defined symbols
	for name := range env.bindings {
		items = append(items, readline.PcItem(name))
	}

	return readline.NewPrefixCompleter(items...)
}

func startREPL(env *Env) {
	homeDIR, _ := os.UserHomeDir()
	historyFile := filepath.Join(homeDIR, "minilisp_history")

	// Create completer
	completer := createCompleter(env)

	repl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      completer,
	})

	if err != nil {
		panic(err)
	}
	defer repl.Close()

	fmt.Println("MiniLisp - Type :help for commands")

	var buffer []string // Accumulate multi-line input

	for {
		// Choose prompt based on whether we're continuing
		prompt := "> "
		if len(buffer) > 0 {
			prompt = "... "
		}
		repl.SetPrompt(prompt)

		line, err := repl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				buffer = nil // Clear buffer on Ctrl+C
				if len(buffer) == 0 {
					break
				}
				continue
			}
			buffer = nil // Clear buffer on Ctrl+C
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		// Skip empty lines (unless in multi-line mode)
		if line == "" {
			if len(buffer) == 0 {
				continue
			}
		}

		// Skip comments
		if strings.HasPrefix(line, ";") {
			continue
		}

		// Handle REPL commands (only when not in multi-line mode)
		if strings.HasPrefix(line, ":") && len(buffer) == 0 {
			handleREPLCommand(line, env, repl)
			continue
		}

		// Accumulate input
		buffer = append(buffer, line)
		input := strings.Join(buffer, "\n")

		// Check if input is complete
		if !isCompleteExpr(input) {
			continue // Need more input
		}

		// Clear buffer and evaluate
		buffer = nil

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Error: %v\n", r)
				}
			}()

			expr := readStr(input)
			if expr == nilExpr {
				return
			}
			result := eval(expr, env)
			printResult(result)
		}()
	}
}
