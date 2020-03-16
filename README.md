# Deregexp

[![GoDoc](https://godoc.org/github.com/Jille/deregexp?status.svg)](https://godoc.org/github.com/Jille/deregexp)
[![Build Status](https://travis-ci.org/Jille/deregexp.png)](https://travis-ci.org/Jille/deregexp)

Deregexp converts a regexp to an expression with substring matches. This expression does not express exactly the same as the regexp, but a superset. It can be used for a cheap first pass filter, or for index lookups.

## Examples

```shell
$ go run cmd/test/test.go "Mary.+had.+a.+little.+(lamb|sheep)"
"little" AND "Mary" AND "had" AND ("sheep" OR "lamb")
```

Note that the "a" has been dropped because it is implied by "Mary" and "had".

```shell
$ go run cmd/test/test.go "Mary.+had.+(a|no).+little.+(lamb|sheep)"
"little" AND "Mary" AND "had" AND ("sheep" OR "lamb")
```

Note that the "(a|no)" has been dropped because the "a" is implied by "Mary" and "had", and `(TRUE OR "no")` is always true.

## Known limitations

Case insensitivity is not supported. You can of course use the resulting expression case insensitively, but controlling case sensitivity with `(?i)` is not supported. Patches are welcome.
