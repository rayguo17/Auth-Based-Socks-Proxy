package user

type Cmd interface {
	Start() error //start executor
	Close()       //close executor
	FormByte() []byte
}
