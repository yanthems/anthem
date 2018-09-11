package main

import (
	"encoding/json"
	"github.com/yanthems/anthem"
	"log"
	"net"
)

func main() {

	originP := "9090"
	targetP := "12345"

	transP := "23333"
	managerP := "23334"

	go listen(originP, originNetChan, "")
	go listen(transP, transNetChan, "")
	go listen(managerP, managerNetChan, "")

	for {
		conn := <-managerNetChan

		go func(manager net.Conn) {

			defer manager.Close()

			for {

				origin := <-originNetChan
				if err := hello(manager, targetP); err != nil {
					log.Println(err)

					originNetChan <- origin
					return
				}
				trans := <-transNetChan
				go anthem.SerToCli(origin, trans)
			}
		}(conn)
	}

}

func hello(conn net.Conn, target string) error {
	hi := anthem.Msg{
		Port: target,
	}

	raw, err := json.Marshal(hi)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(hi, "<===")
	_, err = conn.Write(raw)
	if err != nil {
		return err
	}
	return nil
}

var (
	originNetChan  = make(chan net.Conn, 1024)
	transNetChan   = make(chan net.Conn, 1024)
	managerNetChan = make(chan net.Conn, 1024)
)

func listen(port string, netChan chan net.Conn, target string) {

	ser, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", port))
	if err != nil {
		log.Println("server listen error")
		return
	}

	go func() {
		for {
			conn, err := ser.Accept()
			if err != nil {
				log.Println("server accept error")
				return
			}
			if target != "" {
				hello(conn, target)
			} else {
				log.Println("new connection from", conn.RemoteAddr().String())
			}
			netChan <- conn
		}
	}()
}
