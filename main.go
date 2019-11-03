package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	serverIP       = "10.11.3.1" //Fill in your ip here
	serverPort     = "4242"
	serverAdress   = serverIP + ":" + serverPort
	serverProtocol = "tcp"
	serverPass     = "host"
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
