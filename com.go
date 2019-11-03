package main

import "fmt"

func prefixMessage(sender user, message string) string {
	message = fmt.Sprintf("[%s: %s]: %s", sender.channel, sender.nickname, message)
	return message
}

func prefixWhisper(sender user, message string) string {
	message = fmt.Sprintf("-> %s(%s) whispered: %s", sender.username, sender.nickname, message)
	return message
}

func (user user) receiveMessage() (string, error) {
	message, err := user.connReader.ReadBytes('\n')
	return string(message), err
}

func (user user) sendMessagePrefix(message string) {
	message = prefixMessage(user, message)
	user.sendMessage(message)
}

func (user user) sendMessage(message string) {
	_, err := user.conn.Write([]byte(message))
	if err != nil {
		errorMsg("Failed to send message: "+err.Error(), 1)
	}
}

func (user user) sendMessagePrefixLn(message string) {
	user.sendMessagePrefix(message + "\n")
}

func (user user) sendMessageLn(message string) {
	user.sendMessage(message + "\n")
}

//---- Direct messages
func (clientManager *clientManager) sendServerMessage(receiver *user, message string) {
	message = prefixMessage(*clientManager.serverUser, message)
	clientManager.sendDirectMessageLn(receiver, message)
}

func (clientManager *clientManager) sendDirectMessagePrefix(receiver *user, message string) {
	message = prefixMessage(*clientManager.curUser, message)
	clientManager.sendDirectMessage(receiver, message)
}

func (clientManager *clientManager) sendDirectMessagePrefixWhisper(receiver *user, message string) {
	message = prefixWhisper(*clientManager.curUser, message)
	clientManager.sendDirectMessage(receiver, message)
}

func (clientManager *clientManager) sendDirectMessage(receiver *user, message string) {
	_, err := receiver.conn.Write([]byte(message))
	if err != nil {
		errorMsg("Failed to send message: "+err.Error(), 1)
	}
}

func (clientManager *clientManager) sendDirectMessagePrefixLn(receiver *user, message string) {
	clientManager.sendDirectMessagePrefix(receiver, message+"\n")
}

func (clientManager *clientManager) sendDirectMessagePrefixWhisperLn(receiver *user, message string) {
	clientManager.sendDirectMessagePrefixWhisper(receiver, message+"\n")
}
func (clientManager *clientManager) sendDirectMessageLn(receiver *user, message string) {
	clientManager.sendDirectMessage(receiver, message+"\n")
}

func (clientManager *clientManager) sendMessageToAllUsers(message string) {
	for _, user := range *clientManager.allUsers {
		if user.status == active && user.conn != nil &&
			user.conn != clientManager.curUser.conn &&
			user.channel == clientManager.curUser.channel {
			clientManager.sendDirectMessagePrefix(user, message)
		}
	}
}
