package lang

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

func Read(r io.Reader) ([]interface{}, error) {
	forms, err := ParseReader("<repl>", r)
	return forms.([]interface{}), err
}

func ReadString(s string) ([]interface{}, error) {
	return Read(strings.NewReader(s))
}

func ReadForm(s string) (interface{}, error) {
	forms, err := ReadString(s)
	if err != nil {
		return nil, err
	}
	if len(forms) != 1 {
		return nil, errors.New("ReadForm expects to read one form, found zero or many")
	}
	return forms[0], nil
}

func MustReadForm(s string) interface{} {
	form, err := ReadForm(s)
	if err != nil {
		panic(err)
	}
	return form
}

func ReadList(s string) (*List, error) {
	form, err := ReadForm(s)
	if err != nil {
		return nil, err
	}
	list, ok := form.(*List)
	if !ok {
		return nil, errors.New("Unable to read list from specified input")
	}
	return list, nil
}

func MustReadList(s string) *List {
	list, err := ReadList(s)
	if err != nil {
		panic(err)
	}
	return list
}

func ReadVector(s string) (*Vector, error) {
	form, err := ReadForm(s)
	if err != nil {
		return nil, err
	}
	vector, ok := form.(*Vector)
	if !ok {
		return nil, errors.New("Unable to read vector from specified input")
	}
	return vector, nil
}

func MustReadVector(s string) *Vector {
	vector, err := ReadVector(s)
	if err != nil {
		panic(err)
	}
	return vector
}
