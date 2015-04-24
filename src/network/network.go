package network

import (
	"elevtypes"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
	"strconv"

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
			b, err_Json := json.Marshal(status)
			if err_Json != nil {
				fmt.Println("error with JSON")
				fmt.Println(err_Json)
			}

			for i := 0; i < 4; i++ {
				//_, err1 := status_sender.Write(append([]byte("random ass junk"), b[0:]...))
				_, err1 := status_sender.Write(b)
				if err1 != nil {
					fmt.Println("Error writing data to server. Waiting for 10 seconds before sending again.")
					time.Sleep(10 * time.Second)
					fmt.Println(err1)
					break
				}
			}
		}
	}
}

func ReadStatus(statusmap_send_chan chan<- map[string]elevtypes.Status, netIsAlive chan<- bool, orders_from_unresponsive_elev_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	var statusMap map[string]elevtypes.Status
	statusMap = make(map[string]elevtypes.Status)
	var respondingElevatorList []string 
	var msg elevtypes.Status
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
			b := make([]byte, 1024)
			n, raddr, err_read := status_receiver.ReadFromUDP(b)
			
			if err_read != nil {
				fmt.Println("Error reading from UDP")
			} else {
				err_decoding := json.Unmarshal(b[0:n], &msg)
				if err_decoding != nil {
					fmt.Println("Error decoding client msg")
				} else {
					statusMap[raddr.String()] = msg
					includeElevator(&respondingElevatorList, raddr.String() )
					statusmap_send_chan <- statusMap
					/*for key, _ := range statusMap {
						fmt.Printf(" \nIP: %s  \n", key)
						PrintMatrix(statusMap[key].OrderMatrix)
					}*/
				}
			}
			
		}
	}()
	go func() {
		var OutDatedStatusMap map[string]elevtypes.Status
		OutDatedStatusMap = make(map[string]elevtypes.Status)
		for {
			respondingElevatorList = respondingElevatorList[:0]
			time.Sleep(1 * time.Second)
			for key, _ := range statusMap {
				if !isInList(respondingElevatorList, key) {
					//PrintMatrix(statusMap[key].OrderMatrix)					
					for key, v := range statusMap {
						OutDatedStatusMap[key] = v 
					}
					delete(statusMap, key)
					if IsLowestIP(respondingElevatorList) {
						orders_from_unresponsive_elev_chan <- OutDatedStatusMap[key].OrderMatrix
												
						
					}				
				}
			}
		}
	}()
}

func IsLowestIP(respondingElevatorList []string) bool {
	lowestIP := 1000
	localIP := strings.Split(GetIP(), ".")
	for i,_ := range respondingElevatorList {
		curIPtmp := strings.Split(respondingElevatorList[i], ":")[0]
		curIP := strings.Split(curIPtmp, ".")
		compare, err := strconv.Atoi(curIP[3])
		if err != nil {
			fmt.Println(err)
		}
		if compare < lowestIP {
			lowestIP = compare
		}
/*		if curIP[3] < lowestIP {
			lowestIP = curIP[3]
		}
		fmt.Println("Current IP: ", curIP[3])*/
	}
	//fmt.Println("Lowest IP: ", lowestIP)
	ownIP, err2 := strconv.Atoi(localIP[3])
	if err2 != nil {
			fmt.Println(err2)
		}
	if ownIP==lowestIP {
		return true
	}else {
		return false
	}
}

func isInList(respondingElevators []string, elevatorIP string) bool {
	for i,_ := range respondingElevators {
		if elevatorIP == respondingElevators[i] {
			return true
		}
	}
	return false
}
func includeElevator(respondingElevatorList *[]string, elevatorIP string ) {
	if !isInList(*respondingElevatorList, elevatorIP) {
		*respondingElevatorList = append(*respondingElevatorList, elevatorIP)
	}		 
}

func GetIP() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Printf("Error with dialing Google.com", err)
		return "localhost"
	} else {
		return strings.Split(string(conn.LocalAddr().String()), ":")[0]
	}
}

func PrintMatrix(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	fmt.Println("\nFloor \t UP \t DOWN")
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		fmt.Printf("\n %d \t", i)
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			fmt.Printf("%d \t ", matrix[i][j])
		}
	}
}
func InitMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS; y++ {
			matrix[x][y] = 0
		}
	}
}
