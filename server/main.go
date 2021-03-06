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


	go listen(originP, originNetChan)
	go listen(transP, transNetChan)
	go listen(managerP, managerNetChan)

	for {
		conn := <-managerNetChan
		log.Println("get manager")
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

				go func() {
					defer origin.Close()
					defer trans.Close()
					if err:=anthem.SerToCli(origin, trans);err!=nil{
						log.Println(
							origin.LocalAddr().String(),
							trans.LocalAddr().String(),
							err)
					}
				}()
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

func listen(port string, netChan chan net.Conn) {

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
			log.Println("new connection from", conn.RemoteAddr().String())
			netChan <- conn
		}
	}()
}
