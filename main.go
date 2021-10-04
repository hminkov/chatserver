package main

import (
	"chat-server/server"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = ":8888"
)

func main() {
	s := server.NewServer()
	s.Listen(CONN_HOST + CONN_PORT)

	// start the server
	s.Start()

	// close the server after the program closes
	defer s.CloseServer()

}
