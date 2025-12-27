package autocomplete

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

type CodecraftersAutoCompleter struct {
	trieRoot *TrieNode
}

func NewCodecraftersAutoCompleter() *CodecraftersAutoCompleter {
	root := NewTrieNode()

	for _, executableNameString := range GetExecutablesFromPath() {
		root.AddWord(executableNameString)
	}

	// Add "exit" as it's not a built-in executable
	// by default
	root.AddWord("exit")

	return &CodecraftersAutoCompleter{
		trieRoot: root,
	}
}

func (ccAutoCompleter *CodecraftersAutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	searchResults := ccAutoCompleter.trieRoot.GetPrefixedWords(string(line), false)

	if len(searchResults) > 0 {
		return searchResults, len(string(line))
	}

	return [][]rune{
		{rune(readline.CharBell)},
	}, pos
}

func (ccAutoCompleter *CodecraftersAutoCompleter) CompleteExecutable(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key == readline.CharTab {
		line = line[:len(line)-1]
		pos--

		secondTab := false
		if line[pos-1] == readline.CharBell {
			secondTab = true
			line = line[:pos-1]
		}

		searchResults := ccAutoCompleter.trieRoot.GetPrefixedWords(string(line), true)
		if len(searchResults) == 0 {
			return append(line, readline.CharBell), pos + 1, true
		} else if len(searchResults) == 1 {
			searchResult := searchResults[0]
			return append(searchResult, ' '), len(searchResult) + 1, true
		} else {
			if !secondTab {
				newLine = append(line, readline.CharBell)
				newPos = len(string(newLine))
				ok = true
				return
			}
			var retString strings.Builder
			for _, searchResult := range searchResults {
				retString.WriteString(string(searchResult) + "  ")
			}
			fmt.Println("")
			fmt.Println(retString.String())
			return line, len(line), true
		}
	} else {
		return line, pos, false
	}
}
