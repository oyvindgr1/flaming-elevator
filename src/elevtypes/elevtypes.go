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
type State_enum int
  
const (
	Running State_enum = iota
	Idle 
	Door
	Undefined
)
	
type Status struct {
	CurFloor int
	Dir int
	OrderList []Order//This elevator's orders to execute
	UnprocessedOrders []Order//This elevator's orders, not yet assigned
	//state State_Enum
}

const N_BUTTONS = 3
const N_FLOORS = 4

//ORDER_UP = 0, ORDER_DOWN = 1, ORDER_INTERNAL = 2
type Order struct {
	Floor int
	Dir int
}


//var UnprocessedOrders []Order


