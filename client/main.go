package main

import (
	"net"
	"log"
		"encoding/json"
	"io"
	"bytes"
)


func main(){
	transP:="23333"
	transH:="127.0.0.1"

	target:=hello(transH,transP)
	connect(target)

	trans_con:=<-transNetChan
	target_con:=<-transNetChan

	log.Println(transP,"->",target)
	//trans(trans_con,target_con)

	anthem.SerToCli(trans_con,target_con)
}

var (
	targetNetChan = make(chan net.Conn,1024)
	transNetChan = make(chan net.Conn,1024)
)

func connect(port string){
	local,err:=net.Dial("tcp",net.JoinHostPort("127.0.0.1",port))
	if err!=nil{
		log.Println(err)
		return
	}
	log.Println("connect to local",port)
	transNetChan<-local
}
func hello(host,port string)string{
		trans_ser, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return ""
		}

		raw:=make([]byte,2048)
		_,err= trans_ser.Read(raw)
		if err!=nil{
			log.Println(err)
			return ""
		}
		raw=bytes.TrimRightFunc(raw, func(r rune) bool {
			return r=='\x00'
		})

	hi:=map[string]interface{}{}
		if err:=json.Unmarshal(raw,&hi);err!=nil{
			log.Println(err)
			return ""
		}
		if tp,exist:=hi["target"];exist{
			if str,ok:=tp.(string);ok{
				result:=str
				log.Println("target =",result,"<===")
				transNetChan<- trans_ser
				return result
			}
		}
		return ""
}

func trans(trans_con,target_con net.Conn){
	defer trans_con.Close()
	defer target_con.Close()
	for{
		log.Println("start translate")
		go io.Copy(trans_con,target_con)
		io.Copy(target_con,trans_con)
	}
}