//+build ignore

package pedro

type PetriNet struct {
	tokens      []token
	places      []label        // these may not be a great idea, some transitions need to be marked as silent
	transitions map[label]bool // where the bool says if it is a labelled transition (as opposed to a silent one).
	arcs        []arc
}

type arc struct {
	source      label
	destination label
	tokens      entityMarking //These should be fifos of tokens, no multiplicity, has to be changed in pedro(OCaml)
}

type token string
type label string // the label of a place or a transition

type entityMarking = []token

type marking map[label]entityMarking

type MarkedPetriNet struct {
	pn      PetriNet
	marking marking
}

// splits two arrays a and b, into three a -b, a \cap b, b - a
func tokenSplit(a []token, b []token) ([]token, []token, []token) {
	var anotb []token
	var aandb []token
	var bnota []token

	for _, ael := range a {
		fanotb := true // flag for in a and not in b
		for _, bel := range b {
			if ael == bel {
				fanotb = false
				break
			}
		}
		if fanotb {
			aandb = append(aandb, ael)
		} else {
			anotb = append(anotb, ael)
		}
	}
	for _, bel := range b {
		fbnota := true
		for _, ael := range a {
			if ael == bel {
				fbnota = false
				break
			}
		}
		if fbnota {
			bnota = append(bnota, bel)
		}
	}
	return anotb, aandb, bnota
}

// This is a function to remove from an array, is this really needed?
func remove_idx(s entityMarking, i int) entityMarking {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// remove token from entity marking (removes the first instance and stops)
func remove_token(m entityMarking, t token) (bool, entityMarking) {
	for i, v := range m {
		if v == t {
			return true, remove_idx(m, i)
		}
	}
	return false, m
}

// consumes mreq resources from mavail
func consumeFromEntityMarking(mavail entityMarking, mreq entityMarking) (bool, entityMarking) {
	mres := mavail // copy?

	var res bool
	for _, v := range mreq {
		res, mres = remove_token(mres, v)
		if !res {
			return false, mavail
		}
	}
	return true, mres
}

func (m MarkedPetriNet) findArcsToLabel(label label) []arc {
	var res []arc
	for _, v := range m.pn.arcs {
		if v.destination == label {
			res = append(res, v)
		}
	}

	return res
}

func (m MarkedPetriNet) findArcsFromLabel(label label) []arc {
	var res []arc
	for _, v := range m.pn.arcs {
		if v.source == label {
			res = append(res, v)
		}
	}

	return res
}

// consume the resources needed to execute t and return the marking
func (mn MarkedPetriNet) consume(tr label) (bool, marking) {
	var collect_arcs = mn.findArcsToLabel(tr)
	if len(collect_arcs) == 0 {
		return false, mn.marking
	}

	m := mn.marking // copy ?
	var res bool
	for _, v := range collect_arcs {
		var mavail = m[v.source]
		res, mavail = consumeFromEntityMarking(mavail, v.tokens)
		if !res {
			return false, mn.marking
		} else {
			m[v.source] = mavail
		}
	}

	return true, m
}

func (mn MarkedPetriNet) provide(tr label) (bool, marking) {
	var collect_arcs = mn.findArcsFromLabel(tr)
	if len(collect_arcs) == 0 {
		return false, mn.marking
	}
	m := mn.marking // copy ?
	for _, v := range collect_arcs {
		m[v.destination] = append(m[v.destination], v.tokens...)
	}

	return true, m
}

func (mn MarkedPetriNet) do_transition(tr label) (bool, marking) {
	res, m := mn.consume(tr)

	if res {
		mn.marking = m
		res, mn.marking = mn.provide(tr) // if provide fails what happens to mn.marking
	}
	return res, mn.marking
}
