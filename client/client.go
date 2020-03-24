package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/hadi-ilies/ProtoBuffExample/protobuff"
	"github.com/jroimartin/gocui"
	"google.golang.org/protobuf/proto"
)

//connect: connect to server
func connect(ipAddr string, port string) net.Conn {
	client, err := net.Dial("tcp", ipAddr+":"+port)
	if err != nil {
		log.Fatal("Error during the connection :", err)
	}
	return client
}

//sendText: get and send the text from the stdin
func sendText(client net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		//buffer that store server's response
		text := getLine(scanner)
		if len(text) != 0 {
			clientBDD.Msg = text
			data, err := proto.Marshal(clientBDD)
			if err != nil {
				log.Fatalln("Failed to serialize post in protobuf:", err)
			}
			//send proto struct serialized!
			client.Write(data)
		} else {
			break
		}
	}
}

func startGUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// handle error
	}
	defer g.Close()

	g.Cursor = true
	//set the main func
	g.SetManagerFunc(layout)

	//set key binding
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, input); err != nil {
		log.Panicln(err)
	}
	// if err := g.SetKeybinding("hello", gocui.KeyEnter, gocui.ModNone, input); err != nil {
	// 	log.Panicln(err)
	// }

	//main loop but we need to code above this line
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		// handle error
	}
}

var clientBDD = &protobuff.ClientBDD{
	Name: "Lol",
}

//get line from stdin gocui
func input(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	clientBDD.Name = l
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	//get size window
	maxX, maxY := g.Size()
	//set a "view" and its size
	if v, err := g.SetView("hello", maxX/2-14, maxY/2, maxX/2+14, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "set your Name...")
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("hello"); err != nil {
			return err
		}
	}
	return nil
}

//StartClient : Start a simple client that will connect to the server and send text
func StartClient(ipAddr string, port string) {
	//start gui
	startGUI()
	println("Your name is : ", clientBDD.Name)
	//connect client
	client := connect(ipAddr, port)
	//async text input
	go sendText(client)
	for {
		// Read the incoming connection into the buffer.
		client.SetReadDeadline(time.Now().Add(30 * time.Microsecond))
		buf := make([]byte, 1024)
		_, err := client.Read(buf)
		if err != nil {
			if err == io.EOF {
				// io.EOF, etc
				println("EOF detected !!!")
				return
			} else if err.(*net.OpError).Timeout() {
				//println("TIMEOUT")
				// no status msgs
				// note: TCP keepalive failures also cause this; keepalive is on by default
				continue
			}
			fmt.Println("Error reading:", err.Error())
		}
		//print message sent by the server
		println(string(buf))
	}
}
