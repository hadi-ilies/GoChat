package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

//room is a kind of channel that allow users to make group conv
type room struct {
	nbMaxChatter uint8
	chatters     []net.Conn
	messagesDB   []string
}

//NewRoom is a room constructor
func newRoom(nbChatter uint8) room {
	return room{nbMaxChatter: nbChatter}
}

//isFull: check whether the room is full or not
func (room *room) isFull() bool {
	if uint8(len(room.chatters)) < room.nbMaxChatter {
		return false
	}
	return true
}

//close: close the room by closing all chatters conn
func (room *room) close() {
	//note the _ (ignored) is the index and chatter the element
	for _, chatter := range room.chatters {
		chatter.Close()
	}
}

//getChatterID: get client index from chatters
func (room *room) getChatterID(target *net.Conn) int {
	for index, chatter := range room.chatters {
		if chatter == *target {
			return index
		}
	}
	return -1 //chatter does not exist
}

//deleteChatter: allow to delete/eject client from the room
func (room *room) deleteChatter(target *net.Conn) {
	chatterID := room.getChatterID(target)
	if chatterID == -1 {
		log.Fatal("Error chatter doesn't exist !!!")
	}
	//close chatter
	room.chatters[chatterID].Close()
	//eject chatter from the room
	room.chatters = append(room.chatters[:chatterID], room.chatters[chatterID+1:]...) //chatterID is the index of the element to be removed
}

//addChatter: allow to add a new client in a room
func (room *room) addChatter(newChatter net.Conn) {
	if uint8(len(room.chatters)) < room.nbMaxChatter {
		room.chatters = append(room.chatters, newChatter)
	} else {
		println("Too many chatters in this Room ! Create another one or increase the maximum number of chatters ;)")
	}
}

//this function detect whether read has detected an EOF or just has 'timeouted' or there is no error
//@return bool, bool: the first bool is eof detection and the seconde timeout detection
func (room *room) handleReadError(err error) (bool, bool) {
	if err != nil {
		if err == io.EOF {
			// io.EOF, etc
			println("EOF detected !!!")
			return true, false //tmp
		} else if err.(*net.OpError).Timeout() {
			//this allow us to set an "no-blocking read"
			// no status msgs
			// note: TCP keepalive failures also cause this; keepalive is on by default
			return false, true
		}
		fmt.Println("Error reading:", err.Error())
	}
	return false, false
}

//sendToAll: send text to all clients
func (room *room) sendToAll(text string) {
	for i := range room.chatters {
		room.chatters[i].Write([]byte(text))
	}
}

//sendToAllExcept: send text to all clients except the ones passed in arg
func (room *room) sendToAllExcept(text string, chattersID []uint8) {
	for i := range room.chatters {
		isExcepted := false

		for _, chatterID := range chattersID {
			if uint8(i) == chatterID {
				isExcepted = true
			}
		}
		if !isExcepted {
			room.chatters[i].Write([]byte(text))
		}
	}
}

//sendTo: send text to all chatters specified in param
//Note We considerate that the chattersId slice is always good
func (room *room) sendTo(text string, chattersID []uint8) {
	for _, chatterID := range chattersID {
		room.chatters[chatterID].Write([]byte(text))
	}
}
