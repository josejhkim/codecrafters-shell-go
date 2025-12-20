package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

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
		}

		switch {
		case len(command) < 4:
			fmt.Println(command + ": command not found")
		case command == "exit":
			os.Exit(0)
			return
		case command[:4] == "echo":
			fmt.Println(command[5:])
		case command[:4] == "type":
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
			fmt.Println(command + ": command not found")
		}
	}
}
