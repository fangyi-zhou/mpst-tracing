package pedro

import "errors"

type PetriNet struct {
	tokens              []token
	placesOrTransitions []placeOrTransition
	arcs                []arc
}

type arc struct {
	source      placeOrTransition
	destination placeOrTransition
	tokens      []tokenWithMultiplicity
}

type token string
type placeOrTransition string
type label string

type tokenWithMultiplicity struct {
	token        token
	multiplicity int
}

type marking map[placeOrTransition][]tokenWithMultiplicity

type MarkedPetriNet struct {
	pn      PetriNet
	marking marking
}

func (m MarkedPetriNet) Reduce(label label) (error, marking) {
	return errors.New("I didn't implement"), nil
}
