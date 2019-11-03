package main

import (
	"bufio"
	"net"
)

type userStatus int

const (
	inactive userStatus = 0
	active   userStatus = 1
)

type user struct {
	conn       net.Conn
	connReader *bufio.Reader
	nickname   string
	username   string
	password   string
	status     userStatus
	channel    string
}

func newUser(conn net.Conn) *user {
	newUser := new(user)
	newUser.conn = conn
	newUser.channel = "World"
	if conn != nil {
		newUser.connReader = bufio.NewReader(conn)
	}
	return newUser
}
