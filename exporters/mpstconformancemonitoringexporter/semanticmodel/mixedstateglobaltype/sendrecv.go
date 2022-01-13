package mixedstateglobaltype

import (
	"errors"
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
)

type Send struct {
	origin string
	dest   string
	label  string
	cont   MixedStateGlobalType
}

func (s Send) PossiblePrefixes() []model.Action {
	var prefixes []model.Action
	// First add the send action
	prefixes = append(prefixes, model.Action{
		Label:  s.label,
		Src:    s.origin,
		Dest:   s.dest,
		IsSend: true,
	})
	contPrefixes := s.cont.PossiblePrefixes()
	for _, contPrefix := range contPrefixes {
		if contPrefix.Subject() != s.origin {
			prefixes = append(prefixes, contPrefix)
		}
	}
	return prefixes
}

func (s Send) ConsumePrefix(g *mixedStateGlobalTypeSemanticModel, m model.Action) (MixedStateGlobalType, error) {
	if m.Src == s.origin && m.Dest == s.dest && m.IsSend {
		if m.Label == s.label {
			// Send prefix consumed
			return s.cont, nil
		} else {
			return nil, errors.New("label " + m.Label + " not permitted in the global type " + s.String())
		}
	}
	if s.origin != m.Subject() && s.dest != m.Subject() {
		// Reduction under Prefix
		consumed, err := s.cont.ConsumePrefix(g, m)
		if err != nil {
			return nil, err
		}
		return Send{
			origin: s.origin,
			dest:   s.dest,
			label:  s.label,
			cont:   consumed,
		}, nil
	}
	return nil, errors.New(
		"cannot consume message " + m.String() + " in the global type " + s.String(),
	)
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
	b.WriteString(": <")
	b.WriteString(s.label)
	b.WriteString("> .\n")
	s.cont.stringWithBuilder(b)
}

type Recv struct {
	origin string
	dest   string
	label  string
	cont   MixedStateGlobalType
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

func (r Recv) ConsumePrefix(g *mixedStateGlobalTypeSemanticModel, m model.Action) (MixedStateGlobalType, error) {
	if m.Src == r.origin && m.Dest == r.dest && !m.IsSend {
		return r.cont, nil
	}
	if m.Subject() != r.dest {
		// Reduction under prefix
		newCont, err := r.cont.ConsumePrefix(g, m)
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
	return nil, errors.New(
		"cannot consume message " + m.String() + " in the global type " + r.String(),
	)
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
