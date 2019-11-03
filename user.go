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
	channel    uint8
}

func getUser(clientManager *clientManager, username string) (*user, bool) {
	for _, user := range *clientManager.allUsers {
		if user.username == username {
			return user, true
		}
	}
	return clientManager.curUser, false
}

func newUser(conn net.Conn) *user {
	newUser := new(user)
	newUser.conn = conn
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
		sendMessageLn(conn, "[SERVER] : Welcome,")
		sendMessageLn(conn, "[SERVER] : !init:  Initialize account or login")
		sendMessageLn(conn, "[SERVER] : !help:  See list of possible commands")
		sendMessageLn(conn, "")
		fmt.Println("Connection Established")
		newUserChannel <- newUser
	}
}

func handleRemoveUser(clientManager clientManager, removedUser *user) {
	for _, user := range *clientManager.allUsers {
		if user.conn == removedUser.conn {
			user.status = inactive
			user.connReader = nil
			user.conn.Close()
			fmt.Println("Connection Terminated")
		}
	}
}

func handleUserRequest(clientManager clientManager, removedUserChannel chan *user) {
	for {
		message, err := receiveMessage(clientManager.curUser.connReader)
		if err != nil {
			switch {
			case err.Error() != "EOF":
				fmt.Println("Error receiving message:", err.Error())

			default:
				removedUserChannel <- clientManager.curUser
				return
			}
		}
		handled := handleCommand(&clientManager, message)
		if handled != true {
			switch clientManager.curUser.status {
			case inactive:
				sendMessageLn(clientManager.curUser.conn, "Please initialize")
				sendMessageLn(clientManager.curUser.conn, "!help\tif you feel lost")

			case active:
				sendMessageAllUsers(&clientManager, message)
			}
		}
		fmt.Print(message)
	}
}
