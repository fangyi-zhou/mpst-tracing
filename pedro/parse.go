package pedro

import (
	"bufio"
	"fmt"
	"github.com/nsf/sexp"
	"os"
	"strconv"
)

func LoadFromSexpFile(filename string) (*MarkedPetriNet, error) {
	var ctx sexp.SourceContext
	var err error
	sourceFile := ctx.AddFile(filename, -1)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	loadedSexp, err := sexp.Parse(reader, sourceFile)
	if err != nil {
		return nil, err
	}
	return LoadFromSexp(loadedSexp)
}

func ensureKeyword(node *sexp.Node, keyword string) error {
	if node == nil {
		return fmt.Errorf("expected keyword %s, got nil", keyword)
	}
	if node.Value != keyword {
		return fmt.Errorf("expected keyword %s, got %s", keyword, node.Value)
	}
	return nil
}

func LoadFromSexp(node *sexp.Node) (*MarkedPetriNet, error) {
	if !node.IsList() {
		return nil, fmt.Errorf("expect a list s-expression") // TODO: add location for error reporting
	}
	petriNetKeyWord := node.Children.Children
	if err := ensureKeyword(petriNetKeyWord, "petri-net"); err != nil {
		return nil, err
	}
	tokensExpression := petriNetKeyWord.Next
	placesExpression := tokensExpression.Next
	transitionsExpression := placesExpression.Next
	arcsExpression := transitionsExpression.Next
	placesOrTransitions := make([]placeOrTransition, 0)
	marking := make(marking)

	tokens, err := loadTokensExpression(tokensExpression)
	if err != nil {
		return nil, err
	}

	placesOrTransitions, marking, err = loadPlacesExpression(placesExpression, placesOrTransitions, marking)
	if err != nil {
		return nil, err
	}

	placesOrTransitions, err = loadTransitionsExpression(transitionsExpression, placesOrTransitions)
	if err != nil {
		return nil, err
	}

	arcs, err := loadArcsExpression(arcsExpression)
	if err != nil {
		return nil, err
	}

	pn := PetriNet{
		tokens:              tokens,
		placesOrTransitions: placesOrTransitions,
		arcs:                arcs,
	}
	pnMarked := MarkedPetriNet{
		pn:      pn,
		marking: marking,
	}

	return &pnMarked, nil
}

func loadTokensExpression(node *sexp.Node) ([]token, error) {
	tokens := make([]token, 0)
	node = node.Children
	if err := ensureKeyword(node, "tokens"); err != nil {
		return nil, err
	}
	node = node.Next
	for node != nil {
		tokenExpression := node.Children
		if err := ensureKeyword(tokenExpression, "token"); err != nil {
			return nil, err
		}
		tokenExpression = tokenExpression.Next
		tokenValue := tokenExpression.Value
		// Token sort is ignored
		tokens = append(tokens, token(tokenValue))
		node = node.Next
	}
	return tokens, nil
}

func loadPlacesExpression(node *sexp.Node, placesOrTransitions []placeOrTransition, marking marking) ([]placeOrTransition, marking, error) {
	node = node.Children
	if err := ensureKeyword(node, "places"); err != nil {
		return nil, nil, err
	}
	node = node.Next
	for node != nil {
		placeExpression := node.Children
		if err := ensureKeyword(placeExpression, "place"); err != nil {
			return nil, nil, err
		}
		placeExpression = placeExpression.Next
		placeValue := placeOrTransition(placeExpression.Value)
		initialMarkingExpression := placeExpression.Next.Children
		for initialMarkingExpression != nil {
			markingExpression := initialMarkingExpression.Children
			token := token(markingExpression.Value)
			multiplicity, err := strconv.Atoi(markingExpression.Next.Value)
			if err != nil {
				return nil, nil, err
			}
			marking[placeValue] = append(marking[placeValue], tokenWithMultiplicity{
				token:        token,
				multiplicity: multiplicity,
			})
			initialMarkingExpression = initialMarkingExpression.Next
		}
		placesOrTransitions = append(placesOrTransitions, placeValue)
		node = node.Next
	}
	return placesOrTransitions, marking, nil
}

func loadTransitionsExpression(node *sexp.Node, placesOrTransitions []placeOrTransition) ([]placeOrTransition, error) {
	node = node.Children
	if err := ensureKeyword(node, "transitions"); err != nil {
		return nil, err
	}
	node = node.Next
	for node != nil {
		transitionExpression := node.Children
		if err := ensureKeyword(transitionExpression, "transition"); err != nil {
			return nil, err
		}
		transitionExpression = transitionExpression.Next
		transitionValue := placeOrTransition(transitionExpression.Value)
		placesOrTransitions = append(placesOrTransitions, transitionValue)
		node = node.Next
	}
	return placesOrTransitions, nil
}

func loadArcsExpression(node *sexp.Node) ([]arc, error) {
	node = node.Children
	arcs := make([]arc, 0)
	if err := ensureKeyword(node, "arcs"); err != nil {
		return nil, err
	}
	node = node.Next
	for node != nil {
		arcExpression := node.Children
		if err := ensureKeyword(arcExpression, "arc"); err != nil {
			return nil, err
		}
		srcExpression := arcExpression.Next
		source := placeOrTransition(srcExpression.Value)
		dstExpresion := srcExpression.Next
		destination := placeOrTransition(dstExpresion.Value)
		markingsExpression := dstExpresion.Next.Next.Children
		tokens := make([]tokenWithMultiplicity, 0)
		for markingsExpression != nil {
			markingExpression := markingsExpression.Children
			token := token(markingExpression.Value)
			multiplicity, err := strconv.Atoi(markingExpression.Next.Value)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tokenWithMultiplicity{
				token:        token,
				multiplicity: multiplicity,
			})
			markingsExpression = markingsExpression.Next
		}
		arc := arc{
			source:      source,
			destination: destination,
			tokens:      nil,
		}
		arcs = append(arcs, arc)
		node = node.Next
	}
	return arcs, nil
}
