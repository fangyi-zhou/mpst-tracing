//+build ignore

package pedro

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestParseTwoBuyer(t *testing.T) {
	file := path.Join("..", "twobuyer", "TwoBuyer_petri_net.sexp")
	pn, err := LoadFromSexpFile(file)
	assert.NoError(t, err)
	assert.NotNil(t, pn)
	tokens := []string{"buy", "share", "quote'", "quote", "query", "S", "B", "A"}
	for _, tk := range tokens {
		assert.Contains(t, pn.pn.tokens, token(tk))
	}
}
