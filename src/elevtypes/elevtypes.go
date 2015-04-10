package elevtypes

import (

)


/*const (
	NewOrderType MessageType  = iota
	DeleteOrderType
	CostType
	AddOrderType
)

type Message struct {
	Message MessageType 
	OrderInfo Order
	Cost int
}*/

	
type Status struct {
	CurFloor int
	Dir int
	OrderMatrix [N_FLOORS][N_BUTTONS]int//This elevator's orders to execute
	UnprocessedMatrix [N_FLOORS][N_BUTTONS-1]int//This elevator's orders, not yet assigned
	//state State_Enum
}

const N_BUTTONS = 3
const N_FLOORS = 4

//OrderType: ORDER_UP = 0, ORDER_DOWN = 1, ORDER_INTERNAL = 2
type Order struct {
	Floor int
	OrderType int
}


//var UnprocessedOrders []Order


