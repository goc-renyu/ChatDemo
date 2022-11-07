package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Hub *Hub

	Conn *websocket.Conn

	Uid string

	UserName string

	RoomId string

	Received chan Message
}

func (c *Client) getOnlineMsg() Message {
	var onlineData OnlineData
	onlineData.Uid = c.Uid
	onlineData.RoomId = c.RoomId
	onlineData.UserName = c.UserName
	return Message{
		Status: msgTypeOnline,
		Data: onlineData,
	}
}

func (c *Client) getOfflineMsg() Message {
	return Message{
		Status: msgTypeOffline,
		Data: OfflineData{
			Uid: c.Uid,
			RoomId: c.RoomId,
			UserName: c.UserName,
		},
	}
}

func (c *Client) read() {
	defer func() {
		//捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("read 发生错误", err)
		}
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.Hub.UnRegister <- c
			log.Println("ReadMessage 发生错误")
			return
		}
		var clientMsg Message
		json.Unmarshal(message, &clientMsg)
		fmt.Println(clientMsg)
		if clientMsg.Data != nil {
			if clientMsg.Status == msgTypeOnline {
				c.Uid = clientMsg.Data.(map[string]interface{})["uid"].(string)
				c.UserName = clientMsg.Data.(map[string]interface{})["userName"].(string)
				c.RoomId = clientMsg.Data.(map[string]interface{})["roomId"].(string)
				c.Hub.Register <- c
			}
		}
		c.Hub.Broadcast <- clientMsg
	}
}

func (c *Client) write() {
	defer func() {
		//捕获write抛出的panic
		if err := recover();err!=nil{
			log.Println("write 发生错误",err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Received:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			serverMsg, _ := json.Marshal(message)
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(serverMsg)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Received: make(chan Message),
	}
	go client.read()
	go client.write()
	select {}
}
