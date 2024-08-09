package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func SolveExpressions(raw, flow Literal) (Literal, error) {
	switch raw.kind.name {
	case EXPRESSION:
		result, err := ExecuteCalls(raw.value.([]Call), flow)

		if err != nil {
			return Nada, fmt.Errorf("| inside expression:\n%w", err)
		}

		return result, nil

	case LIST:
		solvedList := []Literal{}

		for _, literal := range raw.value.([]Literal) {
			solved, err := SolveExpressions(literal, flow)

			if err != nil {
				return Nada, nil
			}

			solvedList = append(solvedList, solved)
		}

		spec := DetermineListSpec(solvedList)

		return Literal{solvedList, Kind{name: LIST, spec1: &spec}}, nil

	case PAIR:
		left, err := SolveExpressions(raw.value.(Pair).left, flow)
		if err != nil {
			return Nada, err
		}

		right, err := SolveExpressions(raw.value.(Pair).right, flow)
		if err != nil {
			return Nada, err
		}

		return Literal{
			Pair{left, right},
			Kind{name: PAIR, spec1: &left.kind, spec2: &right.kind},
		}, nil

	default:
		return raw, nil
	}
}

/* Busca encontrar uma função com as exatas constraints passadas*/
func SearchSpecificFunctionVariant(name string, constraints Constraints) (Function, error) {
	for _, function := range FUNCTIONS {
		// if function.name == name &&
		// 	constraints.flow == function.constraints.flow &&
		// 	constraints.parameter == function.constraints.parameter {
		// 	return function, nil
		// }

		if function.name == name &&
			constraints.flow.Format() == function.constraints.flow.Format() &&
			constraints.parameter.Format() == function.constraints.parameter.Format() {
			return function, nil
		}
	}

	return Function{}, fmt.Errorf(
		"\nvariant %s not found for function '%s'",
		FormatConstraint(constraints), name,
	)
}

/* Busca encontrar uma função, podendo esta ser genérica (@any & @any)*/
func SearchCompatibleFunctionVariant(name string, constraints Constraints) (Function, error) {
	found_one := false
	generic, found_generic := Function{}, false

	// fmt.Printf("searching: '%s' against %s\n", name, FormatConstraint(constraints))

	for _, function := range FUNCTIONS {

		if function.name == name {
			// fmt.Printf("found (%s & %s)\n", function.constraints.flow.Format(), function.constraints.parameter.Format())

			found_one = true

			if function.constraints.flow.name == ANY && function.constraints.parameter.name == ANY {
				found_generic = true
				generic = function
				continue
			}

			if function.constraints.flow.Matches(&constraints.flow) &&
				function.constraints.parameter.Matches(&constraints.parameter) {
				return function, nil
			}

			// if function.constraints.flow.MatchesStrict(&constraints.flow) &&
			// 	function.constraints.parameter.MatchesStrict(&constraints.parameter) {
			// 	return function, nil
			// }

			// if function.constraints.flow.Matches(&constraints.flow) &&
			// 	function.constraints.parameter.Matches(&constraints.parameter) {
			// 	generic = function
			// 	found_generic = true
			// }

			// if function.constraints.flow == ANY && function.constraints.parameter == ANY {
			// 	found_generic = true
			// 	generic = function
			// 	continue
			// }

			// if (function.constraints.flow == ANY ||
			// 	function.constraints.flow == constraints.flow) &&
			// 	(function.constraints.parameter == ANY ||
			// 		function.constraints.parameter == constraints.parameter) {

			// 	return function, nil
			// }
		}
	}

	if found_generic {
		return generic, nil
	}

	if found_one {
		return Function{}, fmt.Errorf(
			"\nno variant (%s & %s) found for function '%s'",
			constraints.flow.Format(), constraints.parameter.Format(), name,
		)
	} else {
		return Function{}, fmt.Errorf("\nfunction '%s' not declared", name)
	}
}

func ExecuteFunction(function Function, flow, argument Literal) (Literal, error) {
	// fmt.Printf("calling: %s\n", function.name)

	if function.is_bound {
		// fmt.Printf("bound to: %s\n", FormatLiteral(function.bound_argument))

		ARGUMENT_STACK.Push(function.bound_argument)
		defer ARGUMENT_STACK.Pop()

	} else {
		if argument != Nada {
			// fmt.Printf("pushing argument: %s\n", FormatLiteral(argument))
			ARGUMENT_STACK.Push(argument)
			defer ARGUMENT_STACK.Pop()
		}
	}

	if function.is_recursive {
		FUNCTION_STACK.Push(function)
		defer FUNCTION_STACK.Pop()
	}

	result, err := function.implementation(flow, argument)

	if err != nil {
		return Nada, err
	}

	return result, nil
}

func ExecuteCalls(calls []Call, startingFlow Literal) (Literal, error) {
	flow := startingFlow

	for _, call := range calls {
		if call.functionName == "pipe" || call.functionName == "comment" {
			continue
		}

		solvedArgument, err := SolveExpressions(call.argument, flow)

		if err != nil {
			return Nada, err
		}

		function, err := SearchCompatibleFunctionVariant(call.functionName, Constraints{
			flow.kind, solvedArgument.kind,
		})

		if err != nil {
			return Nada, err
		}

		flow, err = ExecuteFunction(function, flow, solvedArgument)

		if err != nil {
			return Nada, fmt.Errorf("| %s: at %s line %d\n%w",
				call.functionName,
				filepath.Base(call.file),
				call.line,
				err,
			)
		}
	}

	return flow, nil
}

func ExecuteFile(path string) error {
	currentPath := ""

	if len(PATH_STACK.content) == 0 {
		currentPath, _ = os.Getwd()
	}

	currentPath, _ = PATH_STACK.Peek(0)
	absolutePath, _ := filepath.Abs(filepath.Dir(currentPath) + "\\" + path)

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return fmt.Errorf("\ncouldn't find file: %s", absolutePath)
	}

	PATH_STACK.Push(absolutePath)

	file, err := os.ReadFile(absolutePath)
	if err != nil {
		return err
	}

	code, err := Parse(filepath.Base(absolutePath), file)
	if err != nil {
		return err
	}

	_, err = ExecuteCalls(code.([]Call), Nada)

	if err != nil {
		return err
	}

	PATH_STACK.Pop()

	return nil
}
