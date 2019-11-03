package main

import (
	"fmt"
	"net"
)

type clientManager struct {
	serverUser *user
	curUser    *user
	allUsers   *[]*user
}

func createListener() net.Listener {
	listener, err := net.Listen(serverProtocol, serverAdress)
	if err != nil {
		errorMsg("Failed to listen to server: "+err.Error(), 1)
	}
	fmt.Println("Server listening on: ", serverAdress)
	return listener
}

func createServerUser(clientManager *clientManager) {
	serverUser := newUser(nil)
	clientManager.serverUser = serverUser
	clientManager.serverUser.username = "server"
	clientManager.serverUser.nickname = clientManager.serverUser.username
	clientManager.serverUser.password = serverPass
	clientManager.allUsers = &[]*user{clientManager.serverUser}
}

func createServer() {
	serverListener := createListener()
	defer serverListener.Close()
	clientManager := clientManager{}
	createServerUser(&clientManager)
	serverLoop(clientManager, serverListener)
}

func serverLoop(clientManager clientManager, listener net.Listener) {
	newUserChannel := make(chan *user)
	removedUserChannel := make(chan *user)

	go handleNewUser(listener, newUserChannel)

	for {
		select {
		case newUser := <-newUserChannel:
			clientManager.curUser = newUser
			go handleUserRequest(clientManager, removedUserChannel)
		case removedUser := <-removedUserChannel:
			handleRemoveUser(clientManager, removedUser)
		}
	}
}

func sendMessageAllUsers(clientManager *clientManager, message string) {
	fmt.Printf("CurUser %p\nAllusers %v\n", clientManager.curUser, clientManager.allUsers)

	for _, user := range *clientManager.allUsers {
		if user.status == active &&
			user.conn != clientManager.curUser.conn &&
			user.channel == clientManager.curUser.channel {
			sendMessage(user.conn, message)
		}
	}
}
