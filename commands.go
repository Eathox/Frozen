package main

import "strings"

var commands = map[string]func(*user, []string){
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

func handleCommand(curUser *user, message string) bool {
	words := strings.Fields(message)
	for key, function := range commands {
		if key == strings.ToLower(words[0]) {
			function(curUser, words)
			return true
		}
	}
	return false
}

func initUser(curUser *user, words []string) {
	countWords := len(words)
	if countWords < 3 {
		sendMessage(curUser.conn, "usage: !init USERNAME PASSWORD (NICKNAME)\n")
		return
	}

	existingUser := getUser(words[1])
	switch {
	case existingUser == nil:
		curUser.username = words[1]
		curUser.password = words[2]
		if countWords == 3 {
			curUser.nickname = curUser.username
		}
		allUsers = append(allUsers, curUser)

	case existingUser.password == words[2] && existingUser.status == inactive:
		existingUser.conn = curUser.conn
		existingUser.connReader = curUser.connReader
		*curUser = *existingUser

	case existingUser.password == words[2] && existingUser.status == active:
		sendMessage(curUser.conn, "[SERVER] : Failed to init, User already active\n")
		return

	default:
		sendMessage(curUser.conn, "[SERVER] : Failed to init, Username already in use\n")
		return
	}

	curUser.status = active
	if countWords > 3 {
		curUser.nickname = words[3]
	}
	sendMessage(curUser.conn, "[SERVER] : Welcome, "+curUser.nickname+"\n")
}
