package globaltype

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestLoadingStreamingFromProtocol(t *testing.T) {
	file := path.Join("testdata", "Streaming.protocol")
	gtype, err := LoadFromProtobuf(file)
	assert.NoError(t, err)
	assert.NotNil(t, gtype)
}
