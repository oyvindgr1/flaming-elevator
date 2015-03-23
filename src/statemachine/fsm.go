package statemachine

import(
	"driver"
	
	
)

const(
	Upwards   = 1
	Downwards = -1
	Rest = 0
)

func ElevatorInit() {
	if driver.ElevInit() == 0 {
		return 0
	} else {
		if driver.ElevGetFloorSensorSignal() != -1 {
		} else {
			driver.ElevSetSpeed(Downwards*300)
			floor := driver.ElevGetFloorSensorSignal()
			for floor == -1 {
				floor = driver.ElevGetFloorSensorSignal()
			}
			ElevatorBrake()
		}
		fmt.Printf("Initialized\n")
		return 1
	}
}

func EventHandler() {
}

func ElevatorBrake() {
	driver.ElevSetSpeed(Upwards*300)
	time.After(time.Millisecond*20)
}
	
	
