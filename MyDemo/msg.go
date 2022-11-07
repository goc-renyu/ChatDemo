package main

type Message struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// 定义消息类型
const msgTypeOnline = 1      // 上线
const msgTypeOffline = 2     // 离线
const msgTypePrivateChat = 3 // 私聊
const msgTypePublicChat = 4  // 群聊

type OnlineData struct {
	Uid 	   string `json:"uid"`
	RoomId   string `json:"roomId"`
	UserName string `json:"userName"`
}

type OfflineData struct {
	Uid      string `json:"uid"`
	RoomId   string `json:"roomId"`
	UserName string `json:"userName"`
}