package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		InitializeFunctions()

		err := ExecuteFile(os.Args[1])

		if err != nil {
			Panic(err.Error())
		}

	} else {
		fmt.Printf("fluxo v0.1-alpha\nusage: fluxo [source_file]")
	}
}
