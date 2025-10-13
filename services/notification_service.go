package services

import (
	messagemodel "chitchat/models"
	"context"
	"fmt"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// Creating a variable for getting client for sending notification
var client *messaging.Client

// Initializing firebase to send notifications
func InitializeFirebaseApp() {

	//Getting options to create firebase app with credentials
	options := option.WithCredentialsFile("service.json")
	//Initializing firebase
	app, initErr := firebase.NewApp(context.TODO(), nil, options)

	if initErr != nil {
		fmt.Println("Firebase initializing error : ", initErr)
		return
	}
	fClient, clientErr := app.Messaging(context.TODO())

	if clientErr != nil {
		fmt.Println("Client initializing error : ", clientErr)
		return
	}

	client = fClient

	fmt.Println("Firebase initialized")

}

// For sending normal notification eg: when someone sends a request to someone
func SendNotification(deviceToken string, title string, body string, imageUrl string) {
	message := &messaging.Message{
		Token: deviceToken,
		Data: map[string]string{
			"isMessageNotification": "false",
			"imageUrl":              imageUrl,
			"title":                 title,
			"body":                  body,
		},
	}
	//Sending the notification
	_, responseErr := client.Send(context.TODO(), message)

	if responseErr != nil {
		fmt.Println("Response error in firebase : ", responseErr)
		return
	}
}

// For sending when message notification when new message comes
func SendMessageNotification(userId int, anotherUserId int, username string, profilePic string, messageModel messagemodel.MessageModel) {
	fmt.Println("Entered in message notification")
	//Creating the topic to send notifications
	topic := "userMessage" + strconv.Itoa(userId) + strconv.Itoa(anotherUserId)
	newMessage := &messaging.Message{
		Topic: topic,
		Data: map[string]string{
			"title":                 username,
			"body":                  messageModel.Message,
			"imageUrl":              profilePic,
			"isMessageNotification": "true",
			"type":                  messageModel.Type,
		},
	}
	//Sending notification
	response, sendErr := client.Send(context.TODO(), newMessage)

	if sendErr != nil {
		fmt.Println("Sending notification error in message notification : ", sendErr)
		return
	}

	fmt.Println(response)
}

// For subscribing a user to a topic to get another user's message as notifications
func SubscribeToUserMessageTopic(deviceToken string, userId int, anotherUserId int) {
	//Creating the topic to sub
	topic := "userMessage" + strconv.Itoa(userId) + strconv.Itoa(anotherUserId)

	_, subErr := client.SubscribeToTopic(context.TODO(), []string{deviceToken}, topic)

	if subErr != nil {
		fmt.Println("Subscribing user message topic error : ", subErr)
		return
	}
}

// For unsubscribing a user from the topic to not get another user's message as notifications
func UnSubscribeFromUserMessageTopic(deviceToken string, userId int, anotherUserId int) {
	//Creating the topic to unsub
	topic := "userMessage" + strconv.Itoa(userId) + strconv.Itoa(anotherUserId)

	_, unSub := client.UnsubscribeFromTopic(context.TODO(), []string{deviceToken}, topic)

	if unSub != nil {
		fmt.Println("Subscribing user message topic error : ", unSub)
		return
	}
}

// For subscribing a user to a topic to get group messages as notifications
func SubscribeUserToGroupMessageTopic(deviceToken string, groupId string) {
	//Creating a topic for a group to subscribe to get notifications
	topic := "groupMessage" + groupId

	//Subscribing the user with the device to get group messages as notifications
	_, subErr := client.SubscribeToTopic(context.TODO(), []string{deviceToken}, topic)

	if subErr != nil {
		fmt.Println("FCM subscription error : ", subErr)
		return
	}

}

// For unsubscribing a user from a topic when the user is removed or exited
func UnsubscribeFromGroupMessageTopic(deviceToken string, groupId string) {
	//Topic where the user is going to be removed from
	topic := "groupMessage" + groupId

	_, unSubErr := client.UnsubscribeFromTopic(context.TODO(), []string{deviceToken}, topic)

	if unSubErr != nil {

		fmt.Println("FCM unsubscribe error : ", unSubErr)
		return
	}

}

// For subscribing to call topic to get call notification (audio , video)
func SubscribeToGroupCallTopic(deviceToken string, groupId string) {
	//Topic to subscribe
	topic := "groupCall" + groupId

	_, subErr := client.SubscribeToTopic(context.TODO(), []string{deviceToken}, topic)

	if subErr != nil {
		fmt.Println("Subscribing call topic error : ", subErr)
		return
	}

}

// For unsubscribing a user from call topic
func UnsubscribeFromGroupCallTopic(deviceToken string, groupId string) {
	//Topic to unsubscribe
	topic := "groupCall" + groupId

	_, unSubErr := client.UnsubscribeFromTopic(context.TODO(), []string{deviceToken}, topic)

	if unSubErr != nil {
		fmt.Println("Unsubscribing error : ", unSubErr)
		return

	}

}

// For sending group message notification
func SendGroupMessageNotification(groupId string, title string, body string, imageUrl string, messageType string) {
	//Topic to send notification to all offline members of the group
	topic := "groupMessage" + groupId
	//Creating new message to send notification to
	newMessage := &messaging.Message{
		Topic: topic,
		Data: map[string]string{
			"title":                 title,
			"body":                  body,
			"imageUrl":              imageUrl,
			"type":                  messageType,
			"isMessageNotification": "true",
		},
	}

	_, sendErr := client.Send(context.TODO(), newMessage)

	if sendErr != nil {
		fmt.Println("Group notification sending error : ", sendErr)
		return
	}

}

// For sending call notification when someone calls
func SendCallNotification(
	token string,
	callType string,
	callerName string,
	callerId int,
	calleeId int,
	imageUrl string,
) {
	notificationId := uuid.NewString()
	currentTime := time.Now()
	//Creating payload for sending notification with details
	newMessage := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"callType":         callType,
			"type":             "call",
			"callerName":       callerName,
			"callerId":         strconv.Itoa(callerId),
			"currentUserId":    strconv.Itoa(calleeId),
			"imageUrl":         imageUrl,
			"id":               notificationId,
			"notificationTime": currentTime.Format(time.RFC3339Nano),
		},
	}

	fmt.Println(newMessage.Data)

	_, resErr := client.Send(context.TODO(), newMessage)

	if resErr != nil {
		fmt.Println("Send call notification error : ", resErr)
		return
	}

}

// For sending group call notification
func SendGroupCallNotification(groupId string, groupName string, groupProfilePic string, callType string) {
	//Creating the topic to send notification to all who subscribed to this topic
	topic := "groupCall" + groupId
	currentTime := time.Now()
	message := &messaging.Message{
		Topic: topic,
		Data: map[string]string{
			"title":            groupName,
			"imageUrl":         groupProfilePic,
			"groupId":          groupId,
			"callType":         callType,
			"type":             "groupCall",
			"notificationTime": currentTime.Format(time.RFC3339Nano),
		},
	}

	sendRes, sendErr := client.Send(context.TODO(), message)

	if sendErr != nil {
		fmt.Println("Group call notification error : ", sendErr)
		return
	}

	fmt.Println(sendRes)
}
