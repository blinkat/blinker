package spectre

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// header reserve byte number
const header_reserve = 6

type Header struct {
	Serial int
	Total  int
	Length int // current pack byte size
}

type Pack struct {
	Header *Header
	Body   []byte
	Sign   bool
}

func PackReserve() int {
	return header_reserve
}

func Unpack(src []byte, sub int) (*Pack, error) {
	if len(src) <= header_reserve {
		return nil, fmt.Errorf(namespace + ": src length smaller.")
	}
	head := src[:header_reserve]
	body := src[header_reserve:]

	header := &Header{}
	binary.Read(bytes.NewBuffer(head[:2]), binary.BigEndian, &header.Total)
	binary.Read(bytes.NewBuffer(head[3:5]), binary.BigEndian, &header.Serial)
	binary.Read(bytes.NewBuffer(head[5:7]), binary.BigEndian, &header.Length)

	return &Pack{
		Header: header,
		Body:   body[:header.Length],
	}, nil
}

func SendPack(src []byte, sub int, conn net.Conn) error {
	sub -= header_reserve
	leng := len([]byte(src))
	splt := leng / sub
	if leng%sub != 0 {
		splt += 1
	}

	packs := make([]*Pack, splt)
	for k := range packs {
		p := &Pack{
			Header: &Header{
				Serial: k,
				Total:  splt,
			},
			Sign: true,
		}

		if k == splt-1 {
			p.Header.Length = leng - (k * sub)
			p.Body = src[k:]
		} else {
			p.Header.Length = sub
			p.Body = src[k : (k+1)*sub]
		}

		packs[k] = p
	}

	// send

	for _, v := range packs {
		var head bytes.Buffer
		binary.Write(&head, binary.BigEndian, v.Header.Total)
		binary.Write(&head, binary.BigEndian, v.Header.Serial)
		binary.Write(&head, binary.BigEndian, v.Header.Length)

		_, err := conn.Write(append(head.Bytes(), v.Body...))
		if err != nil {
			return err
		}
	}
	return nil
}
