package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func SolveExpressions(raw, flow Literal) Literal {
	switch raw.kind {
	case EXPRESSION:
		return ExecuteCalls(raw.value.([]Call), flow)
	case LIST:
		solvedList := []Literal{}
		for _, literal := range raw.value.([]Literal) {
			solvedList = append(solvedList, SolveExpressions(literal, flow))
		}
		return Literal{solvedList, LIST}
	case PAIR:
		return Literal{Pair{
			SolveExpressions(raw.value.(Pair).left, flow),
			SolveExpressions(raw.value.(Pair).right, flow),
		}, PAIR}
	default:
		return raw
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

func CheckFunctionConstraints(function Function, flow, argument Literal) bool {
	if function.constraints.flow != ANY {
		if flow.kind != function.constraints.flow {
			Panic(fmt.Sprintf(
				"function '%s' expected %s in the flow, got %s\n",
				function.name, function.constraints.flow, flow.kind,
			))
			return false
		}
	}

	if function.constraints.parameter != ANY {
		if argument.kind != function.constraints.parameter {
			Panic(fmt.Sprintf(
				"function '%s' expected %s as argument, got %s\n",
				function.name, function.constraints.parameter, argument.kind,
			))
			return false
		}
	}

	return true
}

func ExecuteFunction(function Function, flow, argument Literal) Literal {
	// fmt.Printf("calling: %s\n", function.name)

	expandedArgument := SolveExpressions(argument, flow)

	if !CheckFunctionConstraints(function, flow, expandedArgument) {
		Panic(fmt.Sprintf("function %s's constraints not met!!!!", function.name))
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

	return function.implementation(flow, expandedArgument)
}

func ExecuteCalls(calls []Call, startingFlow Literal) Literal {
	flow := startingFlow

	for _, call := range calls {
		if call.functionName == "pipe" || call.functionName == "comment" {
			continue
		}

		if function, ok := SearchFunction(call.functionName); ok {
			flow = ExecuteFunction(function, flow, call.argument)
		} else {
			Panic(fmt.Sprintf("function %s not found!!!!", call.functionName))
		}
	}

	return flow
}

func ExecuteFile(path string) {
	currentPath := ""

	if len(PATH_STACK.content) == 0 {
		currentPath, _ = os.Getwd()
	}

	currentPath, _ = PATH_STACK.Peek(0)
	absolutePath, _ := filepath.Abs(filepath.Dir(currentPath) + "\\" + path)

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		Panic("couldn't find file: " + absolutePath)
	}

	PATH_STACK.Push(absolutePath)
	defer PATH_STACK.Pop()

	file, err := os.ReadFile(absolutePath)
	if err != nil {
		Panic(err.Error())
	}

	code, err := Parse(filepath.Base(absolutePath), file)
	if err != nil {
		Panic(err.Error())
	}

	ExecuteCalls(code.([]Call), Nada)
}
