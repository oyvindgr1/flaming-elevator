package order

import (
	"driver"
	"elevtypes"
	"fmt"
	"time"
)

/*type Order_struct struct{
	IP string
	Orders map[string][N_FLOORS][N_BUTTON]bool

} */
//, reset_button chan Order, reset_all chan bool
var orderMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
var unprocessedMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int

func OrderListener(orders_local_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, orders_external_elevator_chan chan<- [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int) {
	var ButtonMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	InitMatrix(&ButtonMatrix)
	go func() {
		for {
			time.Sleep(15 * time.Millisecond)
			for i := 0; i < elevtypes.N_FLOORS; i++ {
				for j := 0; j < elevtypes.N_BUTTONS; j++ {
					ButtonMatrix[i][j] = driver.GetButtonSignal(j, i)
					if ButtonMatrix[i][j] == 1 {
						if j == 2 {
							fmt.Println("New order internal")
							orderMatrix[i][j] = 1
							orders_local_elevator_chan <- orderMatrix
						} else {
							fmt.Println("New order external")
							unprocessedMatrix[i][j] = 1
							orders_external_elevator_chan <- unprocessedMatrix

						}
					}
				}
			}
		}
	}()
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
