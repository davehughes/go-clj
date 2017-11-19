package lang

import "fmt"

func showComponent(t string, c interface{}) (interface{}, error) {
	value := fmt.Sprintf("(%s %s)", t, c)
	fmt.Println(value)
	return value, nil
}

type Keyword struct {
	Position
	Name      string
	Namespace string
}

type Symbol struct {
	Position
	Name string
}

type String struct {
	Position
	Value string
}

func (s String) String() string {
	return fmt.Sprintf("(string %v)", s.Value)
}

type Comment struct {
	Position
	Value string
}

type Position struct {
	Line   int
	Column int
	Offset int
}

type CollectionType int

const (
	CollectionList CollectionType = iota
	CollectionVector
	CollectionMap
	CollectionSet
	CollectionLambda
)

type Collection struct {
	Position
	Type  CollectionType
	Items []interface{}
}

type Integer struct {
	Position
	ParsedValue string
	Value       int64
}

type Float struct {
	Position
	ParsedValue string
	Value       float64
}

type QuotedForm struct {
	Position
	Form interface{}
}
