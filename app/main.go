package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/internal/execute"
	"github.com/codecrafters-io/shell-starter-go/app/internal/history"
	"github.com/codecrafters-io/shell-starter-go/app/internal/listener"
)

func main() {
	history.InitializeHistory()
	listener := listener.NewListener()

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
	})

	rl.Config.SetListener(listener.Listen)
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
		history.AppendToHistory(command)

		stdoutDest := os.Stdout
		stderrDest := os.Stderr

		if strings.Contains(command, "|") {

			// split the pipeline into individual commands
			// e.g. "ls -la | grep foo | wc -l" -> ["ls -la ", " grep foo ", " wc -l"]
			commands := strings.Split(command, "|")

			var cmds []*exec.Cmd

			var prevPipeReader io.Reader = nil

			for i, cmdStr := range commands {
				cmdStr = strings.TrimSpace(cmdStr)

				var stdoutDest io.Writer = os.Stdout
				var stderrDest io.Writer = os.Stderr

				var nextPipeReader *os.File
				var nextPipeWriter *os.File
				var err error

				if i < len(commands)-1 {
					nextPipeReader, nextPipeWriter, err = os.Pipe()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error creating pipe:", err)
						break
					}
					stdoutDest = nextPipeWriter
				}

				cmd := execute.ExecuteUserInput(cmdStr, false, prevPipeReader, stdoutDest, stderrDest)

				if cmd != nil {
					cmds = append(cmds, cmd)
				}

				if prevPipeReader != nil {
					if c, ok := prevPipeReader.(io.Closer); ok {
						c.Close()
					}
				}

				if nextPipeWriter != nil {
					nextPipeWriter.Close()
				}

				prevPipeReader = nextPipeReader
			}

			for _, cmd := range cmds {
				cmd.Wait()
			}
		} else {
			execute.ExecuteUserInput(command, true, nil, stdoutDest, stderrDest)
		}
	}
}
