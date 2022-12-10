package util

const (
	Ipv4Address   = "IPV4"
	DomainAddress = "Domain"
	Ipv6Address   = "IPV6"
)

type Address interface {
	String() string
	Addr() string
	Port() [2]byte
	AddrType() string
	AddrByte() []byte
}
