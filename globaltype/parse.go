package globaltype

import (
	"bufio"
	"fmt"
	"github.com/nsf/sexp"
	"os"
)

func LoadFromSexp(sexp *sexp.Node) (GlobalType, error) {
	if !sexp.IsList() {
		return nil, fmt.Errorf("Expect the top level s-expression to be a list")
	}
	return loadGType(sexp.Children)
}

func loadGType(node *sexp.Node) (GlobalType, error) {
	if node.IsScalar() {
		switch node.Value {
		case "EndG":
			return NewDone(), nil
		default:
			return nil, fmt.Errorf("unrecognised global type constructor %s", node.Children.Value)
		}
	}
	switch node.Children.Value {
	case "MessageG":
		return loadMessageG(node.Children)
	case "ChoiceG":
		return loadChoiceG(node.Children)
	default:
		return nil, fmt.Errorf("unrecognised global type constructor %s", node.Children.Value)
	}
}

func loadChoiceG(node *sexp.Node) (GlobalType, error) {
	// node is the keyword ChoiceG
	// the next s-expression is the choice role

	// The list of continuations follows
	choices := node.Next.Next
	if choices == nil || choices.IsScalar() {
		return nil, fmt.Errorf("Expect a list of choices")
	}
	choice := choices.Children
	gtypes := make([]GlobalType, 0)
	for choice != nil {
		gtype, err := loadGType(choice)
		if err != nil {
			return nil, err
		}
		gtypes = append(gtypes, gtype)
		choice = choice.Next
	}
	return combineGtypes(gtypes)
}

func combineGtypes(gtypes []GlobalType) (GlobalType, error) {
	if len(gtypes) == 0 {
		// nuscr should probably report an error of an empty choice, but we handle this case for completeness
		return NewDone(), nil
	} else if len(gtypes) == 1 {
		// nuscr should probably also report an error in the case of a degenerate choice, we also handle for completeness
		return gtypes[0], nil
	}
	combined := Send{conts: make(map[string]GlobalType)}
	for _, gtype := range gtypes {
		switch gtype := gtype.(type) {
		case Send:
			combined.origin = gtype.origin
			combined.dest = gtype.dest
			for label, k := range gtype.conts {
				combined.conts[label] = k
			}
		default:
			return nil, fmt.Errorf("cannot combine gtypes %s", gtype)
		}
	}
	return combined, nil
}

func loadMessageG(node *sexp.Node) (GlobalType, error) {
	msgExp := node.Next
	if msgExp == nil || msgExp.IsScalar() {
		return nil, fmt.Errorf("Expect a message")
	}
	labelExp := msgExp.Children
	if labelExp.Children.Value != "label" {
		return nil, fmt.Errorf("Expect keyword label, got %s", labelExp.Children.Value)
	}
	label := labelExp.Children.Next.Value
	sendRoleExp := msgExp.Next
	if sendRoleExp == nil {
		return nil, fmt.Errorf("Expect a sending role")
	}
	recvRoleExp := sendRoleExp.Next
	if recvRoleExp == nil {
		return nil, fmt.Errorf("Expect a receiving role")
	}
	contExp := recvRoleExp.Next
	if contExp == nil {
		return nil, fmt.Errorf("Expect a continuation")
	}
	cont, err := loadGType(contExp)
	if err != nil {
		return nil, err
	}
	return Send{
		origin: sendRoleExp.Value,
		dest:   recvRoleExp.Value,
		conts: map[string]GlobalType{
			label: Recv{
				origin: sendRoleExp.Value,
				dest:   recvRoleExp.Value,
				label:  label,
				cont:   cont,
			},
		},
	}, nil
}

func LoadFromSexpFile(filename string) (GlobalType, error) {
	var ctx sexp.SourceContext
	var err error
	sourceFile := ctx.AddFile(filename, -1)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	sexp, err := sexp.Parse(reader, sourceFile)
	if err != nil {
		return nil, err
	}
	return LoadFromSexp(sexp)
}
