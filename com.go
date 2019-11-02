package main

import (
	"bufio"
	"net"
)

func recieveMessage(connReader *bufio.Reader) (string, error) {
	message, err := connReader.ReadBytes('\n')
	return string(message), err
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		errorMsg("Failed to send message: "+err.Error(), 1)
	}
}
