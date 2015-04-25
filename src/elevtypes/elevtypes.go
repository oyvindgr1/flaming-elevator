package elevtypes

import ()

type ElevatorStatus struct {
	CurFloor     	     	int
	ServeDirection       	int
	WorkLoad             	int		
	OrderMatrix          	[N_FLOORS][N_BUTTONS]int     		//Local elevator's orders to execute
	UnassignedOrdersMatrix [N_FLOORS][N_BUTTONS - 1]int 		//Local elevator's orders, not yet assigned
}

const N_BUTTONS = 3
const N_FLOORS = 4

