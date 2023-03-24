package main

import (
	"fmt"
)

func main() {
	createGUI()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logf(format string, a ...any) {
	// we use println instead of fmt.Print because
	// android doesn't pipe stdout/stderr to logcat
	println(fmt.Sprintf(format, a...))
}
