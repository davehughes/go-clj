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

func (s Keyword) String() string {
	return fmt.Sprintf("(keyword %v)", s.Name)
}

type Symbol struct {
	Position
	Name string
}

func (s Symbol) String() string {
	return fmt.Sprintf("(symbol %v)", s.Name)
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

type Container struct {
	Position
	Forms []interface{}
}

type List Container
type Vector Container
type Map Container
type Set Container
type Lambda Container

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

type ReservedWord struct {
	Position
}

type Nil ReservedWord
type True ReservedWord
type False ReservedWord
