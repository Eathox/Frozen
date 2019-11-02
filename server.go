package main

import (
	"fmt"
	"net"
	"os"
)

type user struct {
	conn     net.Conn
	nickname string
	username string
	password string
	channel  uint8
}

func createServer() {
	listener, err := net.Listen(serverProtocol, serverAdress)
	if err != nil {
		fmt.Println("Error Listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on ", serverAdress)
	getNewUser(listener)
}

func getNewUser(listener net.Listener) {
	newUserChannel := make(chan user)
	removedUserChannel := make(chan user)
	allUsers := []user{}

	go handleNewUser(listener, newUserChannel)

	for {
		select {
		case newUser := <-newUserChannel:
			allUsers = append(allUsers, newUser)
			go handleUserRequest(newUser, &allUsers, removedUserChannel)
		case removedUser := <-removedUserChannel:
			handleRemoveUser(removedUser, allUsers)
		}
	}
}

func handleNewUser(listener net.Listener, newUserChannel chan user) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		println("Connection Established")
		newUser := user{conn: conn}
		newUserChannel <- newUser
	}
}

func handleRemoveUser(removedUser user, allUsers []user) {
	for index, user := range allUsers {
		if user.conn == removedUser.conn {
			allUsers = append(allUsers[:index], allUsers[index+1:]...)
		}
	}
	println("Connection Terminated")
}

func handleUserRequest(curUser user, allUsers *[]user, removedUserChannel chan user) {
	for {
		message, err := recieveMessage(curUser.conn)
		if err != nil {
			switch {
			case err.Error() != "EOF":
				fmt.Println("Error recieving message:", err.Error())
			default:
				removedUserChannel <- curUser
				return
			}
		}
		for _, user := range *allUsers {
			if user.conn != curUser.conn {
				sendMessage(user.conn, message)
			}
		}
		println(message)
	}
}
