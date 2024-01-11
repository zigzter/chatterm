package utils

import (
	"strings"
)

type TrieNode struct {
	children map[rune]*TrieNode
	isName   bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

type Trie struct {
	Root         *TrieNode
	Suggestions  []string
	CurrentIndex int
	Prefix       string
}

func (t *Trie) autocomplete(input string) []string {
	if strings.HasPrefix(input, "@") || strings.HasPrefix(input, "/ban ") {
		searchTerm := input[1:]
		if strings.HasPrefix(input, "/ban ") {
			searchTerm = strings.TrimPrefix(input, "/ban ")
		}
		return t.Search(searchTerm)
	}
	return nil
}

func (t *Trie) UpdateSuggestion(input string) string {
	if t.Prefix == "" {
		t.Suggestions = nil
		t.Prefix = input
	}
	if t.Suggestions == nil {
		t.Suggestions = t.autocomplete(t.Prefix)
		t.CurrentIndex = 0
	} else {
		t.CurrentIndex = (t.CurrentIndex + 1) % len(t.Suggestions)
	}
	if len(t.Suggestions) > 0 {
		return t.Suggestions[t.CurrentIndex]
	}
	return input
}

func (t *Trie) Search(prefix string) []string {
	current := t.Root
	for _, c := range prefix {
		node, ok := current.children[c]
		if !ok {
			return nil
		}
		current = node
	}
	return t.collectNames(current, prefix, []string{})
}

func (t *Trie) collectNames(node *TrieNode, prefix string, words []string) []string {
	if node.isName {
		words = append(words, prefix)
	}
	for c, child := range node.children {
		words = t.collectNames(child, prefix+string(c), words)
	}
	return words
}

func (t *Trie) Insert(name string) {
	current := t.Root
	for _, c := range name {
		node, ok := current.children[c]
		if !ok {
			node = NewTrieNode()
			current.children[c] = node
		}
		current = node
	}
	current.isName = true
}

func (t *Trie) Populate(names []string) {
	for _, name := range names {
		t.Insert(name)
	}
}
