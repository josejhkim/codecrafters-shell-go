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

func PrintHistory(out *io.Writer) {
	for i, cmd := range history {
		fmt.Fprintf(*out, "    %d  %s\n", i+1, cmd)
	}
}
