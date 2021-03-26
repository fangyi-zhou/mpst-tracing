package globaltype

import "fmt"

type Message struct {
	Label  string
	Origin string
	Dest   string
	Action string
}

func (m Message) Subject() string {
	if m.Action == "send" {
		return m.Origin
	} else if m.Action == "recv" {
		return m.Dest
	} else {
		panic("invalid action: " + m.Action)
	}
}

func (m Message) String() string {
	return fmt.Sprintf("%s_%s_%s:%s", m.Origin, m.Action, m.Dest, m.Label)
}
