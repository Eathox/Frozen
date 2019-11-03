package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == serverPass {
		fmt.Println("Creating server...")
		createServer()
	} else {
		fmt.Println("Handling client...")
		handleClient()
	}
}
