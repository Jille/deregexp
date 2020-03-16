package deregexp

import (
	"reflect"
	"regexp/syntax"
	"testing"
)

func TestSequences(t *testing.T) {
	tests := []struct {
		input string
		want  [][]string
	}{
		{
			input: "Hello",
			want:  [][]string{
				{"Hello"},
			},
		},
		{
			input: "Hell[o]",
			want:  [][]string{
				{"Hello"},
			},
		},
		{
			input: "Hello+",
			want:  [][]string{
				{"Hello"},
			},
		},
		{
			input: "H[ea]llo",
			want:  [][]string{
				{"Hallo"},
				{"Hello"},
			},
		},
		{
			input: "Hello{3,4}",
			want:  [][]string{
				{"Hellooo"},
				{"Helloooo"},
			},
		},
		{
			input: "Hello{3,4}( 123)?",
			want:  [][]string{
				{"Hellooo 123"},
				{"Hellooo"},
				{"Helloooo 123"},
				{"Helloooo"},
			},
		},
		{
			input: "Hello{3,}( 123)?(456)?",
			want:  [][]string{
				{"Hellooo", " 123456"},
				{"Hellooo", " 123"},
				{"Hellooo", "456"},
				{"Hellooo"},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			re, err := syntax.Parse(tc.input, syntax.Perl)
			if err != nil {
				t.Fatalf("failed to parse regexp: %v", err)
			}
			p := StripBare(re)
			seqs := flatSequences(p)
			// TODO(quis): Ignore ordering
			if !reflect.DeepEqual(seqs, tc.want) {
				t.Errorf("mismatch: %#v; want %#v", seqs, tc.want)
			}
		})
	}
}
