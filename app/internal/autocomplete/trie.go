package autocomplete

import "slices"

type TrieNode struct {
	Char     rune
	IsEnd    bool
	Children map[rune]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Char:     ' ',
		IsEnd:    false,
		Children: make(map[rune]*TrieNode),
	}
}

func (root *TrieNode) AddWord(word string) {
	curr := root
	for _, c := range word {
		if child, okay := curr.Children[c]; okay {
			curr = child
		} else {
			newChild := &TrieNode{
				Char:     c,
				IsEnd:    false,
				Children: make(map[rune]*TrieNode),
			}
			curr.Children[c] = newChild
			curr = newChild
		}
	}
	curr.IsEnd = true
}

func (root *TrieNode) GetPrefixedWords(prefix string, withPrefix bool) [][]rune {
	curr := *root
	for _, c := range prefix {
		if child, okay := curr.Children[c]; okay {
			curr = *child
		} else {
			return nil
		}
	}

	currString := []rune{}

	if withPrefix {
		currString = []rune(prefix)
	}

	ret := [][]rune{}
	if curr.IsEnd {
		wordFound := make([]rune, len(currString))
		copy(wordFound, currString)

		ret = append(ret, wordFound)
	}

	ret = curr.DFS(currString, ret)
	slices.SortFunc(ret, func(a, b []rune) int {
		if string(a) < string(b) {
			return -1
		}
		return 1
	})
	return ret
}

func (node *TrieNode) DFS(curr []rune, rets [][]rune) [][]rune {
	for c, child := range node.Children {
		curr = append(curr, c)

		if child.IsEnd {
			wordFound := make([]rune, len(curr))
			copy(wordFound, curr)
			rets = append(rets, wordFound)
		}

		rets = child.DFS(curr, rets)

		curr = curr[:len(curr)-1]
	}
	return rets
}
