package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/hadi-ilies/GoChat/protobuff"
)

//Client is a user that will connect to my chatserver
type Client struct {
	ipAddr string
	port   string
	packet *protobuff.ClientBDD
	conn   net.Conn
}

//NewClient is the Client's constructor
func NewClient(ipAddr string, port string) Client {
	return Client{ipAddr: ipAddr, port: port, packet: &protobuff.ClientBDD{}}
}

//connect: connect to server //param "tcp"
func (client *Client) connect(connType string) net.Conn {
	conn, err := net.Dial(connType, client.ipAddr+":"+client.port)
	if err != nil {
		log.Fatal("Error during the connection :", err)
	}
	return conn
}

//register: with this function the client has to insert a name in order to connect to our server
func (client *Client) register() {
	scanner := bufio.NewScanner(os.Stdin)
	print("Register your name : ")
	for {
		text := getLine(scanner)
		if len(text) != 0 {
			client.packet.Name = text
			client.conn = client.connect("tcp")
			return
		}
	}
}

//StartClient : Start a simple client that will connect to the server and send text
func (client *Client) StartClient() {
	client.register()
	client.startGUI()
}
