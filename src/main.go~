package main

import (
	"elevtypes"
	"network"
	"order"
	"statemachine"
	"driver"
)

func main() {
	
	elevator_status_map_chan := make(chan map[string]elevtypes.ElevatorStatus, 1)
	elevator_status_chan := make(chan elevtypes.ElevatorStatus, 1)
	orders_from_unresponsive_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_local_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_external_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, 1)

	driver.Init()	
	
	go order.LocalOrderListener(orders_local_elev_chan, orders_external_elev_chan)
	go order.OrdersFromNetworkListener(orders_local_elev_chan, elevator_status_map_chan,orders_external_elev_chan, orders_from_unresponsive_elev_chan)
	go network.SendElevatorStatus(elevator_status_chan)
	go network.ReadElevatorStatus(elevator_status_map_chan, orders_from_unresponsive_elev_chan)
	go statemachine.LocalElevatorController(orders_local_elev_chan, orders_external_elev_chan, elevator_status_chan)

	select {}


}
