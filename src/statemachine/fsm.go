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
	NEW_DIRECTION
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

func StateMachine(current_order chan Order, previous_order chan Order, delete_order chan Order, state_chan chan State_enum, button_update_channel chan[N_BUTTONS][N_FLOORS]bool) {	
	var cur_order Order
	var prev_order Order
	get_previous_floor_chan := make(chan Order, 1)
	delete_order_chan := make(chan Order, 1)
	state_update_chan := make(chan State_t, 1)
	var ButtonMatrix [N_BUTTONS][N_FLOORS]bool


	go func () {
		for {
			select {
			case local_list = <- update_local_list_c:
			case head_order = <-head_order_c:
			case prev_order = <-get_prev_floor_c:
				prev_order_c <- prev_order
			case del_req := <-delete_order:
				del_order <- del_req
			case new_state := <-state:
				state_c <- new_state
			}
		}
	}();

	go func () {
		var state State_enum
		var event Event_enum
		for {
			time.Sleep(10 * time.Millisecond)
			switch event {
			case NEW_ORDER:
				event = Elevator_run(state_update_c, get_prev_floor_c, &state, head_order)
			case FLOOR_REACHED:
				event = Elevator_door(state_update_c, &state)
			case NO_ORDERS:
				event = Elevator_wait(state_update_c, &state)

				
		}
	}();
}
/*
func Elvator_wait(state_update_c chan State_t, state *State_t) {
	if *state != WAIT {
		*state = WAIT
		state_update_c <- RUN
	}
	return NO_ORDERS
}

func Elevator_run(state_update_c chan State_t, get_prev_floor_c chan Order, state *State_t, head_order Order) {
	if *state != RUN {
		*state = RUN
		driver.Elev_set_speed(300 * head_order.Dir)
		state_update_c <- RUN
	}
	current_floor := driver.Elev_get_floor_sensor_signal()
	if current_floor != -1 {
		var current Order
		current.Floor = current_floor
		current.Dir = head_order.Dir
		get_prev_floor_c <- current
		driver.Elev_set_floor_indicator(current_floor)
		if current_floor == head_order.Floor {
			Elevator_break(head_order.dir)
			return FLOOR_REACHED
		}
	}
	if current_floor == head_order.Floor {
		Elevator_break(head_order.Dir)
		return FLOOR_REACHED
	}
	if driver.Elev_get_stop_signal() {
		Elevator_break(head_order.Dir)
		return STOP
	}
	if driver.Elev_get_obstruction_signal() {
		Elevator_break(head_order.Dir)
		return OBSTRUCTION
	}
	return NEW_ORDER
}

func Elevator_door(state_update_c chan State_t, state *State_t) {
	if driver.Elev_get_floor_sensor_signal() != -1 {
		if *state != DOOR {
			*state = DOOR
			state_update_c <- DOOR
			driver.Elev_set_door_open_lamp(1)
		}
		time.Sleep(3000 * time.Millisecond)
		if driver.Elev_get_obstruction_signal() {
			return DOOR
		}
		if driver.Elev_get_stop_signal() {
			return STOP
		}
		driver.Elev_set_door_open_lamp(0)
		return NEW_ORDER
	}else {
		return UNDEF
	}
}*/

func ElevatorBrake(dir int) {
	driver.SetSpeed(dir*300)
	time.Sleep(time.Millisecond*20)
	driver.SetSpeed(0)
}
	
	
