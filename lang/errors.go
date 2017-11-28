package lang

import "github.com/pkg/errors"

var (
	ErrNotImplemented          = errors.New("Not implemented")
	ErrOddNumberOfMapForms     = errors.New("Map literal must contain an even number of forms")
	ErrIfTooFewArgs            = errors.New("Too few args to if")
	ErrIfTooManyArgs           = errors.New("Too many args to if")
	ErrUnresolvedSymbol        = errors.New("Unable to resolve symbol")
	ErrOddNumberOfBindingForms = errors.New("Binding vector must have an even number of forms")
	ErrUnrecognizedBindingForm = errors.New("Unrecognized binding form")
)
