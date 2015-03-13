
package main
import (
"netmod"
)





func main(){
	
	
	state1 := netmod.State{LocalIP() ,"11", "22", "heihei"}
	Client_map := make(map[string]State)	

	go netmod.Read_status(Client_map)
	go netmod.Send_status(state1)
	
	
	
	/**
	driver.IoInit()
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
