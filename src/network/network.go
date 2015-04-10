
package network
import (
	"net"
	"fmt"
	"time"
	"encoding/json"
	"elevtypes"
	"strings"
)




func SendStatus(status_send_chan <-chan elevtypes.Status) {
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
		case status := <-status_send_chan:
			time.Sleep(20 * time.Millisecond)
			b, err_Json := json.Marshal(status)
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


func ReadStatus(statusmap_send_chan chan<- map[string]elevtypes.Status, netIsAlive chan<- bool) {
	var statusMap map[string]elevtypes.Status	
	statusMap = make(map[string]elevtypes.Status)
	var msg elevtypes.Status
	b := make([]byte, 1024)
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	//addr, err := net.ResolveUDPAddr("udp", getIP()+NetworkPort
	/*if err_ != nil {
			fmt.Println("error:", err)
			netIsAlive <- false
			return
		}*/
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
	go func() {	
		for {
			time.Sleep(5 * time.Millisecond)
			n, raddr, err_read := status_receiver.ReadFromUDP(b)
			if err_read != nil {
				fmt.Println("Error reading from UDP")
			} else {
				err_decoding := json.Unmarshal(b[0:n], &msg)
				if err_decoding != nil {
					fmt.Println("Error decoding client msg")
				} else {
					statusMap[raddr.String()] = msg
					for key, _ := range statusMap {
						fmt.Printf(" IP: %s \n", raddr.String())
						PrintMatrix(statusMap[key].UnprocessedMatrix)
					}
				}
			}
		}
	}()
	
	go func() {
		for {
			select {
			case statusmap_send_chan <- statusMap:
				time.Sleep(10*time.Millisecond)	
			}		
		}
	}()
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

func PrintMatrix(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS-1]int) {
	time.Sleep(1*time.Second)
	fmt.Println("Floor \t UP \t DOWN")	
	for i := 0; i < elevtypes.N_FLOORS; i++ { 
		fmt.Printf("\n %d \t", i)
		for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
			fmt.Printf("%d \t ", matrix[i][j])	
		}
	}
}
