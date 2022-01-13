package mixedstateglobaltype

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"strings"
)

type Choice struct {
	choicer string
	choices []MixedStateGlobalType
}

func (c Choice) PossiblePrefixes() []model.Action {
	prefixes := make([]model.Action, 0)
	for _, choice := range c.choices {
		prefixes = append(prefixes, choice.PossiblePrefixes()...)
	}
	return prefixes
}

func (c Choice) ConsumePrefix(g *mixedStateGlobalTypeSemanticModel, message model.Action) (MixedStateGlobalType, error) {
	if message.Subject() == c.choicer {
		// Choicer reduction
		success := false
		successIdx := -1
		var cont MixedStateGlobalType = nil
		for idx, choice := range c.choices {
			next, err := choice.ConsumePrefix(g, message)
			if err == nil {
				success = true
				cont = next
				successIdx = idx
				break
			}
		}
		if success {
			residuals := make([][]model.Action, 0)
			for idx, choice := range c.choices {
				if idx != successIdx {
					for _, residual := range choice.ResidualActions(c.choicer) {
						if len(residual) > 0 {
							residuals = append(residuals, residual)
						}
					}
				}
			}
			g.AddResidualActions(residuals)
			return cont, nil
		}
	} else {
		// Non-choicer reduction
		for idx, choice := range c.choices {
			next, err := choice.ConsumePrefix(g, message)
			if err == nil {
				newChoices := make([]MixedStateGlobalType, len(c.choices))
				copy(newChoices, c.choices)
				newChoices[idx] = next
				return Choice{
					choicer: c.choicer,
					choices: newChoices,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no choice branch can reduce %s", message.String())
}

func (c Choice) IsDone() bool {
	return false
}

func (c Choice) String() string {
	var b strings.Builder
	c.stringWithBuilder(&b)
	return b.String()
}

func (c Choice) stringWithBuilder(b *strings.Builder) {
	b.WriteString(c.choicer)
	b.WriteString(" CHOICE: ")
	b.WriteString("{\n")
	for _, cont := range c.choices {
		cont.stringWithBuilder(b)
		b.WriteString(";\n")
	}
	b.WriteString("}\n")
}

func (c Choice) ResidualActions(choicer string) [][]model.Action {
	res := make([][]model.Action, 0)
	for _, choice := range c.choices {
		res = append(res, choice.ResidualActions(choicer)...)
	}
	return res
}
