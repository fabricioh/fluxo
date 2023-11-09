package main

import "os"

func main() {

	file, err := os.ReadFile("test.fl")
	if err != nil {
		panic(err)
	}

	result, err := Parse("main", file)
	if err != nil {
		panic(err)
	}

	InitializeFunctions()
	ExecuteCalls(result.([]Call), Literal{"Hello, world!", TEXT})
}

/*
- Tagged unions (Enums helpful);
- Function pointers;
- Global mutable array;
- Lambdas/Closures;
*/
