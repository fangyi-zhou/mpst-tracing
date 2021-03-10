package pedro

import (
	"bufio"
	"fmt"
	"github.com/nsf/sexp"
	"os"
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
	placesOrTransitions := make([]label, 0)
	marking := make(marking)

	tokens, tokenSorts, err := loadTokensExpression(tokensExpression)
	if err != nil {
		return nil, err
	}

	placesOrTransitions, marking, err = loadPlacesExpression(placesExpression, tokenSorts, placesOrTransitions, marking)
	if err != nil {
		return nil, err
	}

	placesOrTransitions, err = loadTransitionsExpression(transitionsExpression, placesOrTransitions)
	if err != nil {
		return nil, err
	}

	arcs, err := loadArcsExpression(arcsExpression, tokenSorts)
	if err != nil {
		return nil, err
	}

	pn := PetriNet{
		tokens:              tokens,
		tokenSorts:          tokenSorts,
		placesOrTransitions: placesOrTransitions,
		arcs:                arcs,
	}
	pnMarked := MarkedPetriNet{
		pn:      pn,
		marking: marking,
	}

	return &pnMarked, nil
}

func loadTokensExpression(node *sexp.Node) ([]token, map[token]sort, error) {
	tokenSorts := make(map[token]sort)
	tokens := make([]token, 0)
	node = node.Children
	if err := ensureKeyword(node, "tokens"); err != nil {
		return nil, nil, err
	}
	node = node.Next
	for node != nil {
		tokenExpression := node.Children
		if err := ensureKeyword(tokenExpression, "token"); err != nil {
			return nil, nil, err
		}
		tokenExpression = tokenExpression.Next
		tokenValue := token(tokenExpression.Value)
		sortExpression := tokenExpression.Next
		sortValue := sort(sortExpression.Value)
		// Token sort is ignored
		tokens = append(tokens, tokenValue)
		tokenSorts[tokenValue] = sortValue
		node = node.Next
	}
	return tokens, tokenSorts, nil
}

func loadPlacesExpression(node *sexp.Node, tokenSorts map[token]sort, placesOrTransitions []label, marking marking) ([]label, marking, error) {
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
		placeValue := label(placeExpression.Value)
		if _, exists := marking[placeValue]; exists {
			return nil, nil, fmt.Errorf("redefined place %s", placeValue)
		}
		initialTokenQueue := make(map[sort]tokenQueue)
		initialMarkingExpression := placeExpression.Next.Children
		for initialMarkingExpression != nil {
			token := token(initialMarkingExpression.Value)
			sort, exists := tokenSorts[token]
			if !exists {
				return nil, nil, fmt.Errorf("unbound token %s", token)
			}
			initialTokenQueue[sort] = append(initialTokenQueue[sort], token)
			initialMarkingExpression = initialMarkingExpression.Next
		}
		marking[placeValue] = initialTokenQueue
		placesOrTransitions = append(placesOrTransitions, placeValue)
		node = node.Next
	}
	return placesOrTransitions, marking, nil
}

func loadTransitionsExpression(node *sexp.Node, placesOrTransitions []label) ([]label, error) {
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
		transitionValue := label(transitionExpression.Value)
		placesOrTransitions = append(placesOrTransitions, transitionValue)
		node = node.Next
	}
	return placesOrTransitions, nil
}

func loadArcsExpression(node *sexp.Node, tokenSorts map[token]sort) ([]arc, error) {
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
		source := label(srcExpression.Value)
		dstExpresion := srcExpression.Next
		destination := label(dstExpresion.Value)
		markingsExpression := dstExpresion.Next.Next.Children
		tokens := make(entityMarking)
		for markingsExpression != nil {
			token := token(markingsExpression.Value)
			sort, exists := tokenSorts[token]
			if !exists {
				return nil, fmt.Errorf("unbound token %s", token)
			}
			tokens[sort] = append(tokens[sort], token)
			markingsExpression = markingsExpression.Next
		}
		arc := arc{
			source:      source,
			destination: destination,
			tokens:      tokens,
		}
		arcs = append(arcs, arc)
		node = node.Next
	}
	return arcs, nil
}
