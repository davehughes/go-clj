package lang

import "github.com/pkg/errors"

type EvaluationContext struct {
	Bindings map[string]interface{}
}

func Eval(form interface{}, ctx *EvaluationContext) (interface{}, error) {
	switch t := form.(type) {
	case Symbol:
		binding, ok := ctx.Bindings[t.Name]
		if !ok {
			return nil, errors.New("Unbound symbol")
		}
		return binding, nil
	default:
		return nil, errors.New("Unrecognized form")
	}
}
