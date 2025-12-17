package deregexp

import (
	"regexp/syntax"
	"slices"
	"strings"
)

// Longest returns the longest literal string that must be in an input to have a chance of matching the regexp. Can be used for cheap prefiltering.
// An error will only be returned iff the given regexp is invalid.
// The returned string might end in truncated UTF-8 characters.
func Longest(regex string) (string, error) {
	re, err := syntax.Parse(regex, syntax.Perl)
	if err != nil {
		return "", err
	}
	seqs := flatSequences(stripBare(re))
	if len(seqs) == 0 {
		return "", nil
	}
	slices.SortFunc(seqs, func(a, b []string) int {
		return lengthSum(a) - lengthSum(b)
	})
	var best string
	for _, s := range seqs[0] {
		for offset := 0; len(s)-len(best) > offset; offset++ {
			rem := s[offset:]
			for l := len(best) + 1; len(rem)+1 > l; l++ {
				if !allSequencesContain(seqs[1:], rem[:l]) {
					break
				}
				best = rem[:l]
			}
		}
	}
	return best, nil
}

func lengthSum(ss []string) int {
	var ret int
	for _, s := range ss {
		ret += len(s)
	}
	return ret
}

func anyContains(haystacks []string, needle string) bool {
	for _, s := range haystacks {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

func allSequencesContain(seqs [][]string, needle string) bool {
	for _, seq := range seqs {
		if !anyContains(seq, needle) {
			return false
		}
	}
	return true
}
