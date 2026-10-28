package main

import (
	"bytes"
	"html/template"
	"log"
)

type Message struct {
	ClientID string
	Text     string
}

type WSMessage struct {
	Headers interface{} `json:"HEADERS"`
	Text    string      `json:"text"`
}

type Hub struct {
	clients    map[*Client]bool
	messages   []*Message
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("cli ent registered %s", client.id)
			// for _, msg := range h.messages {
			// 	client.send <- getMessageTemplate(msg)
			// }
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case msg := <-h.broadcast:
			h.messages = append(h.messages, msg)
			for client := range h.clients {
				select {
				case client.send <- getMessageTemplate(msg):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("template/message.html")
	if err != nil {
		log.Printf("template parsing: %v", err)
		return nil
	}

	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Printf("template execution: %v", err)
		return nil
	}

	return renderedMessage.Bytes()
}
