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
		errorMsg("Failed to establish connection to server: "+err.Error(), 1)
	}

	curUser := newUser(conn)
	go handleReceive(curUser)
	handleSend(curUser)
}

func handleSend(curUser *user) {
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
			curUser.sendMessage(message)
		}
	}
}

func handleReceive(curUser *user) {
	for {
		message, err := curUser.receiveMessage()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Failed to receive from server:", err.Error())
			}
			os.Exit(1)
		}
		fmt.Print(message)
	}
}
