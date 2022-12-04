package user

type Cmd interface {
	Start() //start executor
	Close() //close executor
}
