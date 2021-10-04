package client

type commandId int

const (
	CHANGEUSERNAME commandId = iota
	MSG
	QUIT
)

type Command struct {
	ACTION commandId
	User   *TcpUser
	Args   []string
}
