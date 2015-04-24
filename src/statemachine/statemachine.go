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

func StateMachine(orders_local_elevator_chan chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, status_update_chan chan<- elevtypes.Status) {
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
			status.CurFloor = floor
			status.ServeDirection = serveDirection
			status_update_chan <- status
			time.Sleep(50 * time.Millisecond)
			//fmt.Println("\nservedir: ",serveDirection)
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			//fmt.Println("\nSTATE : ", state)
			switch state {
			case WAIT:
				wait(status.OrderMatrix, &state, &serveDirection, &runDirection)
				//order.PrintMatrix(status.OrderMatrix)
			case RUN_UP:
				runUp(status.OrderMatrix, &state, &serveDirection)
				//order.PrintMatrix(status.OrderMatrix)
			case RUN_DOWN:
				runDown(status.OrderMatrix, &state, &serveDirection)
				//order.PrintMatrix(status.OrderMatrix)
			case OPEN:
				open(orders_local_elevator_chan, status.OrderMatrix, &state, runDirection, &serveDirection)
				//order.PrintMatrix(status.OrderMatrix)
			}
		}
	}()
}
func open(orders_local_elevator_chan chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, runDirection int, serveDirection *int) {
	elevatorBrake(runDirection)
	curFloor := driver.GetFloorSensorSignal()
	driver.SetDoorOpenLamp(1)
	time.Sleep(3 * time.Second)
	driver.SetDoorOpenLamp(0)
	order.DeleteOrder(curFloor, 2, orders_local_elevator_chan)
	if *serveDirection == 0 {
		order.DeleteOrder(curFloor, 0, orders_local_elevator_chan)
		for i := curFloor+1; i < elevtypes.N_FLOORS; i++ {
			if orderMatrix[i][0] == 1 || orderMatrix[i][2] == 1 {
				*state = RUN_UP
				return
			}
		}
	} else if *serveDirection == 1 {
		order.DeleteOrder(curFloor, 1, orders_local_elevator_chan)
		for i := 0; i < curFloor-1; i++ {
			if orderMatrix[i][1] == 1 || orderMatrix[i][2] == 1 {
				*state = RUN_DOWN
				return
			}
		}
	}
	*state = WAIT
}

func wait(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int, runDirection *int) {
	*serveDirection = -1	
	curFloor := driver.GetFloorSensorSignal()
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			if orderMatrix[i][0] == 1 { //Serve Order Up
				if curFloor == i { //Order In Current Floor
					*state = OPEN
					*serveDirection = 0
					return
				} else if curFloor < i { //Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = -1
					return
				} else { //Going Down
					*state = RUN_DOWN
					*serveDirection = 0
					*runDirection = 1
					return
				}
			} else if orderMatrix[i][1] == 1 { //Serve Order Down
				if curFloor == i {
					*state = OPEN //Order In Currrent Floor
					*serveDirection = 1
					return
				} else if curFloor < i {
					*state = RUN_UP //Going Up
					*serveDirection = 1
					*runDirection = -1
					return
				} else { //Going Down
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = 1
					return
				}
			} else if orderMatrix[i][2] == 1 { //Serve Internal
				if curFloor == i { //Order In Current Floor
					*state = OPEN
					*serveDirection = -1
					return
				} else if curFloor < i { //Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = -1
					return
				} else { //Going Down
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = 1
					return
				}
			}
		}
	}
}

func runUp(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int) {
	driver.SetSpeed(300)
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
