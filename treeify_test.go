package deregexp

import (
	"fmt"
	"testing"
)

func TestTreeify(t *testing.T) {
	tests := []struct {
		input [][]string
		want  string
	}{
		{
			input: [][]string{
				{"Hello"},
			},
			want: `"Hello"`,
		},
		{
			input: [][]string{
				{"Hallo"},
				{"Hello"},
			},
			want: `"Hallo" OR "Hello"`,
		},
		{
			input: [][]string{
				{"Hellooo"},
				{"Halloooo"},
			},
			want: `"Halloooo" OR "Hellooo"`,
		},
		{
			input: [][]string{
				{"Hellooo 123"},
				{"Hellooo"},
				{"Helloooo 123"},
				{"Helloooo"},
			},
			want: `"Hellooo"`,
		},
		{
			input: [][]string{
				{"Hellooo", " 123456"},
				{"Hellooo", " 123"},
				{"Hellooo", "456"},
				{"Hellooo"},
			},
			want: `"Hellooo"`,
		},
		{
			input: [][]string{
				{"Hellooo"},
				{"Hellooo"},
			},
			want: `"Hellooo"`,
		},
		{
			input: [][]string{
				{"a"},
				{"aa"},
			},
			want: `"a"`,
		},
		{
			input: [][]string{
				{"a", "c"},
				{"aa", "b", "a", "c"},
				{"c"},
			},
			want: `"c"`,
		},
		{
			input: [][]string{
				{"a", "c"},
				{"aa", "b", "a", "c"},
				{"aaa", "c"},
			},
			want: `"a" AND "c"`,
		},
		{
			input: [][]string{
				{"Mary", "had", "a", "little", "lamb"},
				{"Mary", "had", "a", "little", "sheep"},
			},
			want: `"little" AND "Mary" AND "had" AND ("sheep" OR "lamb")`,
		},
		{
			input: [][]string{
				{"Mary", "had", "a", "little", "lamb"},
			},
			want: `"little" AND "Mary" AND "lamb" AND "had"`,
		},
		{
			input: [][]string{
				{"a", "b"},
				{"a", "c"},
				{"a", "d"},
				{"e"},
			},
			want: `"e" OR ("a" AND ("b" OR "c" OR "d"))`,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.input), func(t *testing.T) {
			n := Treeify(tc.input)
			if n.Expr() != tc.want {
				t.Errorf("Treeify: %s (%+v); want %s", n.Expr(), n, tc.want)
			}
		})
	}
}
