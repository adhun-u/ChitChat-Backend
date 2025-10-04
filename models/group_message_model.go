package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

// For getting details for pushing notification
type GroupMessageNotificationModel struct {
	GroupId         string `json:"groupId"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	MessageType     string `json:"messageType"`
	GroupProfilePic string `json:"groupProfilePic"`
}

// For making a client for group chat for sending an indicator for typing and recording audio
type GroupClient struct {
	UserId     int
	GroupId    string
	Connection *websocket.Conn
}

// For reading the indication from json to model
type GroupIndication struct {
	Indication     string `json:"indication"`
	IndicationType string `json:"indicationType"`
	SenderId       int    `json:"senderId"`
	GroupId        string `json:"groupId"`
}

// For handling each groups
type GroupsConnection struct {
	Clients map[int]*GroupClient
	Mutex   sync.RWMutex
}
