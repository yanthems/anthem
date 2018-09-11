package main

import (
	"bytes"
	"encoding/json"
	"github.com/yanthems/anthem"
	"log"
	"net"
)

func main() {

	if !connect(RemoteHost, ManagerPort, managerNetChan) {
		return
	}

	go func() {

		conn := <-managerNetChan
		defer conn.Close()

		for {

			raw := make([]byte, 256)
			_, err := conn.Read(raw)
			if err != nil {
				log.Println(err)
				return
			}
			raw = bytes.TrimRightFunc(raw, func(r rune) bool {
				return r == '\x00'
			})
			hi := anthem.Msg{}
			if err := json.Unmarshal(raw, &hi); err != nil {
				log.Println(err)
				return
			}

			go hello(hi.Port)

		}

	}()
}

type Trans struct {
	Remote net.Conn
	Local  net.Conn
}

var (
	targetNetChan  = make(chan net.Conn, 1024)
	transNetChan   = make(chan net.Conn, 1024)
	managerNetChan = make(chan net.Conn, 1024)

	RemoteHost  = "127.0.0.1"
	TransPort   = "23333"
	ManagerPort = "23334"

	LocalHost = "127.0.0.1"
)

func connect(host, port string, netChan chan net.Conn) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("connect to ", host, port)
	netChan <- conn
	return true
}
func hello(port string) {

	log.Println("target port =", port, "<===")

	go func() {
		for {
			if connect(RemoteHost, TransPort, transNetChan) {
				break
			}
		}
	}()

	go func() {

		for {
			if connect(LocalHost, port, targetNetChan) {
				break
			}
		}

	}()

	trans := <-transNetChan
	target := <-targetNetChan

	anthem.SerToCli(trans, target)

}
