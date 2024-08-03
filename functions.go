package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func InitializeFunctions() {

	/* Funções intrínsecas da linguagem */

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "val",
		constraints: Constraints{ANY, ANY},
		implementation: func(flow, argument Literal) (Literal, error) {
			return argument, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "val",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return flow, nil
		},
	})

	//---------------------------------------------------------------------------- IO

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "print",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			fmt.Printf("%s", FormatLiteral(flow))
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "print",
		constraints: Constraints{ANY, ANY},
		implementation: func(flow, argument Literal) (Literal, error) {
			fmt.Printf("%s", FormatLiteral(argument))
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "println",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			fmt.Printf("%s\n", FormatLiteral(flow))
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "println",
		constraints: Constraints{ANY, ANY},
		implementation: func(flow, argument Literal) (Literal, error) {
			fmt.Printf("%s\n", FormatLiteral(argument))
			return flow, nil
		},
	})

	//---------------------------------------------------------------------------- TEXT

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "conc",
		constraints: Constraints{TEXT, TEXT},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(string) + argument.value.(string), TEXT}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "chars",
		constraints: Constraints{TEXT, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			list := []Literal{}

			for _, c := range flow.value.(string) {
				list = append(list, Literal{string(c), TEXT})
			}

			return Literal{list, LIST}, nil
		},
	})

	//---------------------------------------------------------------------------- OPERAÇÕES MATEMÁTICAS

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "add",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(int64) + argument.value.(int64), NUMBER}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "sub",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(int64) - argument.value.(int64), NUMBER}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "mul",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(int64) * argument.value.(int64), NUMBER}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "div",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(int64) / argument.value.(int64), NUMBER}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "mod",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.value.(int64) % argument.value.(int64), NUMBER}, nil
		},
	})

	//---------------------------------------------------------------------------- VARIÁVEIS

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "set",
		constraints: Constraints{ANY, TEXT},
		implementation: func(flow, argument Literal) (Literal, error) {
			VARIABLES[argument.value.(string)] = flow
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "get",
		constraints: Constraints{ANY, TEXT},
		implementation: func(flow, argument Literal) (Literal, error) {
			val, ok := VARIABLES[argument.value.(string)]

			if !ok {
				return Nada, fmt.Errorf("\nvariable '%s' not defined", argument.value.(string))
			}

			return val, nil
		},
	})

	//---------------------------------------------------------------------------- FUNÇÕES

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "def",
		constraints: Constraints{FUNCTION, TEXT},
		implementation: func(flow, argument Literal) (Literal, error) {
			function := flow.value.(Function)
			name := argument.value.(string)

			function.name = name

			_, err := SearchSpecificFunctionVariant(name, function.constraints)

			// fmt.Printf("\nfound: %v\n", FormatLiteral(Literal{f, FUNCTION}))

			if err == nil {
				return Nada, fmt.Errorf(
					"\nvariant %s for function '%s' already defined",
					FormatConstraint(function.constraints), name,
				)
			}

			FUNCTIONS = append(FUNCTIONS, function)

			return Literal{function, FUNCTION}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "constraints",
		constraints: Constraints{FUNCTION, PAIR},
		implementation: func(flow, argument Literal) (Literal, error) {
			function := flow.value.(Function)
			function.constraints = Constraints{
				argument.value.(Pair).left.value.(string),
				argument.value.(Pair).right.value.(string),
			}

			// fmt.Printf("%v", FormatLiteral(Literal{function, FUNCTION}))

			return Literal{function, FUNCTION}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "do",
		constraints: Constraints{ANY, PAIR},
		implementation: func(flow, argument Literal) (Literal, error) {
			function := argument.value.(Pair).right.value.(Function)
			actualArgument := argument.value.(Pair).left

			return ExecuteFunction(function, flow, actualArgument)
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "arg",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			if arg, ok := ARGUMENT_STACK.Peek(0); ok {
				return arg, nil
			} else {
				return Nada, errors.New("\nno argument in the current scope")
			}
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "self",
		constraints: Constraints{ANY, ANY},
		implementation: func(flow, argument Literal) (Literal, error) {
			if function, ok := FUNCTION_STACK.Peek(0); ok {
				return ExecuteFunction(function, flow, argument)
			} else {
				return flow, errors.New("\ncannot call 'self' if not inside a recursive function")
			}
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "bind",
		constraints: Constraints{ANY, PAIR},
		implementation: func(flow, argument Literal) (Literal, error) {
			function := argument.value.(Pair).right.value.(Function)
			function.is_bound = true
			function.bound_argument = argument.value.(Pair).left
			return Literal{function, FUNCTION}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "idem",
		constraints: Constraints{ANY, FUNCTION},
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
	})

	//---------------------------------------------------------------------------- STACK

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "push",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			STACK.Push(flow)
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "pop",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			if STACK.Pop() {
				return flow, nil
			} else {
				return Nada, errors.New("\ncannot pop empty stack")
			}
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "peek",
		constraints: Constraints{ANY, ANY},
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
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "stack",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{STACK.content, LIST}, nil
		},
	})

	//---------------------------------------------------------------------------- CONDICIONAIS

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "if",
		constraints: Constraints{ANY, PAIR},
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
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "case",
		constraints: Constraints{ANY, LIST},
		implementation: func(flow, argument Literal) (Literal, error) {
			for _, elem := range argument.value.([]Literal) {
				test := elem.value.(Pair).left.value.(Function)
				arm := elem.value.(Pair).right.value.(Function)

				result, err := ExecuteFunction(test, flow, Nada)

				if err != nil {
					return Nada, err
				}

				if result.value.(bool) {
					return ExecuteFunction(arm, flow, Nada)
				}
			}

			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "and",
		constraints: Constraints{ANY, LIST},
		implementation: func(flow, argument Literal) (Literal, error) {
			for _, elem := range argument.value.([]Literal) {
				result, err := ExecuteFunction(elem.value.(Function), flow, Nada)

				if err != nil {
					return Nada, err
				}

				if !result.value.(bool) {
					return Literal{false, LOGICAL}, nil
				}
			}

			return Literal{true, LOGICAL}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "or",
		constraints: Constraints{ANY, LIST},
		implementation: func(flow, argument Literal) (Literal, error) {
			for _, elem := range argument.value.([]Literal) {
				result, err := ExecuteFunction(elem.value.(Function), flow, Nada)

				if err != nil {
					return Nada, err
				}

				if result.value.(bool) {
					return Literal{true, LOGICAL}, nil
				}
			}

			return Literal{false, LOGICAL}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "less",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{
				flow.value.(int64) < argument.value.(int64),
				NUMBER,
			}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "less_eql",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{
				flow.value.(int64) <= argument.value.(int64),
				NUMBER,
			}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "grt",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{
				flow.value.(int64) > argument.value.(int64),
				NUMBER,
			}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "grt_eql",
		constraints: Constraints{NUMBER, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{
				flow.value.(int64) >= argument.value.(int64),
				NUMBER,
			}, nil
		},
	})

	// FUNCTIONS = append(FUNCTIONS, Function{
	// 	name:        "eql",
	// 	constraints: Constraints{LOGICAL, NADA},
	// 	implementation: func(flow, argument Literal) (Literal, error) {

	// 	},
	// })

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "not",
		constraints: Constraints{LOGICAL, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{!flow.value.(bool), LOGICAL}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "else",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{true, LOGICAL}, nil
		},
	})

	//---------------------------------------------------------------------------- CONTROLE DE FLUXO

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "aside",
		constraints: Constraints{ANY, FUNCTION},
		implementation: func(flow, argument Literal) (Literal, error) {
			oldFlow := flow
			_, err := ExecuteFunction(argument.value.(Function), flow, Nada)

			if err != nil {
				return Nada, err
			}

			return oldFlow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "wait",
		constraints: Constraints{ANY, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			time.Sleep(time.Duration(argument.value.(int64)) * time.Millisecond)
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "exit",
		constraints: Constraints{ANY, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			os.Exit(int(argument.value.(int64)))
			return flow, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "panic",
		constraints: Constraints{ANY, TEXT},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Nada, errors.New(argument.value.(string))
		},
	})

	//---------------------------------------------------------------------------- LIST E PAIR

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "left",
		constraints: Constraints{PAIR, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return flow.value.(Pair).left, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "right",
		constraints: Constraints{PAIR, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return flow.value.(Pair).right, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "len",
		constraints: Constraints{LIST, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{int64(len(flow.value.([]Literal))), NUMBER}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "index",
		constraints: Constraints{LIST, NUMBER},
		implementation: func(flow, argument Literal) (Literal, error) {
			if int(argument.value.(int64)) > len(flow.value.([]Literal))-1 {
				return Nada, errors.New("\nindex out of bounds")
			}

			return flow.value.([]Literal)[argument.value.(int64)], nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "append",
		constraints: Constraints{LIST, ANY},
		implementation: func(flow, argument Literal) (Literal, error) {
			list := flow.value.([]Literal)
			list = append(list, argument)
			return Literal{list, LIST}, nil
		},
	})

	//---------------------------------------------------------------------------- TIPOS

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "type",
		constraints: Constraints{ANY, NADA},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.kind, TYPE}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "is",
		constraints: Constraints{ANY, TYPE},
		implementation: func(flow, argument Literal) (Literal, error) {
			return Literal{flow.kind == argument.value.(string), LOGICAL}, nil
		},
	})

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "as",
		constraints: Constraints{ANY, TYPE},
		implementation: func(flow, argument Literal) (Literal, error) {
			switch argument.value.(string) {
			case TEXT:
				return Literal{FormatLiteral(flow), TEXT}, nil

			case NUMBER:
				if flow.kind == TEXT {
					number, err := strconv.ParseInt(flow.value.(string), 10, 64)

					if err != nil {
						return Nada, err
					}

					return Literal{number, NUMBER}, nil

				} else {
					return Nada, fmt.Errorf("\ncan't convert @%s to @number", flow.kind)
				}

			case NADA:
				return Nada, nil

			case ANY, PAIR, TYPE, LIST, LOGICAL, FUNCTION:
				return Nada, fmt.Errorf("\ncan't convert @%s to @%s", flow.kind, argument.value.(string))
			}

			return Nada, fmt.Errorf("\nno such type as @%s", argument.value.(string))
		},
	})

	//---------------------------------------------------------------------------- ARQUIVOS

	FUNCTIONS = append(FUNCTIONS, Function{
		name:        "exec",
		constraints: Constraints{ANY, LIST},
		implementation: func(flow, argument Literal) (Literal, error) {
			for _, literal := range argument.value.([]Literal) {
				if literal.kind != TEXT {
					return Nada, fmt.Errorf(
						"\nfunction 'exec' expected a @list of @text, found element of type @%s",
						literal.kind,
					)
				}

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
	})
}
