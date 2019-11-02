package main

import (
	"fmt"
	"net"
	"os"
)

func createServer() {
	listener, err := net.Listen(serverProtocol, serverAdress)
	if err != nil {
		fmt.Println("Error Listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

}
