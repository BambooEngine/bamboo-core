/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LISENCE>.
 */
package bamboo

import (
	"unicode"
)

const (
	FindResultNotMatch = iota
	FindResultMatchPrefix
	FindResultMatchFull
)

type Node struct {
	Full       bool
	Dictionary bool
	Next       map[rune]*Node
}

func AddTrie(trie *Node, s []rune, dictionary bool, down bool) {
	if trie.Next == nil {
		trie.Next = map[rune]*Node{}
	}

	//add original char
	s0 := s[0]
	if trie.Next[s0] == nil {
		trie.Next[s0] = &Node{}
	}

	if len(s) == 1 {
		if !trie.Next[s0].Full {
			trie.Next[s0].Full = !down
			trie.Next[s0].Dictionary = dictionary
		}
	} else {
		AddTrie(trie.Next[s0], s[1:], dictionary, down)
	}

	//add down 1 level char
	var r0 = AddToneToChar(RemoveMarkFromChar(s0), uint8(TONE_NONE))
	if r0 != s0 {
		if trie.Next[r0] == nil {
			trie.Next[r0] = &Node{}
		}

		if len(s) == 1 {
			trie.Next[r0].Full = true
		} else {
			AddTrie(trie.Next[r0], s[1:], false, true)
		}
	}
}

func TestString(trie *Node, s []rune, dictionary bool) uint8 {

	if len(s) == 0 {
		if dictionary {
			if trie.Full && trie.Dictionary {
				return FindResultMatchFull
			}
			return FindResultNotMatch
		}
		if trie.Full {
			return FindResultMatchFull
		}
		return FindResultMatchPrefix
	}

	c := unicode.ToLower(s[0])

	if trie.Next[c] != nil {
		r := TestString(trie.Next[c], s[1:], dictionary)
		if r != FindResultNotMatch {
			return r
		}
	}

	return FindResultNotMatch
}

func dfs(trie *Node, lookup map[string]bool, s string) {
	if trie.Full {
		lookup[s] = true
	}
	for chr, t := range trie.Next {
		var key = s + string(chr)
		dfs(t, lookup, key)
	}
}

func FindNode(trie *Node, s []rune) *Node {
	if len(s) == 0 {
		return trie
	}
	c := s[0]
	if trie.Next[c] != nil {
		return FindNode(trie.Next[c], s[1:])
	}
	// not match
	return nil
}

func FindWords(trie *Node, s string) []string {
	var words []string
	var node = FindNode(trie, []rune(s))
	if node == nil {
		return nil
	}
	var lookup = map[string]bool{}
	dfs(node, lookup, s)
	for w := range lookup {
		words = append(words, w)
	}
	return words
}
