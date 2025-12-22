package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	for {
		fmt.Print("$ ")
		//var input string
		//fmt.Scanln(&input)
		//fmt.Printf("%s: command not found \n", input)

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = command[:len(command)-1]

		keywords := map[string]int{
			"type": 1,
			"echo": 1,
			"exit": 1,
			"pwd":  1,
			"cd":   1,
		}

		switch {
		case len(command) <= 0:
			continue
		case command == "pwd":
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(cwd)
		case command[:2] == "cd":
			// change directory
			// to the provided absolute path
			absPath := command[3:]
			if len(absPath) >= 1 && absPath[:1] == "~" {
				homePath, err := os.UserHomeDir()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				} else {
					absPath = homePath + absPath[1:]
				}
			}
			err := os.Chdir(absPath)
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", absPath)
			}
		case command == "exit":
			os.Exit(0)
			return
		case len(command) >= 4 && command[:4] == "echo":
			if len(command) <= 4 {
				fmt.Println("Usage: $ echo [command]")
				break
			}
			args := parseArgsWithQuotes(command, len("echo")+1)
			for idx, arg := range args {
				fmt.Print(arg)
				if idx < len(arg)-1 {
					fmt.Print(" ")
				}
			}
			fmt.Println("")

		case len(command) >= 4 && command[:4] == "type":
			if len(command) <= 4 {
				fmt.Printf(": not found\n")
				break
			}
			keyword := command[5:]
			if _, ok := keywords[keyword]; ok {
				fmt.Printf("%s is a shell builtin\n", keyword)
			} else if isExecutable, fileName := isExecutableFromPath(keyword); isExecutable {
				fmt.Printf("%s is %s\n", keyword, fileName)
			} else {
				fmt.Printf("%s: not found\n", keyword)
			}
		default:
			cmdName := strings.Fields(command)[0]
			isExecutable, _ := isExecutableFromPath(cmdName)
			if !isExecutable {
				fmt.Println(cmdName + ": command not found")
				continue
			}
			var out strings.Builder
			args := parseArgsWithQuotes(command, len(cmdName)+1)
			cmd := exec.Command(cmdName, args...)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(out.String())
		}
	}
}

func isExecutableFromPath(commandName string) (bool, string) {
	path, _ := exec.LookPath(commandName)
	if path != "" {
		return true, path
	}

	pathString := os.Getenv("PATH")
	pathDirs := strings.Split(pathString, string(os.PathListSeparator))

	for _, pathDir := range pathDirs {
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

func parseArgsWithQuotes(command string, index int) []string {
	args := make([]string, 0)

	currArg := ""
	for index < len(command) {
		curr := index
		switch command[curr] {
		case '"':
			curr++
			index = curr
			for index < len(command) && command[index] != '"' {
				index++
			}
			if index < len(command) {
				currArg += command[curr:index]
			}
			index++
		case '\'':
			curr++
			index = curr
			for index < len(command) && command[index] != '\'' {
				index++
			}
			if index < len(command) {
				currArg += command[curr:index]
			}
			index++
		default:
			for index < len(command) && (command[index] != ' ' && command[index] != '\'' && command[index] != '"') {
				index++
			}
			currArg += command[curr:index]
		}
		if (index >= len(command) || command[index] == ' ') && len(currArg) > 0 {
			args = append(args, currArg)
			currArg = ""
		}
		for index < len(command) && command[index] == ' ' {
			index++
		}
	}

	return args
}
