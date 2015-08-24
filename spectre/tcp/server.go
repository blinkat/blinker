package tcp

import (
	"fmt"
	"github.com/blinkat/blinker/invoker"
	"github.com/blinkat/blinker/spectre"
	"net"
)

const namespace = "blinker/spectre/tcp"
const subpackage_size = 1024

type Handler interface {
	Connected(conn spectre.Connectors) // connected event.
	Error(err error)
}

var clients map[string]*connecter

var listener net.Listener
var handler Handler
var get_invoker func() *invoker.Invoker

func ListenAndServer(port int) error {
	var err error
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	go listen_loop()
	return nil
}

func listen_loop() {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		var ctr *invoker.Invoker
		if get_invoker != nil {
			ctr = get_invoker()
		}
		cer := newConnecter(conn, ctr)
		cer.listen()
		cer.put_key()
		cer.handler.Connected(cer)
	}
}

func init() {
	clients = make(map[string]*connecter)
}

func SetInvoker(i func() *invoker.Invoker) {
	get_invoker = i
}

func SetHandler(h Handler) {
	handler = h
}
