package main

const (
	ANY        = "ANY"
	NADA       = "NADA"
	TEXT       = "TEXT"
	PAIR       = "PAIR"
	LIST       = "LIST"
	TYPE       = "TYPE"
	NUMBER     = "NUMBER"
	LOGICAL    = "LOGICAL"
	FUNCTION   = "FUNCTION"
	EXPRESSION = "EXPRESSION"
)

type Literal struct {
	value any
	kind  string
}

type Pair struct {
	left, right Literal
}

type Call struct {
	functionName string
	argument     Literal
}

type Function struct {
	name           string
	body           []Call
	constraints    Constraint
	is_recursive   bool
	is_bound       bool
	bound_argument Literal
	implementation func(Literal, Literal) Literal
}

type Constraint struct {
	flow, parameter string
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
var STACK = Stack[Literal]{}
var ARGUMENT_STACK = Stack[Literal]{}
var FUNCTION_STACK = Stack[Function]{}
var Nada = Literal{"nada", NADA}
