package globaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"strings"
)

type Send struct {
	origin string
	dest   string
	conts  map[string]GlobalType
}

func (s Send) PossiblePrefixes() []model.Action {
	var prefixes []model.Action
	for label, cont := range s.conts {
		// First add the send action
		prefixes = append(prefixes, model.Action{
			Label:  label,
			Src:    s.origin,
			Dest:   s.dest,
			IsSend: true,
		})
		contPrefixes := cont.PossiblePrefixes()
		for _, contPrefix := range contPrefixes {
			if contPrefix.Subject() != s.origin {
				prefixes = append(prefixes, contPrefix)
			}
		}
	}
	return prefixes
}

func (s Send) ConsumePrefix(m model.Action) (GlobalType, error) {
	if m.Src == s.origin && m.Dest == s.dest && m.IsSend {
		cont, exists := s.conts[m.Label]
		if exists {
			// Send prefix consumed
			return cont, nil
		} else {
			return nil, errors.New("label " + m.Label + " not permitted in the global type " + s.String())
		}
	}
	if s.origin != m.Subject() && s.dest != m.Subject() {
		// Reduction under Prefix
		var newCont = map[string]GlobalType{}
		for label, cont := range s.conts {
			consumed, err := cont.ConsumePrefix(m)
			if err != nil {
				return nil, err
			}
			newCont[label] = consumed
		}
		return Send{
			origin: s.origin,
			dest:   s.dest,
			conts:  newCont,
		}, nil
	}
	return nil, errors.New("cannot consume message " + m.String() + " in the global type " + s.String())
}

func (s Send) IsDone() bool {
	return false
}

func (s Send) String() string {
	var b strings.Builder
	s.stringWithBuilder(&b)
	return b.String()
}

func (s Send) stringWithBuilder(b *strings.Builder) {
	b.WriteString(s.origin)
	b.WriteString(" --> ")
	b.WriteString(s.dest)
	b.WriteString(": {\n")
	for label, cont := range s.conts {
		b.WriteString(label)
		b.WriteString(": ")
		cont.stringWithBuilder(b)
		b.WriteString("\n")
	}
	b.WriteString("}\n")
}

type Recv struct {
	origin string
	dest   string
	label  string
	cont   GlobalType
}

func (r Recv) PossiblePrefixes() []model.Action {
	prefixes := []model.Action{{
		Label:  r.label,
		Src:    r.origin,
		Dest:   r.dest,
		IsSend: false,
	}}
	contPrefixes := r.cont.PossiblePrefixes()
	for _, prefix := range contPrefixes {
		if prefix.Subject() != r.dest {
			prefixes = append(prefixes, prefix)
		}
	}
	return prefixes
}

func (r Recv) ConsumePrefix(m model.Action) (GlobalType, error) {
	if m.Src == r.origin && m.Dest == r.dest && !m.IsSend {
		return r.cont, nil
	}
	if m.Subject() != r.dest {
		// Reduction under prefix
		newCont, err := r.cont.ConsumePrefix(m)
		if err != nil {
			return nil, err
		} else {
			return Recv{
				origin: r.origin,
				dest:   r.dest,
				label:  r.label,
				cont:   newCont,
			}, nil
		}
	}
	return nil, errors.New("cannot consume message " + m.String() + " in the global type " + r.String())
}

func (r Recv) IsDone() bool {
	return false
}

func (r Recv) String() string {
	var b strings.Builder
	r.stringWithBuilder(&b)
	return b.String()
}

func (r Recv) stringWithBuilder(b *strings.Builder) {
	b.WriteString(r.origin)
	b.WriteString(" -~> ")
	b.WriteString(r.dest)
	b.WriteString(" ")
	b.WriteString(r.label)
	b.WriteString(": ")
	r.cont.stringWithBuilder(b)
}
