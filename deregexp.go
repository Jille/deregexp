package deregexp

import "regexp/syntax"

func Deregexp(regex string) (Node, error) {
	re, err := syntax.Parse(regex, syntax.Perl)
	if err != nil {
		return nil, err
	}
	return Treeify(flatSequences(stripBare(re))), nil
}
