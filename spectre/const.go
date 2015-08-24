package spectre

import (
	"math/rand"
	"net"
	"time"
)

type ProtocolType int

const (
	PROTOCOL_TCP = ProtocolType(iota)
	PROTOCOL_UDP
	PROTOCOL_HTTP

	namespace = "blinker/spectre"
)

var chars = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
	'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
	'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
}

//var verify_key = aes.NewCipher([]byte("58ad451ea0ed5f82f13efa4d4cb49a8c"))

func NewGUID(s int) string {
	rs := []rune{}
	leng := len(chars)
	for i := 0; i < s; i++ {
		rs = append(rs, chars[rand.Intn(leng)])
	}
	return string(rs)
}

type Connectors interface {
	UUID() string
	Write(src []byte) error
	Addr() net.Addr
	RegisterTime() time.Time
	OnlineMS() int64
}
