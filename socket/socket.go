package socket

import (
	"chitchat/config"
	"chitchat/helpers"
	messagemodel "chitchat/models"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Initializing the upgrader for upgrading http to websockets
var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// For sending messages and receiving
var chatConn = make(map[messagemodel.ChatKey]*websocket.Conn)

// For checking online
var onlineConn = make(map[int]*websocket.Conn)

// For adding receiver's id for indication online and seen info
var receiverMap = make(map[int]int)

var mutex = sync.Mutex{}

// For connecting websockets
func ConnectMessageSocket(ctx *gin.Context) {
	//Getting current user id using query
	currentUserIdStr := ctx.Query("currentUserId")
	//Getting current user name
	currentUsername := ctx.Query("currentUsername")
	//Getting current user profile pic
	currentUserProfilePic := ctx.Query("currentUserProfilePic")
	//Converting string current user id into int current user id
	currentUserId, currentUserIdconvErr := strconv.Atoi(currentUserIdStr)

	if currentUserIdconvErr != nil {
		fmt.Println("Conversion error in websocket : ", currentUserIdconvErr)
		helpers.SendMessageAsJson(ctx, "Provide valid user id", http.StatusNotAcceptable)
		return
	}
	fmt.Println("Current user id in connection : ", currentUserId)
	//Upgrading http to websockets
	connection, upgraderErr := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if upgraderErr != nil {
		fmt.Println("Socket Error  : ", upgraderErr)
		return
	}

	mutex.Lock()
	onlineConn[currentUserId] = connection
	mutex.Unlock()

	//Indicating opposite user when current user enters in online
	if _, ok := onlineConn[currentUserId]; ok {
		if receiverMap == nil {
			return
		}
		if receiverId, ok := receiverMap[currentUserId]; ok {
			indicateOnline(onlineConn[receiverId], true)
		}
	}
	fmt.Println("Connected socket")
	defer func() {
		//Deleting the connection with current user when current user exits
		mutex.Lock()
		connection.Close()
		delete(onlineConn, currentUserId)
		//To indicate receiver when current user exits
		if receiverId, ok := receiverMap[currentUserId]; ok {
			indicateOnline(onlineConn[receiverId], false)
		}
		mutex.Unlock()
	}()
	for {
		var chatConnModel messagemodel.ChatConnectionModel
		chatConnModel.Mutext.RLock()
		readErr := connection.ReadJSON(&chatConnModel)
		if readErr != nil {
			fmt.Println("Reading error in socket")
			return
		}
		chatConnModel.Mutext.RUnlock()
		//--------- ONE TO ONE CHAT ----------------
		if chatConnModel.IsInChat {
			fmt.Println("Entered")
			//Adding current user's receiver
			chatConnModel.Mutext.Lock()
			receiverMap[chatConnModel.ReceiverId] = currentUserId
			chatConnModel.Mutext.Unlock()
			fmt.Println("Current user id : ", currentUserId)
			//Adding current user to chat connection
			chatConnModel.Mutext.Lock()
			chatConn[messagemodel.ChatKey{SenderId: currentUserId, ReceiverId: chatConnModel.ReceiverId}] = connection
			chatConnModel.Mutext.Unlock()
			//Key to identify current user
			chatConnModel.Mutext.RLock()
			senderConnKey := messagemodel.ChatKey{SenderId: currentUserId, ReceiverId: chatConnModel.ReceiverId}
			chatConnModel.Mutext.RUnlock()
			//Checking if receiver is in online
			chatConnModel.Mutext.Lock()
			if _, ok := onlineConn[chatConnModel.ReceiverId]; ok {
				fmt.Println("Online")
				//Sending online indication
				indicateOnline(chatConn[senderConnKey], true)
			} else {
				fmt.Println("Offline")
				indicateOnline(chatConn[senderConnKey], false)
			}
			chatConnModel.Mutext.Unlock()
			//To send seen indication to receiver's connection when current user is in chat connection
			chatConnModel.Mutext.Lock()
			seenIndication(chatConnModel.ReceiverId, currentUserId)
			chatConnModel.Mutext.Unlock()
			for {
				var chat messagemodel.MessageModel
				//Reading the message as json
				chatConnModel.Mutext.RLock()
				messageErr := connection.ReadJSON(&chat)
				if messageErr != nil {
					fmt.Println("Socket reading error : ", messageErr)
					break
				}
				chatConnModel.Mutext.RUnlock()
				if chat.Exit {
					fmt.Println("Chat connection exited userId : ", currentUserId)
					chatConnModel.Mutext.Lock()
					delete(chatConn, senderConnKey)
					delete(receiverMap, chatConnModel.ReceiverId)
					chatConnModel.Mutext.Unlock()
					break
				}

				chatConnModel.Mutext.RLock()
				//Sending and receiving message
				handleErr := handleMessages(
					chatConn,
					&chat,
					currentUsername,
					currentUserProfilePic,
				)
				chatConnModel.Mutext.RUnlock()
				if handleErr != nil {
					fmt.Println("Everything is deleted")
					//Deleting the user from chat connection when the user exists
					chatConnModel.Mutext.Lock()
					delete(chatConn, senderConnKey)
					delete(receiverMap, chatConnModel.ReceiverId)
					chatConnModel.Mutext.Unlock()
					fmt.Println(handleErr)
					break
				}

			}
		}

	}

}

// For handling messages
func handleMessages(chatConn map[messagemodel.ChatKey]*websocket.Conn, chat *messagemodel.MessageModel, currentUsername string, currentUserProfilePic string) error {
	//Sending the message back to sender
	if senderConnection, ok := chatConn[messagemodel.ChatKey{SenderId: chat.SenderId, ReceiverId: chat.ReceiverId}]; ok {
		fmt.Println("entered int sender connec ")
		fmt.Println("Message type : ", chat.Type)
		//If type of the message is media type , then not sending to sender
		if chat.Type != "image" && chat.Type != "audio" && chat.Type != "video" && chat.Type != "voice" {
			sendingErr := sendMessage(senderConnection, chat, chat.SenderId)
			if sendingErr != nil {
				return sendingErr
			}
		}

	} else {
		fmt.Println("No connection")
	}
	_, isReceiverInChat := chatConn[messagemodel.ChatKey{SenderId: chat.ReceiverId, ReceiverId: chat.SenderId}]
	chat.IsRead = isReceiverInChat
	chat.IsSeen = isReceiverInChat
	//Sending the message to receiver
	if receiverConn, ok := chatConn[messagemodel.ChatKey{SenderId: chat.ReceiverId, ReceiverId: chat.SenderId}]; ok {
		//Sending the message to receiver
		sendingErr := sendMessage(receiverConn, chat, chat.ReceiverId)
		if sendingErr != nil {
			return sendingErr
		}
		//Sending seen indication when receiver is in chat connection
		seenIndication(chat.SenderId, chat.ReceiverId)

	} else if chat.Type != "Typing" && chat.Type != "Not typing" && chat.Type != "Recording" && chat.Type != "Not recording" {
		//If receiver is not in chat connection , then sending to online connection
		if OnlineReceiverConn, ok := onlineConn[chat.ReceiverId]; ok {
			sendingErr := sendMessage(OnlineReceiverConn, chat, chat.ReceiverId)
			if sendingErr != nil {
				return sendingErr
			}
		} else {
			//Saving message
			saveMessageAsTemp((*chat))
		}
		//Sending notification to receiver's device
		notifyReceiver(currentUsername, currentUserProfilePic, (*chat))
	}
	return nil

}

// For sending notification to receiver's device
func notifyReceiver(senderName string, senderProfile string, model messagemodel.MessageModel) {
	services.SendMessageNotification(
		model.SenderId,
		model.ReceiverId,
		senderName,
		senderProfile,
		model,
	)

}

// For saving messages when receiver is not in connection
func saveMessageAsTemp(chat messagemodel.MessageModel) {
	//Inserting to tempMessage collection for storing temporary messages
	collection := config.MongoDB.Collection("tempMessages")
	_, insertionErr := collection.InsertOne(context.TODO(), chat)
	fmt.Println("chat read status : ", chat.IsRead)
	fmt.Println("Chat seen status : ", chat.IsSeen)
	if insertionErr != nil {
		fmt.Println("Temp message insertion error in socket : ", insertionErr)
	}
}

// For sending the message to receiver's connection
func sendMessage(conn *websocket.Conn, chat *messagemodel.MessageModel, userId int) error {
	//Checking if the message type is not a message
	if chat.Type == "Typing" || chat.Type == "Not typing" || chat.Type == "Recording" || chat.Type == "Not recording" {
		//Sending the indication to receiver
		if chat.ReceiverId == userId {
			fmt.Println("Chat type : ", chat.Type)
			writeMessageErr := conn.WriteJSON(messagemodel.InteractionIndicatorModel{
				ReceiverId: chat.ReceiverId,
				Indication: chat.Type,
			})
			if writeMessageErr != nil {
				fmt.Println("Indicator writing error in message socket : ", writeMessageErr)
				return writeMessageErr
			}
		}

	} else {
		writeMessageErr := conn.WriteJSON(chat)
		if writeMessageErr != nil {
			fmt.Println("Sending message Socket error : ", writeMessageErr)
			return writeMessageErr
		}

	}
	return nil
}

// To indicate current user when the receiver is in online
func indicateOnline(senderConn *websocket.Conn, isOnline bool) {
	fmt.Println("Sent online indication")
	senderConn.WriteJSON(gin.H{"isOnline": isOnline, "type": "status"})
}

// To send "seen" indication
func seenIndication(senderId int, receiverId int) {
	if conn, isInChatConn := chatConn[messagemodel.ChatKey{SenderId: senderId, ReceiverId: receiverId}]; isInChatConn {
		conn.WriteJSON(gin.H{"isSeen": true, "type": "seen", "senderId": senderId, "receiverId": receiverId})
	} else if onlineConn, isInOnlineConn := onlineConn[senderId]; isInOnlineConn {
		onlineConn.WriteJSON(gin.H{"isSeen": true, "type": "seen", "senderId": senderId, "receiverId": receiverId})
	}
}
