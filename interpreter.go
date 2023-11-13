package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func SolveExpressions(raw, flow Literal) (Literal, error) {
	switch raw.kind {
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
		return Literal{solvedList, LIST}, nil

	case PAIR:
		left, err := SolveExpressions(raw.value.(Pair).left, flow)
		if err != nil {
			return Nada, err
		}

		right, err := SolveExpressions(raw.value.(Pair).right, flow)
		if err != nil {
			return Nada, err
		}

		return Literal{Pair{left, right}, PAIR}, nil

	default:
		return raw, nil
	}
}

func SearchFunction(name string) (Function, bool) {
	for _, function := range FUNCTIONS {
		if function.name == name {
			return function, true
		}
	}
	return Function{}, false
}

func CheckFunctionConstraints(function Function, flow, argument Literal) error {
	if function.constraints.flow != ANY {
		if flow.kind != function.constraints.flow {
			return fmt.Errorf("\nfunction '%s' expected %s in the flow, got %s",
				function.name, function.constraints.flow, flow.kind)
		}
	}

	if function.constraints.parameter != ANY {
		if argument.kind != function.constraints.parameter {
			return fmt.Errorf("\nfunction '%s' expected %s as argument, got %s",
				function.name, function.constraints.parameter, argument.kind)
		}
	}

	return nil
}

func ExecuteFunction(function Function, flow, argument Literal) (Literal, error) {
	// fmt.Printf("calling: %s\n", function.name)

	expandedArgument, err := SolveExpressions(argument, flow)
	if err != nil {
		return Nada, err
	}

	err = CheckFunctionConstraints(function, flow, expandedArgument)
	if err != nil {
		return Nada, err
	}

	if function.is_bound {
		// fmt.Printf("bound to: %s\n", FormatLiteral(function.bound_argument))

		ARGUMENT_STACK.Push(function.bound_argument)
		defer ARGUMENT_STACK.Pop()
	} else {
		if expandedArgument != Nada {
			// fmt.Printf("pushing argument: %s\n", FormatLiteral(expandedArgument))
			ARGUMENT_STACK.Push(expandedArgument)
			defer ARGUMENT_STACK.Pop()
		}
	}

	if function.is_recursive {
		FUNCTION_STACK.Push(function)
		defer FUNCTION_STACK.Pop()
	}

	result, err := function.implementation(flow, expandedArgument)

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

		if function, ok := SearchFunction(call.functionName); ok {
			var err error = nil

			flow, err = ExecuteFunction(function, flow, call.argument)

			if err != nil {
				return Nada, fmt.Errorf("| %s: at %s line %d\n%w",
					call.functionName,
					filepath.Base(call.file),
					call.line,
					err,
				)
			}

		} else {
			return Nada, fmt.Errorf("function %s not found", call.functionName)
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
