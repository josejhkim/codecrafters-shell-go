package execute

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/internal/history"
	"github.com/codecrafters-io/shell-starter-go/app/internal/parsing"
)

var Builtins = map[string]int{
	"type":    1,
	"echo":    1,
	"exit":    1,
	"pwd":     1,
	"history": 1,
	"cd":      1,
}

func ExecuteUserInput(command string, waitForFinish bool, stdin io.Reader, stdoutDest, stderrDest io.Writer) *exec.Cmd {
	command = strings.TrimSpace(command)

	var append bool = false
	var toErr bool = false

	var outputFile string

	if strings.Contains(command, ">") {
		toErr = strings.Contains(command, "2>")
		if strings.Contains(command, ">>") {
			append = true
			separationIndex := strings.Index(command, ">>")
			outputFile = strings.TrimSpace(command[separationIndex+2:])
			command = strings.TrimSpace(command[:separationIndex-1])
		} else {
			separationIndex := strings.Index(command, ">")
			outputFile = strings.TrimSpace(command[separationIndex+1:])
			command = strings.TrimSpace(command[:separationIndex-1])
		}

		var flag int
		if append {
			flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		} else {
			flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
		}
		f, err := os.OpenFile(outputFile, flag, 0777)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening file:", err)
			return nil
		}
		defer f.Close()

		if toErr {
			stderrDest = f
		} else {
			stdoutDest = f
		}
	}

	cmdAndArgs := parsing.ParseArgsWithQuotes(command, 0)

	cmd := RunCommand(cmdAndArgs, waitForFinish, stdin, stdoutDest, stderrDest)

	return cmd
}

func RunCommand(cmdAndArgs []string, waitForFinish bool, stdin io.Reader, stdout, stderr io.Writer) *exec.Cmd {
	if len(cmdAndArgs) <= 0 {
		return nil
	}

	command := cmdAndArgs[0]

	switch command {
	case "":
		return nil

	case "pwd":
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(stdout, cwd)
	case "history":
		history.PrintHistory(&stdout)
	case "cd":
		// change directory
		// to the provided absolute path
		absPath := cmdAndArgs[1]
		if len(absPath) >= 1 && absPath[:1] == "~" {
			homePath, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(stderr, err.Error())
			} else {
				absPath = homePath + absPath[1:]
			}
		}
		err := os.Chdir(absPath)
		if err != nil {
			fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", absPath)
		}

	case "exit":
		os.Exit(0)

	case "echo":
		if len(cmdAndArgs) <= 1 {
			fmt.Fprint(stdout, "Usage: $ echo [command]")
			break
		}
		length := len(cmdAndArgs[1:])
		for idx, arg := range cmdAndArgs[1:] {
			if idx < length-1 {
				fmt.Fprint(stdout, arg)
				fmt.Fprint(stdout, " ")
			} else {
				fmt.Fprintln(stdout, arg)
			}
		}

	case "type":
		if len(cmdAndArgs) <= 1 {
			stdout.Write([]byte(": not found\n"))
			break
		}
		keyword := cmdAndArgs[1]
		if _, ok := Builtins[keyword]; ok {
			fmt.Fprintf(stdout, "%s is a shell builtin\n", keyword)
		} else if isExecutable, fileName := isExecutableFromPath(keyword); isExecutable {
			fmt.Fprintf(stdout, "%s is %s\n", keyword, fileName)
		} else {
			fmt.Fprintf(stdout, "%s: not found\n", keyword)
		}

	default:
		isExecutable, _ := isExecutableFromPath(command)
		if !isExecutable {
			fmt.Fprintln(stdout, command+": command not found")
		}
		cmd := exec.Command(command, cmdAndArgs[1:]...)
		if stdin != nil {
			cmd.Stdin = stdin
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if waitForFinish {
			cmd.Run()
		} else {
			cmd.Start()
		}
		return cmd
	}
	return nil
}

func isExecutableFromPath(commandName string) (bool, string) {
	path, _ := exec.LookPath(commandName)
	if path != "" {
		return true, path
	}

	pathString := os.Getenv("PATH")
	pathDirs := strings.SplitSeq(pathString, string(os.PathListSeparator))

	for pathDir := range pathDirs {
		entries, err := os.ReadDir(pathDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(pathDir, entry.Name())
			if !entry.IsDir() && entry.Name() == commandName && entry.Type()&0111 != 0 {
				return true, fullPath
			}
		}
	}
	return false, ""
}
