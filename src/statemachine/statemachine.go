package statemachine

import(
	"driver"
	"time"
	"fmt"
	"elevtypes"
	"order"
	
)

type State_enum int
  
const (
	WAIT State_enum = iota
	RUN 
	OPEN
)

	
func ElevatorInit() int {
	init := driver.Init()
	if init == 0 {
		return 0
	} else {
		if driver.GetFloorSensorSignal() != -1 {
		} else {
			driver.SetSpeed(-1*300)
			floor := driver.GetFloorSensorSignal()
			for floor == -1 {
				floor = driver.GetFloorSensorSignal()
			}
			ElevatorBrake(1)
		}
		fmt.Printf("Initialized\n")
		return 1
	}
}

func StateMachine(orders_local_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int,orders_external_elevator_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS-1]int, 	status_update_chan chan<- elevtypes.Status) {
	var status elevtypes.Status
	var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int	
	var unprocessedMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS-1]int	
	var state State_enum
	var floor int
	var dir int
	state = WAIT
	
	order.InitMatrix(&orderMatrix)
	InitMatrix(&unprocessedMatrix)
	
	go func () {
		for {
			select {
			case orderMatrix := <-orders_local_elevator_chan:
				fmt.Println("oppdaterer Ordermatrix...")
				status.OrderMatrix = orderMatrix
			case unprocessedMatrix := <-orders_external_elevator_chan:
				fmt.Println("Oppdaterer unprocessedmatrix..")  
				status.UnprocessedMatrix = unprocessedMatrix
			}
		}
	}()

	go func () {
		for {
			
			status.CurFloor	= floor		
			status.Dir = dir
			status_update_chan <- status	
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func () {
		for {
			time.Sleep(10 * time.Millisecond)
			switch state {
			case WAIT:
				wait(orderMatrix, &state)
			/*case RUN:
				run(orderList, &state)
			case OPEN:
				open()*/
			}
		}
	}()
}

func wait(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, state *State_enum) {

	for i := 0; i < elevtypes.N_FLOORS; i++ { 
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			if orderMatrix[i][j] == 1 {
				time.Sleep(1 * time.Second)
				fmt.Println("Matrix not empty...")
			}
		}
	}		
}

/*func run(orderList []Order, state *elevtypes.State_enum) {
			
			
	
}*/











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
func ElevatorBrake(dir int) {
	driver.SetSpeed(dir*300)
	time.Sleep(time.Millisecond*20)
	driver.SetSpeed(0)
}

func InitMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS-1]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ { 
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			matrix[x][y] = 0
		}
	}
}

