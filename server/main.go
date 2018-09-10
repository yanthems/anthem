package main

import (
	"net"
	"log"
	"encoding/json"
	"io"
	"github.com/yanthems/anthem"

)



func main(){

	originP:="9090"
	targetP:="12345"

	transP:="23333"

	go listen(originP,originNetChan,"")
	go listen(transP,transNetChan,targetP)

	trans_con:=<-transNetChan
	origin_con:=<-originNetChan

	//trans(trans_con,origin_con)

	anthem.SerToCli(origin_con,trans_con)
}

func hello(conn net.Conn,target string){
	hi:=map[string]interface{}{
		"target":target,
	}
	raw,err:=json.Marshal(hi)
	if err!=nil{
		log.Println(err)
		return
	}
	log.Println(hi,"<===")
	conn.Write(raw)
}
var (
	originNetChan = make(chan net.Conn,1024)
	transNetChan = make(chan net.Conn,1024)
)
func listen(port string,netChan chan net.Conn,target string){

	ser, err := net.Listen("tcp",net.JoinHostPort("127.0.0.1", port))
	if err != nil {
		log.Println("server listen error")
		return
	}

	go func(){
		for {
			conn, err := ser.Accept()
			if err != nil {
				log.Println("server accept error")
				return
			}
			if target!="" {
				hello(conn, target)
			}else{
				log.Println("new connection from",conn.RemoteAddr().String())
			}
			netChan<-conn
		}
	}()
}

func trans(trans_con,origin_con net.Conn){
	defer trans_con.Close()
	defer origin_con.Close()
		log.Println("start translate")
		go io.Copy(origin_con,trans_con)
		io.Copy(trans_con,origin_con)
		origin_con.Close()
}