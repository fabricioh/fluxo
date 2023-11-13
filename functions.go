package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

/* Funções intrínsecas da linguagem */

func InitializeFunctions() {
	FUNCTIONS = []Function{
		{
			name:        "val",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) (Literal, error) {
				return argument, nil
			},
		},

		//-------------------------------------------------------- IO

		{
			name:        "print",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				fmt.Printf("%s", FormatLiteral(flow))
				return flow, nil
			},
		},

		{
			name:        "println",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				fmt.Printf("%s\n", FormatLiteral(flow))
				return flow, nil
			},
		},

		//-------------------------------------------------------- TEXT

		{
			name:        "conc",
			constraints: Constraint{TEXT, TEXT},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(string) + argument.value.(string), TEXT}, nil
			},
		},

		{
			name:        "chars",
			constraints: Constraint{TEXT, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				list := []Literal{}

				for _, c := range flow.value.(string) {
					list = append(list, Literal{string(c), TEXT})
				}

				return Literal{list, LIST}, nil
			},
		},

		//-------------------------------------------------------- MATH

		{
			name:        "add",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(int64) + argument.value.(int64), NUMBER}, nil
			},
		},

		{
			name:        "sub",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(int64) - argument.value.(int64), NUMBER}, nil
			},
		},

		{
			name:        "mul",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(int64) * argument.value.(int64), NUMBER}, nil
			},
		},

		{
			name:        "div",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(int64) / argument.value.(int64), NUMBER}, nil
			},
		},

		{
			name:        "mod",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{flow.value.(int64) % argument.value.(int64), NUMBER}, nil
			},
		},

		//-------------------------------------------------------- FUNCTIONS

		{
			name:        "def",
			constraints: Constraint{FUNCTION, TEXT},
			implementation: func(flow, argument Literal) (Literal, error) {
				function := flow.value.(Function)
				function.name = argument.value.(string)
				FUNCTIONS = append(FUNCTIONS, function)

				return Literal{function, FUNCTION}, nil
			},
		},

		{
			name:        "takes",
			constraints: Constraint{FUNCTION, PAIR},
			implementation: func(flow, argument Literal) (Literal, error) {
				function := flow.value.(Function)
				function.constraints = Constraint{
					argument.value.(Pair).left.value.(string),
					argument.value.(Pair).right.value.(string),
				}

				return Literal{function, FUNCTION}, nil
			},
		},

		{
			name:        "do",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) (Literal, error) {
				function := argument.value.(Pair).right.value.(Function)
				actualArgument := argument.value.(Pair).left

				return ExecuteFunction(function, flow, actualArgument)
			},
		},

		{
			name:        "arg",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				if arg, ok := ARGUMENT_STACK.Peek(0); ok {
					return arg, nil
				} else {
					return Nada, errors.New("\nno argument in the current scope")
				}
			},
		},

		{
			name:        "self",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) (Literal, error) {
				if function, ok := FUNCTION_STACK.Peek(0); ok {
					return ExecuteFunction(function, flow, argument)
				} else {
					return flow, errors.New("\ncannot call 'self' if not inside a recursive function")
				}
			},
		},

		{
			name:        "bind",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) (Literal, error) {
				function := argument.value.(Pair).right.value.(Function)
				function.is_bound = true
				function.bound_argument = argument.value.(Pair).left
				return Literal{function, FUNCTION}, nil
			},
		},

		{
			name:        "idem",
			constraints: Constraint{ANY, FUNCTION},
			implementation: func(flow, argument Literal) (Literal, error) {
				if arg, ok := ARGUMENT_STACK.Peek(1); ok {
					// fmt.Printf("current arg: %s", FormatLiteral(arg))
					function := argument.value.(Function)
					function.is_bound = true
					function.bound_argument = arg
					return Literal{function, FUNCTION}, nil
				} else {
					return flow, errors.New("\nno argument in current scope")
				}
			},
		},

		//-------------------------------------------------------- STACK

		{
			name:        "push",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				STACK.Push(flow)
				return flow, nil
			},
		},

		{
			name:        "pop",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				if STACK.Pop() {
					return flow, nil
				} else {
					return Nada, errors.New("\ncannot pop empty stack")
				}
			},
		},

		{
			name:        "peek",
			constraints: Constraint{ANY, ANY},
			implementation: func(flow, argument Literal) (Literal, error) {
				index := 0

				if argument.kind == NUMBER {
					index = int(argument.value.(int64))
				}

				if top, ok := STACK.Peek(index); ok {
					return top, nil
				} else {
					return flow, errors.New("\ncannot peek empty stack")
				}
			},
		},

		{
			name:        "stack",
			constraints: Constraint{ANY, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{STACK.content, LIST}, nil
			},
		},

		//-------------------------------------------------------- CONDITIONALS

		{
			name:        "if",
			constraints: Constraint{ANY, PAIR},
			implementation: func(flow, argument Literal) (Literal, error) {
				test := argument.value.(Pair).left.value.(Function)
				arm := argument.value.(Pair).right.value.(Function)

				result, err := ExecuteFunction(test, flow, Nada)

				if err != nil {
					return Nada, err
				}

				if result.value.(bool) {
					return ExecuteFunction(arm, flow, Nada)
				} else {
					return flow, nil
				}
			},
		},

		{
			name:        "less",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{
					flow.value.(int64) < argument.value.(int64),
					NUMBER,
				}, nil
			},
		},

		{
			name:        "less_eql",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{
					flow.value.(int64) <= argument.value.(int64),
					NUMBER,
				}, nil
			},
		},

		{
			name:        "grt",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{
					flow.value.(int64) > argument.value.(int64),
					NUMBER,
				}, nil
			},
		},

		{
			name:        "grt_eql",
			constraints: Constraint{NUMBER, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{
					flow.value.(int64) >= argument.value.(int64),
					NUMBER,
				}, nil
			},
		},

		{
			name:        "not",
			constraints: Constraint{LOGICAL, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{!flow.value.(bool), LOGICAL}, nil
			},
		},

		//-------------------------------------------------------- CONTROL FLOW

		{
			name:        "aside",
			constraints: Constraint{ANY, FUNCTION},
			implementation: func(flow, argument Literal) (Literal, error) {
				oldFlow := flow
				_, err := ExecuteFunction(argument.value.(Function), flow, Nada)

				if err != nil {
					return Nada, err
				}

				return oldFlow, nil
			},
		},

		{
			name:        "wait",
			constraints: Constraint{ANY, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				time.Sleep(time.Duration(argument.value.(int64)) * time.Millisecond)
				return flow, nil
			},
		},

		{
			name:        "exit",
			constraints: Constraint{ANY, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				os.Exit(int(argument.value.(int64)))
				return flow, nil
			},
		},

		{
			name:        "panic",
			constraints: Constraint{ANY, TEXT},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Nada, errors.New(argument.value.(string))
			},
		},

		//-------------------------------------------------------- COLLECTIONS

		{
			name:        "left",
			constraints: Constraint{PAIR, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				return flow.value.(Pair).left, nil
			},
		},

		{
			name:        "right",
			constraints: Constraint{PAIR, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				return flow.value.(Pair).right, nil
			},
		},

		{
			name:        "len",
			constraints: Constraint{LIST, NADA},
			implementation: func(flow, argument Literal) (Literal, error) {
				return Literal{int64(len(flow.value.([]Literal))), NUMBER}, nil
			},
		},

		{
			name:        "index",
			constraints: Constraint{LIST, NUMBER},
			implementation: func(flow, argument Literal) (Literal, error) {
				if int(argument.value.(int64)) > len(flow.value.([]Literal))-1 {
					return Nada, errors.New("\nindex out of bounds")
				}

				return flow.value.([]Literal)[argument.value.(int64)], nil
			},
		},

		{
			name:        "append",
			constraints: Constraint{LIST, ANY},
			implementation: func(flow, argument Literal) (Literal, error) {
				list := flow.value.([]Literal)
				list = append(list, argument)
				return Literal{list, LIST}, nil
			},
		},

		//-------------------------------------------------------- TYPE SYSTEM

		// {
		// 	name:        "model",
		// 	constraints: Constraint{ANY, ANY},
		// 	implementation: func(flow, argument Literal) (Literal, error) {
		// 		return Literal{DetermineModel(flow), TEXT}
		// 	},
		// },

		{
			name:        "as",
			constraints: Constraint{ANY, TEXT},
			implementation: func(flow, argument Literal) (Literal, error) {
				switch argument.value.(string) {
				case TEXT:
					return Literal{FormatLiteral(flow), TEXT}, nil

				default:
					return Nada, errors.New("\ninvalid type name passed to function 'as'")
				}
			},
		},

		//-------------------------------------------------------- FILES

		{
			name:        "exec",
			constraints: Constraint{ANY, LIST},
			implementation: func(flow, argument Literal) (Literal, error) {
				for _, literal := range argument.value.([]Literal) {
					fileName := literal.value.(string)

					if !strings.HasSuffix(fileName, ".fl") {
						fileName += ".fl"
					}

					err := ExecuteFile(fileName)

					if err != nil {
						return Nada, err
					}
				}

				return flow, nil
			},
		},
	}
}
