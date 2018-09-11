package anthem

import (
	"io"
	"log"
	"net"
)

func SerToCli(ser, cli net.Conn) {

	log.Println("start translate")

	go io.Copy(cli, ser)
	io.Copy(ser, cli)
}

type Msg struct {
	Type string `json:"type"` // beats,msg,
	Port string `json:"port"`
}
