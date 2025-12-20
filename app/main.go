package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
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
		case command == "exit":
			os.Exit(0)
			return
		case command[:4] == "echo":
			fmt.Println(command[5:])
		case command[:4] == "type":
			keyword := command[5:]
			if _, ok := keywords[keyword]; ok {
				fmt.Printf("%s is a shell builtin\n", keyword)
			} else {
				fmt.Printf("%s: not found\n", keyword)
			}
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
