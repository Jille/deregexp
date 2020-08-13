package deregexp

import (
	"reflect"
	"regexp/syntax"
	"testing"
)

func TestStripBare(t *testing.T) {
	tests := []struct {
		input string
		want  part
	}{
		{
			input: "Hello",
			want:  word("Hello"),
		},
		{
			input: "Hell[o]",
			want:  word("Hello"),
		},
		{
			input: "Hello+",
			want:  concatenation{word("Hell"), concatenation{word("o"), separator{}}},
		},
		{
			input: "H[ea]llo",
			want:  concatenation{word("H"), orPart{word("a"), word("e")}, word("llo")},
		},
		{
			input: "Hello{3,4}",
			want:  concatenation{word("Hell"), orPart{concatenation{word("o"), word("o"), word("o")}, concatenation{word("o"), word("o"), word("o"), word("o")}}},
		},
		{
			input: "a[0-9]b",
			want:  concatenation{word("a"), separator{}, word("b")},
		},
		{
			input: "[0-3]",
			want:  orPart{word("0"), word("1"), word("2"), word("3")},
		},
		{
			input: "[0124]",
			want:  orPart{word("0"), word("1"), word("2"), word("4")},
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			re, err := syntax.Parse(tc.input, syntax.Perl)
			if err != nil {
				t.Fatalf("failed to parse regexp: %v", err)
			}
			p := stripBare(re)
			if !reflect.DeepEqual(p, tc.want) {
				t.Errorf("mismatch: %#v; want %#v", p, tc.want)
			}
		})
	}
}
