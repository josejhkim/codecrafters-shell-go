package autocomplete

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {

	t.Run("test prefix search", func(t *testing.T) {
		rootNode := NewTrieNode()

		prefix := "str"

		words := map[string]int{
			"string": 1,
			"strung": 1,
			"strang": 1,
		}

		wordsWithoutPrefix := map[string]int{}

		for key := range words {
			rootNode.AddWord(key)
			wordsWithoutPrefix[key[len(prefix)+1:]] = 1
		}

		res := rootNode.GetPrefixedWords("str", false)

		if len(res) != len(words) {
			t.Errorf("Only %d prefixes instead of %d", len(res), len(words))
		}

		for _, s := range res {
			if _, ok := wordsWithoutPrefix[string(s)]; !ok {
				fmt.Printf("%s missing from the word search retrieval", prefix+string(s))
			}
		}
	})
}
