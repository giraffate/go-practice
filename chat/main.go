package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type Room struct {
	Clients map[*Client]struct{}
	forward chan []byte
	join    chan *Client
	leave   chan *Client
}

func NewRoom() *Room {
	return &Room{
		Clients: make(map[*Client]struct{}),
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
	}
}

func (r *Room) ServeHTTP(ws *websocket.Conn) {
	client := Client{
		socket: ws,
		room:   r,
		send:   make(chan []byte),
	}
	r.join <- &client
	defer func() { r.leave <- &client }()
	go client.Read()
	client.Write()
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.Clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.Clients, client)
		case msg := <-r.forward:
			for client := range r.Clients {
				select {
				case client.send <- msg:
				default:
					delete(r.Clients, client)
				}
			}
		}
	}
}

type Client struct {
	socket *websocket.Conn
	room   *Room
	send   chan []byte
}

func (c *Client) Read() {
	for {
		var b []byte
		if err := websocket.Message.Receive(c.socket, &b); err == nil {
			c.room.forward <- b
		} else {
			break
		}
	}
}

func (c *Client) Write() {
	for {
		for b := range c.send {
			if err := websocket.Message.Send(c.socket, b); err != nil {
				break
			}
		}
	}
}

func main() {
	r := NewRoom()
	go r.Run()

	http.Handle("/chat", websocket.Handler(r.ServeHTTP))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
}
