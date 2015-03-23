package elevtypes

import (
	"net"

)



type State struct {
	IP net.IP
	CurFloor string
	HeisNr string 	 
	Astring string
	Floor int
	Dir int
	UnprocessedOrders []Order
	OrdersToExecute []Order  
}

const N_BUTTONS = 3
const N_FLOORS = 4

//ORDER_UP = 0, ORDER_DOWN = 1, ORDER_INTERNAL = 2
type Order struct {
	Floor int
	Dir int
}

var UnprocessedOrders []Order

var lampChannelMatrix= [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannelMatrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}	