package util

type Address interface {
	String() string
	Addr() string
	Port() [2]byte
}
