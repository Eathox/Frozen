package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleClient() {
	conn, err := net.Dial(serverProtocol, serverAdress)
	if err != nil {
		errorMsg("Failed to establish connection to server: " + err.Error(), 1)
	}
	go handleRecieve(conn)
	handleSend(conn)
}

func handleSend(conn net.Conn) {
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		message, err := stdinReader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Failed to send to server:", err.Error())
			}
			os.Exit(1)
		}
		if len(strings.TrimSpace(message)) != 0 {
			sendMessage(conn, message)
		}
	}
}

func handleRecieve(conn net.Conn) {
	connReader := bufio.NewReader(conn)
	for {
		message, err := recieveMessage(connReader)
		if err != nil {
			errorMsg("Failed to recieve from server: " + err.Error(), 1)
		}
		fmt.Print(message)
	}
}
