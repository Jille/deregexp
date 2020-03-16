package deregexp

import (
	"fmt"
	"regexp/syntax"
	"strings"
)

type part interface {
	describePart() string
}

type word string
type separator struct{}
type orPart []part
type concatenation []part

func (w word) describePart() string    { return string(w) }
func (separator) describePart() string { return "." }
func (g orPart) describePart() string {
	var ret []string
	for _, w := range g {
		ret = append(ret, w.describePart())
	}
	return "(" + strings.Join(ret, "|") + ")"
}
func (c concatenation) describePart() string {
	var ret []string
	for _, w := range c {
		ret = append(ret, w.describePart())
	}
	return "(" + strings.Join(ret, ", ") + ")"
}

func StripBare(re *syntax.Regexp) (retPart part) {
	switch re.Op {
	case syntax.OpNoMatch: // matches no strings
		return word("__no_matches")
	case syntax.OpEmptyMatch: // matches empty string
		return word("")
	case syntax.OpLiteral: // matches Runes sequence
		return word(re.Rune)
	case syntax.OpCharClass: // matches Runes interpreted as range pair list
		rs := uniqueInt32(re.Rune)
		if len(rs) > 5 {
			return separator{}
		}
		var ret orPart
		for _, r := range rs {
			ret = append(ret, word(fmt.Sprintf("%c", r)))
		}
		return ret
	case syntax.OpAnyCharNotNL: // matches any character except newline
		return separator{}
	case syntax.OpAnyChar: // matches any character
		return separator{}
	case syntax.OpBeginLine: // matches empty string at beginning of line
		return separator{}
	case syntax.OpEndLine: // matches empty string at end of line
		return separator{}
	case syntax.OpBeginText: // matches empty string at beginning of text
		return separator{}
	case syntax.OpEndText: // matches empty string at end of text
		return separator{}
	case syntax.OpWordBoundary: // matches word boundary `\b`
		return word("")
	case syntax.OpNoWordBoundary: // matches word non-boundary `\B`
		return word("")
	case syntax.OpCapture: // capturing subexpression with index Cap, optional name Name
		return StripBare(re.Sub[0])
	case syntax.OpStar: // matches Sub[0] zero or more times
		return separator{}
	case syntax.OpPlus: // matches Sub[0] one or more times
		return concatenation{StripBare(re.Sub[0]), separator{}}
	case syntax.OpQuest: // matches Sub[0] zero or one times
		return orPart{StripBare(re.Sub[0]), word("")}
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		s := StripBare(re.Sub[0])
		if re.Max == -1 || re.Max-re.Min > 5 {
			var ret concatenation
			for i := 0; re.Min > i; i++ {
				ret = append(ret, s)
			}
			if re.Min != re.Max {
				ret = append(ret, separator{})
			}
			return ret
		} else {
			var ret orPart
			for j := re.Min; re.Max >= j; j++ {
				var c concatenation
				for i := 0; j > i; i++ {
					c = append(c, s)
				}
				ret = append(ret, c)
			}
			return ret
		}
	case syntax.OpConcat: // matches concatenation of Subs
		var ret concatenation
		for _, s := range re.Sub {
			ret = append(ret, StripBare(s))
		}
		return ret
	case syntax.OpAlternate: // matches alternation of Subs
		var ret orPart
		for _, s := range re.Sub {
			ret = append(ret, StripBare(s))
		}
		return ret
	default:
		panic(fmt.Errorf("unknown opcode %d", re.Op))
	}
}

func uniqueInt32(a []int32) []int32 {
	seen := map[int32]bool{}
	var ret []int32
	for _, e := range a {
		if !seen[e] {
			ret = append(ret, e)
			seen[e] = true
		}
	}
	return ret
}
