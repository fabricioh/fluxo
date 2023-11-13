package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Panic(message string) {
	currentFile, _ := PATH_STACK.Peek(0)

	fmt.Printf("\n===========%s===========\n", color.MagentaString("ERROR"))

	if len(PATH_STACK.content) > 0 {
		fmt.Printf("\nwhile executing file: %s\n\n", color.CyanString(currentFile))
	}

	fmt.Printf("trace:\n%s\n", message)
	println()

	os.Exit(1)
}
