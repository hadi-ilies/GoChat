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

//scrollView: scroll my view
func (client *Client) scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

//scrollDown: use scrollView in order to scroll down
func (client *Client) scrollDown(g *gocui.Gui, v *gocui.View) error {
	client.scrollView(v, 2)
	return nil
}

//scrollUp: use scrollView in order to scroll up
func (client *Client) scrollUp(g *gocui.Gui, v *gocui.View) error {
	client.scrollView(v, -2)
	return nil
}

//keybindings: this function manage keybindings
//todo: replace "if forest" by a map of functions
func (client *Client) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("input", gocui.KeyCtrlV, gocui.ModNone, client.switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("chatRoom", gocui.KeyCtrlV, gocui.ModNone, client.switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("chatRoom", gocui.KeyArrowUp, gocui.ModNone, client.scrollUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("chatRoom", gocui.KeyArrowDown, gocui.ModNone, client.scrollDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, client.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, client.input); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyCtrlL, gocui.ModNone, client.clearInput); err != nil {
		return err
	}
	return nil
}

//clear: allow to clear the input by pressing (ctrl + L)
func (client *Client) clearInput(g *gocui.Gui, v *gocui.View) error {
	//clear the view
	v.Clear()
	//set the cursor in the begining of the view
	v.SetCursor(0, 0)
	return nil
}

//startGUI: start the ui
func (client *Client) startGUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		//handle error
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

//layout: the "main" function of my ui
func (client *Client) layout(g *gocui.Gui) error {
	//get size window
	maxX, maxY := g.Size()
	//set a "view" and its size 	//create a view and set it
	v, err := client.setView(g, "input", newDimension(1, maxY-9, maxX-1, maxY-1))
	//g.CurrentView().Title
	if err != nil {
		fmt.Fprintln(v, "write text...")
		g.Cursor = true
		v.Editable = true
		v.Frame = true
		v.Title = " New message - <" + client.packet.Name + "> "
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	//set and create a view
	v, err = client.setView(g, "chatRoom", newDimension(1, 1, maxX-1, maxY-14))
	if err != nil {
		v.Autoscroll = true
		v.Frame = true
		//todo add dynamic room name
		v.Title = " Messages - <Room> "
		go client.chatRoom(g, v)
	}
	return nil
}

//setView: create a gocui view
func (client *Client) setView(g *gocui.Gui, name string, coords dimension) (*gocui.View, error) {
	v, err := g.SetView(name, coords.x0, coords.y0, coords.x1, coords.y1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Fatal("View error:", err)
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
			fmt.Fprintln(v, "\033[31;2m"+string(buf)+"\033[0m")
			return nil
		})
	}
}

//input: get and send the text from the input view and send it to the server
func (client *Client) input(g *gocui.Gui, v *gocui.View) error {
	var l = ""
	var err error

	//set the cursor in the begining of the view
	v.SetCursor(0, 0)
	//get 'y' point of the cursor in order to get the line
	_, cy := v.Cursor()
	//get line
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	//if the in put is empty we do not send our data
	if l == "" {
		return nil
	}
	//save the data wrote in the input view into the client's packet
	client.packet.Msg = l
	//serialize our packet with protobuff
	data, err := proto.Marshal(client.packet)
	if err != nil {
		log.Fatalln("Failed to serialize post in protobuf:", err)
	}
	//send proto struct serialized!
	client.conn.Write(data)
	//clear input view after data has been sent
	v.Clear()
	return nil
}

//switchView: switch the focus between input/chatroom view by pressing (ctrl + v)
func (client *Client) switchView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "input" {
		_, err := g.SetCurrentView("chatRoom")
		return err
	}
	_, err := g.SetCurrentView("input")
	return err
}

//quit: quit the ui by clicking on (ctrl + c)
func (client *Client) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
