package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Jille/deregexp"
)

func main() {
	r := os.Args[1]
	n, err := deregexp.Deregexp(r)
	if err != nil {
		log.Fatalf("Regexp %q is invalid: %v", r, n)
	}
	fmt.Println(n.Expr())
}
