package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/hadi-ilies/ProtoBuffExample/protobuff"
	"google.golang.org/protobuf/proto"
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

//sendText: get and send the text from the stdin
func (client *Client) sendText(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		//buffer that store server's response
		text := getLine(scanner)
		if len(text) != 0 {
			client.packet.Msg = text
			data, err := proto.Marshal(client.packet)
			if err != nil {
				log.Fatalln("Failed to serialize post in protobuf:", err)
			}
			//send proto struct serialized!
			conn.Write(data)
		} else {
			break
		}
	}
}
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
