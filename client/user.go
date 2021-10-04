package client

import (
	"bufio"
	"net"
	"strings"
)

type TcpUser struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Commands  chan<- Command
	Conn      net.Conn
}

func (u *TcpUser) CommandParser() {

	for {
		msg, err := bufio.NewReader(u.Conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/changename":
			u.Commands <- Command{
				ACTION: CHANGEUSERNAME,
				User:   u,
				Args:   args,
			}
		case "/msg":
			u.Commands <- Command{
				ACTION: MSG,
				User:   u,
				Args:   args,
			}
		case "/quit":
			u.Commands <- Command{
				ACTION: QUIT,
				User:   u,
			}
		default:
			u.MsgCurrentUser("Unknown command")
		}
	}
}

func (u *TcpUser) MsgCurrentUser(msg string) {
	u.Conn.Write([]byte("> " + msg + "\n"))
}
