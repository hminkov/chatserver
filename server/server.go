package server

import (
	"bufio"
	client "chat-server/client"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CONN_TYPE          = "tcp"
	USERNAMEVALIDATOR  = "[a-zA-Z0-9]{3,}"
	FIRSTNAMEVALIDATOR = "[a-zA-Z]{3,}"
	LASTNAMEVALIDATOR  = "[a-zA-Z]{3,}"
	EMAILVALIDATOR     = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+[a-zA-Z0-9-.]+$"
	PASSWORDVALIDATOR  = "[a-zA-Z]{8,}"
)

type ChatServer struct {
	listener net.Listener
	Commands chan client.Command
	users    []*client.TcpUser //---> this way we save the data for all users and unables us to send them messages
	mutex    *sync.Mutex
}

func NewServer() *ChatServer {
	return &ChatServer{
		Commands: make(chan client.Command),
		mutex:    &sync.Mutex{},
	}
}

func (s *ChatServer) Listen(address string) {

	//with net.Listen I start the server
	l, err := net.Listen(CONN_TYPE, address)

	if err != nil {
		log.Printf("Unable to start server: %s:", err.Error())
		log.Println()
		os.Exit(1)
	}

	s.listener = l
	log.Printf("Started server on %v", address)
}

func (s *ChatServer) CloseServer() {
	s.listener.Close()
}

func (s *ChatServer) run() {
	for cmd := range s.Commands {
		switch cmd.ACTION {
		case client.CHANGEUSERNAME:
			s.changeUsername(cmd.User, cmd.Args)
		case client.MSG:
			s.msg(cmd.User, cmd.Args)
		case client.QUIT:
			s.quit(cmd.User)
		}
	}
}

func (s *ChatServer) Start() {
	go s.run()
	for {
		// I need a way to break the loop
		conn, err := s.listener.Accept()

		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		} else {
			// handle connection
			user := s.accept(conn)
			go s.newUser(user) // this is a blocking operation so we need Goroutine. This is also how we handle the connection.
			continue
		}
	}
}

func (s *ChatServer) accept(conn net.Conn) *client.TcpUser {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Println("A new user has connected", conn.RemoteAddr())
	u := &client.TcpUser{
		Commands: s.Commands,
		Conn:     conn,
	}
	s.users = append(s.users, u)
	return u
}

func (s *ChatServer) newUser(u *client.TcpUser) {

	userDataInputs := []string{"username", "firstname", "lastname", "email"}

	for _, dataType := range userDataInputs {
		u.MsgCurrentUser("Enter your " + dataType + ": ")
		scanner := bufio.NewScanner(u.Conn)
		counter := 1
		for {
			scanner.Scan()
			data := scanner.Text()
			if !userDataValidator(u, data, dataType) {
				triesLeft := strconv.Itoa(3 - counter)
				u.MsgCurrentUser("You have " + triesLeft + " more tries left")
				counter = counter + 1
			} else {
				break
			}
			if counter == 3 {
				u.MsgCurrentUser("You have no more tries left. Please come back later!\n Bye!")
				u.Conn.Close()
				break
			}
		}
	}

	u.MsgCurrentUser("User Validated! Welcome to our chat world!")
	u.CommandParser()
}

func userDataValidator(u *client.TcpUser, input string, inputType string) bool {
	usernamevalidation, _ := regexp.Compile(USERNAMEVALIDATOR)
	firstNameValidation, _ := regexp.Compile(FIRSTNAMEVALIDATOR)
	lastNameValidation, _ := regexp.Compile(LASTNAMEVALIDATOR)
	emailValidation, _ := regexp.Compile(EMAILVALIDATOR)
	// PasswordValidation, _ := regexp.Compile(PASSWORDVALIDATOR) //-->> to be implemented later on
	switch inputType {
	case "username":
		if !usernamevalidation.MatchString(input) {
			u.MsgCurrentUser("The username is not in a appropiate format.")
			return false
		} else {
			u.Username = input
			return true
		}
	case "firstname":
		if !firstNameValidation.MatchString(input) {
			u.MsgCurrentUser("First Name is not valid")
			return false
		} else {
			return true
		}
	case "lastname":
		if !lastNameValidation.MatchString(input) {
			u.MsgCurrentUser("Last Name is not valid")
			return false
		} else {
			return true
		}
	default:
		if !emailValidation.MatchString(input) {
			u.MsgCurrentUser("Email is not valid")
			return false
		} else {
			return true
		}
	}
	//Later on we can use the same for password
	// if PasswordValidation.MatchString(input) {
	// 	return "Password is not valid"
	// }
}

func (s *ChatServer) changeUsername(u *client.TcpUser, args []string) {
	oldUsername := u.Username
	u.Username = args[1]
	u.MsgCurrentUser("Username is changed.")
	u.MsgCurrentUser(fmt.Sprintf("Old username: %s", oldUsername))
	u.MsgCurrentUser(fmt.Sprintf("New username: %s", u.Username))
}

func (s *ChatServer) msg(u *client.TcpUser, args []string) {
	if len(args) < 2 {
		u.MsgCurrentUser("message is required, usage: /msg MSG")
		return
	}

	msg := u.Username + ": " + strings.Join(args[1:], " ")
	s.msgAllUsers(msg) //printout for the users
	log.Println(msg)   //printout to the console
}
func (s *ChatServer) msgAllUsers(msg string) {
	for _, user := range s.users {
		user.MsgCurrentUser(msg) //printout for the users
	}
}

func (s *ChatServer) quit(u *client.TcpUser) {
	defer u.Conn.Close()
	time := (string)(time.Now().Local().String())
	msg := time + ": " + u.Username + " left the chat"
	s.msgAllUsers(msg)
}
