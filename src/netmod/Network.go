


func Send_status(state1 State) {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	Check_error(err_conv_ip)
	status_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
		time.Sleep(1000 * time.Millisecond)
		b, err_Json := json.Marshal(state1)
        if err_Json != nil {
	        fmt.Println("error with JSON")
	        fmt.Println(err1)
        }
		n, err1 := status_sender.Write(b)
        if err1 != nil {
                fmt.Println("error writing data to server")
                fmt.Println(err)
                return
        }
	}
}


func Read_status(Client_map map[State]int) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	status_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		time.Sleep(1000 * time.Millisecond)
		b := make([]byte, 1024)
		n, raddr, _ := status_receiver.ReadFromUDP(b)
		err_decoding := json.Unmarshal(b[0:n]
		if err_decoding != nil {
			fmt.Println("error decoding client msg")
		}
		Client_map[b] = raddr
		
		for key := range m {
		    fmt.Println(key.IP.String())
		}
	}
}