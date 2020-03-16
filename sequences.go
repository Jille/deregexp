package deregexp

// sequenceable parts are parts that are allowed in a sequence.
type sequenceable interface {
	part
	isSequenceable()
}

func (word) isSequenceable()      {}
func (separator) isSequenceable() {}

// sequence is a simplified concatenation in which all concatenations and orParts have been resolved.
type sequence []sequenceable

// flatSequences converts a part into all possible options described by that part.
// concatenation{word("Hi"), separator{}, word(" there")} -> [][]string{[]string{"Hi", " there"}}
// orPart{word("Hi"), word("Bye")} -> [][]string{[]string{"Hi"}, []string{"Bye"}}
func flatSequences(p part) [][]string {
	seqs := toSequences(concatenation{p})
	ret := make([][]string, len(seqs))
	for i, s := range seqs {
		ret[i] = s.flatten()
	}
	return ret
}

// toSequence returns all possible combinations described by a concatenation. If there are no orParts in c, this will only return a single sequence.
func toSequences(c concatenation) []sequence {
	ret := []sequence{sequence{}}
	var todo []part = c
	for len(todo) > 0 {
		p := todo[0]
		todo = todo[1:]

		switch tp := p.(type) {
		case sequenceable: // word, separator
			// Add this sequenceable to all possible results.
			for i, o := range ret {
				ret[i] = append(o, tp)
			}
		case orPart:
			// Multiply all possible combinations we had so far with the options from this orPart.
			newRet := make([]sequence, 0, len(ret)*len(tp))
			for _, o1 := range ret {
				for _, o2 := range tp {
					o3 := toSequences(concatenation{o1.toConcatenation(), o2})
					newRet = append(newRet, o3...)
				}
			}
			ret = newRet
		case concatenation:
			// "Recurse" into tp by prepending it to our remaining todo list.
			todo = append(append([]part{}, tp...), todo...)
		default:
			panic("wtf")
		}
	}
	return ret
}

func (s sequence) toConcatenation() concatenation {
	ret := make(concatenation, len(s))
	for i, e := range s {
		ret[i] = e
	}
	return ret
}

// flatten reduces a sequences to a list of words. Consecutive words get merged together, separators are boundaries between words.
func (s sequence) flatten() []string {
	var ret []string
	var w string
	for _, p := range s {
		switch t := p.(type) {
		case word:
			w += string(t)
		case separator:
			if w != "" {
				ret = append(ret, w)
				w = ""
			}
		default:
			panic("wtf")
		}
	}
	if w != "" {
		ret = append(ret, w)
	}
	return ret
}
