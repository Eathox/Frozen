package main

import "os"

func errorMsg(message string, code int) {
	println("Error Frozen:", message)
	os.Exit(code)
}
