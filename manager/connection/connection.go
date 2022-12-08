package connection

//these interface is handle by socks or light protocol, they use it in there main routine.
type FrontConnection interface {
	Close() error
	ProtocolClose() error
	CommandHandle() error
}
