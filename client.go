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
		fmt.Println("Error establishing connection:", err.Error())
		os.Exit(1)
	}
	go handleRecieve(conn)
	handleSend(conn)
}

func handleSend(conn net.Conn) {
	for {
		stdinReader := bufio.NewReader(os.Stdin)
		message, err := stdinReader.ReadString('\n')
		if err != nil {
			if (err.Error() != "EOF") {
				fmt.Println("Error reading message:", err.Error())
			}
			os.Exit(1)
		}
		if len(strings.TrimSpace(message)) != 0 {
			sendMessage(conn, message)
		}
	}
}

func handleRecieve(conn net.Conn) {
	for {
		message, err := recieveMessage(conn)
		if (err != nil) {
			fmt.Println("Error recieving message:", err.Error())
			os.Exit(1)
		}
		fmt.Print(message)
	}
}
