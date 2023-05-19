package main

import (
	"fmt"
)

// hi there
// at points you will see both logf and a logger function
// logger will log to the Output window in the GUI
// logf will log to stdout
// this is so that logf will be sort of like a debug log
// and logger will be like a normal log for the user
// :+1:

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
