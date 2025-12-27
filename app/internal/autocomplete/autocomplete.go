package autocomplete

import (
	"github.com/chzyer/readline"
)

type CodecraftersAutoCompleter struct {
	trieRoot *TrieNode
}

func NewCodecraftersAutoCompleter() *CodecraftersAutoCompleter {
	root := NewTrieNode()

	for _, executableNameString := range GetExecutablesFromPath() {
		root.AddWord(executableNameString + " ")
	}

	// Add "exit" as it's not a built-in executable
	// by default
	root.AddWord("exit ")

	return &CodecraftersAutoCompleter{
		trieRoot: root,
	}
}

func (ccAutoCompleter *CodecraftersAutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	searchResults := ccAutoCompleter.trieRoot.GetPrefixedWords(string(line))

	if len(searchResults) > 0 {
		return searchResults, len(string(line))
	}

	return [][]rune{
		{rune(readline.CharBell)},
	}, pos
}
