package tcp

import (
	"bytes"
	"fmt"
	"github.com/blinkat/blinker/invoker"
	"github.com/blinkat/blinker/spectre"
	"net"
	"reflect"
	"time"
)

const end_byte = byte(0)
const uuid_size = 24

type connecter struct {
	conn          net.Conn
	uuid          string
	running       bool
	buffer        []*spectre.Pack
	invoker       *invoker.Invoker
	register_time time.Time
	handler       Handler
	key           []byte // client pub key
}

func (c *connecter) listen() error {
	c.running = true
	go c.loop()
	return nil
}

func (c *connecter) loop() {
	buf := make([]byte, subpackage_size)
	data := make(chan []byte, 0)
	go c.read_done(data)

	for c.running {
		_, err := c.conn.Read(buf)
		pack, err := spectre.Unpack(buf, subpackage_size)
		if err != nil {
			fmt.Println(err)
			c.handler.Error(err)
			pack.Sign = false
			continue
		}
		if len(c.buffer) == pack.Header.Total {
			c.unpack(data)
		} else {
			c.buffer = append(c.buffer, pack)
		}
	}
}

func (c *connecter) unpack(chand chan []byte) {
	var buf bytes.Buffer
	for _, v := range c.buffer {
		if !v.Sign {
			return
		}
		buf.Write(v.Body)
	}

	chand <- buf.Bytes()
	c.buffer = []*spectre.Pack{}
}

func (c *connecter) read_done(chand chan []byte) {
	for c.running {
		select {
		case data := <-chand:
			var err error
			if data, err = c.decrypto(data); err != nil {
				c.handler.Error(err)
				continue
			}

			if c.key == nil {
				c.key = data
			} else {
				//fmt.Println(string(data))
				data, err = c.invoker.Decrypt(data)
				if err != nil {
					c.handler.Error(err)
					continue
				}
				c.handle_message(data)
			}
		}
	}
}

func (c *connecter) handle_message(msg []byte) error {
	fmt.Println(msg)
	return nil
}

func (c *connecter) decrypto(data []byte) ([]byte, error) {
	if c.invoker != nil {
		return c.invoker.Decrypt(data)
	}
	return data, nil
}

func (c *connecter) encrypto(data []byte) ([]byte, error) {
	if c.invoker != nil {
		return c.invoker.Encrypt(data)
	}
	return data, nil
}

func (c *connecter) put_key() error {
	k, e := c.invoker.PutKey()
	if e != nil {
		return e
	}
	return spectre.SendPack(k, subpackage_size, c.conn)
}

// func (c *connecter) write(data []byte) error {
// 	leng := len(data)
// 	max := subpackage_size - 1

// 	if leng >= max {
// 		sub := leng / max
// 		if leng%max != 0 {
// 			sub += 1
// 		}

// 		for i := 0; i < sub; i++ {
// 			last := (i + 1) * max
// 			if last > leng {
// 				_, err := c.conn.Write(append([]byte{1}, data[i*max:]...))
// 				if err != nil {
// 					return err
// 				}
// 				break
// 			}

// 			_, err := c.conn.Write(append([]byte{0}, data[i*max:last]...))
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	} else {
// 		_, err := c.conn.Write(append([]byte{1}, data...))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// interface
func (c *connecter) UUID() string {
	return c.uuid
}

func (c *connecter) Addr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connecter) RegisterTime() time.Time {
	return c.register_time
}

func (c *connecter) OnlineMS() int64 {
	return (time.Now().UnixNano() - c.register_time.UnixNano()) / int64(time.Millisecond)
}

func (c *connecter) Write(data []byte) error {
	var err error
	if data, err = c.invoker.EncryptByKey(data, c.key); err != nil {
		return err
	}
	fmt.Println(data)

	return spectre.SendPack(data, subpackage_size, c.conn)
}

func newConnecter(conn net.Conn, iv *invoker.Invoker) *connecter {
	c := &connecter{
		conn:          conn,
		uuid:          spectre.NewGUID(24),
		running:       false,
		buffer:        []*spectre.Pack{},
		invoker:       iv,
		register_time: time.Now(),
	}

	if h, err := copy_handler(); err != nil {
		fmt.Println(err)
		return nil
	} else {
		c.handler = h
	}

step:
	if _, ok := clients[c.uuid]; ok {
		c.uuid = spectre.NewGUID(24)
		goto step
	}
	fmt.Println(namespace+": connect client:", c.conn.LocalAddr())
	clients[c.uuid] = c
	return c
}

func copy_handler() (Handler, error) {
	t := reflect.TypeOf(handler).Elem()
	return reflect.New(t).Interface().(Handler), nil
}
