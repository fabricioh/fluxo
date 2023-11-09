package main

import "fmt"

func FormatLiteral(literal Literal) string {
	switch literal.kind {
	case LIST:
		result := "[ "
		for _, elem := range literal.value.([]Literal) {
			result += FormatLiteral(elem) + " "
		}
		return result + "]"
	case PAIR:
		return fmt.Sprintf(
			"(%s -> %s)",
			FormatLiteral(literal.value.(Pair).left),
			FormatLiteral(literal.value.(Pair).right),
		)
	case FUNCTION:
		function := literal.value.(Function)
		result := "("

		if function.is_recursive {
			result += "recursive "
		}

		result += fmt.Sprintf("function (%s -> %s)",
			literal.value.(Function).constraints.flow,
			literal.value.(Function).constraints.parameter,
		)

		if function.is_bound {
			result += fmt.Sprintf(
				" bound to %s",
				FormatLiteral(function.bound_argument),
			)
		}

		return result + ")"
	default:
		return fmt.Sprintf("%v", literal.value)
	}
}

func CreateFunction(body []Call, is_recursive bool) Function {
	newFunction := Function{
		name:         "ANON",
		body:         body,
		is_bound:     false,
		is_recursive: is_recursive,
		constraints:  Constraint{ANY, ANY},
	}

	newFunction.implementation = func(flow, argument Literal) Literal {
		return ExecuteCalls(newFunction.body, flow)
	}

	return newFunction
}
