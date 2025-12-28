package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/internal/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/internal/execute"
)

func main() {
	autoCompleter := autocomplete.NewCodecraftersAutoCompleter()

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
	})

	rl.Config.SetListener(autoCompleter.CompleteExecutable)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		command, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		stdoutDest := os.Stdout
		stderrDest := os.Stderr

		if strings.Contains(command, "|") {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Fprintln(stdoutDest, err)
				os.Exit(1)
			}
			before, after, _ := strings.Cut(command, "|")
			firstCmd := strings.TrimSpace(before)
			secondCmd := strings.TrimSpace(after)

			cmd1 := execute.ExecuteUserInput(firstCmd, false, nil, w, w)
			w.Close()

			execute.ExecuteUserInput(secondCmd, true, r, stdoutDest, stderrDest)
			r.Close()

			if cmd1 != nil {
				cmd1.Wait()
			}
		} else {
			execute.ExecuteUserInput(command, true, nil, stdoutDest, stderrDest)
		}
	}
}
