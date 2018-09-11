package anthem

import (
	"io"
	"log"
	"net"
)


type errContain chan error
var(
	errList = make( chan errContain,1024 ) //管道池
)
func generalErrChan()errContain{
	if len(errList)==0{
		return make(chan error,1)
	}
	return <-errList
}
func recycleErrChan(c errContain){
	errList<-c
}

func SerToCli(ser, cli net.Conn)error{

	log.Println("start translate")

	ch:=generalErrChan()
	defer recycleErrChan(ch)

	go func() {
		_,err:=io.Copy(cli,ser)
		if err!=nil{
			ch<-err
		}
	}()
	if len(ch)!=0{
		return <-ch
	}
	_,err:=io.Copy(ser,cli)
	if err!=nil{
		ch<-err
	}
	if len(ch)!=0{
		return <-ch
	}
	return nil
}

type Msg struct {
	Type string `json:"type"` // beats,msg,
	Port string `json:"port"`
}
