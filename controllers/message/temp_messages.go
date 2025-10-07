package message

import (
	"chitchat/config"
	"chitchat/helpers"
	messagemodel "chitchat/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// --------- FETCH TEMPORARY MESSAGE CONTROLLER ------------------
// For fetching temporary messages when receiver is not in online
func FetchTempMessagesController(ctx *gin.Context) {
	contxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	//Getting curren userId from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("User id does not exist in temp messages")
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Then fetching all temporary messages using current user id as receiverId
	collection := config.MongoDB.Collection("tempMessages")
	//Condition
	filter := bson.M{"receiverId": currentUserId}

	cursor, fetchErr := collection.Find(contxt, filter)

	if fetchErr != nil {
		fmt.Println("Temp messages fetching error in temp message : ", fetchErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var messages []messagemodel.MessageModel
	//Adding all temporary messages to array of messages
	addingErr := cursor.All(contxt, &messages)

	if addingErr != nil {
		fmt.Println("Adding data error in temp message : ", addingErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if len(messages) == 0 {
		empty := []string{}
		ctx.JSON(http.StatusOK, gin.H{"messages": empty})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"messages": messages})

}

// ------- DELETE TEMPORARY MESSAGES CONTROLLER ----------------
// For deleting temporary messages when user get all messages
func DeleteTempMessageController(ctx *gin.Context) {
	contxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("Current user id does not exist in delete/tempMessages")
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusInternalServerError)
		return
	}

	//Condition for deleting temporary messages
	delCondition := bson.M{"receiverId": currentUserId}

	collection := config.MongoDB.Collection("tempMessages")

	//Deleting all temporary messages using above condition
	_, delErr := collection.DeleteMany(contxt, delCondition)

	if delErr != nil {
		fmt.Println("Deletion error in delete/tempMessages : ", delErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while deleting", http.StatusInternalServerError)
		return
	}
	helpers.SendMessageAsJson(ctx, "Deleted successfully", http.StatusOK)
}

// --------- DELETE SINGLE MESSAGE CONTROLLER ------------
// For deleting single message from tempMessages collection
func DeleteSingleMessage(ctx *gin.Context) {
	fmt.Println("Entered")
	//Querying chat id from request
	chatId := ctx.Query("chatId")

	if chatId == "" {
		helpers.SendMessageAsJson(ctx, "Provide chatId", http.StatusNotAcceptable)
	}
	//Condition to delete the chat
	conditionToDelete := bson.M{
		"chatId": chatId,
	}

	fmt.Println("ChatId : ", chatId)
	tempMessagesCollec := config.MongoDB.Collection("tempMessages")

	delRes, delErr := tempMessagesCollec.DeleteOne(context.TODO(), conditionToDelete)

	if delErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while deleting", http.StatusInternalServerError)
		return
	}

	if delRes.DeletedCount != 0 {
		helpers.SendMessageAsJson(ctx, "Deleted successfully", http.StatusOK)
	} else {
		helpers.SendMessageAsJson(ctx, "Chat not found", http.StatusNotFound)
	}
}
