package main

import "fmt"

func FormatLiteral(literal Literal) string {
	switch literal.kind.name {
	case LIST:
		result := "[ "
		for _, elem := range literal.value.([]Literal) {
			result += FormatLiteral(elem) + " "
		}
		return result + "]"

	case PAIR:
		return fmt.Sprintf(
			"(%s & %s)",
			FormatLiteral(literal.value.(Pair).left),
			FormatLiteral(literal.value.(Pair).right),
		)

	case FUNCTION:
		function := literal.value.(Function)
		result := "{"

		if function.is_recursive {
			result += "recursive "
		}

		result += fmt.Sprintf("function: (@%s & @%s)",
			literal.value.(Function).constraints.flow.Format(),
			literal.value.(Function).constraints.parameter.Format(),
		)

		if function.is_bound {
			result += fmt.Sprintf(
				" bound to %s",
				FormatLiteral(function.bound_argument),
			)
		}

		return result + "}"

	case TYPE:
		return literal.value.(Kind).Format()

	default:
		return fmt.Sprintf("%v", literal.value)
	}
}

func FormatConstraint(constraints Constraints) string {
	return fmt.Sprintf("(%s & %s)", constraints.flow.Format(), constraints.parameter.Format())
}

func CreateFunction(body []Call, is_recursive bool) Function {
	newFunction := Function{
		name:         "ANON",
		body:         body,
		is_bound:     false,
		is_recursive: is_recursive,
		constraints:  Constraints{Kind{name: ANY}, Kind{name: ANY}},
	}

	newFunction.implementation = func(flow, argument Literal) (Literal, error) {
		return ExecuteCalls(newFunction.body, flow)
	}

	return newFunction
}

func DetermineListSpec(list []Literal) Kind {
	spec := Kind{name: ANY}

	for i, elem := range list {
		if i == 0 {
			spec = elem.kind
			continue
		}

		// fmt.Printf("%s == %s = %v\n", spec.Format(), elem.kind.Format(), spec.Matches(&elem.kind))

		if !spec.Matches(&elem.kind) {
			spec = Kind{name: ANY}
			break
		}
	}

	return spec
}
