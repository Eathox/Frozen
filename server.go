package main

import (
	"fmt"
	"net"
)

const (
	serverIP       = "10.11.3.2" //Fill in your ip here
	serverPort     = "4242"
	serverAdress   = serverIP + ":" + serverPort
	serverProtocol = "tcp"
	serverPass     = "host"
)

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
	clientManager.serverUser.username = "SERVER"
	clientManager.serverUser.nickname = clientManager.serverUser.username
	clientManager.serverUser.password = serverPass
	clientManager.allUsers = &[]*user{clientManager.serverUser}
}

func createServer() {
	serverListener := createListener()
	defer serverListener.Close()
	clientManager := clientManager{}
	createServerUser(&clientManager)
	clientManager.channels = &[]string{
		"World",
		"Random",
		"Games",
		"Books",
		"Sports",
		"Cry_Corner",
		"Laugh_Corner",
		"Venting_Corner",
		"Super_Cereal",
		"Secret_Gossip",
	}
	serverLoop(clientManager, serverListener)
}

func serverLoop(clientManager clientManager, listener net.Listener) {
	newUserChannel := make(chan *user)

	go handleNewUser(listener, newUserChannel)

	for {
		newUser := <-newUserChannel
		clientManager.welcomeUser(newUser)
		go clientManager.handleUserRequest(newUser)
	}
}

func handleNewUser(listener net.Listener, newUserChannel chan *user) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			errorMsg("Failed to accept connection: "+err.Error(), 1)
		}

		newUser := newUser(conn)
		fmt.Println("Connection Established")
		newUserChannel <- newUser
	}
}
