
package netmod
import (
	"net"
	"fmt"
	"time"
	"encoding/json"
	"declarations"
)



func Send_status(status_chan chan State) {
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
			_, err1 := status_sender.Write(b)
		        if err1 != nil {
                		fmt.Println("error writing data to server")
                		fmt.Println(err1)
       			         return
       			 }
		default:
		}
	}
}



func Read_status(Client_map map[string]State) {
	var m State
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	if err_conv_ip_listen != nil {
			fmt.Println("error:", err_conv_ip_listen)
		}
	status_receiver, err_listen := net.ListenUDP("udp", laddr)
	if err_listen != nil {
			fmt.Println("error:", err_listen)
		}
	for {
		time.Sleep(1000 * time.Millisecond)
		b := make([]byte, 1024)
		n, raddr, _ := status_receiver.ReadFromUDP(b)
		err_decoding := json.Unmarshal(b[0:n], &m)
		if err_decoding != nil {
			fmt.Println("error decoding client msg")
		}
		Client_map[raddr.String()] = m;
		
		for key := range Client_map {
		    fmt.Println("%s", key)
		}
	}
}
