package main

import (
	"elevtypes"
	"network"
	"order"
	"statemachine"
	"driver"
)

func main() {
	
	statusmap_chan := make(chan map[string]elevtypes.Status, 1)
	status_chan := make(chan elevtypes.Status, 1)
	
	orders_from_unresponsive_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_local_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS]int, 1)
	orders_external_elev_chan := make(chan [elevtypes.N_FLOORS][elevtypes.N_BUTTONS - 1]int, 1)

	netIsAlive := make(chan bool)
	driver.Init()	
	
	go order.OrderListener(orders_local_elev_chan, orders_external_elev_chan)
	go order.OrdersFromNetwork(orders_local_elev_chan, statusmap_chan,orders_external_elev_chan, orders_from_unresponsive_elev_chan)
	go order.ErrorRecovery()
	
	go network.SendStatus(status_chan)
	go network.ReadStatus(statusmap_chan, netIsAlive, orders_from_unresponsive_elev_chan)


	go statemachine.StateMachine(orders_local_elev_chan, orders_external_elev_chan, status_chan)

	/*go func() {
		for {
			select{
			case aBool := <- netIsAlive:
				fmt.Printf(" Net alive: %t", aBool)
			}
		}
	}()*/

	select {}


}
