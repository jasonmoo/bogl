package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type (
	Matrix struct {
		rs   [][]rune
		xmax int
		ymax int

		seen [][]bool

		walks int
	}

	Word struct {
		ps []point
		w  string
	}
	point struct {
		x, y int
	}
)

func (p point) String() string {
	return fmt.Sprintf("%d:%d", p.x, p.y)
}

func (w Word) String() string {
	return fmt.Sprintf("%s: %v", w.w, w.ps)
}

func NewMatrix(x, y int) *Matrix {
	m := &Matrix{
		rs:   make([][]rune, y),
		seen: make([][]bool, y),
		xmax: x,
		ymax: y,
	}
	for y, _ := range m.rs {
		m.rs[y] = make([]rune, x)
		m.seen[y] = make([]bool, x)
	}
	return m
}

func (m *Matrix) Randomize(chars string) {
	rs := []rune(chars)
	for y := 0; y < m.ymax; y++ {
		for x := 0; x < m.xmax; x++ {
			m.rs[y][x] = rs[rand.Intn(len(rs))]
		}
	}
}

func (m *Matrix) FindWords(trie Trie) []Word {

	var (
		node *TrieNode

		words []Word
		word  Word

		walk func(x, y int)
	)

	walk = func(x, y int) {

		// just for stats
		m.walks++

		// skip this path if we hit a seen block
		if m.seen[y][x] {
			return
		}
		m.seen[y][x] = true
		defer func() { m.seen[y][x] = false }()

		c := m.rs[y][x]

		if node == nil {
			n, exists := trie[c]
			if !exists {
				fmt.Println(string(c), "not found in root of trie?")
				return
			}
			node = n
		} else {
			n, exists := node.next[c]
			if !exists {
				return
			}
			prev := node
			defer func() { node = prev }()
			node = n
		}

		word.ps = append(word.ps, point{x, y})
		defer func() { word.ps = word.ps[:len(word.ps)-1] }()

		if node.word != "" {
			word.w = node.word
			words = append(words, word)
		}

		// if row above
		if y > 0 {
			if x > 0 {
				walk(x-1, y-1)
			}
			walk(x, y-1)
			if x < m.xmax-1 {
				walk(x+1, y-1)
			}
		}
		// current row on either side
		if x > 0 {
			walk(x-1, y)
		}
		if x < m.xmax-1 {
			walk(x+1, y)
		}
		// if row below
		if y < m.ymax-1 {
			if x > 0 {
				walk(x-1, y+1)
			}
			walk(x, y+1)
			if x < m.xmax-1 {
				walk(x+1, y+1)
			}
		}

	}

	for y := 0; y < m.ymax; y++ {
		for x := 0; x < m.xmax; x++ {
			walk(x, y)
			node = nil
		}
	}

	return words

}

func (m *Matrix) String() string {
	var buf bytes.Buffer
	buf.WriteByte('+')
	fmt.Fprint(&buf, strings.Repeat("-", m.xmax))
	buf.WriteByte('+')
	buf.WriteByte('\n')
	for y := 0; y < m.ymax; y++ {
		buf.WriteByte('|')
		for x := 0; x < m.xmax; x++ {
			buf.WriteRune(m.rs[y][x])
		}
		buf.WriteByte('|')
		buf.WriteByte('\n')
	}
	buf.WriteByte('+')
	fmt.Fprint(&buf, strings.Repeat("-", m.xmax))
	buf.WriteByte('+')
	buf.WriteByte('\n')
	return buf.String()
}

const (
	fullAlphabet = `abcdefghijklmnopqrstuvwxyz`
)

var (
	size = flag.Int("size", 4, "matrix = size x size")
)

func main() {

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	trie := LoadTrie("/usr/share/dict/words")
	fmt.Println("loaded trie in", time.Since(start))

	m := NewMatrix(*size, *size)
	m.Randomize(fullAlphabet + fullAlphabet)
	fmt.Println(m)

	start = time.Now()
	words := m.FindWords(trie)
	fmt.Println(m.walks, "walks found", len(words), "words in", time.Since(start))

	for _, word := range words {
		fmt.Println(word)
	}

}
