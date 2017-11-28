package lang

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Go code goes here...

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 7, col: 1, offset: 31},
			expr: &actionExpr{
				pos: position{line: 7, col: 10, offset: 40},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 7, col: 10, offset: 40},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 7, col: 10, offset: 40},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 7, col: 16, offset: 46},
								expr: &ruleRefExpr{
									pos:  position{line: 7, col: 16, offset: 46},
									name: "Form",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 7, col: 22, offset: 52},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Form",
			pos:  position{line: 11, col: 1, offset: 83},
			expr: &actionExpr{
				pos: position{line: 11, col: 9, offset: 91},
				run: (*parser).callonForm1,
				expr: &seqExpr{
					pos: position{line: 11, col: 9, offset: 91},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 11, col: 9, offset: 91},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 11, col: 11, offset: 93},
							label: "form",
							expr: &choiceExpr{
								pos: position{line: 11, col: 18, offset: 100},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 11, col: 18, offset: 100},
										name: "LineComment",
									},
									&ruleRefExpr{
										pos:  position{line: 11, col: 32, offset: 114},
										name: "List",
									},
									&ruleRefExpr{
										pos:  position{line: 11, col: 39, offset: 121},
										name: "Lambda",
									},
									&ruleRefExpr{
										pos:  position{line: 11, col: 48, offset: 130},
										name: "Vector",
									},
									&ruleRefExpr{
										pos:  position{line: 11, col: 57, offset: 139},
										name: "Set",
									},
									&ruleRefExpr{
										pos:  position{line: 11, col: 63, offset: 145},
										name: "Map",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 13, offset: 164},
										name: "QuotedForm",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 26, offset: 177},
										name: "AtomDeref",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 38, offset: 189},
										name: "Numeric",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 48, offset: 199},
										name: "Symbol",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 57, offset: 208},
										name: "Keyword",
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 67, offset: 218},
										name: "String",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 13, col: 15, offset: 239},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 18, col: 1, offset: 306},
			expr: &notExpr{
				pos: position{line: 18, col: 7, offset: 312},
				expr: &anyMatcher{
					line: 18, col: 8, offset: 313,
				},
			},
		},
		{
			name: "EOL",
			pos:  position{line: 19, col: 1, offset: 315},
			expr: &seqExpr{
				pos: position{line: 19, col: 8, offset: 322},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 19, col: 8, offset: 322},
						expr: &ruleRefExpr{
							pos:  position{line: 19, col: 8, offset: 322},
							name: "LineComment",
						},
					},
					&choiceExpr{
						pos: position{line: 19, col: 23, offset: 337},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 19, col: 23, offset: 337},
								val:        "\r\n",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 19, col: 32, offset: 346},
								val:        "\n\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 19, col: 41, offset: 355},
								val:        "\r",
								ignoreCase: false,
							},
							&litMatcher{
								pos:        position{line: 19, col: 48, offset: 362},
								val:        "\n",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 19, col: 55, offset: 369},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 20, col: 1, offset: 374},
			expr: &zeroOrMoreExpr{
				pos: position{line: 20, col: 19, offset: 392},
				expr: &charClassMatcher{
					pos:        position{line: 20, col: 19, offset: 392},
					val:        "[ \\n\\t\\r,]",
					chars:      []rune{' ', '\n', '\t', '\r', ','},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "LineComment",
			pos:  position{line: 22, col: 1, offset: 405},
			expr: &actionExpr{
				pos: position{line: 22, col: 16, offset: 420},
				run: (*parser).callonLineComment1,
				expr: &seqExpr{
					pos: position{line: 22, col: 16, offset: 420},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 22, col: 16, offset: 420},
							expr: &ruleRefExpr{
								pos:  position{line: 22, col: 16, offset: 420},
								name: "_",
							},
						},
						&litMatcher{
							pos:        position{line: 22, col: 19, offset: 423},
							val:        ";",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 22, col: 23, offset: 427},
							label: "comment",
							expr: &zeroOrMoreExpr{
								pos: position{line: 22, col: 31, offset: 435},
								expr: &charClassMatcher{
									pos:        position{line: 22, col: 31, offset: 435},
									val:        "[^\\r\\n]",
									chars:      []rune{'\r', '\n'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "List",
			pos:  position{line: 29, col: 1, offset: 608},
			expr: &actionExpr{
				pos: position{line: 29, col: 9, offset: 616},
				run: (*parser).callonList1,
				expr: &seqExpr{
					pos: position{line: 29, col: 9, offset: 616},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 29, col: 9, offset: 616},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 29, col: 13, offset: 620},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 29, col: 19, offset: 626},
								expr: &ruleRefExpr{
									pos:  position{line: 29, col: 19, offset: 626},
									name: "Form",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 29, col: 25, offset: 632},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Vector",
			pos:  position{line: 36, col: 1, offset: 800},
			expr: &actionExpr{
				pos: position{line: 36, col: 11, offset: 810},
				run: (*parser).callonVector1,
				expr: &seqExpr{
					pos: position{line: 36, col: 11, offset: 810},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 36, col: 11, offset: 810},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 36, col: 15, offset: 814},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 36, col: 21, offset: 820},
								expr: &ruleRefExpr{
									pos:  position{line: 36, col: 21, offset: 820},
									name: "Form",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 36, col: 27, offset: 826},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Map",
			pos:  position{line: 43, col: 1, offset: 995},
			expr: &actionExpr{
				pos: position{line: 43, col: 8, offset: 1002},
				run: (*parser).callonMap1,
				expr: &seqExpr{
					pos: position{line: 43, col: 8, offset: 1002},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 43, col: 8, offset: 1002},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 43, col: 12, offset: 1006},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 43, col: 18, offset: 1012},
								expr: &ruleRefExpr{
									pos:  position{line: 43, col: 18, offset: 1012},
									name: "Form",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 43, col: 24, offset: 1018},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Set",
			pos:  position{line: 54, col: 1, offset: 1303},
			expr: &actionExpr{
				pos: position{line: 54, col: 8, offset: 1310},
				run: (*parser).callonSet1,
				expr: &seqExpr{
					pos: position{line: 54, col: 8, offset: 1310},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 54, col: 8, offset: 1310},
							val:        "#",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 54, col: 12, offset: 1314},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 54, col: 16, offset: 1318},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 54, col: 22, offset: 1324},
								expr: &ruleRefExpr{
									pos:  position{line: 54, col: 22, offset: 1324},
									name: "Form",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 54, col: 28, offset: 1330},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Lambda",
			pos:  position{line: 61, col: 1, offset: 1496},
			expr: &actionExpr{
				pos: position{line: 61, col: 11, offset: 1506},
				run: (*parser).callonLambda1,
				expr: &seqExpr{
					pos: position{line: 61, col: 11, offset: 1506},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 61, col: 11, offset: 1506},
							val:        "#",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 61, col: 15, offset: 1510},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 61, col: 19, offset: 1514},
							label: "forms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 61, col: 25, offset: 1520},
								expr: &ruleRefExpr{
									pos:  position{line: 61, col: 25, offset: 1520},
									name: "Form",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 61, col: 31, offset: 1526},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AtomDeref",
			pos:  position{line: 68, col: 1, offset: 1691},
			expr: &actionExpr{
				pos: position{line: 68, col: 14, offset: 1704},
				run: (*parser).callonAtomDeref1,
				expr: &seqExpr{
					pos: position{line: 68, col: 14, offset: 1704},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 68, col: 14, offset: 1704},
							val:        "@",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 68, col: 18, offset: 1708},
							label: "form",
							expr: &ruleRefExpr{
								pos:  position{line: 68, col: 23, offset: 1713},
								name: "Form",
							},
						},
					},
				},
			},
		},
		{
			name: "Keyword",
			pos:  position{line: 80, col: 1, offset: 1965},
			expr: &choiceExpr{
				pos: position{line: 80, col: 14, offset: 1978},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 80, col: 14, offset: 1978},
						name: "NamespacedKeyword",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 34, offset: 1998},
						name: "BareKeyword",
					},
				},
			},
		},
		{
			name: "NamespacedKeyword",
			pos:  position{line: 82, col: 1, offset: 2012},
			expr: &actionExpr{
				pos: position{line: 82, col: 22, offset: 2033},
				run: (*parser).callonNamespacedKeyword1,
				expr: &seqExpr{
					pos: position{line: 82, col: 22, offset: 2033},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 82, col: 22, offset: 2033},
							val:        ":",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 82, col: 26, offset: 2037},
							val:        ":",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 82, col: 30, offset: 2041},
							label: "sym",
							expr: &ruleRefExpr{
								pos:  position{line: 82, col: 34, offset: 2045},
								name: "Symbol",
							},
						},
					},
				},
			},
		},
		{
			name: "BareKeyword",
			pos:  position{line: 90, col: 1, offset: 2251},
			expr: &actionExpr{
				pos: position{line: 90, col: 16, offset: 2266},
				run: (*parser).callonBareKeyword1,
				expr: &seqExpr{
					pos: position{line: 90, col: 16, offset: 2266},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 90, col: 16, offset: 2266},
							val:        ":",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 90, col: 20, offset: 2270},
							label: "sym",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 24, offset: 2274},
								name: "Symbol",
							},
						},
					},
				},
			},
		},
		{
			name: "Symbol",
			pos:  position{line: 97, col: 1, offset: 2441},
			expr: &actionExpr{
				pos: position{line: 97, col: 11, offset: 2451},
				run: (*parser).callonSymbol1,
				expr: &oneOrMoreExpr{
					pos: position{line: 97, col: 11, offset: 2451},
					expr: &charClassMatcher{
						pos:        position{line: 97, col: 11, offset: 2451},
						val:        "[a-z0-9-+*?!&^%#?<>%/\\\\]i",
						chars:      []rune{'-', '+', '*', '?', '!', '&', '^', '%', '#', '?', '<', '>', '%', '/', '\\'},
						ranges:     []rune{'a', 'z', '0', '9'},
						ignoreCase: true,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "QuotedForm",
			pos:  position{line: 114, col: 1, offset: 2883},
			expr: &actionExpr{
				pos: position{line: 114, col: 15, offset: 2897},
				run: (*parser).callonQuotedForm1,
				expr: &seqExpr{
					pos: position{line: 114, col: 15, offset: 2897},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 114, col: 17, offset: 2899},
							val:        "'",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 114, col: 23, offset: 2905},
							label: "form",
							expr: &ruleRefExpr{
								pos:  position{line: 114, col: 28, offset: 2910},
								name: "Form",
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 126, col: 1, offset: 3162},
			expr: &actionExpr{
				pos: position{line: 126, col: 11, offset: 3172},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 126, col: 11, offset: 3172},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 126, col: 11, offset: 3172},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 126, col: 15, offset: 3176},
							label: "value",
							expr: &zeroOrMoreExpr{
								pos: position{line: 126, col: 21, offset: 3182},
								expr: &charClassMatcher{
									pos:        position{line: 126, col: 21, offset: 3182},
									val:        "[^\"]",
									chars:      []rune{'"'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
						&litMatcher{
							pos:        position{line: 126, col: 27, offset: 3188},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Numeric",
			pos:  position{line: 133, col: 1, offset: 3346},
			expr: &choiceExpr{
				pos: position{line: 133, col: 14, offset: 3359},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 133, col: 14, offset: 3359},
						name: "SciNot",
					},
					&ruleRefExpr{
						pos:  position{line: 133, col: 23, offset: 3368},
						name: "Float",
					},
					&ruleRefExpr{
						pos:  position{line: 133, col: 30, offset: 3375},
						name: "Integer",
					},
				},
			},
		},
		{
			name: "SciNotSuffix",
			pos:  position{line: 135, col: 1, offset: 3386},
			expr: &seqExpr{
				pos: position{line: 135, col: 17, offset: 3402},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 135, col: 17, offset: 3402},
						val:        "[eE]",
						chars:      []rune{'e', 'E'},
						ignoreCase: false,
						inverted:   false,
					},
					&ruleRefExpr{
						pos:  position{line: 135, col: 22, offset: 3407},
						name: "Integer",
					},
				},
			},
		},
		{
			name: "SciNot",
			pos:  position{line: 137, col: 1, offset: 3416},
			expr: &actionExpr{
				pos: position{line: 137, col: 11, offset: 3426},
				run: (*parser).callonSciNot1,
				expr: &seqExpr{
					pos: position{line: 137, col: 11, offset: 3426},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 137, col: 11, offset: 3426},
							label: "coefficient",
							expr: &choiceExpr{
								pos: position{line: 137, col: 25, offset: 3440},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 137, col: 25, offset: 3440},
										name: "Float",
									},
									&ruleRefExpr{
										pos:  position{line: 137, col: 33, offset: 3448},
										name: "Integer",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 137, col: 43, offset: 3458},
							label: "exponent",
							expr: &ruleRefExpr{
								pos:  position{line: 137, col: 52, offset: 3467},
								name: "SciNotSuffix",
							},
						},
					},
				},
			},
		},
		{
			name: "Float",
			pos:  position{line: 143, col: 1, offset: 3598},
			expr: &actionExpr{
				pos: position{line: 143, col: 10, offset: 3607},
				run: (*parser).callonFloat1,
				expr: &seqExpr{
					pos: position{line: 143, col: 10, offset: 3607},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 143, col: 10, offset: 3607},
							expr: &litMatcher{
								pos:        position{line: 143, col: 10, offset: 3607},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 143, col: 15, offset: 3612},
							expr: &charClassMatcher{
								pos:        position{line: 143, col: 15, offset: 3612},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 143, col: 22, offset: 3619},
							val:        ".",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 143, col: 26, offset: 3623},
							expr: &charClassMatcher{
								pos:        position{line: 143, col: 26, offset: 3623},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 155, col: 1, offset: 3913},
			expr: &actionExpr{
				pos: position{line: 155, col: 12, offset: 3924},
				run: (*parser).callonInteger1,
				expr: &seqExpr{
					pos: position{line: 155, col: 12, offset: 3924},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 155, col: 12, offset: 3924},
							expr: &litMatcher{
								pos:        position{line: 155, col: 12, offset: 3924},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 155, col: 17, offset: 3929},
							expr: &charClassMatcher{
								pos:        position{line: 155, col: 17, offset: 3929},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
	},
}

func (c *current) onInput1(forms interface{}) (interface{}, error) {
	return forms, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["forms"])
}

func (c *current) onForm1(form interface{}) (interface{}, error) {
	// fmt.Printf("Form: %v\n", form)
	return form, nil
}

func (p *parser) callonForm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onForm1(stack["form"])
}

func (c *current) onLineComment1(comment interface{}) (interface{}, error) {
	return &Comment{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Value:    fmt.Sprint(comment),
	}, nil
}

func (p *parser) callonLineComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineComment1(stack["comment"])
}

func (c *current) onList1(forms interface{}) (interface{}, error) {
	return &List{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms:    forms.([]interface{}),
	}, nil
}

func (p *parser) callonList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList1(stack["forms"])
}

func (c *current) onVector1(forms interface{}) (interface{}, error) {
	return &Vector{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms:    forms.([]interface{}),
	}, nil
}

func (p *parser) callonVector1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVector1(stack["forms"])
}

func (c *current) onMap1(forms interface{}) (interface{}, error) {
	numForms := len(forms.([]interface{}))
	if numForms%2 != 0 {
		return nil, ErrOddNumberOfMapForms
	}
	return &Map{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms:    forms.([]interface{}),
	}, nil
}

func (p *parser) callonMap1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMap1(stack["forms"])
}

func (c *current) onSet1(forms interface{}) (interface{}, error) {
	return &Set{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms:    forms.([]interface{}),
	}, nil
}

func (p *parser) callonSet1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSet1(stack["forms"])
}

func (c *current) onLambda1(forms interface{}) (interface{}, error) {
	return &Lambda{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms:    forms.([]interface{}),
	}, nil
}

func (p *parser) callonLambda1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLambda1(stack["forms"])
}

func (c *current) onAtomDeref1(form interface{}) (interface{}, error) {
	return &List{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms: []interface{}{
			&Symbol{
				Name: "deref",
			},
			form,
		},
	}, nil
}

func (p *parser) callonAtomDeref1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAtomDeref1(stack["form"])
}

func (c *current) onNamespacedKeyword1(sym interface{}) (interface{}, error) {
	return &Keyword{
		Position:  Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Name:      (sym.(*Symbol)).Name,
		Namespace: "TODO: namespaced",
	}, nil
}

func (p *parser) callonNamespacedKeyword1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespacedKeyword1(stack["sym"])
}

func (c *current) onBareKeyword1(sym interface{}) (interface{}, error) {
	return &Keyword{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Name:     (sym.(*Symbol)).Name,
	}, nil
}

func (p *parser) callonBareKeyword1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBareKeyword1(stack["sym"])
}

func (c *current) onSymbol1() (interface{}, error) {
	pos := Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset}
	switch string(c.text) {
	case "nil":
		return &Nil{Position: pos}, nil
	case "true":
		return &True{Position: pos}, nil
	case "false":
		return &False{Position: pos}, nil
	default:
		return &Symbol{
			Position: pos,
			Name:     string(c.text),
		}, nil
	}
}

func (p *parser) callonSymbol1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSymbol1()
}

func (c *current) onQuotedForm1(form interface{}) (interface{}, error) {
	return &List{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Forms: []interface{}{
			&Symbol{
				Name: "quote",
			},
			form,
		},
	}, nil
}

func (p *parser) callonQuotedForm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedForm1(stack["form"])
}

func (c *current) onString1(value interface{}) (interface{}, error) {
	return &String{
		Position: Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		Value:    string(c.text),
	}, nil
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1(stack["value"])
}

func (c *current) onSciNot1(coefficient, exponent interface{}) (interface{}, error) {
	fl := fmt.Sprintf("(sci-not %s :exponent %s)", coefficient, exponent)
	fmt.Println(fl)
	return fl, nil
}

func (p *parser) callonSciNot1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSciNot1(stack["coefficient"], stack["exponent"])
}

func (c *current) onFloat1() (interface{}, error) {
	value, err := strconv.ParseFloat(string(c.text), 64)
	if err != nil {
		panic(err)
	}
	return &Float{
		Position:    Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		ParsedValue: string(c.text),
		Value:       value,
	}, nil
}

func (p *parser) callonFloat1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFloat1()
}

func (c *current) onInteger1() (interface{}, error) {
	value, err := strconv.ParseInt(string(c.text), 0, 64)
	if err != nil {
		panic(err)
	}
	return &Integer{
		Position:    Position{Line: c.pos.line, Column: c.pos.col, Offset: c.pos.offset},
		ParsedValue: string(c.text),
		Value:       value,
	}, nil
}

func (p *parser) callonInteger1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}



