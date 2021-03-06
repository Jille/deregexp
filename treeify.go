package deregexp

import (
	"fmt"
	"strings"
)

// Node is either an AndNode or OrNode in the tree.
// The tree describes an expression of ANDs and ORs and words.
type Node interface {
	Expr() string
}

// AndNode expresses a set of words and OrNodes that must all match for this node to match.
type AndNode struct {
	Words    []string
	Children []OrNode
}

// OrNode expresses a set of words and AndNodes of which at least one must match for this node to match.
type OrNode struct {
	Words    []string
	Children []AndNode
}

// Expr converts this node to a string.
func (n AndNode) Expr() string {
	var parts []string
	for _, w := range n.Words {
		parts = append(parts, fmt.Sprintf("%q", w))
	}
	for _, o := range n.Children {
		parts = append(parts, fmt.Sprintf("(%s)", o.Expr()))
	}
	return strings.Join(parts, " AND ")
}

// Expr converts this node to a string.
func (n OrNode) Expr() string {
	var parts []string
	for _, w := range n.Words {
		parts = append(parts, fmt.Sprintf("%q", w))
	}
	for _, a := range n.Children {
		parts = append(parts, fmt.Sprintf("(%s)", a.Expr()))
	}
	return strings.Join(parts, " OR ")
}

// append adds other to n and returns the new result.
// The new node is ANDed together with the existing node.
func (n AndNode) append(other Node) AndNode {
	if other == nil {
		return n
	}
	if a, ok := other.(AndNode); ok {
		return AndNode{
			Words:    append(n.Words, a.Words...),
			Children: append(n.Children, a.Children...),
		}
	}
	o := other.(OrNode)
	if len(o.Words)+len(o.Children) == 1 {
		if len(o.Words) == 1 {
			return AndNode{
				Words:    append(n.Words, o.Words[0]),
				Children: n.Children,
			}
		}
		return n.append(o.Children[0])
	}
	return AndNode{
		Words:    n.Words,
		Children: append(n.Children, o),
	}
}

// append adds other to n and returns the new result.
// The new node is ORed together with the existing node.
func (n OrNode) append(other Node) OrNode {
	if other == nil {
		return n
	}
	if o, ok := other.(OrNode); ok {
		return OrNode{
			Words:    append(n.Words, o.Words...),
			Children: append(n.Children, o.Children...),
		}
	}
	o := other.(AndNode)
	if len(o.Words)+len(o.Children) == 1 {
		if len(o.Words) == 1 {
			return OrNode{
				Words:    append(n.Words, o.Words[0]),
				Children: n.Children,
			}
		}
		return n.append(o.Children[0])
	}
	return OrNode{
		Words:    n.Words,
		Children: append(n.Children, o),
	}
}

// Treeify simplifies all possible options passed in into a node tree.
// Treeify tries to make the tree as simple as possible.
func Treeify(sequences [][]string) Node {
	for _, o := range sequences {
		if len(o) == 0 {
			return nil
		}
	}
	bestWord := mostCommon(sequences)
	if bestWord == "" {
		panic("no more words?")
	}
	var with, without [][]string
	for _, o := range sequences {
		if containsWord(o, bestWord) {
			with = append(with, filterWordAndImplications(o, bestWord))
		} else {
			without = append(without, o)
		}
	}
	wn := Treeify(with)
	wn = AndNode{
		Words: []string{bestWord},
	}.append(wn)
	if len(without) == 0 {
		return wn
	}
	won := Treeify(without)
	return OrNode{}.append(wn).append(won)
}

// mostCommon figures out which word is the most common among all options.
// If a word is a substring of another word, that makes the other word considered even more common.
func mostCommon(sequences [][]string) string {
	words := map[string][]string{}
	// First find all unique words.
	for _, o := range sequences {
		for _, w := range o {
			words[w] = nil
		}
	}
	// Now figure out for each word which words contain it as a substring.
	for w1 := range words {
		for w2 := range words {
			if strings.Contains(w2, w1) {
				words[w2] = append(words[w2], w1)
			}
		}
	}
	// Now score each word a point for every possible option it is in, giving at most one point per possibility to each word.
	scores := map[string]int{}
	for _, o := range sequences {
		scored := map[string]bool{}
		for _, w := range o {
			for _, sw := range words[w] {
				if !scored[sw] {
					scored[sw] = true
					scores[sw]++
				}
			}
		}
	}
	// Find the highest scoring word in the map, breaking ties by longest word, then lexographically.
	bestScore := 0
	bestWord := ""
	for w, s := range scores {
		if s > bestScore || (s == bestScore && (len(w) > len(bestWord) || (len(w) == len(bestWord) && w < bestWord))) {
			bestScore = s
			bestWord = w
		}
	}
	return bestWord
}

// filterWordAndImplications removes needle from haystack, and all words from haystack that are substrings of needle.
func filterWordAndImplications(haystack []string, needle string) []string {
	ret := make([]string, 0, len(haystack))
	for _, e := range haystack {
		if !strings.Contains(needle, e) {
			ret = append(ret, e)
		}
	}
	return ret
}

// containsWord returns true if needle is a substring of any word in the haystack.
func containsWord(haystack []string, needle string) bool {
	for _, w := range haystack {
		if strings.Contains(w, needle) {
			return true
		}
	}
	return false
}
