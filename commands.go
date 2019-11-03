package main

import (
	"fmt"
	"strings"
)

type command struct {
	function    func(*clientManager, []string)
	description string
}

var commands = map[string]command{
	"!init":     command{initUser, "Logs in user, params: USERNAME PASSWORD (NICKNAME)"},
	"!nickname": command{changeNickname, "Change nickname, params: NICKNAME"},
	"!join":     command{joinChannel, "Join channel or create new one, params: CHANNEL"},
	"!leave":    command{leaveChannel, "Move back to defualt channel"},
	"!users":    command{listUsers, "List all users"},
	"!channels": command{listChannels, "List all channels"},
	"!whisper":  command{privateMessage, "Send private message to person, params: USERNAME MESSAGE"},
	"!msg":      command{privateMessage, "Send private message to person, params: USERNAME MESSAGE"},
	"!where":    command{currentChannel, "Show current channel"},
}

func handleCommand(clientManager *clientManager, message string) bool {
	words := strings.Fields(message)

	if strings.ToLower(words[0]) != "!init" && strings.ToLower(words[0]) != "!help" &&
		clientManager.curUser.status == inactive {
		return false
	}
	if len(words) == 0 {
		return false
	}

	if strings.ToLower(words[0]) == "!help" {
		printHelp(clientManager)
		return true
	}

	for key, command := range commands {
		if key == strings.ToLower(words[0]) {
			command.function(clientManager, words)
			return true
		}
	}

	if words[0][0] == '!' {
		clientManager.sendServerMessage(clientManager.curUser, "Invalid command")
		clientManager.sendServerMessage(clientManager.curUser, "!help\tif you feel lost")
		return true
	}
	return false
}

func currentChannel(clientManager *clientManager, words []string) {
	clientManager.showCurrentChannel()
	_ = words
}

func privateMessage(clientManager *clientManager, words []string) {
	if len(words) < 3 {
		clientManager.sendServerMessage(clientManager.curUser, "usage: !whisper/msg USERNAME MESSAGE")
		return
	}
	var message string
	for _, word := range words[2:] {
		message = message + word + " "
	}
	reciever, found := clientManager.getUser(words[1])
	if found == false || reciever.status == inactive {
		clientManager.sendServerMessage(clientManager.curUser, "Failed to send message, User does not exist")
		return
	}
	clientManager.sendDirectMessagePrefixWhisperLn(reciever, message)
}

func printHelp(clientManager *clientManager) {
	for key, command := range commands {
		clientManager.curUser.sendMessageLn(fmt.Sprintf("%s: %s", key, command.description))
	}
}

func changeNickname(clientManager *clientManager, words []string) {
	countWords := len(words)
	if countWords < 2 {
		clientManager.sendServerMessage(clientManager.curUser, "usage: !nickname NICKNAME")
		return
	}

	clientManager.curUser.nickname = words[1]
}

func leaveChannel(clientManager *clientManager, words []string) {
	clientManager.leaveChannel()
	_ = words
}

func joinChannel(clientManager *clientManager, words []string) {
	if len(words) > 1 {
		clientManager.switchChannel(words[1])
	} else {
		clientManager.sendServerMessage(clientManager.curUser, "usage: !init CHANNEL_NAME/ID")
	}
}

func listChannels(clientManager *clientManager, words []string) {
	clientManager.listAllChannels()
	_ = words
}

func listUsers(clientManager *clientManager, words []string) {
	clientManager.listActiveUsers()
	_ = words
}

func initUser(clientManager *clientManager, words []string) {
	countWords := len(words)
	if clientManager.curUser.status == active {
		clientManager.sendServerMessage(clientManager.curUser, "You are allreadt active")
		return
	}

	if countWords < 3 {
		clientManager.sendServerMessage(clientManager.curUser, "Usage: !init USERNAME PASSWORD (NICKNAME)")
		return
	}

	existingUser, userFound := clientManager.getUser(words[1])
	switch {
	case userFound == false:
		clientManager.curUser.username = words[1]
		clientManager.curUser.password = words[2]
		if countWords == 3 {
			clientManager.curUser.nickname = clientManager.curUser.username
		}
		clientManager.addCurUser()

	case existingUser.password == words[2] && existingUser.status == inactive:
		existingUser.conn = clientManager.curUser.conn
		existingUser.connReader = clientManager.curUser.connReader
		clientManager.curUser = existingUser
		fmt.Println("User:", clientManager.curUser.username, "logged in")

	default:
		clientManager.sendServerMessage(clientManager.curUser, "Failed to init, Username already in use")
		return
	}

	clientManager.curUser.status = active
	if countWords > 3 {
		clientManager.curUser.nickname = words[3]
	}
	clientManager.sendServerMessage(clientManager.curUser, "Welcome, "+clientManager.curUser.nickname)
}
