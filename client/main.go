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

	func() {

		conn := <-managerNetChan
		defer conn.Close()

		for {
			//log.Println("start read")
			raw := make([]byte, 256)
			_, err := conn.Read(raw)
			if err != nil {
				log.Println(err)
				return
			}
			raw = bytes.TrimRightFunc(raw, func(r rune) bool {
				return r == '\x00'
			})
			//log.Println(raw)
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
				log.Println("connect trans success")
				break
			}
		}
	}()

	go func() {

		for {
			if connect(LocalHost, port, targetNetChan) {
				log.Println("connect local success")
				break
			}
		}

	}()

	trans := <-transNetChan
	target := <-targetNetChan
	log.Println("start ",port)
	anthem.SerToCli(trans, target)

}
