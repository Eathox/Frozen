package main

import (
	"fmt"
	"net"
	"os"
)

func handleClient() {
	conn, err := net.Dial(serverProtocol, serverAdress)
	if err != nil {
		fmt.Println("Error establishing connection:", err.Error())
		os.Exit(1)
	}
	_ = conn
}
