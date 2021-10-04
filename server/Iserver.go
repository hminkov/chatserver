package server

//the interface actually helps us to define the behavior clearly.
type IChatServer interface {
	Listen(address string) error
	Start()
	Close()
	CloseServer()
}
