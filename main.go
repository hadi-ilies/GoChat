package main

import (
	"os"

	"./client"
	"./server"
)

func usage(exitValue int) {
	println("USAGE")
	println("\tgo run main.go [options]")
	println("DESCRIPTION")
	println("\tThis Program is an update of this project:\"https://github.com/hadi-ilies/SimpleGolangChatServer\" that works following a room concept, for instance the server create/manage rooms and users connect to it.")
	println("\tEach users can connect to a room or join one if there is an available spot. The server checks if there is an available spot in active rooms, if not, it will create a new one.")
	println("OPTIONS")
	println("\t-S, --server <[ipaddr], port>")
	println("\t\tstart server using the ip and port pass in argv")
	println("\t-C, --client <ipaddr, port>")
	println("\t\tstart client on the ip and port pass in argv")
	println("\t-h, --help")
	println("\t\tDisplay the program usage")
	os.Exit(exitValue)
}

func main() {
	if len(os.Args) < 3 {
		usage(1)
	} else {
		if os.Args[1] == "--server" || os.Args[1] == "-S" {
			println("start server here")
			if len(os.Args) == 3 {
				//create server. The IPV4 will be automaticly chosen whether we have internet or not
				server := server.NewServer("tcp")
				//start server on the port chosen in argv
				server.Start(os.Args[2])
			} else {
				//create server with a specific IPV4
				server := server.NewServer("tcp", os.Args[2])
				//start server on the port chosen in argv
				server.Start(os.Args[3])
			}
		}
		if (os.Args[1] == "--client" && len(os.Args) == 4) || (os.Args[1] == "-C" && len(os.Args) == 4) {
			println("Start client here")
			client := client.NewClient(os.Args[2], os.Args[3])
			client.StartClient()
		}
	}
}
