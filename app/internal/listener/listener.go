package listener

import (
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/internal/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/internal/history"
)

type Listener struct {
	autoCompleter *autocomplete.CodecraftersAutoCompleter
}

func NewListener() *Listener {
	autoCompleter := autocomplete.NewCodecraftersAutoCompleter()

	return &Listener{
		autoCompleter: autoCompleter,
	}
}

func (listener *Listener) Listen(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	switch key {
	case readline.CharTab:
		return listener.autoCompleter.AutoComplete(line, pos)
	case readline.CharNext:
		history.IncrementIndex()
		return history.GetHistoryAtIndex()
	case readline.CharPrev:
		history.DecrementIndex()
		return history.GetHistoryAtIndex()
	default:
		return line, pos, false
	}
}
