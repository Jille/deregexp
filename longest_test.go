package deregexp

import (
	"testing"
)

func TestLongest(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "Hello",
			want:  "Hello",
		},
		{
			input: "H[ae]llo",
			want:  "llo",
		},
		{
			input: "Hello?",
			want:  "Hell",
		},
		{
			input: "Hello{3,4}",
			want:  "Hellooo",
		},
		{
			input: "1[2b](3|c)",
			want:  "1",
		},
		{
			input: "1.(3|c)",
			want:  "1",
		},
		{
			input: "1.?(3|c)",
			want:  "1",
		},
		{
			input: "hi (alligator|elevator)",
			want:  "ator",
		},
		{
			input: "(a|b)",
			want:  "",
		},
		{
			input: ".",
			want:  "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			l, err := Longest(tc.input)
			if err != nil {
				t.Fatalf("Longest(): %v", err)
			}
			if l != tc.want {
				t.Errorf("Longest(): %q; want %q", l, tc.want)
			}
		})
	}
}
