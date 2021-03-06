package statemachine

import (
	"driver"
	"elevtypes"
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

//Controls the local elevator. Updates elevator status and sends it to network module. 
func LocalElevatorController(orders_local_elevator_chan chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, elevator_status_update_chan chan<- elevtypes.ElevatorStatus) {
	var elevatorStatus elevtypes.ElevatorStatus
	var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	var unassignedOrdersMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int
	var state State_enum
	var floor int
	var serveDirection int
	var runDirection int
	state = WAIT
	go order.ErrorRecovery()
	order.InitOrderMatrix(&orderMatrix)
	order.InitUnassignedOrdersMatrix(&unassignedOrdersMatrix)

	//Receive Orders and update elevatorStatus
	go func() {
		for {
			select {
			case orderMatrix := <-orders_local_elevator_chan:
				elevatorStatus.OrderMatrix = orderMatrix
				elevatorStatus.WorkLoad = sumOfOrders(orderMatrix)
			case unassignedOrdersMatrix := <-orders_external_elevator_chan:
				elevatorStatus.UnassignedOrdersMatrix = unassignedOrdersMatrix
			}
		}
	}()
	
	//Update elevator status
	go func() {
		for {
			if driver.GetFloorSensorSignal() != -1 {
				floor = driver.GetFloorSensorSignal()
			}
			driver.SetLightFloorIndicator(floor)
			elevatorStatus.CurFloor = floor
			elevatorStatus.ServeDirection = serveDirection
			elevator_status_update_chan <- elevatorStatus
			time.Sleep(100 * time.Millisecond)
		}
	}()

	//STATEMACHINE
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			switch state {
			case WAIT:
				wait(elevatorStatus.OrderMatrix, &state, &serveDirection, &runDirection)
			case RUN_UP:
				runDirection = 1
				runUp(elevatorStatus.OrderMatrix, &state, &serveDirection)
			case RUN_DOWN:
				runDirection = -1
				runDown(elevatorStatus.OrderMatrix, &state, &serveDirection)
			case OPEN:
				open(orders_local_elevator_chan, elevatorStatus.OrderMatrix, &state, runDirection, &serveDirection)
			}
		}
	}()
}

func open(orders_local_elevator_chan chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, runDirection int, serveDirection *int) {
	elevatorBrake(runDirection)
	time.Sleep(500 * time.Millisecond)
	curFloor := driver.GetFloorSensorSignal()
	if curFloor != -1 {
		driver.SetDoorOpenLamp(1)
		time.Sleep(2 * time.Second)
		driver.SetDoorOpenLamp(0)
		order.DeleteOrder(curFloor, 2, orders_local_elevator_chan)
		if *serveDirection == 0 {
			order.DeleteOrder(curFloor, 0, orders_local_elevator_chan)
			for i := curFloor + 1; i < elevtypes.N_FLOORS; i++ {
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
}

func wait(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int, runDirection *int) {
	*serveDirection = -1
	curFloor := driver.GetFloorSensorSignal()
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			if orderMatrix[i][0] == 1 { 				//Serve Order Up
				if curFloor == i { 				//Order In Current Floor
					*state = OPEN
					*serveDirection = 0
					return
				} else if curFloor < i { 			//Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = 1
					return
				} else { 					//Going Down
					*state = RUN_DOWN
					*serveDirection = 0
					*runDirection = -1
					return
				}
			} else if orderMatrix[i][1] == 1 { 			//Serve Order Down
				if curFloor == i {
					*state = OPEN 				//Order In Currrent Floor
					*serveDirection = 1
					return
				} else if curFloor < i {
					*state = RUN_UP 			//Going Up
					*serveDirection = 1
					*runDirection = 1
					return
				} else { 					//Going Down
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = -1
					return
				}
			} else if orderMatrix[i][2] == 1 { 			//Serve Internal
				if curFloor == i { 				//Order In Current Floor
					*state = OPEN
					*serveDirection = -1
					return
				} else if curFloor < i { 			//Going Up
					*state = RUN_UP
					*serveDirection = 0
					*runDirection = 1
					return
				} else { 					//Going Down
					*state = RUN_DOWN
					*serveDirection = 1
					*runDirection = -1
					return
				}
			}
		}
	}
}

func runUp(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int) {
	driver.SetSpeed(300)
	curFloor := driver.GetFloorSensorSignal()
	curFloorIsLastRemainingOrder := true
	if curFloor != -1 {							//HIT FLOOR
		if *serveDirection == 0 {					//Serving orders UP
			if orderMatrix[curFloor][0] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
		if *serveDirection == 1 {					//Serving orders DOWN
			for i := curFloor + 1; i < elevtypes.N_FLOORS; i++ {
				if orderMatrix[i][1] == 1 {
					curFloorIsLastRemainingOrder = false
				}
			}
			if orderMatrix[curFloor][1] == 1 && curFloorIsLastRemainingOrder || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
	}
}

func runDown(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum, serveDirection *int) {
	driver.SetSpeed(-300)
	curFloor := driver.GetFloorSensorSignal()
	curFloorIsLastRemainingOrder := true
	if curFloor != -1 {							//HIT FLOOR
		if *serveDirection == 0 {					//Serving orders UP
			for i := 0; i < curFloor; i++ {
				if orderMatrix[i][0] == 1 {
					curFloorIsLastRemainingOrder = false
				}
			}
			if orderMatrix[curFloor][0] == 1 && curFloorIsLastRemainingOrder || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
		if *serveDirection == 1 {					//Serving orders DOWN
			if orderMatrix[curFloor][1] == 1 || orderMatrix[curFloor][2] == 1 {
				*state = OPEN
			}
		}
	}
}

func elevatorBrake(dir int) {
	driver.SetSpeed(-dir * 100)
	time.Sleep(time.Millisecond * 20)
	driver.SetSpeed(0)
}

func sumOfOrders(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) int {
	counter := 0
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS; y++ {
			if orderMatrix[x][y] == 1 {
				counter = counter + 1
			}
		}
	}
	return counter
}

