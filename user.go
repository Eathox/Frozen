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

var allUsers = []*user{}

func getUser(username string) *user {
	for _, user := range allUsers {
		if user.username == username {
			return user
		}
	}
	return nil
}

func getNewUser(listener net.Listener) {
	newUserChannel := make(chan *user)
	removedUserChannel := make(chan *user)

	go handleNewUser(listener, newUserChannel)

	for {
		select {
		case newUser := <-newUserChannel:
			go handleUserRequest(newUser, removedUserChannel)
		case removedUser := <-removedUserChannel:
			handleRemoveUser(removedUser)
		}
	}
}

func newUser(conn net.Conn) user {
	newUser := user{
		conn:       conn,
		connReader: bufio.NewReader(conn),
	}
	return (newUser)
}

func handleNewUser(listener net.Listener, newUserChannel chan *user) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			errorMsg("Failed to accept connection: "+err.Error(), 1)
		}
		newUser := newUser(conn)
		sendMessage(conn, "[SERVER] : Welcome,\n")
		sendMessage(conn, "[SERVER] : !init:  Initialize account or login\n")
		sendMessage(conn, "[SERVER] : !help:  See list of possible commands\n")
		sendMessage(conn, "\n")
		fmt.Println("Connection Established")
		newUserChannel <- &newUser
	}
}

func handleRemoveUser(removedUser *user) {
	for _, user := range allUsers {
		if (*user).conn == (*removedUser).conn {
			user.status = inactive
			user.connReader = nil
			user.conn.Close()
			fmt.Println("Connection Terminated")
		}
	}
}

func handleUserRequest(curUser *user, removedUserChannel chan *user) {
	for {
		message, err := receiveMessage((*curUser).connReader)
		if err != nil {
			switch {
			case err.Error() != "EOF":
				fmt.Println("Error receiving message:", err.Error())

			default:
				removedUserChannel <- curUser
				return
			}
		}
		handled := handleCommand(curUser, message)
		if handled != true {
			switch (*curUser).status {
			case inactive:
				sendMessage((*curUser).conn, "Please initialize\n")
				sendMessage((*curUser).conn, "!help\tif you feel lost\n")

			case active:
				sendMessageAllUsers(curUser, message)
			}
		}
		fmt.Print(message)
	}
}
