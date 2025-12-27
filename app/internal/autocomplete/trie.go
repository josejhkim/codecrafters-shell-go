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

func (root *TrieNode) GetPrefixedWords(prefix string, withPrefix bool) ([]rune, [][]rune) {
	curr := *root
	longestPrefix := prefix

	for _, c := range prefix {
		if child, okay := curr.Children[c]; okay {
			curr = *child
		} else {
			return nil, nil
		}
	}

	currString := []rune{}

	if withPrefix {
		currString = []rune(longestPrefix)
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

	if len(ret) > 1 {
		minLen := len(ret[0])
		for _, l := range ret {
			minLen = min(minLen, len(l))
		}

		prefixString := ""

		for i := 0; i < minLen; i++ {
			c := ret[0][i]
			allSame := true
			for j := 0; j < len(ret); j++ {
				if ret[j][i] != c {
					allSame = false
					break
				}
			}
			if allSame {
				prefixString += string(c)
			} else {
				break
			}
		}

		if withPrefix {
			longestPrefix = prefixString
		} else {
			longestPrefix = prefix + prefixString
		}
	}

	return []rune(longestPrefix), ret
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
