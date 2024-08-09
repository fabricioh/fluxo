package main

const (
	ANY        = "any"
	NADA       = "nada"
	TEXT       = "text"
	PAIR       = "pair"
	LIST       = "list"
	TYPE       = "type"
	NUMBER     = "number"
	LOGICAL    = "logical"
	FUNCTION   = "function"
	EXPRESSION = "expression"
)

type Literal struct {
	value any
	kind  Kind
}

type Pair struct {
	left, right Literal
}

type Call struct {
	functionName string
	argument     Literal
	file         string
	line         int
}

type Function struct {
	name           string
	body           []Call
	constraints    Constraints
	is_recursive   bool
	is_bound       bool
	bound_argument Literal
	file           string
	line           int
	implementation func(Literal, Literal) (Literal, error)
}

type Constraints struct {
	flow, parameter Kind
}

type Kind struct {
	name         string
	spec1, spec2 *Kind
}

func (k Kind) Format() string {
	result := "@" + k.name

	if k.spec1 != nil {
		result += "(" + k.spec1.Format()

		if k.spec2 != nil {
			result += " & " + k.spec2.Format()
		}

		result += ")"
	}

	return result
}

func (k *Kind) Matches(other *Kind) bool {
	if k.name == ANY {
		return true
	}

	switch k.name {
	case LIST:
		return other.name == LIST && k.spec1.Matches(other.spec1)
	case PAIR:
		return other.name == PAIR && k.spec1.Matches(other.spec1) && k.spec2.Matches(other.spec2)
	default:
		return k.name == other.name
	}
}

type Stack[T any] struct {
	content []T
}

func (s *Stack[T]) Push(new T) {
	s.content = append([]T{new}, s.content...)
}

func (s *Stack[T]) Pop() bool {
	if len(s.content) == 0 {
		return false
	} else {
		s.content = s.content[1:]
		return true
	}
}

func (s *Stack[T]) Peek(index int) (T, bool) {
	if len(s.content) > 0 && index <= len(s.content)-1 {
		return s.content[index], true
	} else {
		var temp T
		return temp, false
	}
}

var FUNCTIONS = []Function{}
var VARIABLES = map[string]Literal{}
var STACK = Stack[Literal]{}
var ARGUMENT_STACK = Stack[Literal]{}
var FUNCTION_STACK = Stack[Function]{}
var PATH_STACK = Stack[string]{}
var Nada = Literal{"nada", Kind{name: NADA}}
