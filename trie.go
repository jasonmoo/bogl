package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type (
	Trie map[rune]*TrieNode

	TrieNode struct {
		char rune
		word string
		next map[rune]*TrieNode
	}
)

func NewTrieNode(c rune) *TrieNode {
	return &TrieNode{
		char: c,
		next: make(map[rune]*TrieNode),
	}
}

func (tn *TrieNode) String() string {
	var cs []rune
	for c, _ := range tn.next {
		cs = append(cs, c)
	}
	sort.Slice(cs, func(i, j int) bool {
		return cs[i] < cs[j]
	})
	return fmt.Sprintf("%q: %q", tn.char, cs)
}

func (t Trie) AddWord(word string) {

	// drop single char words
	if len(word) < 2 {
		return
	}

	word = strings.ToLower(word)

	var node *TrieNode

	for _, c := range word {
		if node == nil {
			var exists bool
			node, exists = t[c]
			if !exists {
				node = NewTrieNode(c)
				t[c] = node
			}
			continue
		}

		next, exists := node.next[c]
		if !exists {
			next = NewTrieNode(c)
			node.next[c] = next
		}

		node = next
	}

	node.word = word

}

func LoadTrie(rootPath string) map[rune]*TrieNode {

	trie := make(Trie)

	file, err := os.OpenFile(rootPath, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		trie.AddWord(sc.Text())
	}

	return trie

}
