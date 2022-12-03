package socks

const (
	HandShakeRequest  int = 1
	HandShakeResponse     = 2
	AuthRequest           = 3
	AuthResponse          = 4
)

type SocksHandshake struct {
}

func fromByte(buf []byte, messageType int) interface{} {
	//
}
