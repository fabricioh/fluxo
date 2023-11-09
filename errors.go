package main

import (
	"fmt"
	"os"
)

func Panic(message string) {
	fmt.Printf("===========ERROR===========\ninfo: %s", message)
	os.Exit(1)
}
