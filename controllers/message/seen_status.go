package message

import (
	"chitchat/config"
	"chitchat/helpers"
	messagemodel "chitchat/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// --------- FETCH SEEN STATUS CONTROLLER ------------
// To fetch seen status
func FetchSeenStatusController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	//Querying receiver id
	receiverIdStr := ctx.Query("receiverId")
	//Converting string receiver id into int receiver id
	receiverId, convErr := strconv.Atoi(receiverIdStr)

	if convErr != nil {
		fmt.Println("Could not convert id while fetching seen status : ", convErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if !exist {
		fmt.Println("current user id does not exist in while fetching seen status")
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Fetching read status
	collection := config.MongoDB.Collection("seenStatus")

	//Condition for fetching data
	condition := bson.M{
		"$and": []bson.M{
			{"senderId": currentUserId},
			{"receiverId": receiverId},
		},
	}
	var readModel messagemodel.ReadStatusModel

	fetchErr := collection.FindOne(context.TODO(), condition).Decode(&readModel)

	if fetchErr != nil {
		fmt.Println("Could not fetch seen status : ", fetchErr)
		helpers.SendMessageAsJson(
			ctx,
			"Something went wrong",
			http.StatusNoContent)
		return
	}

	fmt.Println("Seen data : ", readModel.IsSeen)
	fmt.Println("Receiver id : ", readModel.ReceiverId)
	fmt.Println("Sender id : ", readModel.SenderId)

	//Sending response
	ctx.JSON(http.StatusOK,
		gin.H{"message": "success",
			"senderId":   readModel.SenderId,
			"receiverId": readModel.ReceiverId,
			"isSeen":     readModel.IsSeen})
}

// -------------- DELETE SEEN STATUS CONTROLLER -----------
// For deleting seen status
func DeleteSeenStatusController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exists := ctx.Get("userId")

	if !exists {
		fmt.Println("Current user id does not exist in delete seen controller ")
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Querying receiver id
	receiverIdStr := ctx.Query("receiverId")
	//Converting receiverid string to receiver id int
	receiverId, convErr := strconv.Atoi(receiverIdStr)

	if convErr != nil {
		fmt.Println("Conv error in delete seen status : ", convErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Condition for deleting seen status
	condition := bson.M{
		"$and": []bson.M{
			{"senderId": currentUserId},
			{"receiverId": receiverId},
		},
	}

	collection := config.MongoDB.Collection("seenStatus")
	//Deleting it
	_, delErr := collection.DeleteMany(context.TODO(), condition)

	if delErr != nil {
		fmt.Println("Deletiong error in delete seen status : ", delErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	helpers.SendMessageAsJson(ctx, "Deleted successfully", http.StatusOK)
}

// ----------------- SAVE SEEN INDICATION -----------
// For saving seen info when current user sees receiver's message
func SaveSeenInfoController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	//Querying sender id
	senderId, convErr := strconv.Atoi(ctx.Query("senderId"))

	if !exist {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if convErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Saving the seen status
	collection := config.MongoDB.Collection("seenStatus")

	//Seen info
	seenInfo := bson.M{
		"senderId":   senderId,
		"receiverId": currentUserId,
		"isSeen":     true,
	}

	_, inseErr := collection.InsertOne(context.TODO(), seenInfo)

	if inseErr != nil {
		fmt.Println("Could not insert in seen status : ", inseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	helpers.SendMessageAsJson(ctx, "Saved successfully", http.StatusOK)

}
