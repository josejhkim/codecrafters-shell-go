package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var history []string
var index int = 0
var lastAppendedHistory int = 0

func InitializeHistory() {
	fileName := os.Getenv("HISTFILE")
	if len(fileName) > 0 {
		AppendToHistoryFromFile(fileName)
	}
}

func SaveHistory() {
	fileName := os.Getenv("HISTFILE")
	if len(fileName) > 0 {
		SaveHistoryToFile(fileName, false)
	}
}

func AppendToHistory(command string) {
	history = append(history, command)
	index++
}

func GetHistory() []string {
	return history
}

func PrintHistory(out *io.Writer, limit int) {
	length := len(history)
	for i, cmd := range history {
		if i < length-limit {
			continue
		}
		fmt.Fprintf(*out, "    %d  %s\n", i+1, cmd)
	}
}

func GetHistoryLength() int {
	return len(history)
}

func IncrementIndex() {
	if index == len(history) {
		return
	}
	index++
}

func DecrementIndex() {
	if index == 0 {
		return
	}
	index--
}

func GetHistoryAtIndex() (newLine []rune, newPos int, ok bool) {
	if index < len(history) {
		ret := []rune(history[index])
		return ret, len(ret), true
	}
	return []rune(""), 0, true
}

func AppendToHistoryFromFile(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		AppendToHistory(scanner.Text())
	}
	return true
}

func SaveHistoryToFile(fileName string, append bool) bool {
	var flag int
	if append {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(fileName, flag, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return false
	}
	defer f.Close()
	for i, h := range history {
		if i < lastAppendedHistory {
			continue
		}
		f.WriteString(h + "\n")
	}
	lastAppendedHistory = len(history)
	return true
}
