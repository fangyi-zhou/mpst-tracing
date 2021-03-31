package model

import (
	"fmt"
	"regexp"
)

type Action struct {
	Src    string
	Dest   string
	Label  string
	IsSend bool
}

func (a Action) Subject() string {
	if a.IsSend {
		return a.Src
	} else {
		return a.Dest
	}
}

func (a Action) String() string {
	var action string
	if a.IsSend {
		action = "!"
	} else {
		action = "?"
	}
	return fmt.Sprintf("%s%s%s<%s>", a.Src, action, a.Dest, a.Label)
}

func NewActionFromString(actionString string) (Action, error) {
	actionRegex := regexp.MustCompile(`^(?P<src>\w+)(?P<action>[\!\?])(?P<dest>\w+)\<(?P<label>\w+)\>$`)
	matches := actionRegex.FindStringSubmatch(actionString)
	if matches == nil {
		return Action{}, fmt.Errorf("unrecognised action string %s", actionString)
	}
	var isSend bool
	if matches[2] == "!" {
		isSend = true
	} else {
		isSend = false
	}
	return Action{
		Src:    matches[1],
		Dest:   matches[3],
		Label:  matches[4],
		IsSend: isSend,
	}, nil
}
