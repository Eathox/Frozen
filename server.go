package main

import (
	"fmt"
	"net"
)

func createServer() {
	listener, err := net.Listen(serverProtocol, serverAdress)
	if err != nil {
		errorMsg("Failed to listen to server: "+err.Error(), 1)
	}
	defer listener.Close()
	fmt.Println("Listening on ", serverAdress)
	getNewUser(listener)
}

func sendMessageAllUsers(curUser *user, message string) {
	for _, user := range allUsers {
		if (*user).status == active &&
			(*user).conn != (*curUser).conn &&
			(*user).channel == (*curUser).channel {
			sendMessage((*user).conn, message)
		}
	}
}
