package deregexp

type sequenceable interface {
	part
	isSequenceable()
}

func (word) isSequenceable()      {}
func (separator) isSequenceable() {}

type sequence []sequenceable

func  flatSequences(p part) [][]string {
	seqs := toSequences(concatenation{p})
	ret := make([][]string, len(seqs))
	for i, s := range seqs {
		ret[i] = s.flatten()
	}
	return ret
}

func  toSequences(c concatenation) []sequence {
	ret := []sequence{sequence{}}
	var todo []part = c
	for len(todo) > 0 {
		p := todo[0]
		todo = todo[1:]

		switch tp := p.(type) {
		case sequenceable: // word, separator
			for i, o := range ret {
				ret[i] = append(o, tp)
			}
		case orPart:
			newRet := make([]sequence, 0, len(ret)*len(tp))
			for _, o1 := range ret {
				for _, o2 := range tp {
					o3 := toSequences(concatenation{o1.toConcatenation(), o2})
					newRet = append(newRet, o3...)
				}
			}
			ret = newRet
		case concatenation:
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
