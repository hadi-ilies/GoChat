package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hadi-ilies/ProtoBuffExample/protobuff"
	"google.golang.org/protobuf/proto"
)

//Server is a struct that handle and manage all the chatrooms
type Server struct {
	port           string
	ipAddr         string
	connectionType string
}

//NewServer is the server's constructor
func NewServer(connectionType string, ipAddr ...string) Server {
	//get my IP addr
	ip := getOutboundIP().String()
	if len(ipAddr) > 0 {
		ip = ipAddr[0]
	}
	println("IP equal: ", ip)
	return Server{ipAddr: ip, connectionType: connectionType}
}

func (server *Server) createListenner(port string) net.Listener {
	//create listenner.
	l, err := net.Listen(server.connectionType, server.ipAddr+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//save the current port
	server.port = port
	return l
}

//Start allow to start/run a server in tcp on the localhost:8080
func (server *Server) Start(port string) {
	//create listenner
	l := server.createListenner(port)
	//close server when server stop
	defer l.Close()
	//define nb max chatter per rooms
	const nbMaxChatterInRoom = uint8(3)
	//create a slice of room/channel
	var rooms []room
	//Listen for an incoming connection.
	for {
		clientHasFoundRoom := false
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
		}
		println("A new client has been detected")
		//add client in a free Room
		for index := range rooms {
			println("nb chatters BEFORE = ", len(rooms[index].chatters))
			//todo this cond is not required we can use the return value of addChatter but I am too lazy
			if rooms[index].isFull() {
				continue
			} else {
				rooms[index].addChatter(conn)
				rooms[index].sendToAll("\nA new client has joined the room\n")
				clientHasFoundRoom = true
			}
			println("nb chatters AFTER = ", len(rooms[index].chatters))
		}
		//if there is no free room, I create a new one
		if !clientHasFoundRoom {
			println("Create a new room")
			//create a new room
			rooms = append(rooms, newRoom(nbMaxChatterInRoom))
			//get the last index of the slice
			newRoomIndex := len(rooms) - 1
			//add client in new room
			rooms[newRoomIndex].addChatter(conn)
			//load room in a goroutine (thread)
			go manageRoom(&rooms[newRoomIndex])
		}
	}
}

//todo use channels instead of pointers
//manageRoom: handle and manage all the room logic. In our case the client send a message in the room and the server resend the message for all other room's clients
func manageRoom(room *room) {
	defer room.close()
	println("Start to manage room")
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		for i, chatter := range room.chatters {
			//todo find a way to set read in non blocking mode WITHOUT timeout
			chatter.SetReadDeadline(time.Now().Add(1 * time.Microsecond))
			//the first return value of read correspond to the number of byte read
			_, err := chatter.Read(buf)
			eof, timeout := room.handleReadError(err)
			if eof {
				room.deleteChatter(&chatter)
				room.sendToAll("\nClient " + strconv.Itoa(i+1) + " has quitted the room\n")
				continue
			}
			if timeout {
				continue
			}
			//unserialize
			clientbdd := protobuff.ClientBDD{}

			err = proto.Unmarshal(buf, &clientbdd)
			if err != nil {
				println("Unmarshall error detected: ", err.Error()) //todo check why unmarshall func return an err but still works
			}

			//send text to all clients except the one who sent the message to the server
			room.sendToAllExcept(clientbdd.GetName()+": "+clientbdd.GetMsg(), []uint8{uint8(i)})
			//resend the text to the client in order to "save" his messages
			//Note/todo: it better to display the text directly in the client and not resend the message to the client
			room.sendTo("\nMe: "+string(buf), []uint8{uint8(i)})
		}
	}
}
