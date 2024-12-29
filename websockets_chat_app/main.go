package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var ws websocket.Upgrader

var Clients map[string]*Client

type Message struct {
	Data string
	From string
}

type Client struct {
	Name string
	io.Writer
	io.Reader
	conn      *websocket.Conn
	readChan  chan Message // messages received over read chan need to be written to friends of client (supports 1 friend only)
	writeChan chan Message // messages received over write chan need to be written back to client
	friends   []string
}

func (c *Client) ReadMessage() {
	for {
		_, msg, err := c.conn.ReadMessage() // from the context of the server i.e. server reads message from client
		if err != nil {
			if err == io.EOF {
				fmt.Println("received EOF when reading message.")
				break
			}
			fmt.Println("error reading message : ", err)
			continue
		}
		// fmt.Println(c.Name, " : sending message : ", string(msg))
		c.readChan <- Message{string(msg), c.Name} // all read messages go to readChan
		fmt.Println("sent message")
	}
}

func (c *Client) WriteMessage() {
	for {
		msg := <-c.writeChan
		// fmt.Println(c.Name, ": Writing message : ", string(msg))       // all messages to be written to screen go to writeChan
		err := c.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(msg.From, " says : ", msg.Data))) // from the context of the server i.e. server writes message to client
		if err != nil {
			if err == io.EOF {
				fmt.Println("received EOF when writing message.")
				break
			}
			fmt.Println("error reading message : ", err)
			continue
		}
	}
}

func chatEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received chat endpoint request from ", r.RemoteAddr)
	ws.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := ws.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client := &Client{
		Name:      r.Header.Get("username"),
		friends:   strings.Split(r.Header.Get("friends"), ","),
		conn:      ws,
		readChan:  make(chan Message),
		writeChan: make(chan Message),
	}
	fmt.Printf("Creating client for %v with friends %v\n", client.Name, client.friends)
	Clients[client.Name] = client
	go Clients[client.Name].ReadMessage()
	go Clients[client.Name].WriteMessage()
}

func healthCheckEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the /health_check endpoint response"))
}

func init() {
	Clients = make(map[string]*Client)
}

func manageFriends() {
	for {
		for _, client := range Clients {
			go func(client *Client) {
				msg := <-client.readChan
				for _, friend := range client.friends {
					// sending out the message to all of my friends :)
					go func(msg Message) {
						// fmt.Println("FROM: ", client.Name, "TO:", Clients[friend].Name, " sending message ", msg)
						Clients[friend].writeChan <- msg
					}(msg)
				}
			}(client)
		}
	}
}

func main() {
	srvr := http.Server{Addr: ":9000"}
	http.HandleFunc("/chat", chatEndpoint)
	http.HandleFunc("/health_check", healthCheckEndpoint)
	fmt.Println("STARTING SERVER")
	go manageFriends() // go over the friends
	fmt.Println(srvr.ListenAndServe())
}
