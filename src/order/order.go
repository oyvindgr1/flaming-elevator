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
			PrintMatrix(orderMatrix)
		}
	}()
}

func OrdersFromNetwork(statusmap_send_chan <-chan map[string]elevtypes.Status) {
	for {
		select {
		case statusMap := <-statusmap_send_chan:
			for key, _ := range statusMap {
				if !MatrixIsEmpty(statusMap[key].UnprocessedOrdersMatrix) {
					CostFunction(statusMap)
					checkUnprocessedMatrix(statusMap)
				}
			}
		}

	}
}

func checkUnprocessedMatrix(statusMap map[string]elevtypes.Status) {
	for key, _ := range statusMap {
		for i := 0; i < elevtypes.N_FLOORS; i++ {
			for j := 0; j < elevtypes.N_BUTTONS-1; j++ {
				if unprocessedOrdersMatrix[i][j] == statusMap[key].UnprocessedOrdersMatrix[i][j] {
					unprocessedOrdersMatrix[i][j] = 0
				}
			}
		}
	}
}

func CostFunction(statusMap map[string]elevtypes.Status) {
	fmt.Println("In Costfunction.")
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
		if orderType != statusMap[key].Dir {
			penaltyMap[key] = penaltyMap[key] + 4
		}
	}
	lowestPenalty = 100
	for key, _ := range penaltyMap {
		if penaltyMap[key] < lowestPenalty {
			lowestPenalty = penaltyMap[key]
			lowestPenaltyIP = key
			fmt.Println("Penalty: ", lowestPenalty)
		}
	}
	if strings.Split(lowestPenaltyIP, ":")[0] == network.GetIP() {
		orderMatrix[orderFloor][orderType] = 1
	}

	//PrintMatrix(orderMatrix)
}

func AbsoluteValue(value int) int {
	if value < 0 {
		return -value
	} else {
		return value
	}
}

func MatrixIsEmpty(matrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) bool {
	for x := 0; x < elevtypes.N_FLOORS; x++ {
		for y := 0; y < elevtypes.N_BUTTONS-1; y++ {
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

/*
func OrderAppend(orderList_chan <-chan []Order, orders_local_elev_chan chan elevtypes.Order) {
	for {
		select {
		case newOrder := <-orders_local_elev_chan:
			isInList = false
			if !IsInOrderList(orderList, newOrder) {
				orderList = append(orderList, newOrder)
				orderList_chan <- orderList
				fmt.Printf("Floor of new order: %d, Type of new order: %d\n", newOrder.Floor, newOrder.Dir)
			}
		}
	}
}()
//ORDER_UP = 0, ORDER_DOWN = 1, ORDER_INTERNAL = 2
func OrderDelete(orderType int, order Order) {



}*/

/*

	go func() {
		for {
			select {
			case resetAll := <- reset_all:
				for i := 0; i < N_FLOORS; i++ {
						list[ORDER_DOWN][i] = false
						list[ORDER_UP][i] = false
					}
			case resetFloor := <- reset_button:
				ButtonMatrix[resetFloor.Dir][resetFloor.Floor] = false
			}
		}
	}()




func Init() {


	orderList := make(map[string][N_FLOORS][N_BUTTON]bool)

	Order_struct =
	orderChannel = make(chan ordersToExecute)
}

func UnprocessedOrderListGenerator( OrderMatrix [][]bool) {

	for {
		select {
			case internalOrder <- orderInternal



}

func IsInOrderList(orderList []Order, order elevtypes.Order) bool {
	for i,_ := range orderList {
		if order == orderList[i] {
			return true
		}
	}
}*/
