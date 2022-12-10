package connection

type Cmd interface {
	Start() error //start executor
	Close()       //close executor
	FormByte() []byte
	RemoteAddress() string
	Status() int
}
