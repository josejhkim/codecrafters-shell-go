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
	keywords := map[string]int{
		"type": 1,
		"echo": 1,
		"exit": 1,
		"pwd":  1,
		"cd":   1,
	}

	for {
		fmt.Print("$ ")

		// for stdout
		var out strings.Builder

		// for stderr
		var errs strings.Builder

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = command[:len(command)-1]

		var append bool = false
		var appendToErr bool = false
		var appendOutputFile string

		var redirect bool = false
		var redirectToErr bool = false
		var redirectOutputFile string

		if strings.Contains(command, ">>") {
			append = true
			appendToErr = strings.Contains(command, "2>>")
			separationIndex := strings.Index(command, ">>")
			appendOutputFile = strings.TrimSpace(command[separationIndex+2:])
			command = strings.TrimSpace(command[:separationIndex-1])
		} else if strings.Contains(command, ">") {
			redirect = true
			redirectToErr = strings.Contains(command, "2>")
			separationIndex := strings.Index(command, ">")
			redirectOutputFile = strings.TrimSpace(command[separationIndex+1:])
			command = strings.TrimSpace(command[:separationIndex-1])
		}

		switch {
		case len(command) <= 0:
			continue

		case strings.Contains(command, ">"):

		case command == "pwd":
			cwd, err := os.Getwd()
			if err != nil {
				errs.WriteString(fmt.Sprintln(err))
				os.Exit(1)
			}
			out.WriteString(fmt.Sprintln(cwd))

		case command[:2] == "cd":
			// change directory
			// to the provided absolute path
			absPath := command[3:]
			if len(absPath) >= 1 && absPath[:1] == "~" {
				homePath, err := os.UserHomeDir()
				if err != nil {
					errs.WriteString(fmt.Sprintln(err))
				} else {
					absPath = homePath + absPath[1:]
				}
			}
			err := os.Chdir(absPath)
			if err != nil {
				errs.WriteString(fmt.Sprintf("cd: %s: No such file or directory\n", absPath))
			}

		case command == "exit":
			os.Exit(0)
			return

		case len(command) >= 4 && command[:4] == "echo":
			if len(command) <= 4 {
				out.WriteString(fmt.Sprintln("Usage: $ echo [command]"))
				break
			}
			args := parseArgsWithQuotes(command, len("echo")+1)
			for idx, arg := range args {
				out.WriteString(fmt.Sprint(arg))
				if idx < len(arg)-1 {
					out.WriteString(" ")
				}
			}
			out.WriteString(fmt.Sprintln(""))

		case len(command) >= 4 && command[:4] == "type":
			if len(command) <= 4 {
				fmt.Printf(": not found\n")
				break
			}
			keyword := command[5:]
			if _, ok := keywords[keyword]; ok {
				out.WriteString(fmt.Sprintf("%s is a shell builtin\n", keyword))
			} else if isExecutable, fileName := isExecutableFromPath(keyword); isExecutable {
				out.WriteString(fmt.Sprintf("%s is %s\n", keyword, fileName))
			} else {
				out.WriteString(fmt.Sprintf("%s: not found\n", keyword))
			}

		default:
			args := parseArgsWithQuotes(command, 0)
			cmdName := args[0]
			isExecutable, _ := isExecutableFromPath(cmdName)
			if !isExecutable {
				out.WriteString(fmt.Sprintln(cmdName + ": command not found"))
			} else {
				cmd := exec.Command(cmdName, args[1:]...)
				cmd.Stdout = &out
				cmd.Stderr = &errs
				cmd.Run()
			}

		}

		if !redirect && !append {
			if errs.Len() > 0 {
				fmt.Print(errs.String())
			} else {
				fmt.Print(out.String())
			}
		} else if append {
			f, err := os.OpenFile(appendOutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				log.Fatal(err)
			}
			if appendToErr {
				_, err = f.Write([]byte(errs.String()))
			} else {
				_, err = f.Write([]byte(out.String()))
			}
			if err != nil {
				log.Fatal(err)
			}
		} else {
			var err error
			if redirectToErr {
				err = os.WriteFile(redirectOutputFile, []byte(errs.String()), 0777)
				fmt.Print(out.String())
			} else {
				err = os.WriteFile(redirectOutputFile, []byte(out.String()), 0777)
				fmt.Print(errs.String())
			}
			if err != nil {
				log.Fatal(err)
			}
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

var escapableInDoubleQuotes = map[rune]int{
	'"':  1,
	'\\': 1,
	'$':  1,
	'`':  1,
	'\n': 1,
}

func parseArgsWithQuotes(command string, index int) []string {
	runes := []rune(command)
	args := make([]string, 0)

	currArg := ""
	for index < len(command) {
		curr := index
		switch command[curr] {
		case '"':
			curr++
			index = curr
			for index < len(command) && command[index] != '"' {
				if command[index] == '\\' {
					if _, ok := escapableInDoubleQuotes[runes[index+1]]; ok {
						index++
					}
				}
				currArg += string(command[index])
				index++
			}
			index++
		case '\'':
			curr++
			index = curr
			for index < len(command) && command[index] != '\'' {
				// if command[index] == '\\' {
				// 	index++
				// }
				currArg += string(command[index])
				index++
			}
			index++
		default:
			for index < len(command) && (command[index] != ' ' && command[index] != '\'' && command[index] != '"') {
				if command[index] == '\\' {
					index++
				}
				currArg += string(command[index])
				index++
			}
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
