package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		InitializeFunctions()
		ExecuteFile(os.Args[1])

	} else {
		fmt.Printf("fluxo v0.1-alpha\nusage: fluxo [source_file]")
	}
}
