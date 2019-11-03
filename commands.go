package main

import (
	"fmt"
	"strings"
)

type commandData struct {
	function    func(*clientManager, []string)
	description string
}

var commands = map[string]commandData{
	"!init":     commandData{initUser, "Logs in user, params: USERNAME PASSWORD (NICKNAME)"},
	"!nickname": commandData{changeNickname, "Change nickname, params: NICKNAME"},
	"!join":     commandData{joinChannel, "Join channel or create new one, params: CHANNEL"},
	"!leave":    commandData{leaveChannel, "Move back to defualt channel"},
	"!users":    commandData{listUsers, "List all users"},
	"!channels": commandData{listChannels, "List all channels"},
	"!whisper":  commandData{privateMessage, "Send private message to person, params: USERNAME MESSAGE"},
	"!msg":      commandData{privateMessage, "Send private message to person, params: USERNAME MESSAGE"},
	"!where":    commandData{currentChannel, "Show current channel"},
}

func handleCommand(clientManager *clientManager, message string) bool {
	words := strings.Fields(message)
	userCommand := strings.ToLower(words[0])

	if len(words) == 0 ||
		(userCommand != "!init" && userCommand != "!help" &&
			clientManager.curUser.status == inactive) {
		return false
	}

	if userCommand == "!help" {
		printHelp(clientManager)
		return true
	}
	for command, commandData := range commands {
		if command == userCommand {
			commandData.function(clientManager, words)
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

func printHelp(clientManager *clientManager) {
	for command, commandData := range commands {
		clientManager.curUser.sendMessageLn(fmt.Sprintf("%s: %s", command, commandData.description))
	}
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

func changeNickname(clientManager *clientManager, words []string) {
	countWords := len(words)
	if countWords < 2 {
		clientManager.sendServerMessage(clientManager.curUser, "usage: !nickname NICKNAME")
		return
	}

	clientManager.curUser.nickname = words[1]
}

func joinChannel(clientManager *clientManager, words []string) {
	if len(words) > 1 {
		clientManager.switchChannel(words[1])
		return
	}

	clientManager.sendServerMessage(clientManager.curUser, "usage: !init CHANNEL_NAME/ID")
}

func leaveChannel(clientManager *clientManager, words []string) {
	clientManager.leaveChannel()
	_ = words
}

func currentChannel(clientManager *clientManager, words []string) {
	clientManager.showCurrentChannel()
	_ = words
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
	if clientManager.curUser.status == active {
		clientManager.sendServerMessage(clientManager.curUser, "You are already active")
		return
	}

	countWords := len(words)
	if countWords < 3 {
		clientManager.sendServerMessage(clientManager.curUser, "Usage: !init USERNAME PASSWORD (NICKNAME)")
		return
	}

	existingUser, userFound := clientManager.getUser(words[1])
	if userFound == false {
		clientManager.registerNewUser(words[1:])
	} else if existingUser.password == words[2] && existingUser.status == inactive {
		clientManager.loginUser(existingUser)
	} else {
		clientManager.sendServerMessage(clientManager.curUser, "Failed to init, Username already in use")
		return
	}
	if countWords > 3 {
		clientManager.curUser.nickname = words[3]
	}

	clientManager.curUser.status = active
	clientManager.sendServerMessage(clientManager.curUser, "Welcome, "+clientManager.curUser.nickname)
}
