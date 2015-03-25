
package main
import (
        "network"
	"time"
	"driver"
	"order"
	"fmt"
	"elevtypes"
)
è



func main(){
	var orderList []elevtypes.Order
	statusmap_chan := make(chan map[string]elevtypes.Status,1)
	status_chan := make(chan elevtypes.Status,1)
	
	orders_local_elev_chan := make(chan elevtypes.Order, 1)
	
		
	netIsAlive	:= make(chan bool)
	var 

	//var ip = net.IPv4(129,241,187,153)
	state1 := Status{1, 1,nil , nil,2}
	init := driver.Init()
	fmt.Printf("Driver initiated: %d\n", init)
	go func() {
		for {
			select {
			case status_chan <- state1:
			}
		}
	}()
		
	go network.ReadStatus(statusmap_chan, netIsAlive)
	go network.SendStatus(status_chan)
	
	go order.OrderListener(orders_local_elev_chan)
	
	go func() {
		for {
			select {
			case newOrder := <-orders_local_elev_chan:
				orderList = append(orderList, newOrder)
				fmt.Printf("Floor of new order: %d, Type of new order: %d\n", newOrder.Floor, newOrder.Dir)
			}
		}
	}()
	
	go func() {
		for {
			select{
			case aBool := <- netIsAlive:
				fmt.Printf(" Net alive: %t", aBool)
			}
		}
	}()

	
	time.Sleep(1000*time.Second)
	








































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
