package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func recieveMessage(conn net.Conn) (string, error) {
	connReader := bufio.NewReader(conn)
	message, err := connReader.ReadBytes('\n')
	return string(message), err
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message: ", err.Error())
		os.Exit(1)
	}
}
