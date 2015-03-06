package main
import ( 
	"net" 
	"fmt" 
	"time"
	"encoding/json" 
)

type state struct {
    curfloor int 
    heisnr int 
    aString string
}

var m state

func main() {
	
	udpAddr, err := net.ResolveUDPAddr("udp", ":35000")
	fmt.Println(udpAddr)	
	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("error resolving UDP address on ", ":30000")
                fmt.Println(err)
                return
	}	
	defer ln.Close()
	

	var buf []byte = make([]byte, 1500)

        for {

                time.Sleep(100 * time.Millisecond)

                n, address, err := ln.ReadFromUDP(buf)

		err := json.Unmarshal(buf,&m)

                if err != nil {
                        fmt.Println("error reading data from connection")
                        fmt.Println(err)
                        return
                }

                if address != nil {

                        fmt.Println("got message from ", address, " with n = ", n)
	
                        if n > 0 {
                                fmt.Println("from address", address, "got message:", string(m.heisnr)
                        }
                }
        }
}
