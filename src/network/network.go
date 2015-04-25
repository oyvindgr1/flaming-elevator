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

func SendElevatorStatus(elevator_status_send_chan <-chan elevtypes.ElevatorStatus) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	if err_conv_ip != nil {
		fmt.Println("error:", err_conv_ip)
	}
	elevator_status_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	if err_dialudp != nil {
		fmt.Println("error:", err_dialudp)
	}
	for {
		select {
		case elevatorStatus := <-elevator_status_send_chan:
			b, err_Json := json.Marshal(elevatorStatus)
			if err_Json != nil {
				fmt.Println("error with JSON")
				fmt.Println(err_Json)
			}

			for i := 0; i < 4; i++ {
				_, err1 := elevator_status_sender.Write(b)
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

func ReadElevatorStatus(elevator_status_map_send_chan chan<- map[string]elevtypes.ElevatorStatus, orders_from_unresponsive_elev_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	var elevatorStatusMap map[string]elevtypes.ElevatorStatus
	elevatorStatusMap = make(map[string]elevtypes.ElevatorStatus)
	var respondingElevatorList []string 
	var msg elevtypes.ElevatorStatus
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	if err_conv_ip_listen != nil {
		fmt.Println("error:", err_conv_ip_listen)
		return
	}
	elevator_status_receiver, err_listen := net.ListenUDP("udp", laddr)
	if err_listen != nil {
		fmt.Println("error:", err_listen)
		return
	}

	//Read elevator status from network and place in elevatorStatusMap
	go func() {
		for {
			b := make([]byte, 1024)
			n, raddr, err_read := elevator_status_receiver.ReadFromUDP(b)
			
			if err_read != nil {
				fmt.Println("Error reading from UDP")
			} else {
				err_decoding := json.Unmarshal(b[0:n], &msg)
				if err_decoding != nil {
					fmt.Println("Error decoding client msg")
				} else {
					elevatorStatusMap[raddr.String()] = msg
					includeRespondingElevator(&respondingElevatorList, raddr.String())
					elevator_status_map_send_chan <- elevatorStatusMap
				}
			}
			
		}
	}()
	
	//Check for nonresponding elevators
	go func() {
		var OutDatedElevatorStatusMap map[string]elevtypes.ElevatorStatus
		OutDatedElevatorStatusMap = make(map[string]elevtypes.ElevatorStatus)
		for {
			respondingElevatorList = respondingElevatorList[:0]
			time.Sleep(1 * time.Second)
			for key, _ := range elevatorStatusMap {
				if !isInList(respondingElevatorList, key) {
					for key, v := range elevatorStatusMap {
						OutDatedElevatorStatusMap[key] = v 
					}
					delete(elevatorStatusMap, key)
					if LocalIPIsLowest(respondingElevatorList) {
						orders_from_unresponsive_elev_chan <- OutDatedElevatorStatusMap[key].OrderMatrix
												
						
					}				
				}
			}
		}
	}()
}

func LocalIPIsLowest(respondingElevatorList []string) bool {
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
	}

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
func includeRespondingElevator(respondingElevatorList *[]string, elevatorIP string ) {
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
