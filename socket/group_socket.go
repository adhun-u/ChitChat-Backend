package socket

import (
	groupCallController "chitchat/controllers/group"
	"chitchat/helpers"
	"chitchat/models"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Creating an upgrader for upgrading normal http connection to websocket connection
var groupSocketUpgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// For handling each user with in a group
var groups = make(map[string]*models.GroupsConnection)

// For preventing concurrent writes and reads
var globalGroupMutex sync.RWMutex

// ---------------- GROUP SOCKET -------------
// For handling typing and recording indication for a group
func ConnectGroupSocket(ctx *gin.Context) {

	//Querying user id for identify the user
	userId, intConvErr := strconv.Atoi(ctx.Query("userId"))
	//Querying group id for identify the group
	groupId := ctx.Query("groupId")
	//Websocket connection for each user
	conn, connErr := groupSocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if connErr != nil {
		helpers.SendMessageAsJson(ctx, "Websocket connection error while connecting", http.StatusInternalServerError)
		return
	}

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide an integer", http.StatusInternalServerError)
		return
	}

	//Initializing the a group if the group is not created already
	globalGroupMutex.Lock()
	if groups[groupId] == nil {
		groups[groupId] = &models.GroupsConnection{
			Clients: make(map[int]*models.GroupClient),
		}
	}
	//Adding all clients to the group
	group := groups[groupId]

	globalGroupMutex.Unlock()

	group.Mutex.Lock()
	group.Clients[userId] = &models.GroupClient{
		UserId:     userId,
		GroupId:    groupId,
		Connection: conn,
	}
	group.Mutex.Unlock()
	//For clearing all resources related to current user when current user exists
	defer func() {
		group.Mutex.Lock()
		conn.Close()
		delete(group.Clients, userId)
		group.Mutex.Unlock()
		//If group has no clients , then deleting the group
		globalGroupMutex.Lock()
		group, ok := groups[groupId]

		if ok && group != nil {
			if len(groups[groupId].Clients) == 0 {
				delete(groups, groupId)
			}
		}

		globalGroupMutex.Unlock()
	}()
	for {
		fmt.Println("Connected group chat socket")
		var groupIndication models.GroupIndication
		readErr := conn.ReadJSON(&groupIndication)

		if readErr != nil {
			fmt.Println("Reading error in group socket : ", readErr)
			return
		}
		//For sending the info if the group is already in call
		if groupIndication.Indication == "call" {
			if sender, ok := group.Clients[groupIndication.SenderId]; ok {
				if isInCall, ok := groupCallController.GroupCalls[groupId]; ok {
					sender.Connection.WriteJSON(gin.H{"isInCall": isInCall, "indication": "call", "senderId": userId, "indicationType": "call", "groupId": groupId})
				} else {
					sender.Connection.WriteJSON(gin.H{"isInCall": false, "indication": "call", "senderId": userId, "indicationType": "call", "groupId": groupId})
				}
			}
		}
		//For the group to exist from calling state
		if groupIndication.Indication == "close" {
			groupCallController.GroupCalls[groupId] = false
		}
		//For sending seen indication to sender only
		if groupIndication.Indication == "seen" {
			group.Mutex.Lock()
			if senderConn, ok := group.Clients[groupIndication.SenderId]; ok {
				fmt.Println(groupIndication.SenderId)
				handleGroupIndications(senderConn.Connection, groupIndication)
			}
			group.Mutex.Unlock()
		}

		group.Mutex.RLock()

		//For getting all members except current user of a group
		members := make([]*models.GroupClient, 0, len(group.Clients))
		for _, client := range group.Clients {
			//Not adding current user to members to prevent sending indication back to sender
			if client.UserId != userId {
				members = append(members, client)
			}

		}
		group.Mutex.RUnlock()
		for _, client := range members {
			//Sending the indication to each client
			handleGroupIndications(client.Connection, groupIndication)
		}

	}

}

// For sending  the indication to all online members of the group
func handleGroupIndications(connection *websocket.Conn, indication models.GroupIndication) {

	writeErr := connection.WriteJSON(indication)

	if writeErr != nil {
		fmt.Println("Write error : ", writeErr)
	}
}
