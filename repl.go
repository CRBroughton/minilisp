package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

func handleREPLCommand(cmd string, env *Env) {
	switch cmd {
	case ":help":
		printHelp()
	case ":quit", ":q":
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s (try :help)\n", cmd)
	}
}

func printHelp() {
	fmt.Print(`
Available commands:
  :help     - Show this help
  :quit     - Exit REPL (or press Ctrl+D)

Keyboard shortcuts:
  Ctrl+A    - Move to beginning of line
  Ctrl+E    - Move to end of line
  Ctrl+U    - Delete from cursor to beginning
  Ctrl+K    - Delete from cursor to end
  Ctrl+R    - Search command history
  ↑/↓       - Navigate command history
`)
}

func startREPL(env *Env) {
	homeDIR, _ := os.UserHomeDir()
	historyFile := filepath.Join(homeDIR, "minilisp_history")

	repl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       historyFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}
	defer repl.Close()

	fmt.Println("MiniLisp - Type :help for commands")

	for {
		line, err := repl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, ":") {
			handleREPLCommand(line, env)
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Error: %v\n", r)
				}
			}()

			expr := readStr(line)
			if expr == nilExpr {
				return
			}
			result := eval(expr, env)
			fmt.Println("=>", printExpr(result))
		}()
	}
}
