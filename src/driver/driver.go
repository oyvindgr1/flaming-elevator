package driver

import (
	"elevtypes"
	"fmt"
	"math"
	"time"
)

//ORDER_UP = 0, ORDER_DOWN = 1, ORDER_INTERNAL = 2

/*type OrderDirection int
const (
	ORDER_UP OrderDirection = iota
	ORDER_DOWN
	ORDER_INTERNAL
)*/
var lampChannelMatrix = [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannelMatrix = [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Init() int {
	if !IoInit() {
		fmt.Printf("IO initiated\n")
	} else {
		fmt.Printf("IO not initiated\n")
		return 0
	}
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		if i != 0 {
			SetButtonLamp(1, i, 0)
		}
		if i != elevtypes.N_FLOORS-1 {
			SetButtonLamp(0, i, 0)
		}
		SetButtonLamp(2, i, 0)
	}
	SetStopLamp(0)
	SetDoorOpenLamp(0)
	SetLightFloorIndicator(0)
	if GetFloorSensorSignal() != -1 {
		} else {
			SetSpeed(-1 * 300)
			floor := GetFloorSensorSignal()
			for floor == -1 {
				floor = GetFloorSensorSignal()
			}
		SetSpeed(100)
		time.Sleep(time.Millisecond * 20)
		SetSpeed(0)
		}
	return 1
}

func SetStopLamp(i int) {
	if i == 1 {
		Set_bit(LIGHT_STOP)
	} else {
		Clear_bit(LIGHT_STOP)
	}
}

func SetDoorOpenLamp(i int) {
	if i == 1 {
		Set_bit(LIGHT_DOOR_OPEN)
	} else {
		Clear_bit(LIGHT_DOOR_OPEN)
	}
}

func GetButtonSignal(dir int, floor int) int {
	// Need error handling before proceeding$
	if Read_bit(buttonChannelMatrix[floor][dir]) {
		return 1
	} else {
		return 0
	}
}

func SetSpeed(speed int) {
	last_speed := 0
	if speed > 0 {
		Clear_bit(MOTORDIR)
	} else if speed < 0 {
		Set_bit(MOTORDIR)
	} else if last_speed < 0 {
		Clear_bit(MOTORDIR)
	} else if last_speed > 0 {
		Set_bit(MOTORDIR)
	}
	last_speed = speed
	Write_analog(MOTOR, int(2048+4*math.Abs(float64(speed))))
}

func SetButtonLamp(dir int, floor int, val int) {
	/*if floor >= 0 && floor < N_FLOORS {
	if dir == "up" && floor == 3 {
		fmt.Printf("The current direction and floor does not exist (up and 4)")
		return
	}
	if dir == "down" && floor == 0 {
		fmt.Printf("The current direction and floor does not exist (down and 0)")
		return
	}*/
	if val == 1 {
		Set_bit(lampChannelMatrix[floor][dir])
	} else {
		Clear_bit(lampChannelMatrix[floor][dir])
	}
	/*} else {
		fmt.Printf("Floor and direction is out of bounds")
		return
	}*/
}

func GetObstructionSignal() bool {
	return Read_bit(OBSTRUCTION)
}

func GetStopSignal() bool {
	return Read_bit(STOP)
}

func GetFloorSensorSignal() int {
	if Read_bit(SENSOR_FLOOR1) {
		return 0
	} else if Read_bit(SENSOR_FLOOR2) {
		return 1
	} else if Read_bit(SENSOR_FLOOR3) {
		return 2
	} else if Read_bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func SetLightFloorIndicator(floor int) {
	// Binary encoding. One light must always be on.
	if floor >= 0 && floor < elevtypes.N_FLOORS {
		switch floor {
		case 0:
			Clear_bit(LIGHT_FLOOR_IND1)
			Clear_bit(LIGHT_FLOOR_IND2)
		case 1:
			Clear_bit(LIGHT_FLOOR_IND1)
			Set_bit(LIGHT_FLOOR_IND2)
		case 2:
			Set_bit(LIGHT_FLOOR_IND1)
			Clear_bit(LIGHT_FLOOR_IND2)
		case 3:
			Set_bit(LIGHT_FLOOR_IND1)
			Set_bit(LIGHT_FLOOR_IND2)
		}
	} else {
		fmt.Printf("SetFloorIndicator: Elevator out of range.")
	}
}

/*
func ClearAllLights() {
	fmt.Printf("All lights cleared")
	SetDoorOpenLamp(0)
	SetStopLamp(0)
	for i := 0; i < elevtypes.N_FLOORS; i++ {

		if i > 0 {
			SetButtonLamp(1, i, 0)
		}
		if i < elevtypes.N_FLOORS-1 {
			SetButtonLamp(0, i, 0)
		}
	SetButtonLamp(2, i, 0)
	}
}
*/
