package lang

import (
	"fmt"

	"github.com/pkg/errors"
)

type EvaluationContext struct {
	bindings map[string]interface{}
}

func (ctx *EvaluationContext) clone() *EvaluationContext {
	bindings := make(map[string]interface{}, len(ctx.bindings))
	for k, v := range ctx.bindings {
		bindings[k] = v
	}

	return &EvaluationContext{
		bindings: bindings,
	}
}

func NewEvaluationContext() *EvaluationContext {
	return &EvaluationContext{
		bindings: map[string]interface{}{},
	}
}

func NewBaseEvaluationContext() *EvaluationContext {
	return NewEvaluationContext().
		Bind("if", &IfExpression{}).
		Bind("let", &LetExpression{}).
		Bind("fn", &FnExpression{})
}

func (ctx *EvaluationContext) Bind(name string, value interface{}) *EvaluationContext {
	// TODO: naive cloning will be a performance problem, rework with persistent vectors
	nextCtx := ctx.clone()
	nextCtx.bindings[name] = value
	return nextCtx
}

func (ctx *EvaluationContext) Resolve(name string) (interface{}, bool) {
	value, ok := ctx.bindings[name]
	return value, ok
}

func (ctx *EvaluationContext) BindForm(bindingForm interface{}, bindingValue interface{}) (*EvaluationContext, error) {
	// TODO: add support for destructured binding
	switch t := bindingForm.(type) {
	case *Symbol:
		evaluatedValue, err := Eval(bindingValue, ctx)
		if err != nil {
			return nil, err
		}
		return ctx.Bind(t.Name, evaluatedValue), nil
	default:
		return nil, ErrUnrecognizedBindingForm
	}
}

func (ctx *EvaluationContext) BindForms(bindingForms ...interface{}) (*EvaluationContext, error) {
	if len(bindingForms)%2 != 0 {
		return nil, ErrOddNumberOfMapForms
	}

	var err error
	for i := 0; i < len(bindingForms); i += 2 {
		ctx, err = ctx.BindForm(bindingForms[i], bindingForms[i+1])
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func (ctx *EvaluationContext) BindArgs(args []interface{}, positionalBindings []interface{}, variadicBinding interface{}) (*EvaluationContext, error) {
	var err error

	// Check for args-signature mismatch
	if (variadicBinding == nil && len(args) != len(positionalBindings)) ||
		(variadicBinding != nil && len(args) < len(positionalBindings)) {
		return nil, errors.New("Wrong number of args passed")
	}

	// Bind positional arguments
	for idx, binding := range positionalBindings {
		ctx, err = ctx.BindForm(binding, args[idx])
		if err != nil {
			return nil, err
		}
	}

	// Bind variadic arguments if appropriate
	if variadicBinding != nil {
		var bindingValue interface{} = &Nil{}
		if len(args) > len(positionalBindings) {
			bindingValue = &List{Forms: args[len(positionalBindings):]}
		}
		ctx, err = ctx.BindForm(variadicBinding, bindingValue)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

type IFunction interface {
	Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error)
}

type IMacro interface {
	Expand(forms []interface{}) (interface{}, error)
}

type NopFunction struct{}

func (f NopFunction) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	return nil, nil
}

// Special forms;
// def
// if
// let
// fn
// loop/recur

type IfExpression struct{}

func (e IfExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	// Should have 2 or 3 args (if 2, the third is implicitly nil
	if len(args) < 2 {
		return nil, ErrIfTooFewArgs
	}
	if len(args) > 3 {
		return nil, ErrIfTooManyArgs
	}
	predicate, trueValue := args[0], args[1]
	var falseValue interface{} = &Nil{}
	if len(args) == 3 {
		falseValue = args[2]
	}

	predicateResult, err := Eval(predicate, ctx)
	if err != nil {
		return nil, err
	}
	if CoerceToBoolean(predicateResult) {
		return Eval(trueValue, ctx)
	} else {
		return Eval(falseValue, ctx)
	}
}

type LetExpression struct{}

func (e LetExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	var err error

	bindings, expressions := args[0], args[1:]
	bindingsVec, ok := bindings.(*Vector)
	if !ok {
		return nil, errors.New("Let requires a vector for its binding")
	}
	ctx, err = ctx.BindForms(bindingsVec.Forms...)
	if err != nil {
		return nil, err
	}

	var returnValue interface{} = &Nil{}
	for _, expression := range expressions {
		returnValue, err = Eval(expression, ctx)
		if err != nil {
			return nil, err
		}
	}
	return returnValue, nil
}

type FnExpression struct{}

// Returns an Fn object that can be used to apply the defined function to a set of
// args at a callsite.
func (e FnExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	// if args is a [vec] followed by expressions, index as a single 'overload'
	// if args is a list of ([vec] body) forms, index as multiple overloads

	// representation of arity?
	// + # positional + varargs (?)
	// + [] -> 0, false
	// + [& xs] -> 0, true
	// + [x y] -> 2, false
	// + [x & xs] -> 1, true
	// + varargs is set to nil if none are provided, or a List of forms if provided
	// + if a previous overload has the same arity, raise an exception
	//   -> (fn ([x] 1) ([x] 2))
	//   -> "can't have 2 overloads with same arity"
	//   -> (fn ([x] 1) ([x y] 2) ([x & args] 3))
	//      "can't have fixed arity function with more params than variadic functions"
	//   -> (fn ([x & args] 3) ([x y & args] 4))
	//      "can't have more than 1 variadic overload"

	// TODO: support multiple overloads; for now, assume there's just one
	if len(args) < 1 {
		return nil, errors.New("Function overload must provide an args vector")
	}
	signatureVec, ok := args[0].(*Vector)
	if !ok {
		return nil, errors.New("Function overload must provide an args vector")
	}
	signature, err := functionSignature(signatureVec)
	if err != nil {
		return nil, err
	}

	return &Fn{
		overloads: []*fnOverload{
			&fnOverload{
				signature:   signature,
				expressions: args[1:],
			},
		},
	}, nil
}

type Fn struct {
	overloads []*fnOverload // TODO: variants of different arities
}

func (f Fn) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	var err error

	// TODO: identify variant based on args arity
	overload := f.overloads[0]

	// Bind args, creating merged context
	ctx, err = overload.signature.bindArgs(ctx, args)
	if err != nil {
		return nil, err
	}

	// Evaluate forms, returning last value
	var returnValue interface{} = nil
	for _, expr := range overload.expressions {
		returnValue, err = Eval(expr, ctx)
	}
	return returnValue, nil
}

type fnOverload struct {
	signature   *fnSignature
	docstring   string
	expressions []interface{}
}

type fnSignature struct {
	positionalBindings []interface{}
	variadicBinding    interface{}
}

func (s *fnSignature) bindArgs(ctx *EvaluationContext, args []interface{}) (*EvaluationContext, error) {
	return ctx.BindArgs(args, s.positionalBindings, s.variadicBinding)
}

func functionSignature(form interface{}) (*fnSignature, error) {
	vector, ok := form.(*Vector)
	if !ok {
		return nil, errors.New("Function arity requires a vector of forms")
	}

	sig := &fnSignature{
		positionalBindings: []interface{}{},
		variadicBinding:    nil,
	}
	for idx, form := range vector.Forms {
		// TODO: implement destructuring
		switch t := form.(type) {
		case *Symbol:
			if t.Name == "&" {
				if idx == (len(vector.Forms) - 2) {
					sig.variadicBinding = vector.Forms[idx+1]
					return sig, nil
				} else {
					// This is an odd edge case, but Clojure seems to treat trailing ampersand the same
					// as if it were omitted.
					// TODO: handle cases like (fn [x & foo bar] ...)
					return sig, nil
				}
			} else {
				sig.positionalBindings = append(sig.positionalBindings, t)
			}
		case *Vector, *List, *Map:
			sig.positionalBindings = append(sig.positionalBindings, t)
		default:
			return nil, errors.New("Unsupported binding form")
		}
	}
	return sig, nil
}

type DefExpression struct{}

func (e DefExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	// (def symbol doc-string? init?)
	return nil, ErrNotImplemented
}

type LoopExpression struct{}

func (e LoopExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	return nil, ErrNotImplemented
}

type RecurExpression struct{}

func (e RecurExpression) Apply(ctx *EvaluationContext, args []interface{}) (interface{}, error) {
	return nil, ErrNotImplemented
}

func Eval(form interface{}, ctx *EvaluationContext) (interface{}, error) {
	switch t := form.(type) {

	case *Symbol:
		binding, ok := ctx.Resolve(t.Name)
		if !ok {
			return nil, ErrUnresolvedSymbol
		}
		return binding, nil

	case *List:
		if len(t.Forms) == 0 {
			// empty list
			return []interface{}{}, nil
		}
		first, rest := t.Forms[0], t.Forms[1:]
		fmt.Printf("Processing list eval, first: %v, rest %v\n", first, rest)
		firstEval, err := Eval(first, ctx)
		if err != nil {
			return nil, err
		}

		// first item should be a symbol that resolves to an IFn
		// (function, macro, or special form)
		// remaining items are args/kwargs

		fmt.Printf("Evaluated first form %v -> %v\n", first, firstEval)
		fun, ok := firstEval.(IFunction)
		if !ok {
			return nil, errors.New("Cannot cast form to IFunction in list eval")
		}
		fmt.Printf("TODO: apply function %v to args %v\n", fun, rest)
		result, err := fun.Apply(ctx, rest)
		if err != nil {
			return nil, errors.New("Error encountered in function invocation")
		}
		fmt.Printf("Result: %v\n", result)
		return result, err

	// Literals eval as themselves
	case *String, *Integer, *Keyword, *True, *False, *Nil:
		return t, nil

	default:
		return nil, ErrNotImplemented
	}
}

func CoerceToBoolean(x interface{}) bool {
	switch t := x.(type) {
	case *True, True:
		return true
	case *False, False:
		return false
	case *Nil, Nil:
		return false
	case *Integer:
		return t.Value == 1
	default:
		return true
	}
}

// (defn function-name doc-string? args body)
// + macro:
//   (def function-name (fn doc-string? args body))
