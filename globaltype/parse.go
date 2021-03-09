package globaltype

import (
	"bufio"
	"errors"
	"github.com/nsf/sexp"
	"os"
)

func LoadFromSexp(sexp *sexp.Node) (GlobalType, error) {
	return nil, errors.New("unimplemented")
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
