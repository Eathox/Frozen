package main

import (
	"bufio"
	"fmt"
	"net"
)

type userStatus int

const (
	inactive userStatus = 0
	active   userStatus = 1
)

type user struct {
	conn       net.Conn
	connReader *bufio.Reader
	nickname   string
	username   string
	password   string
	status     userStatus
	channel    string
}

func newUser(conn net.Conn) *user {
	newUser := new(user)
	newUser.conn = conn
	newUser.channel = "World"
	if conn != nil {
		newUser.connReader = bufio.NewReader(conn)
	}
	return newUser
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

func (clientManager clientManager) handleUserRequest(curUser *user) {
	clientManager.curUser = curUser

	for {
		message, err := clientManager.curUser.receiveMessage()
		if err != nil {
			switch {
			case err.Error() != "EOF":
				errorMsg("Error receiving message: "+err.Error(), 1)

			default:
				clientManager.logoutCurUser()
				return
			}
		}
		handled := handleCommand(&clientManager, message)
		if handled != true {
			switch clientManager.curUser.status {
			case inactive:
				clientManager.sendServerMessage(clientManager.curUser, "Please initialize")
				clientManager.sendServerMessage(clientManager.curUser, "!help\tif you feel lost")

			case active:
				clientManager.sendMessageToAllUsers(message)
			}
		}
		fmt.Print(message)
	}
}
