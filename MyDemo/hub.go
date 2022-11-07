package main

type Hub struct {
	Clients map[string]map[string]*Client

	Broadcast chan Message

	Register chan *Client

	UnRegister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Clients:    make(map[string]map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Clients[client.RoomId]; !ok{
				h.Clients[client.RoomId] = make(map[string]*Client)
			}
			h.Clients[client.RoomId][client.Uid] = client
		case client := <-h.UnRegister:
			if _, ok := h.Clients[client.RoomId][client.Uid]; ok {
				delete(h.Clients[client.RoomId], client.Uid)
			}
			if len(h.Clients[client.RoomId]) == 0 {
				delete(h.Clients, client.RoomId)
			}
		case message := <-h.Broadcast:
			switch message.Status {
			case msgTypeOnline:
				roomId := message.Data.(map[string]interface{})["roomId"].(string)
				for _, client := range h.Clients[roomId] {
					client.Received <- message
				}
			case msgTypeOffline:
				roomId := message.Data.(map[string]interface{})["roomId"].(string)
				uid := message.Data.(map[string]interface{})["uid"].(string)
				if _,ok := h.Clients[roomId][uid]; ok {
					close(h.Clients[roomId][uid].Received)
					delete(h.Clients[roomId],uid)
				}
				for _, client := range h.Clients[roomId] {
					client.Received <- message
				}
			case msgTypePrivateChat:
				roomId := message.Data.(map[string]interface{})["roomId"].(string)
				uid := message.Data.(map[string]interface{})["uid"].(string)
				toUid := message.Data.(map[string]interface{})["toUid"].(string)
				if _,ok := h.Clients[roomId][toUid]; ok {
					h.Clients[roomId][toUid].Received <- message
				}
				if _,ok := h.Clients[roomId][uid]; ok {
					h.Clients[roomId][uid].Received <- message
				}
			case msgTypePublicChat:
				roomId := message.Data.(map[string]interface{})["roomId"].(string)
				for _, client := range h.Clients[roomId] {
					client.Received <- message
				}
			default:
			}
		}
	}
}