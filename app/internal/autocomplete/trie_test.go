package autocomplete

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {

	t.Run("test prefix search", func(t *testing.T) {
		rootNode := NewTrieNode()

		prefix := "st"

		words := map[string]int{
			"string": 1,
			"strung": 1,
			"strang": 1,
		}

		actualLongestPrefix := "str"

		wordsWithoutPrefix := map[string]int{}

		for key := range words {
			rootNode.AddWord(key)
			wordsWithoutPrefix[key[len(prefix):]] = 1
		}

		longestPrefix, res := rootNode.GetPrefixedWords(prefix, false)

		if string(longestPrefix) != actualLongestPrefix {
			t.Errorf("Longest prefix is %s whereas it should've been %s", string(longestPrefix), actualLongestPrefix)
		}

		if len(res) != len(words) {
			t.Errorf("Only %d prefixes instead of %d", len(res), len(words))
		}

		for _, s := range res {
			if _, ok := wordsWithoutPrefix[string(s)]; !ok {
				fmt.Printf("%s missing from the word search retrieval\n", prefix+string(s))
			}
		}
	})
}
