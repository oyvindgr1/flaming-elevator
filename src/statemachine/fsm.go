package statemachine

import(
	"driver"
	"time"
	"fmt"
	
	
)

const(
	Upwards   = 1
	Downwards = -1
	Rest = 0
)

func ElevatorInit() {
	if driver.Init() == 0 {
		return 0
	} else {
		if driver.GetFloorSensorSignal() != -1 {
		} else {
			driver.SetSpeed(Downwards*300)
			floor := driver.GetFloorSensorSignal()
			for floor == -1 {
				floor = driver.GetFloorSensorSignal()
			}
			ElevatorBrake(Upwards)
		}
		fmt.Printf("Initialized\n")
		return 1
	}
}

func EventHandler() {
}

func ElevatorBrake(dir int) {
	driver.SetSpeed(dir*300)
	time.Sleep(time.Millisecond*20)
	driver.SetSpeed(0)
}
	
	
