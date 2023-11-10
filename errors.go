package main

import (
	"fmt"
	"os"
)

func Panic(message string) {
	currentFile, _ := PATH_STACK.Peek(0)

	println("\n===========ERROR===========\n")
	if len(PATH_STACK.content) > 0 {
		fmt.Printf("while executing %s\n\n", currentFile)
	}
	fmt.Printf("info: %s\n", message)
	println()
	os.Exit(1)
}
