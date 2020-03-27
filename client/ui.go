package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/jroimartin/gocui"
	"google.golang.org/protobuf/proto"
)

func (client *Client) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone, client.switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("chatRoom", gocui.KeyArrowDown, gocui.ModNone, client.switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, client.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, client.input); err != nil {
		return err
	}
	return nil
}

func (client *Client) startGUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		// handle error
	}
	defer g.Close()

	//set the main func
	g.SetManagerFunc(client.layout)

	//set key binding
	if err := client.keybindings(g); err != nil {
		log.Panicln(err)
	}

	//main loop but we need to code above this line
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (client *Client) layout(g *gocui.Gui) error {
	//get size window
	maxX, maxY := g.Size()
	//set a "view" and its size 	//create a view and set it
	v, err := setView(g, "input", newDimension(1, maxY-14, maxX-1, maxY-1))
	//g.CurrentView().Title
	if err != nil {
		fmt.Fprintln(v, "write text...")
		g.Cursor = true
		v.Editable = true
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	//set and create a view
	v, err = setView(g, "chatRoom", newDimension(1, 1, maxX-1, maxY-14))
	if err != nil {
		v.Autoscroll = true
		go client.chatRoom(g, v)
	}
	return nil
}

//create a gocui view
func setView(g *gocui.Gui, name string, coords dimension) (*gocui.View, error) {
	v, err := g.SetView(name, coords.x0, coords.y0, coords.x1, coords.y1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Fatalf("View error:", err)
	}
	return v, err
}

func (client *Client) chatRoom(g *gocui.Gui, v *gocui.View) error {
	// Read the incoming connection into the buffer.
	for {
		client.conn.SetReadDeadline(time.Now().Add(30 * time.Microsecond))
		buf := make([]byte, 1024)
		_, err := client.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// io.EOF, etc
				println("EOF detected !!!")
				return nil
			} else if err.(*net.OpError).Timeout() {
				//println("TIMEOUT")
				// no status msgs
				// note: TCP keepalive failures also cause this; keepalive is on by default
				continue
			}
			fmt.Println("Error reading:", err.Error())
		}

		//update views I almost rage quit ^^ if we do not call Update method its will store data in buffer just like print without flush in C ;)
		g.Update(func(g *gocui.Gui) error {
			//print message sent by the server
			fmt.Fprintln(v, string(buf))
			return nil
		})
	}
	return nil
}

//get line from stdin gocui
func (client *Client) input(g *gocui.Gui, v *gocui.View) error {
	var l = ""
	var err error

	//g.Cursor = false
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	if l == "" {
		return nil
	}

	client.packet.Msg = l
	data, err := proto.Marshal(client.packet)
	if err != nil {
		log.Fatalln("Failed to serialize post in protobuf:", err)
	}
	//send proto struct serialized!
	client.conn.Write(data)
	v.Clear()
	return nil
}

func (client *Client) switchView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "input" {
		_, err := g.SetCurrentView("chatRoom")
		//v.Clear() //tmp
		return err
	}
	v, err := g.SetCurrentView("input")
	v.Clear()
	return err
}

//quit
func (client *Client) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
