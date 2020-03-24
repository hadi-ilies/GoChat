package server

import (
	"net"
)

//getOutboundIP: get my ip addr
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	//if there is no internet we get the localhost
	if err != nil {
		//log.Fatal(err)
		println("No network !")
		return net.IPv4(127, 0, 0, 1)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
