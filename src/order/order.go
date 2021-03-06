package order

import (
	"driver"
	"elevtypes"
	"fmt"
	"network"
	"strings"
	"time"
	"sort"
)

var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
var unassignedOrdersMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int

//Listen to local order buttons 
func LocalOrderListener(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	var ButtonMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	InitOrderMatrix(&ButtonMatrix)
	InitOrderMatrix(&orderMatrix)
	InitUnassignedOrdersMatrix(&unassignedOrdersMatrix)
	go func() {
		for {
			time.Sleep(15 * time.Millisecond)
			for i := 0; i < elevtypes.N_FLOORS; i++ {
				for j := 0; j < elevtypes.N_BUTTONS; j++ {
					ButtonMatrix[i][j] = driver.GetButtonSignal(j, i)
					if ButtonMatrix[1][1] == 1 {
					}
					if ButtonMatrix[i][j] == 1 {
						if j == 2 {
							orderMatrix[i][j] = 1
							orders_local_elevator_chan <- orderMatrix
						} else {
							unassignedOrdersMatrix[i][j] = 1
							orders_external_elevator_chan <- unassignedOrdersMatrix
						}
					}
				}
			}
		}
	}()
}

//Listen to orders from network and set lights
func OrdersFromNetworkListener(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, elevator_status_map_send_chan <-chan map[string]elevtypes.ElevatorStatus, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, orders_from_unresponsive_elev_chan <-chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	var elevatorStatusMap map[string]elevtypes.ElevatorStatus
	elevatorStatusMap = make(map[string]elevtypes.ElevatorStatus)
	var prevElevatorStatusMap map[string]elevtypes.ElevatorStatus
	prevElevatorStatusMap = make(map[string]elevtypes.ElevatorStatus)
	printCounter := 0
	
	//Listen to orders from network
	go func() {
		for {
			for key, v := range elevatorStatusMap {
				prevElevatorStatusMap[key] = v	
			}
			select {
			case elevatorStatusMap = <-elevator_status_map_send_chan:
				for key2, _ := range elevatorStatusMap {
					if !unassignedOrdersMatrixIsEmpty(elevatorStatusMap[key2].UnassignedOrdersMatrix) {
						costFunction(orders_local_elevator_chan, elevatorStatusMap)
						confirmOrderAssignment(elevatorStatusMap, orders_external_elevator_chan)
					}
				}
				for key3, _ := range elevatorStatusMap {
					if !orderMatricesEqual(elevatorStatusMap[key3].OrderMatrix, prevElevatorStatusMap[key3].OrderMatrix) || elevatorStatusMap[key3].CurFloor != prevElevatorStatusMap[key3].CurFloor {
						fmt.Printf("Printnumber: %d \n", printCounter)
						printStatusMap(elevatorStatusMap)
						printCounter = printCounter +1
					}
				}
			case newOrderMatrix := <-orders_from_unresponsive_elev_chan:
				time.Sleep(1 * time.Second)
				addOrdersToUnprocessedMatrix(newOrderMatrix)
				orders_external_elevator_chan <- unassignedOrdersMatrix
			}
		}
	}()

	//Set the Order button lights and internal lights, locally, as the Status map has been received from the network
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			setLights(elevatorStatusMap)
		}
	}()
}

func printStatusMap(elevatorStatusMap map[string]elevtypes.ElevatorStatus) {
	
	var keys []string
	for k := range elevatorStatusMap {
	    keys = append(keys, k)
	}
	sort.Strings(keys)





	for range elevatorStatusMap {
		fmt.Printf("----------------------------------------")
	}
	fmt.Printf("\n")
	for i, _ := range keys {
		fmt.Printf("IP: %s \t\t\t ", keys[i])
	}
	fmt.Printf("\n\n")
	for  range keys {
		fmt.Printf("Floor     UP   DOWN   INTERNAL  \t\t ")
	}
	
	fmt.Printf("\n\n")
	for i := elevtypes.N_FLOORS-1; i >= 0  ; i-- {
		for _, key2 := range keys {
			fmt.Printf(" %d \t %d \t %d \t %d   \t  ", i +1, elevatorStatusMap[key2].OrderMatrix[i][0], elevatorStatusMap[key2].OrderMatrix[i][1], elevatorStatusMap[key2].OrderMatrix[i][2])
			if elevatorStatusMap[key2].CurFloor == i {
				fmt.Printf("   [] ")
			} else{
				fmt.Printf("   -  ")
			}
			fmt.Printf("         ")
		}
		fmt.Printf("\n")
	}
	for _, key2 := range keys {
		if elevatorStatusMap[key2].ServeDirection == 0 {
			fmt.Printf("ServeDirection:  UP     \t\t\t")
		}else if elevatorStatusMap[key2].ServeDirection == 1 {
			fmt.Printf("ServeDirection: DOWN    \t\t\t")
		}else {
			fmt.Printf("ServeDirection: UNDEFINED    \t\t\t")
		}
	}
	fmt.Printf("\n")
	for _, key2 := range keys {
		fmt.Printf("WorkLoad:       %d \t\t\t\t", elevatorStatusMap[key2].WorkLoad)
	}
	fmt.Printf("\n\n\n")
}

//An external order arrives from network. Evaluate cost. Decide if this elevator should take it or not.  
func costFunction(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, elevatorStatusMap map[string]elevtypes.ElevatorStatus) {
	var orderFloor int
	var orderType int
	var penaltyMap map[string]int
	penaltyMap = make(map[string]int)
	var lowestPenalty int
	var lowestPenaltyIP string
	var equalPenaltyList []string
	for key, _ := range elevatorStatusMap {
		for x := 0; x < elevtypes.N_FLOORS; x++ {
			for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
				if elevatorStatusMap[key].UnassignedOrdersMatrix[x][y] == 1 {
					orderFloor = x
					orderType = y
				}
			}
		}
	}
	for key, _ := range elevatorStatusMap {
		penaltyMap[key] = AbsoluteValue(orderFloor - elevatorStatusMap[key].CurFloor + elevatorStatusMap[key].WorkLoad)
		if orderType != elevatorStatusMap[key].ServeDirection {
			penaltyMap[key] = penaltyMap[key] + 4 
		}
	}
	lowestPenalty = 100
	for key, _ := range penaltyMap {
		if penaltyMap[key] < lowestPenalty {
			lowestPenalty = penaltyMap[key]
			lowestPenaltyIP = key
		}
	}
	for key, _ := range penaltyMap {
		if penaltyMap[key] > lowestPenalty {
			delete(penaltyMap, key)
		}
	}
	
	if len(penaltyMap) > 1 {
		for key, _ := range penaltyMap {
			equalPenaltyList = append(equalPenaltyList, key)
		}
		if network.LocalIPIsLowest(equalPenaltyList) {
			orderMatrix[orderFloor][orderType] = 1
			orders_local_elevator_chan <- orderMatrix		
		}
	}else if strings.Split(lowestPenaltyIP, ":")[0] == network.GetIP() {
		orderMatrix[orderFloor][orderType] = 1
		orders_local_elevator_chan <- orderMatrix
	}
}

//Confirm that an elevator has taken external orders.
func confirmOrderAssignment(elevatorStatusMap map[string]elevtypes.ElevatorStatus, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	for key, _ := range elevatorStatusMap {
		for i := 0; i < elevtypes.N_FLOORS; i++ {
			for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
				if unassignedOrdersMatrix[i][j] == elevatorStatusMap[key].OrderMatrix[i][j] {
					unassignedOrdersMatrix[i][j] = 0
					orders_external_elevator_chan <- unassignedOrdersMatrix
				}
			}
		}
	}
}

//In case of emergency; run down to closest floor. New state: open. 
func ErrorRecovery() {
	for {
		prevOrderMatrix := orderMatrix
		time.Sleep(10 * time.Second)
		if !orderMatrixIsEmpty(orderMatrix) && orderMatricesEqual(orderMatrix, prevOrderMatrix) {
			driver.Init()
		}
	}
}

func addOrdersToUnprocessedMatrix(newOrderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			if newOrderMatrix[x][y] == 1 {
				unassignedOrdersMatrix[x][y] = 1
			}
		}
	}
}

func AbsoluteValue(value int) int {
	if value < 0 {
		return -value
	} else {
		return value
	}
}

func orderMatricesEqual(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, prevOrderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) bool {
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			if orderMatrix[i][j] != prevOrderMatrix[i][j] {
				return false
			}
		}
	}
	return true
}

func unassignedOrdersMatrixIsEmpty(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) bool {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			if matrix[x][y] == 1 {
				return false
			}
		}
	}
	return true
}

func orderMatrixIsEmpty(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) bool {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS; y++ {
			if matrix[x][y] == 1 {
				return false
			}
		}
	}
	return true
}

func DeleteOrder(floor int, orderType int, orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	orderMatrix[floor][orderType] = 0
	orders_local_elevator_chan <- orderMatrix
}


func PrintMatrix(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	fmt.Println("\nFloor \t UP \t DOWN \t INTERNAL")
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		fmt.Printf("\n %d \t", i+1)
		for j := 0; j < elevtypes.N_BUTTONS; j++ {
			fmt.Printf("%d \t ", matrix[i][j])
		}
	}
}

func PrintUnassignedOrdersMatrix(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	fmt.Println("\nFloor \t UP \t DOWN \t INTERNAL")
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		fmt.Printf("\n %d \t", i+1)
		for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
			fmt.Printf("%d \t ", matrix[i][j])
		}
	}
}

func setLights(elevatorStatusMap map[string]elevtypes.ElevatorStatus) {
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS -1; j++ {
			lightOn := 0
			for key, _ := range elevatorStatusMap {
				if elevatorStatusMap[key].OrderMatrix[i][j] == 1 {
					lightOn = 1
				} 		
			}
			driver.SetButtonLamp(j, i, lightOn)
		}
	}
	
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		if orderMatrix[i][2] == 1 {
			driver.SetButtonLamp(2, i, 1)
		} else{
			driver.SetButtonLamp(2, i, 0)
		}
	}
}
func InitOrderMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS; y++ {
			matrix[x][y] = 0
		}
	}
}

func InitUnassignedOrdersMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			matrix[x][y] = 0
		}
	}
}
