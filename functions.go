package main

import (
	"fmt"
	"strings"
)

/* Funções intrínsecas da linguagem */

func InitializeFunctions() {
	FUNCTIONS = []Function{
		{
			name:        "val",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) Literal {
				return argument
			},
		},

		//-------------------------------------------------------- IO

		{
			name:        "print",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				fmt.Printf("%s", FormatLiteral(flow))
				return flow
			},
		},

		{
			name:        "println",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				fmt.Printf("%s\n", FormatLiteral(flow))
				return flow
			},
		},

		//-------------------------------------------------------- TEXT

		{
			name:        "conc",
			constraints: Constraint{TEXT, TEXT},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(string) + argument.value.(string), TEXT}
			},
		},

		{
			name:        "chars",
			constraints: Constraint{TEXT, NADA},
			implementation: func(flow, argument Literal) Literal {
				list := []Literal{}

				for _, c := range flow.value.(string) {
					list = append(list, Literal{string(c), TEXT})
				}

				return Literal{list, LIST}
			},
		},

		//-------------------------------------------------------- MATH

		{
			name:        "add",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(int64) + argument.value.(int64), NUMBER}
			},
		},

		{
			name:        "sub",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(int64) - argument.value.(int64), NUMBER}
			},
		},

		{
			name:        "mul",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(int64) * argument.value.(int64), NUMBER}
			},
		},

		{
			name:        "div",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(int64) / argument.value.(int64), NUMBER}
			},
		},

		{
			name:        "mod",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{flow.value.(int64) % argument.value.(int64), NUMBER}
			},
		},

		//-------------------------------------------------------- FUNCTIONS

		{
			name:        "def",
			constraints: Constraint{FUNCTION, TEXT},
			implementation: func(flow, argument Literal) Literal {
				function := flow.value.(Function)
				function.name = argument.value.(string)
				FUNCTIONS = append(FUNCTIONS, function)

				return Literal{function, FUNCTION}
			},
		},

		{
			name:        "takes",
			constraints: Constraint{FUNCTION, PAIR},
			implementation: func(flow, argument Literal) Literal {
				function := flow.value.(Function)
				function.constraints = Constraint{
					argument.value.(Pair).left.value.(string),
					argument.value.(Pair).right.value.(string),
				}

				return Literal{function, FUNCTION}
			},
		},

		{
			name:        "do",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) Literal {
				function := argument.value.(Pair).right.value.(Function)
				actualArgument := argument.value.(Pair).left

				return ExecuteFunction(function, flow, actualArgument)
			},
		},

		{
			name:        "arg",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				if arg, ok := ARGUMENT_STACK.Peek(0); ok {
					return arg
				} else {
					Panic("no argument in the current scope")
					return flow
				}
			},
		},

		{
			name:        "self",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) Literal {
				if function, ok := FUNCTION_STACK.Peek(0); ok {
					return ExecuteFunction(function, flow, argument)
				} else {
					Panic("currently not inside a recursive function")
					return flow
				}
			},
		},

		{
			name:        "bind",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) Literal {
				function := argument.value.(Pair).right.value.(Function)
				function.is_bound = true
				function.bound_argument = argument.value.(Pair).left
				return Literal{function, FUNCTION}
			},
		},

		{
			name:        "idem",
			constraints: Constraint{ANY, FUNCTION},
			implementation: func(flow, argument Literal) Literal {
				if arg, ok := ARGUMENT_STACK.Peek(1); ok {
					// fmt.Printf("current arg: %s", FormatLiteral(arg))
					function := argument.value.(Function)
					function.is_bound = true
					function.bound_argument = arg
					return Literal{function, FUNCTION}
				} else {
					Panic("no argument in current scope")
					return flow
				}
			},
		},

		//-------------------------------------------------------- STACK

		{
			name:        "push",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				STACK.Push(flow)
				return flow
			},
		},

		{
			name:        "pop",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				if STACK.Pop() {
					return flow
				} else {
					Panic("tried to pop empty stack")
					return flow
				}
			},
		},

		{
			name:        "peek",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) Literal {
				index := 0

				if argument.kind == NUMBER {
					index = int(argument.value.(int64))
				}

				if top, ok := STACK.Peek(index); ok {
					return top
				} else {
					Panic("tried to peek empty stack")
					return flow
				}
			},
		},

		{
			name:        "stack",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) Literal {
				return Literal{STACK.content, LIST}
			},
		},

		//-------------------------------------------------------- CONDITIONALS

		{
			name:        "if",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) Literal {
				test := argument.value.(Pair).left.value.(Function)
				arm := argument.value.(Pair).right.value.(Function)

				result := ExecuteFunction(test, flow, Nada)

				if result.value.(bool) {
					return ExecuteFunction(arm, flow, Nada)
				} else {
					return flow
				}
			},
		},

		{
			name:        "less",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{
					flow.value.(int64) < argument.value.(int64),
					NUMBER,
				}
			},
		},

		{
			name:        "greater",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return Literal{
					flow.value.(int64) > argument.value.(int64),
					NUMBER,
				}
			},
		},

		//-------------------------------------------------------- CONTROL FLOW

		{
			name:        "aside",
			constraints: Constraint{ANY, FUNCTION},
			implementation: func(flow, argument Literal) Literal {
				oldFlow := flow
				ExecuteFunction(argument.value.(Function), flow, Nada)
				return oldFlow
			},
		},

		//-------------------------------------------------------- COLLECTIONS

		{
			name:        "left",
			constraints: Constraint{PAIR, NADA},
			implementation: func(flow, argument Literal) Literal {
				return flow.value.(Pair).left
			},
		},

		{
			name:        "right",
			constraints: Constraint{PAIR, NADA},
			implementation: func(flow, argument Literal) Literal {
				return flow.value.(Pair).right
			},
		},

		{
			name:        "len",
			constraints: Constraint{LIST, NADA},
			implementation: func(flow, argument Literal) Literal {
				return Literal{int64(len(flow.value.([]Literal))), NUMBER}
			},
		},

		{
			name:        "index",
			constraints: Constraint{LIST, NUMBER},
			implementation: func(flow, argument Literal) Literal {
				return flow.value.([]Literal)[argument.value.(int64)]
			},
		},

		{
			name:        "append",
			constraints: Constraint{LIST, ANY},
			implementation: func(flow, argument Literal) Literal {
				list := flow.value.([]Literal)
				list = append(list, argument)
				return Literal{list, LIST}
			},
		},

		//-------------------------------------------------------- FILES

		{
			name:        "execute_files",
			constraints: Constraint{ANY, LIST},
			implementation: func(flow, argument Literal) Literal {
				for _, literal := range argument.value.([]Literal) {
					fileName := literal.value.(string)

					if !strings.HasSuffix(fileName, ".fl") {
						fileName += ".fl"
					}

					ExecuteFile(fileName)
				}

				return flow
			},
		},
	}
}
