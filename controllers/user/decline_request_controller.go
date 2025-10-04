package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// -------------- DECLINE REQUEST CONTROLLER --------------
// For declining a request
func DeclineRequestController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	//Querying the id of users who is going to declined
	declinedUserId, intConvErr := strconv.Atoi(ctx.Query("userId"))

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if !exist {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	reqCollec := config.MongoDB.Collection("requestedUsers")
	sentCollec := config.MongoDB.Collection("sentUsers")

	//Condition to delete from both requestedUsers and sentUsers collections
	conditionToDelete := bson.M{
		"$and": []bson.M{
			{
				"sentUserId": currentUserId,
			}, {
				"requestedUserId": declinedUserId,
			},
		},
	}
	//Deleting from requestedUsers collection
	delReqResult, delReqErr := reqCollec.DeleteOne(context.TODO(), conditionToDelete)
	//Deleting from sentUsers collection
	delSentResult, delSentErr := sentCollec.DeleteOne(context.TODO(), conditionToDelete)

	if delSentErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if delReqErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Println(delReqResult.DeletedCount)
	fmt.Println(delSentResult.DeletedCount)
	if delReqResult.DeletedCount != 0 && delSentResult.DeletedCount != 0 {
		fmt.Println("True")
		helpers.SendMessageAsJson(ctx, "Declined successfully", http.StatusOK)
		return
	} else {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)

	}

}
