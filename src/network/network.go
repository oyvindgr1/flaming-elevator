
package network
import (
	"net"
	"fmt"
	"time"
	"encoding/json"
	"elevtypes"
	"strings"
)




func SendInfo(status_chan <-chan elevtypes.Message) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	if err_conv_ip != nil {
			fmt.Println("error:", err_conv_ip)
			}
	status_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	if err_dialudp != nil {
			fmt.Println("error:", err_dialudp)
		}
	for {
		select {
		case status := <-status_chan:
			time.Sleep(1000 * time.Millisecond)
			b, err_Json := json.Marshal(state1)
        		if err_Json != nil {
			        fmt.Println("error with JSON")
			        fmt.Println(err_Json)
        		}
			for i := 0; i < 4; i++ {
				_, err1 := status_sender.Write(b)
			        if err1 != nil {
	                		fmt.Println("Error writing data to server. Waiting for 10 seconds before sending again.")
	                		time.Sleep(10*time.Second)
	                		fmt.Println(err1)
	       			        break
	       			 }
			}
		}
	}
}


func ReadInfo(status_chan chan<- elevtypes.Message, netIsAlive chan<- bool) {
	var msg elevtypes.Message
	b := make([]byte, 1024)
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	addr, err := net.ResolveUDPAddr("udp", getIP()+NetworkPort
	if err_ != nil {
			fmt.Println("error:", err)
			netIsAlive <- false
			return
		}
	if err_conv_ip_listen != nil {
			fmt.Println("error:", err_conv_ip_listen)
			netIsAlive <- false
			return
		}
	status_receiver, err_listen := net.ListenUDP("udp", laddr)
	if err_listen != nil {
			fmt.Println("error:", err_listen)
			netIsAlive <- false
			return
		}
	for {
		time.Sleep(1000 * time.Millisecond)
		n, raddr, _ := status_receiver.ReadFromUDP(b)
		err_decoding := json.Unmarshal(b[0:n], &msg)
		if err_decoding != nil {
			fmt.Println("error decoding client msg")
		}
		Client_map[raddr.String()] = msg;
		
		for key := range Client_map {
		    fmt.Println("%s", key)
		}
	}
}

func getIP() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Printf("Error with dialing Google.com", err)
		return "localhost"
	} else {
		return strings.Split(string(conn.LocalAddr().String()), ":")[0]
	}
}
