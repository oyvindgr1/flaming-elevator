package order

import(
	"elevtypes"
	"driver"
	"time"
	
)

/*type Order_struct struct{
	IP string
	Orders map[string][N_FLOORS][N_BUTTON]bool
	
} */

func OrderListener(orders_local_elevator chan Order, ) {
	var newOrder Order
	var ButtonMatrix [N_FLOORS][N_BUTTONS]bool


	go func() {
		for {
			time.Sleep(15 * time.Millisecond)
			for i := 0; i < N_FLOORS; i++ { 
				for j := 0; j < N_BUTTONS; j++ {
					ButtonMatrix[i][j] = GetButtonSignal(j, i)
					if ButtonMatrix[i][j] == 1 {
						newOrder = Order{i, j}
						if j == 2 {
							orders_local_elevator <- newOrder
						} else {
							orders_local_elevator <- newOrder
						}	
					}
				}
			}
		}
	}
}

//func Init() {


	//orderList := make(map[string][N_FLOORS][N_BUTTON]bool)

	//Order_struct =
	//orderChannel = make(chan ordersToExecute) 
//}

func UnprocessedOrderListGenerator( OrderMatrix [][]bool) {

	for {
		select {
			case internalOrder <- orderInternal  
	
	

}
