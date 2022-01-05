package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionStringParsing1(t *testing.T) {
	actionStr := "A!B<hello>"
	action, err := NewActionFromString(actionStr)
	assert.NoError(t, err)
	assert.Equal(t, "A", action.Src)
	assert.Equal(t, "B", action.Dest)
	assert.Equal(t, "hello", action.Label)
	assert.Equal(t, true, action.IsSend)
}

func TestActionStringParsing2(t *testing.T) {
	actionStr := "A?B<hello>"
	action, err := NewActionFromString(actionStr)
	assert.NoError(t, err)
	assert.Equal(t, "A", action.Src)
	assert.Equal(t, "B", action.Dest)
	assert.Equal(t, "hello", action.Label)
	assert.Equal(t, false, action.IsSend)
}

func TestActionStringify1(t *testing.T) {
	action := Action{
		Src:    "A",
		Dest:   "B",
		Label:  "hello",
		IsSend: true,
	}
	actionString := action.String()
	assert.Equal(t, "A!B<hello>", actionString)
	assert.Equal(t, "A", action.Subject())
}

func TestActionParsingSilentTransition(t *testing.T) {
	actionStr := "_12" // Pedro returns such formats for silent transitions
	_, err := NewActionFromString(actionStr)
	assert.Error(t, err)
}

func TestActionStringify2(t *testing.T) {
	action := Action{
		Src:    "A",
		Dest:   "B",
		Label:  "hello",
		IsSend: false,
	}
	actionString := action.String()
	assert.Equal(t, "A?B<hello>", actionString)
	assert.Equal(t, "B", action.Subject())
}
