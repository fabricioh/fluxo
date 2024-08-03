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
	kind  string
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
var VARIABLES = map[string]Literal{}
var STACK = Stack[Literal]{}
var ARGUMENT_STACK = Stack[Literal]{}
var FUNCTION_STACK = Stack[Function]{}
var PATH_STACK = Stack[string]{}
var Nada = Literal{"nada", NADA}
