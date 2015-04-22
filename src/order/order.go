package order

import (
	"driver"
	"elevtypes"
	"fmt"
	"network"	
	"strings"
	"time"
)

/*type Order_struct struct{
	IP string
	Orders map[string][N_FLOORS][N_BUTTON]bool

} */
//, reset_button chan Order, reset_all chan bool
var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
var unprocessedOrdersMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int

func OrderListener(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	var ButtonMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	InitMatrix(&ButtonMatrix)
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
							unprocessedOrdersMatrix[i][j] = 1
							orders_external_elevator_chan <- unprocessedOrdersMatrix
						}
					}
				}
			}
			//fmt.Println("\norderMatrix local order module: ")
			//PrintMatrix(orderMatrix)
		}
	}()
}

func OrdersFromNetwork(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, statusmap_send_chan <-chan map[string]elevtypes.Status, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	for {
		select {
		case statusMap := <-statusmap_send_chan:
			for key, _ := range statusMap {
				if !MatrixIsEmpty(statusMap[key].UnprocessedOrdersMatrix) {
					CostFunction(orders_local_elevator_chan, statusMap)
					CheckUnprocessedMatrix(statusMap, orders_external_elevator_chan)
				}
			}
		}

	}
}

func CheckUnprocessedMatrix(statusMap map[string]elevtypes.Status, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	for key, _ := range statusMap {
		for i := 0; i < elevtypes.N_FLOORS; i++ {
			for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
				if unprocessedOrdersMatrix[i][j] == statusMap[key].UnprocessedOrdersMatrix[i][j] {
					unprocessedOrdersMatrix[i][j] = 0
					orders_external_elevator_chan <- unprocessedOrdersMatrix
				}
			}
		}
	}
}

func errorRecovery() {
	for {	
		prevOrderMatrix := orderMatrix
		time.Sleep(5 * time.Second)
		if !orderMatrixIsEmpty(orderMatrix) && orderMatricesEqual(orderMatrix, prevOrderMatrix) {
			//disconnect 2 secs from network
			//phoenix
		}
	}
}

orderMatricesEqual(orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, prevOrderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) bool{
	for i := 0; i < elevtypes.N_FLOORS; i++ {
		for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
			if orderMatrix[i][j] != prevOrderMatrix[i][j] {
				return false
			}
		}
	}
	return true
}
	
func CostFunction(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int,statusMap map[string]elevtypes.Status) {
	//fmt.Println("In Costfunction.")
	var orderFloor int
	var orderType int
	var penaltyMap map[string]int
	penaltyMap = make(map[string]int)
	var lowestPenalty int
	var lowestPenaltyIP string
	for key, _ := range statusMap {
		for x := 0; x < elevtypes.N_FLOORS; x++ {
			for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
				if statusMap[key].UnprocessedOrdersMatrix[x][y] == 1 {
					orderFloor = x
					orderType = y
				}
			}
		}
	}
	for key, _ := range statusMap {
		penaltyMap[key] = AbsoluteValue(orderFloor - statusMap[key].CurFloor)
		if orderType != statusMap[key].ServeDirection {
			penaltyMap[key] = penaltyMap[key] + 4
		}
	}
	lowestPenalty = 100
	for key, _ := range penaltyMap {
		if penaltyMap[key] < lowestPenalty {
			lowestPenalty = penaltyMap[key]
			lowestPenaltyIP = key
			//fmt.Println("Penalty: ", lowestPenalty)
		}
	}
	if strings.Split(lowestPenaltyIP, ":")[0] == network.GetIP() {
		orderMatrix[orderFloor][orderType] = 1
		orders_local_elevator_chan <- orderMatrix
	}
}

func AbsoluteValue(value int) int {
	if value < 0 {
		return -value
	} else {
		return value
	}
}

func UnprocessedOrdersMatrixIsEmpty(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) bool {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
			if matrix[x][y] == 1 {
				return false
			}
		}
	}
	return true
}
func OrderMatrixIsEmpty(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) bool {
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

func InitMatrix(matrix *[elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int) {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS; y++ {
			matrix[x][y] = 0
		}
	}
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


