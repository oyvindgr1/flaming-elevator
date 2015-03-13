func Send_status() {
	baddr, err_conv_ip := net.ResolveUDPAddr("udp", "129.241.187.255:20020")
	Check_error(err_conv_ip)
	status_sender, err_dialudp := net.DialUDP("udp", nil, baddr)
	Check_error(err_dialudp)
	for {
			status_sender.Write("IsAlive")
		}
}

func Read_status(lost_orders_c chan driver.Client, all_ips map[string]time.Time, all_clients map[string]driver.Client, localIP net.IP) {
	laddr, err_conv_ip_listen := net.ResolveUDPAddr("udp", ":20020")
	Check_error(err_conv_ip_listen)
	status_receiver, err_listen := net.ListenUDP("udp", laddr)
	Check_error(err_listen)
	for {
		time.Sleep(25 * time.Millisecond)
		b := make([]byte, 1024)
		n, raddr, _ := status_receiver.ReadFromUDP(b)
		err_decoding := json.Unmarshal(b[n:n]
		if err_decoding != nil {
			fmt.Println("error decoding client msg")
		}
		fmt.println("Got message from %s", raddr)
	}
}