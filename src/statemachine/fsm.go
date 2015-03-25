package statemachine

import(
	"driver"
	"time"
	"fmt"
	
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
)

const (
	Running State_enum = iota
	Idle 
	atFloor
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

func StateMachine(current_order <-chan Order, previous_order <-chan Order, delete_order <-chan Order, state_chan <-chan State_enum, button_update_channel <-chan[N_BUTTONS][N_FLOORS]bool) {	
	var cur_order Order
	var prev_order Order
	previous_floor_chan := make(chan Order, 1)
	delete_order_chan := make(chan Order, 1)
	state_update_chan := make(chan State_t, 1)
	var ButtonMatrix [N_BUTTONS][N_FLOORS]bool


	go func () {
		for {
			select {
			case 
			case 
			case 
			case 
		
			}
		}
	}()

	go func () {
		var state State_enum
		var event Event_enum
		for {
			time.Sleep(10 * time.Millisecond)
			switch event {
			case NEW_ORDER:
				event = run(prev_floor_chan, state_update_chan, cur_order, &state)
			case FLOOR_REACHED:
				event = door(state_update_chan, &state)
			case NO_ORDERS:
				event = wait(state_update_chan, &state)
		}
	}()
}

func run(thisOrder Order,................) {
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
		previous_floor_chan

 	
}

func wait(...........) {

}

func door(.............) {

}

func ElevatorBrake(dir int) {
	driver.SetSpeed(dir*300)
	time.Sleep(time.Millisecond*20)
	driver.SetSpeed(0)
}
