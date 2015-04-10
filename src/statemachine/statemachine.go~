package statemachine

import(
	"driver"
	"time"
	"fmt"
	"elevtypes"
	
)

const(
	UP   = 1
	DOWN = -1
	REST = 0
)

const (
	FLOOR_REACHED Event_enum = iota
	NEW_ORDER
	NO_ORDERS
	UNDEFINED
)


	
func ElevatorInit() {
	init = driver.Init()
	if init == 0 {
		return 0
	} else {
		if driver.GetFloorSensorSignal() != -1 {
		} else {
			driver.SetSpeed(DOWN*300)
			floor := driver.GetFloorSensorSignal()
			for floor == -1 {
				floor = driver.GetFloorSensorSignal()
			}
			ElevatorBrake(UP)
		}
		fmt.Printf("Initialized\n")
		return 1
	}
}

func StateMachine(orderList []Order) {	
	var cur_order Order
	var prev_order Order
	var status Status{}
	var state State_enum
	previous_floor_chan := make(chan Order, 1)
	delete_order_chan := make(chan Order, 1)
	state_update_chan := make(chan State_t, 1)
	var ButtonMatrix [N_BUTTONS][N_FLOORS]bool

	go func () {
		//Oppdatere state-variablen etter kjøring av funksjonene i state-switch
		for {
			select {
			case   
			case 
			case 
			}
		}
			


	}()

	go func () {
		

		

		
		for {
			time.Sleep(10 * time.Millisecond)
			switch state {
			case WAIT:
				event = wait()
			case RUN:
				event = (state_update_chan, &state)
			case OPEN:
				event = (state_update_chan, &state)
		}
	}()
}

func wait() {
		


}





func run(thisOrder Order, previous_floor_chan chan Order, state *elevtypes.State_enum, state_update_chan chan elevtypes.State_enum) {
	if *state != Running {
		*state = Running
		driver.SetSpeed(300 * thisOrder.Dir)
		state_update_chan <- Running	
	}
	cur_floor = driver.GetFloorSensorSignal()
	if cur_floor != -1 {
		selfOrder := Order{}
		selfOrder.Dir = thisOrder.Dir
		selfOrder.Floor = cur_floor
		previous_floor_chan <- selfOrder
		driver.SetLightFloorIndicator(cur_floor)
		if cur_floor == thisOrder.Floor {
			ElevatorBrake(thisOrder.Dir)
			return FLOOR_REACHED
		}
	if cur_floor == thisOrder.Floor {
		ElevatorBrake(thisOrder.Floor)
		return FLOOR_REACHED
		}
	}
	return NEW_ORDER
 }



func door(state_update_chan chan elevtypes.State_enum, state *elevtypes.State_enum) {
	if driver.GetFloorSensorSignal() != -1 {
		if *state != Door {
			*state = Door
			state_update_chan <- Door
			driver.SetDoorOpenLamp(1)
		}
		time.Sleep(3*time.Second)
		driver.SetDoorOpenLamp(0)
		return NEW_ORDER
	} else {
		return UNDEFINED
	}

}

func undefined(state_update_chan chan elevtypes.State_enum, state *elevtypes.State_enum)
	if *state != Undefined {
		*state = Undefined
		state_update_chan <- Undefined
	}
	return UNDEFINED
}

func ElevatorBrake(dir int) {
	driver.SetSpeed(dir*300)
	time.Sleep(time.Millisecond*20)
	driver.SetSpeed(0)
}


