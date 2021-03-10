package pedro

type PetriNet struct {
	tokens              []token
	placesOrTransitions []label // these may not be a great idea, some transitions need to be marked as silent
	arcs                []arc
}

type arc struct {
	source      label
	destination label
	tokens      []entityMarking //These should be fifos of tokens, no multiplicity, has to be changed in pedro(OCaml)
}

type token string
type label string // the label of a place or a transition

/*type tokenWithMultiplicity struct {
	token        token
	multiplicity int
}*/

type tokenQueue = []token

type entityMarking = map[label]tokenQueue

type marking map[label]entityMarking

type MarkedPetriNet struct {
	pn      PetriNet
	marking marking
}

func (m MarkedPetriNet) findArcsToLabel(label label) (error, []arc) {
	var res []arc
	for _, v := range m.pn.arcs {
		if v.destination == label {
			res = append(res, v)
		}
	}

	return nil, res
	//return errors.New("I didn't implement"), nil
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

// consumes from a the tokens in b, returns first the extra resources that a has, the consumed resources, and the resources from b that remain to be provided
// this function is a bit scary
//func consumeResources (a []tokenWithMultiplicity, b []tokenWithMultiplicity) ([]tokenWithMultiplicity, []tokenWithMultiplicity, []tokenWithMultiplicity){
//	var extra []tokenWithMultiplicity
//	var consumed []tokenWithMultiplicity
//	var remaining []tokenWithMultiplicity
//
//	for _,ael := range a {
//		toConsume := 0 // flag for in a and not in b
//		for _, bel := range b {
//			if ael.token == bel.token { toConsume = bel.multiplicity ; break }
//		}
//		if toConsume > 0 {
//
//			if ael.multiplicity > toConsume { //if there are strictly more tokens than required
//				eel := ael ; eel.multiplicity = ael.multiplicity - toConsume
//				extra = append(extra, eel)
//
//				cel := ael ; cel.multiplicity = toConsume
//				consumed = append(consumed, cel)
//			} else if ael.multiplicity == toConsume {
//				consumed = append(consumed, ael) // all resources were consumed
//			} else {
//				cel := ael ; cel.multiplicity = ael.multiplicity
//				consumed = append(consumed, cel)
//
//				rel := ael ; cel.multiplicity = toConsume - ael.multiplicity
//				remaining = append(remaining, rel)
//			}
//		} else {
//			extra = append(extra, ael)
//		}
//	}
//	for _,bel := range b {
//		stillRemain := true
//		for _, ael := range a {
//			if ael.token == bel.token { stillRemain = false ; break}
//		}
//		if stillRemain { remaining = append(remaining, bel) }
//	}
//	return extra, consumed, remaining
//}

//func mem(tk tokenWithMultiplicity, tks []tokenWithMultiplicity) (bool, *tokenWithMultiplicity) {
//	for _,v := range tks {
//		if tk.token == v.token { return true, nil}
//	}
//
//	return false, nil
//}
//
//// collects the silent transitions that may bring
//func (m MarkedPetriNet) collectSilentTransition (dst label, tk tokenWithMultiplicity){
//	for _, a := range m.pn.arcs {
//		if a.destination == dst && a.source == "" { //this is an arc of interest (is this how we know if it is silent?)
//			//var xx = mem(tk, a.tokens)
//
//
//		}
//	}
//
//}

//func (m MarkedPetriNet) consumeResourcesOnMarking(act []arc) (error, []arc){
//	mark := m.marking // copy the marking (hopefully, I don't undertand Go semantics well enough)
//
//	for _,v := range act {
//		var extra, _, _ /* remaining */ = consumeResources(mark[v.source], v.tokens)
//		mark[v.source] = extra //resources that stay in the marking after firing
//
//		// at this point we have to process the remaining tokens to see if we can get them all.
//
//	}
//return errors.New("I didn't implement"), nil
//}
//
//
//func (m MarkedPetriNet) Reduce(label label) (error, marking) {
//
//	return errors.New("I didn't implement"), nil
//}
