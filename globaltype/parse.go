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
	default:
		return nil, fmt.Errorf("unrecognised global type constructor %s", node.Children.Value)
	}
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
