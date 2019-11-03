package main

import (
	"fmt"
	"net"
	"strings"
)

type clientManager struct {
	serverUser *user
	curUser    *user
	allUsers   *[]*user
	channels   *[]string
}

func (clientManager *clientManager) addCurUser() {
	*clientManager.allUsers = append(*clientManager.allUsers, clientManager.curUser)
	fmt.Println("User:", clientManager.curUser.username, "registerd")
}

func (clientManager *clientManager) getUser(username string) (*user, bool) {
	for _, user := range *clientManager.allUsers {
		if user.username == username {
			return user, true
		}
	}
	return clientManager.curUser, false
}

func (clientManager *clientManager) showCurrentChannel() {
	curChannel := clientManager.curUser.channel
	clientManager.curUser.sendMessageLn("Current channel: "+curChannel)
}

func (clientManager *clientManager) leaveChannel() {
	clientManager.curUser.channel = "World"
	clientManager.showCurrentChannel()
}

func (clientManager *clientManager) switchChannel(newChannel string) {
	found := false
	newChannel = strings.ToLower(newChannel)
	for _, channel := range *clientManager.channels {
		if newChannel == strings.ToLower(channel) {
			found = true
			break
		}
	}
	if found == false {
		*clientManager.channels = append(*clientManager.channels, newChannel)
	}
	clientManager.curUser.channel = newChannel
	clientManager.showCurrentChannel()
}

func (clientManager *clientManager) listAllChannels() {
	for _, channel := range *clientManager.channels {
		clientManager.curUser.sendMessageLn(channel)
	}
}

func (clientManager *clientManager) listActiveUsers() {
	for _, user := range *clientManager.allUsers {
		if user.status == active {
			clientManager.curUser.sendMessageLn(user.nickname)
		}
	}
}

func (clientManager *clientManager) logoutCurUser() {
	clientManager.curUser.status = inactive
	clientManager.curUser.connReader = nil
	clientManager.curUser.conn.Close()
	clientManager.curUser.conn = nil
	fmt.Println("User:", clientManager.curUser.username, "logged out")
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

func welcomeNewUser(clientManager clientManager, newUser *user) {
	clientManager.sendServerMessage(newUser, "Welcome,")
	clientManager.sendServerMessage(newUser, "!init:  Initialize account or login")
	clientManager.sendServerMessage(newUser, "!help:  See list of possible commands")
	newUser.sendMessageLn("")
}

func serverLoop(clientManager clientManager, listener net.Listener) {
	newUserChannel := make(chan *user)

	go handleNewUser(listener, newUserChannel)

	for {
		newUser := <-newUserChannel
		welcomeNewUser(clientManager, newUser)
		go clientManager.handleUserRequest(newUser)
	}
}
