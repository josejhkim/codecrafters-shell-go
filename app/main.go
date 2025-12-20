package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		command = strings.TrimSpace(command)
		if strings.ToLower(command) == "exit" {
			os.Exit(0)
		}
		fmt.Println(command + ": command not found")
	}
}
