package pedro

import "errors"

type PetriNet interface{}

type petriNet struct {
	tokens              []token
	placesOrTransitions []placeOrTransition
	arcs                []arc
}

type arc struct {
	source      placeOrTransition
	destination placeOrTransition
	tokens      []token
}

type token string
type placeOrTransition string
type label string

type tokenWithMultiplicity struct {
	token        token
	multiplicity int
}

type marking map[placeOrTransition]tokenWithMultiplicity

type MarkedPetriNet struct {
	pn      PetriNet
	marking marking
}

func LoadFromSexpFile(filename string) (error, *MarkedPetriNet) {
	return errors.New("I didn't implement"), nil
}

func (m MarkedPetriNet) Reduce(label label) (error, marking) {
	return errors.New("I didn't implement"), nil
}
