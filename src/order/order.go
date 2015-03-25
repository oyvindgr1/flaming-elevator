package order

import(
	"elevtypes"
	"driver"
	"time"
	"fmt"
	
)

/*type Order_struct struct{
	IP string
	Orders map[string][N_FLOORS][N_BUTTON]bool
	
} */
//, reset_button chan Order, reset_all chan bool 
func OrderListener(orders_local_elevator_chan chan<- elevtypes.Order) {
	var newOrder elevtypes.Order
	var ButtonMatrix [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int
	
	for x := 0; x < elevtypes.N_FLOORS; x++ { 
				for y := 0; y < elevtypes.N_BUTTONS; y++ {
					ButtonMatrix[x][y] = 0
				}
	}
	
	go func() {
		for {
			time.Sleep(15 * time.Millisecond)
			for i := 0; i < elevtypes.N_FLOORS; i++ { 
				for j := 0; j < elevtypes.N_BUTTONS; j++ {
					ButtonMatrix[i][j] = driver.GetButtonSignal(j, i)
					if ButtonMatrix[i][j] == 1 {
						newOrder = elevtypes.Order{i, j}
						if j == 2 {
							fmt.Println("New order internal")
							orders_local_elevator_chan <- newOrder
						} else {
							fmt.Println("New order external")
							orders_local_elevator_chan <- newOrder
						}	
					}
				}
			}
		}
	}()
}
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
	
	

}*/
