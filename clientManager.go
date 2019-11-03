package main

import (
	"fmt"
	"strings"
)

type clientManager struct {
	serverUser *user
	curUser    *user
	allUsers   *[]*user
	channels   *[]string
}

func (clientManager *clientManager) showCurrentChannel() {
	curChannel := clientManager.curUser.channel
	clientManager.curUser.sendMessageLn("Current channel: " + curChannel)
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

func (clientManager *clientManager) registerNewUser(userData []string) {
	clientManager.curUser.username = userData[0]
	clientManager.curUser.password = userData[1]
	if len(userData) > 2 {
		clientManager.curUser.nickname = userData[2]
	} else {
		clientManager.curUser.nickname = clientManager.curUser.username
	}
	clientManager.addCurUser()
	fmt.Println("User:", clientManager.curUser.username, "registerd")
}

func (clientManager *clientManager) addCurUser() {
	*clientManager.allUsers = append(*clientManager.allUsers, clientManager.curUser)
}

func (clientManager *clientManager) getUser(username string) (*user, bool) {
	for _, user := range *clientManager.allUsers {
		if user.username == username {
			return user, true
		}
	}
	return clientManager.curUser, false
}

func (clientManager *clientManager) loginUser(existingUser *user) {
	existingUser.conn = clientManager.curUser.conn
	existingUser.connReader = clientManager.curUser.connReader
	clientManager.curUser = existingUser
	fmt.Println("User:", clientManager.curUser.username, "logged in")
}

func (clientManager *clientManager) logoutCurUser() {
	clientManager.curUser.status = inactive
	clientManager.curUser.connReader = nil
	clientManager.curUser.conn.Close()
	clientManager.curUser.conn = nil
	fmt.Println("User:", clientManager.curUser.username, "logged out")
}

func (clientManager *clientManager) welcomeUser(user *user) {
	clientManager.sendServerMessage(user, "Welcome,")
	clientManager.sendServerMessage(user, "!init:  Initialize account or login")
	clientManager.sendServerMessage(user, "!help:  See list of possible commands")
	user.sendMessageLn("")
}

func (clientManager *clientManager) listActiveUsers() {
	for _, user := range *clientManager.allUsers {
		if user.status == active {
			clientManager.curUser.sendMessageLn(user.nickname)
		}
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
