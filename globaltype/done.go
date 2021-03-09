package globaltype

import (
	"errors"
	"strings"
)

type Done struct{}

func (Done) PossiblePrefixes() []Message {
	return nil
}

func (Done) ConsumePrefix(message Message) (GlobalType, error) {
	return nil, errors.New("end cannot consume prefix")
}

func (Done) IsDone() bool {
	return true
}

func (Done) String() string {
	return "end"
}

func (Done) stringWithBuilder(b *strings.Builder) {
	b.WriteString("end")
}
