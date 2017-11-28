package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO:...
// + Function application/binding?
// + clj forms to bindings?
//   (= foo bar)
//   (and true false)  // or, xor, nand, ...
//
//   (defn and [& forms]
//     (if (boolean (first forms))
//       (apply and (rest forms))
//       false))
//   (defn or [& forms]
//     (if (boolean (first forms))
//       true
//       (apply or (rest forms)))
// + Macros
//   (cond condition1 result1
//         condition2 result2
//         ...
//         :default resultN)
//   (if condition1 result1
//     (if condition2 result2
//       ...
//         resultN))

func TestBasicContextOps(t *testing.T) {
	ctx := NewEvaluationContext()
	ctx = ctx.Bind("test", "hello")
	binding, found := ctx.Resolve("test")
	assert.True(t, found)
	assert.Equal(t, binding, "hello")

	binding, found = ctx.Resolve("unset")
	assert.False(t, found)
}

func TestEval(t *testing.T) {
	ctx := NewEvaluationContext().
		Bind("test", "hello").
		Bind("myfn", NopFunction{})

	testmap := []struct {
		Form   interface{}
		Result interface{}
	}{{
		Form: &Symbol{
			Name: "test",
		},
		Result: "hello",
		// }, {
		// 	Form: &List{
		// 		Forms: []interface{}{
		// 			&Symbol{
		// 				Name: "myfn",
		// 			},
		// 			&String{
		// 				Value: "hello",
		// 			},
		// 		},
		// 	},
		// 	Result: nil,
	}}
	for _, test := range testmap {
		result, err := Eval(test.Form, ctx)
		assert.NoError(t, err)
		assert.Equal(t, result, test.Result)
	}
}

func TestBindOne(t *testing.T) {
	var err error
	ctx := NewEvaluationContext()
	ctx, err = ctx.BindForm(&Symbol{Name: "test"}, &String{Value: "hello"})
	assert.NoError(t, err)

	binding, found := ctx.Resolve("test")
	assert.True(t, found)
	assert.Equal(t, (binding.(*String)).Value, "hello")
}

func TestBindMulti(t *testing.T) {
	var err error
	ctx := NewEvaluationContext()
	ctx, err = ctx.BindForms(
		&Symbol{Name: "foo"}, &String{Value: "hello"},
		&Symbol{Name: "bar"}, &Integer{Value: 123},
	)
	assert.NoError(t, err)

	foo, found := ctx.Resolve("foo")
	assert.True(t, found)
	assert.Equal(t, (foo.(*String)).Value, "hello")

	bar, found := ctx.Resolve("bar")
	assert.True(t, found)
	assert.Equal(t, (bar.(*Integer)).Value, int64(123))
}

func TestLet(t *testing.T) {
	var err error
	var args []interface{}
	var result interface{}
	let := LetExpression{}
	ctx := NewEvaluationContext()

	// (let []) -> nil
	args = []interface{}{
		// bindings
		&Vector{
			Forms: []interface{}{},
		},

		// expressions (empty)
	}
	result, err = let.Apply(ctx, args)
	assert.NoError(t, err)
	assert.IsType(t, &Nil{}, result)

	// (let [x 1]) -> nil
	args = []interface{}{
		// bindings
		&Vector{
			Forms: []interface{}{
				&Symbol{Name: "x"},
				&Integer{Value: 1},
			},
		},

		// expressions (empty)
	}
	result, err = let.Apply(ctx, args)
	assert.NoError(t, err)
	assert.IsType(t, &Nil{}, result)

	// (let [x 1 y "hello"] y) -> "hello"
	args = []interface{}{
		// bindings
		&Vector{
			Forms: []interface{}{
				&Symbol{Name: "x"},
				&Integer{Value: 1},
				&Symbol{Name: "y"},
				&String{Value: "hello"},
			},
		},

		// expressions
		&Symbol{Name: "y"},
	}
	result, err = let.Apply(ctx, args)
	assert.NoError(t, err)
	assert.IsType(t, &String{}, result)
	assert.Equal(t, "hello", (result.(*String)).Value)
}

func TestIf(t *testing.T) {
	var err error
	var args []interface{}
	var result interface{}
	ifExpr := &IfExpression{}
	ctx := NewEvaluationContext()

	// (if true "ok") -> "ok"
	args = []interface{}{
		&True{},
		&String{Value: "ok"},
	}
	result, err = ifExpr.Apply(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, "ok", (result.(*String)).Value)

	// (if false "ok") -> nil
	args = []interface{}{
		&False{},
		&String{Value: "ok"},
	}
	result, err = ifExpr.Apply(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, &Nil{}, result)

	// (if true "yes" "no") -> "yes"
	args = []interface{}{
		&True{},
		&String{Value: "yes"},
		&String{Value: "no"},
	}
	result, err = ifExpr.Apply(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, "yes", (result.(*String)).Value)

	// (if false "yes" "no") -> "no"
	args = []interface{}{
		&False{},
		&String{Value: "yes"},
		&String{Value: "no"},
	}
	result, err = ifExpr.Apply(ctx, args)
	assert.NoError(t, err)
	assert.Equal(t, "no", (result.(*String)).Value)
}

func TestEvalIf(t *testing.T) {
	ctx := NewEvaluationContext().
		Bind("if", &IfExpression{})

	ast, err := ReadString("(if true :yes :no)")
	assert.NoError(t, err)
	assert.Len(t, ast, 1)

	result, err := Eval(ast[0], ctx)
	assert.NoError(t, err)
	assert.IsType(t, &Keyword{}, result)
	assert.Equal(t, "yes", (result.(*Keyword)).Name)
}

func TestNestedLetIf(t *testing.T) {
	ctx := NewEvaluationContext().
		Bind("if", &IfExpression{}).
		Bind("let", &LetExpression{})

	ast, err := ReadString("(let [predicate? (let [x 1] true)] (if predicate? :yes :no))")
	assert.NoError(t, err)
	assert.Len(t, ast, 1)

	result, err := Eval(ast[0], ctx)
	assert.NoError(t, err)
	assert.IsType(t, &Keyword{}, result)
	assert.Equal(t, "yes", (result.(*Keyword)).Name)
}

func TestParsing(t *testing.T) {

	var ast []interface{}
	var err error

	ast, err = ReadString("nil")
	assert.NoError(t, err)
	assert.IsType(t, &Nil{}, ast[0])

	ast, err = ReadString(":foo")
	assert.NoError(t, err)
	assert.IsType(t, &Keyword{}, ast[0])

	ast, err = ReadString("foo")
	assert.NoError(t, err)
	assert.IsType(t, &Symbol{}, ast[0])

	// Quoted forms are converted to lists equivalent to '(quote *form*)'
	ast, err = ReadString("'(foo bar baz)")
	assert.NoError(t, err)
	assert.IsType(t, ast[0], &List{})
	list := ast[0].(*List)
	assert.Len(t, list.Forms, 2)
	assert.Equal(t, &Symbol{Name: "quote"}, list.Forms[0])
	assert.IsType(t, &List{}, list.Forms[1])

	ast, err = ReadString("'baz")
	assert.NoError(t, err)
	assert.IsType(t, &List{}, ast[0])

	ast, err = ReadString("()")
	assert.NoError(t, err)
	assert.IsType(t, &List{}, ast[0])

	ast, err = ReadString("(foo bar baz)")
	assert.NoError(t, err)
	assert.IsType(t, &List{}, ast[0])

	// Map literals
	ast, err = ReadString("{foo bar baz quux}")
	assert.NoError(t, err)
	assert.IsType(t, &Map{}, ast[0])

	// Maps need to have even numbers of forms
	ast, err = ReadString("{foo bar baz}")
	assert.Error(t, err)

	// Vector literals
	ast, err = ReadString("[]")
	assert.NoError(t, err)
	assert.IsType(t, &Vector{}, ast[0])

	ast, err = ReadString("[foo, bar, baz]")
	assert.NoError(t, err)
	vector, ok := ast[0].(*Vector)
	assert.True(t, ok)
	assert.Len(t, vector.Forms, 3)
	assert.IsType(t, &Symbol{}, vector.Forms[0])
	assert.Equal(t, "foo", ((vector.Forms[0]).(*Symbol)).Name)
	assert.IsType(t, &Symbol{}, vector.Forms[1])
	assert.Equal(t, "bar", ((vector.Forms[1]).(*Symbol)).Name)
	assert.IsType(t, &Symbol{}, vector.Forms[2])
	assert.Equal(t, "baz", ((vector.Forms[2]).(*Symbol)).Name)

	// Set literals
	ast, err = ReadString("#{foo bar baz}")
	assert.NoError(t, err)
	assert.IsType(t, &Set{}, ast[0])

	ast, err = ReadString("#(func-of %)")
	assert.NoError(t, err)
	assert.IsType(t, &Lambda{}, ast[0])

	ast, err = ReadString("#(func-of %1)")
	assert.NoError(t, err)
	assert.IsType(t, &Lambda{}, ast[0])

	ast, err = ReadString("#(func-of %2 %1)")
	assert.NoError(t, err)
	assert.IsType(t, &Lambda{}, ast[0])

	// ast, err = ReadString("1.29e8")
	// assert.NoError(t, err)
	// assert.IsType(t, &Float{}, ast[0])

	ast, err = ReadString("3")
	assert.NoError(t, err)
	assert.IsType(t, &Integer{}, ast[0])

	// String literals
	ast, err = ReadString("\"\"")
	assert.NoError(t, err)
	assert.IsType(t, &String{}, ast[0])

	ast, err = ReadString("\"hello world!\"")
	assert.NoError(t, err)
	assert.IsType(t, &String{}, ast[0])

	ast, err = ReadString("\"multi\nline\nstring\"")
	assert.NoError(t, err)
	assert.IsType(t, &String{}, ast[0])

	ast, err = ReadString("\"multi\r\nline\n\rstring\nwith\rmultiple line endings\"")
	assert.NoError(t, err)
	assert.IsType(t, &String{}, ast[0])

	// Comments
	ast, err = ReadString("; comment1\n\t; comment2\n")
	assert.NoError(t, err)
	assert.Len(t, ast, 2)
	assert.IsType(t, &Comment{}, ast[0])
	assert.IsType(t, &Comment{}, ast[1])

	ast, err = ReadString("(+ 1 ; comment until EOL\n 2)")
	assert.NoError(t, err)
	assert.Len(t, ast, 1)
	assert.IsType(t, &List{}, ast[0])
	forms := ((ast[0]).(*List)).Forms
	assert.Len(t, forms, 4)
	assert.IsType(t, &Symbol{}, forms[0])
	assert.IsType(t, &Integer{}, forms[1])
	assert.IsType(t, &Comment{}, forms[2])
	assert.IsType(t, &Integer{}, forms[3])

	// Mixed usage
	ast, err = ReadString("nil foo '(quoted list)\n; here's a line comment\n")
	assert.NoError(t, err)
	assert.Len(t, ast, 4)
	assert.IsType(t, &Nil{}, ast[0])
	assert.IsType(t, &Symbol{}, ast[1])
	assert.IsType(t, &List{}, ast[2])
	assert.IsType(t, &Comment{}, ast[3])
}

func TestSymbolParsing(t *testing.T) {
	symbols := []string{
		"foo",
		"kebab-case",
		"containing:colon",
		"containing/slash",
		"in-tha-dogg#pound",
		"puppies+rainbows",
		"is-it?",
		"it-is!",
		"<!",
		"!>",
		"foo-bar/baz3?<!%>*",
		"#YOLO",
		// "2%milk", TODO: support leading numerals?
	}

	for _, symbol := range symbols {
		ast, err := ReadString(symbol)
		assert.NoError(t, err)
		assert.IsType(t, &Symbol{}, ast[0])
	}
}

func TestArityMatching(t *testing.T) {
	testCases := []struct {
		Input             string
		NumPositionalArgs int
		HasVariadicArgs   bool
	}{
		{
			Input:             "[]",
			NumPositionalArgs: 0,
			HasVariadicArgs:   false,
		},
		{
			Input:             "[x]",
			NumPositionalArgs: 1,
			HasVariadicArgs:   false,
		},
		{
			Input:             "[x y]",
			NumPositionalArgs: 2,
			HasVariadicArgs:   false,
		},
		{
			Input:             "[& varargs]",
			NumPositionalArgs: 0,
			HasVariadicArgs:   true,
		},
		{
			Input:             "[x & varargs]",
			NumPositionalArgs: 1,
			HasVariadicArgs:   true,
		},
		// Advanced destructuring
		{
			Input:             "[{:keys [a b c]}]",
			NumPositionalArgs: 1,
			HasVariadicArgs:   false,
		},
	}

	for _, testCase := range testCases {
		form, err := ReadForm(testCase.Input)
		assert.NoError(t, err)

		sig, err := functionSignature(form)
		require.NoError(t, err)
		assert.Equal(t, testCase.NumPositionalArgs, len(sig.positionalBindings))
		assert.Equal(t, testCase.HasVariadicArgs, sig.variadicBinding != nil)
	}
}

func TestArgumentBinding(t *testing.T) {
	testCases := []struct {
		signatureInput string
		argsInput      string
		expectMismatch bool
		checkBindings  map[string]interface{}
	}{
		{
			signatureInput: "[]",
			argsInput:      "[]",
			checkBindings:  map[string]interface{}{},
		},
		// {
		// 	signatureInput: "[x]",
		// 	argsInput:      "[3]",
		// 	checkBindings: map[string]interface{}{
		// 		"x": &Integer{Value: 1},
		// 	},
		// },
		// {
		// 	signatureInput: "[& v]",
		// 	argsInput:      "[1 2 3]",
		// 	checkBindings: map[string]interface{}{
		// 		"v": &List{
		// 			Forms: []interface{}{
		// 				&Integer{Value: 1},
		// 				&Integer{Value: 2},
		// 				&Integer{Value: 3},
		// 			},
		// 		},
		// 	},
		// },
	}

	for _, testCase := range testCases {
		signatureVec, err := ReadVector(testCase.signatureInput)
		require.NoError(t, err)

		argsVec, err := ReadVector(testCase.argsInput)
		require.NoError(t, err)

		ctx := NewEvaluationContext()
		signature, err := functionSignature(signatureVec)
		require.NoError(t, err)

		ctx, err = signature.bindArgs(ctx, argsVec.Forms)
		if testCase.expectMismatch {
			assert.Error(t, err)
			continue
		} else {
			require.NoError(t, err)
		}

		if testCase.checkBindings != nil {
			for k, v := range testCase.checkBindings {
				value, ok := ctx.Resolve(k)
				assert.True(t, ok)
				assert.Equal(t, v, value)
			}
		}
	}
}

func TestFnDefinition(t *testing.T) {
	fnInput := "((fn [x y] y) 1 2)"
	ctx := NewBaseEvaluationContext()

	form := MustReadForm(fnInput)
	result, err := Eval(form, ctx)
	assert.NoError(t, err)

	resultInt, ok := result.(*Integer)
	assert.True(t, ok)
	assert.Equal(t, int64(2), resultInt.Value)
}
