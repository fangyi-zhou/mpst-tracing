package mixedstateglobaltype

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nsf/sexp"
)

func LoadFromSexp(sexp *sexp.Node) (MixedStateGlobalType, error) {
	if !sexp.IsList() {
		return nil, fmt.Errorf("expect the top level s-expression to be a list")
	}
	return loadGType(sexp.Children)
}

func loadGType(node *sexp.Node) (MixedStateGlobalType, error) {
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

func loadChoiceG(node *sexp.Node) (MixedStateGlobalType, error) {
	// node is the keyword ChoiceG
	// the next s-expression is the choice role
	choicer := node.Next.Value
	// The list of continuations follows
	choices := node.Next.Next
	if choices == nil || choices.IsScalar() {
		return nil, fmt.Errorf("expect a list of choices")
	}
	choice := choices.Children
	gtypes := make([]MixedStateGlobalType, 0)
	for choice != nil {
		gtype, err := loadGType(choice)
		if err != nil {
			return nil, err
		}
		gtypes = append(gtypes, gtype)
		choice = choice.Next
	}
	return Choice{
		choicer: choicer,
		choices: gtypes,
	}, nil
}

func loadMessageG(node *sexp.Node) (MixedStateGlobalType, error) {
	msgExp := node.Next
	if msgExp == nil || msgExp.IsScalar() {
		return nil, fmt.Errorf("expect a message")
	}
	labelExp := msgExp.Children
	if labelExp.Children.Value != "label" {
		return nil, fmt.Errorf("expect keyword label, got %s", labelExp.Children.Value)
	}
	label := labelExp.Children.Next.Value
	sendRoleExp := msgExp.Next
	if sendRoleExp == nil {
		return nil, fmt.Errorf("expect a sending role")
	}
	recvRoleExp := sendRoleExp.Next
	if recvRoleExp == nil {
		return nil, fmt.Errorf("expect a receiving role")
	}
	contExp := recvRoleExp.Next
	if contExp == nil {
		return nil, fmt.Errorf("expect a continuation")
	}
	cont, err := loadGType(contExp)
	if err != nil {
		return nil, err
	}
	return Send{
		origin: sendRoleExp.Value,
		dest:   recvRoleExp.Value,
		label:  label,
		cont: Recv{
			origin: sendRoleExp.Value,
			dest:   recvRoleExp.Value,
			label:  label,
			cont:   cont,
		},
	}, nil
}

func LoadFromSexpFile(filename string) (MixedStateGlobalType, error) {
	var ctx sexp.SourceContext
	var err error
	sourceFile := ctx.AddFile(filename, -1)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	parsed, err := sexp.Parse(reader, sourceFile)
	if err != nil {
		return nil, err
	}
	return LoadFromSexp(parsed)
}
