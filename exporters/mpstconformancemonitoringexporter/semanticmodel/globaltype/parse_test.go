package globaltype

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestParseTwoBuyer(t *testing.T) {
	file := path.Join("..", "..", "..", "..", "twobuyer", "TwoBuyer_global_type.sexp")
	pn, err := LoadFromSexpFile(file)
	assert.NoError(t, err)
	assert.NotNil(t, pn)
}
