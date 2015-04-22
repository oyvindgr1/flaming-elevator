package main

import (
	"driver"
	"elevtypes"
	"network"
	"order"
	"statemachine"
	//"strconv"
	"fmt"
	"strings"
)

func main() {
	
	statusmap_chan := make(chan map[string]elevtypes.Status, 1)
	status_chan := make(chan elevtypes.Status, 1)
	
	orders_todo_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_local_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_external_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, 1)

	netIsAlive := make(chan bool)
	driver.Init()
	//var ip = net.IPv4(129,241,187,153)
	/*var state1 elevtypes.Status
	state1.CurFloor = 1
	state1.Dir = 1
	state1.OrderList = nil
	state1.UnprocessedOrders = orderList
	


	go func() {
		for {
			state1.UnprocessedOrders = orderList
			select {
			case status_chan <- state1:
			}
		}
	}()
	*/
	//var localIP int
	var anIP string
	anIP = "129.241.187.154"
	astring := strings.Split(anIP, ".")
	
	var anIP2 string
	anIP2 = "129.241.187.141"
	astring2 := strings.Split(anIP2, ".")

	if astring[3] > astring2[3] {
		fmt.Println("Success ")
	}
	 
/*	localIP, e := strconv.Atoi(astring)
	if e != nil {
		fmt.Println(e)
	}
	
	fmt.Printf("IP AS INT: %d \n", localIP)
*/	
	
	
	go order.OrderListener(orders_local_elev_chan, orders_external_elev_chan)
	go order.OrdersFromNetwork(orders_local_elev_chan, statusmap_chan,orders_external_elev_chan, orders_todo_elev_chan)
	go order.ErrorRecovery()
	
	go network.SendStatus(status_chan)
	go network.ReadStatus(statusmap_chan, netIsAlive, orders_todo_elev_chan)


	go statemachine.StateMachine(orders_local_elev_chan, orders_external_elev_chan, status_chan)

	/*go func() {
		for {
			select{
			case aBool := <- netIsAlive:
				fmt.Printf(" Net alive: %t", aBool)
			}
		}
	}()*/

	select {}

	/**	driver.IoInit()

	time.Sleep(2*time.Second)
	fmt.Printf("Init:\n")
	driver.Init()
	fmt.Printf("Current floor: ", driver.GetFloorSensorSignal())
	driver.SetLightFloorIndicator(driver.GetFloorSensorSignal())
	fmt.Printf("\nDrive!")
	driver.SetSpeed(300)
	time.Sleep(2*time.Second)
	fmt.Printf("\nStop!")
	driver.SetSpeed(0)

	fmt.Printf("\nButton lamps on")
	driver.SetButtonLamp(0, 0, 1)
	driver.SetButtonLamp(0, 1, 1)
	driver.SetButtonLamp(0, 2, 1)

	driver.SetButtonLamp(1, 1, 1)
	driver.SetButtonLamp(1, 2, 1)
	driver.SetButtonLamp(1, 3, 1)
	time.Sleep(1*time.Second)
	fmt.Printf("\nButton lamps off")
	driver.SetButtonLamp(0, 0, 0)
	driver.SetButtonLamp(0, 1, 0)
	driver.SetButtonLamp(0, 2, 0)

	driver.SetButtonLamp(1, 1, 0)
	driver.SetButtonLamp(1, 2, 0)
	driver.SetButtonLamp(1, 3, 0)


	fmt.Printf("\nOpen door")
	driver.SetDoorOpenLamp(1)
	time.Sleep(2*time.Second)
	fmt.Printf("\nStop lamp")
	driver.SetStopLamp(1)
	time.Sleep(2*time.Second)
	fmt.Printf("\nDoor open")
	driver.SetDoorOpenLamp(1)
	time.Sleep(2*time.Second)
	fmt.Printf("\nObstruction signal")
	driver.GetObstructionSignal()
	time.Sleep(2*time.Second)
	fmt.Printf("\nStop signal")
	driver.GetStopSignal()
	time.Sleep(2*time.Second)
	fmt.Printf("\nFloor sensor signal")
	driver.GetFloorSensorSignal()
	time.Sleep(2*time.Second)
	fmt.Printf("\nButton signal")
	driver.GetButtonSignal(1,1)

	*/

}
