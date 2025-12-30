package history

import (
	"fmt"
	"io"
)

var history []string

func AppendToHistory(command string) {
	history = append(history, command)
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
