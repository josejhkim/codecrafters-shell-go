package autocomplete

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

func (root *TrieNode) GetPrefixedWords(prefix string) [][]rune {
	curr := *root
	for _, c := range prefix {
		if child, okay := curr.Children[c]; okay {
			curr = *child
		} else {
			return nil
		}
	}

	currString := []rune{}
	ret := [][]rune{}
	ret = curr.DFS(currString, ret)

	return ret
}

func (node *TrieNode) DFS(curr []rune, rets [][]rune) [][]rune {
	for c, child := range node.Children {
		curr = append(curr, c)

		if child.IsEnd {
			rets = append(rets, curr)
		}

		rets = child.DFS(curr, rets)
		curr = curr[:len(curr)-1]
	}
	return rets
}
