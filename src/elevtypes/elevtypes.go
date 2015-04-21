package elevtypes

import ()

type Status struct {
	CurFloor          int
	ServeDirection               int
	OrderMatrix       [N_FLOORS][N_BUTTONS]int     //This elevator's orders to execute
	UnprocessedOrdersMatrix [N_FLOORS][N_BUTTONS - 1]int //This elevator's orders, not yet assigned
	//state State_Enum
}

const N_BUTTONS = 3
const N_FLOORS = 4


