// Package deregexp converts a regexp to an expression with substring matches. This expression does not express exactly the same as the regexp, but a superset. It can be used for a cheap first pass filter, or for index lookups.
package deregexp

import "regexp/syntax"

func Deregexp(regex string) (Node, error) {
	re, err := syntax.Parse(regex, syntax.Perl)
	if err != nil {
		return nil, err
	}
	return Treeify(flatSequences(stripBare(re))), nil
}
