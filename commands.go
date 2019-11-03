package main

import (
	"fmt"
	"strings"
)

var commands = map[string]func(*clientManager, []string){
	"!init": initUser,
	// "!nickname": changeNickname,
	// "!join":     joinChannel,
	// "!leave":    leaveChannel,
	// "!users":    listUsers,
	// "!channels": listChannels,
	// "!whisper":  privateMessage,
	// "!msg":      privateMessage,
	// "!help":     printHelp,
}

func handleCommand(clientManager *clientManager, message string) bool {
	words := strings.Fields(message)
	for key, function := range commands {
		if key == strings.ToLower(words[0]) {
			function(clientManager, words)
			return true
		}
	}
	return false
}

func initUser(clientManager *clientManager, words []string) {
	countWords := len(words)
	if countWords < 3 {
		sendMessageLn(clientManager.curUser.conn, "usage: !init USERNAME PASSWORD (NICKNAME)")
		return
	}

	existingUser, userFound := getUser(clientManager, words[1])
	switch {
	case userFound == false:
		clientManager.curUser.username = words[1]
		clientManager.curUser.password = words[2]
		if countWords == 3 {
			clientManager.curUser.nickname = clientManager.curUser.username
		}
		fmt.Printf("CurUser: %p\n", clientManager.curUser)
		*clientManager.allUsers = append(*clientManager.allUsers, clientManager.curUser)

	case existingUser.password == words[2] && existingUser.status == inactive:
		existingUser.conn = clientManager.curUser.conn
		existingUser.connReader = clientManager.curUser.connReader
		clientManager.curUser = existingUser

	default:
		sendMessageLn(clientManager.curUser.conn, "[SERVER] : Failed to init, Username already in use")
		return
	}

	clientManager.curUser.status = active
	if countWords > 3 {
		clientManager.curUser.nickname = words[3]
	}
	sendMessageLn(clientManager.curUser.conn, "[SERVER] : Welcome, "+clientManager.curUser.nickname)
}
