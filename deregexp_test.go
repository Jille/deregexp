package deregexp

import (
	"testing"
)

func TestDeregexp(t *testing.T) {
	tests := []struct {
		input string
		want string
	}{
		{
			input: "Hello",
			want: `"Hello"`,
		},
		{
			input: "H[ae]llo",
			want: `"Hallo" OR "Hello"`,
		},
		{
			input: "Hello?",
			want: `"Hell"`,
		},
		{
			input: "Hello{3,4}",
			want:  `"Hellooo"`,
		},
		{
			input: "1[2b](3|c)",
			want: `"123" OR "12c" OR "1b3" OR "1bc"`,
		},
		{
			input: "1.(3|c)",
			want: `"1" AND ("3" OR "c")`,
		},
		{
			input: "1.?(3|c)",
			want: `"1" AND ("3" OR "c")`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			n, err := Deregexp(tc.input)
			if err != nil {
				t.Fatalf("Deregexp(): %v", err)
			}
			if n.Expr() != tc.want {
				t.Errorf("Deregexp(): %s (%+v); want %s", n.Expr(), n, tc.want)
			}
		})
	}
}
