package statemachine

import (
	"driver"
	"elevtypes"
	"fmt"
	"order"
	"time"
)

type State_enum int

const (
	WAIT State_enum = iota
	RUN_UP
	RUN_DOWN
	OPEN
)

func ElevatorInit() int {
	init := driver.Init()
	if init == 0 {
		return 0
	} else {
		if driver.GetFloorSensorSignal() != -1 {
		} else {
			driver.SetSpeed(-1 * 300)
			floor := driver.GetFloorSensorSignal()
			for floor == -1 {
				floor = driver.GetFloorSensorSignal()
			}
			elevatorBrake(1)
		}
		fmt.Printf("Initialized\n")
		return 1
	}
}

func StateMachine(orders_local_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, status_update_chan chan<- elevtypes.Status) {
	var status elevtypes.Status
	var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	var unprocessedOrdersMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int
	var state State_enum
	var floor int
	var serveDirection int
	var runDirection int
	state = WAIT
	order.InitMatrix(&orderMatrix)
	InitMatrix(&unprocessedOrdersMatrix)

	go func() {
		for {
			select {
			case orderMatrix := <-orders_local_elevator_chan:
				status.OrderMatrix = orderMatrix
			case unprocessedOrdersMatrix := <-orders_external_elevator_chan:
				status.UnprocessedOrdersMatrix = unprocessedOrdersMatrix
			}
		}
	}()

	go func() {
		for {
			if driver.GetFloorSensorSignal() != -1 {
				floor = driver.GetFloorSensorSignal()
			}
			fmt.Println(floor)
			status.CurFloor = floor
			status.Dir = serveDirection
			status_update_chan <- status
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			fmt.Println("\nSTATE : ", state)
			switch state {
			case WAIT:
				wait(status.OrderMatrix, &state, &serveDirection, &runDirection)	
				order.PrintMatrix(status.OrderMatrix)
			case RUN_UP:
				runUp(status.OrderMatrix, &state, &serveDirection)				
				order.PrintMatrix(status.OrderMatrix)
				fmt.Println("HIEHIEHEIEHIEHIE")
			case RUN_DOWN:
				runDown(status.OrderMatrix, &state, &serveDirection)
				order.PrintMatrix(status.OrderMatrix)
			case OPEN:
				open(runDirection)
				order.PrintMatrix(status.OrderMatrix)
			}
		}
	}()
}
func open(runDirection int) {
	elevatorBrake(runDirection)
}

func wait(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int, runDirection *int) {
	curFloor := driver.GetFloorSensorSignal()
	fmt.Printf("FLOOR = %d", curFloor)
	
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			if orderMatrix[i][0] == 1 { //Serve Order Up
				fmt.Println("UP ")
				if curFloor == i { //Order In Current Floor
					*state = OPEN
					break
				} else if curFloor < i { //Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = -1
					break
				} else { //Going Down
					*state = RUN_DOWN
					*serveDirection = 0
					*runDirection = 1
					break
				}
			} else if orderMatrix[i][1] == 1 { //Serve Order Down
				fmt.Println("DOWN ")				
				if curFloor == i {
					*state = OPEN //Order In Currrent Floor
					break
				} else if curFloor < i {
					*state = RUN_UP //Going Up
					*serveDirection = 1
					*runDirection = -1
					break
				} else { //Going Down
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = 1
					break
				}
			} else if orderMatrix[i][2] == 1 { //Serve Internal
				fmt.Println("INTERNAL ")
				if curFloor == i { //Order In Current Floor
					*state = OPEN
					break
				} else if curFloor < i { //Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = -1
					break
				} else { //Going Down
					fmt.Println("GOING DOWN ")
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = 1
					break
				}
			}
		}
	}
}

func runUp(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int) {
	driver.SetSpeed(300)
	fmt.Println("GOING UP")
	curFloor := driver.GetFloorSensorSignal()
	//HIT FLOOR
	if curFloor != -1 {
		//Serving orders UP
		if *serveDirection == 0 {
			if orderMatrix[curFloor][0] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
		//Serving orders DOWN
		if *serveDirection == 1 {
			if orderMatrix[curFloor][1] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
	}
}

func runDown(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int) {
	driver.SetSpeed(-300)
	fmt.Println("GOING DOWN ")
	curFloor := driver.GetFloorSensorSignal()
	//HIT FLOOR
	if curFloor != -1 {
		//Serving orders UP
		if *serveDirection == 0 {
			if orderMatrix[curFloor][0] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}

		}
		//Serving orders DOWN
		if *serveDirection == 1 {
			if orderMatrix[curFloor][1] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
	}
}

/*
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
*/
func elevatorBrake(dir int) {
		driver.SetSpeed(dir * 100)
		time.Sleep(time.Millisecond * 20)
		driver.SetSpeed(0)
}

func InitMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			matrix[x][y] = 0
		}
	}
}
