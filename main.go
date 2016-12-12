package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sensorelay/requests"

	"github.com/hashicorp/mdns"
)

func main() {
	// Setup our service export
	host, _ := os.Hostname()
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	info := []string{"My awesome service"}
	service, _ := mdns.NewMDNSService(host, "_sensorelay._tcp", "", "", port, getLocalIPS(), info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	http.HandleFunc("/reading", requests.Reading)
	log.Println("Listening on port ", port)
	log.Fatal(http.Serve(listener, nil))
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIPS() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var ips []net.IP
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.String() != "127.0.1.1" {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips
}
