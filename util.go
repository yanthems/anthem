package anthem

import (
	"io"
	"log"
	"net"
)

func SerToCli(ser, cli net.Conn) {
	defer ser.Close()
	defer cli.Close()

	log.Println("start translate")

	go io.Copy(cli, ser)
	io.Copy(ser, cli)
}

type Msg struct {
	Port string `json:"port"`
}
