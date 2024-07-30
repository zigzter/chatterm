package utils

import (
	"strings"
)

type TrieNode struct {
	children      map[rune]*TrieNode
	isName        bool
	originalNames []string // Store the case-sensitive names that end here
}

func NewTrie() *Trie {
	return &Trie{Root: NewTrieNode()}
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children:      make(map[rune]*TrieNode),
		originalNames: make([]string, 0),
	}
}

type Trie struct {
	Root         *TrieNode
	Suggestions  []string
	CurrentIndex int
	Prefix       string
}

func (t *Trie) UpdateSuggestion(input string) string {
	newPrefix := strings.ToLower(input)
	if !strings.HasPrefix(newPrefix, t.Prefix) || t.Suggestions == nil {
		t.Prefix = newPrefix
		t.Suggestions = t.Search(t.Prefix)
		t.CurrentIndex = 0
	} else if len(t.Suggestions) > 0 {
		t.CurrentIndex = (t.CurrentIndex + 1) % len(t.Suggestions)
	}
	if len(t.Suggestions) > 0 {
		return t.Suggestions[t.CurrentIndex]
	}
	return input
}

func (t *Trie) Search(prefix string) []string {
	current := t.Root
	for _, c := range strings.ToLower(prefix) {
		node, ok := current.children[c]
		if !ok {
			return nil
		}
		current = node
	}
	return t.getNames(current, prefix, []string{})
}

func (t *Trie) getNames(node *TrieNode, prefix string, words []string) []string {
	if node.isName {
		words = append(words, node.originalNames...)
	}
	for c, child := range node.children {
		words = t.getNames(child, prefix+string(c), words)
	}
	return words
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (t *Trie) Insert(name string) {
	current := t.Root
	for _, c := range strings.ToLower(name) {
		node, ok := current.children[c]
		if !ok {
			node = NewTrieNode()
			current.children[c] = node
		}
		current = node
	}
	current.isName = true
	if !contains(current.originalNames, name) {
		current.originalNames = append(current.originalNames, name)
	}
}

func (t *Trie) Populate(names []string) {
	for _, name := range names {
		t.Insert(name)
	}
}

func (t *Trie) Reset() {
	t.Prefix = ""
	t.Suggestions = nil
	t.CurrentIndex = 0
}
