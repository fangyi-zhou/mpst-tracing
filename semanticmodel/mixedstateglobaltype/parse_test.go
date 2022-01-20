package mixedstateglobaltype

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func testParseFile(t *testing.T, file string) {
	file = path.Join("testdata", file)
	gtype, err := LoadFromSexpFile(file)
	assert.NoError(t, err)
	assert.NotNil(t, gtype)
}

func TestParseTwoBuyer(t *testing.T) {
	testParseFile(t, "TwoBuyer_global_type.sexp")
}

func TestParseTwoButtons(t *testing.T) {
	testParseFile(t, "RedButtons.sexp")
}
