package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ----------- CHANGE LAST MESSAGE TIME CONTROLLER --------------
// For changing last message time to sort
func ChangeLastMessageController(ctx *gin.Context) {

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Querying the user id to get receiver
	userId, userIdIntConvErr := strconv.Atoi(ctx.Query("userId"))

	if userIdIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide opposite user id  as userId and integer", http.StatusInternalServerError)
		return
	}

	//Querying time
	strTime := ctx.Query("time")

	if strTime == "" {
		helpers.SendMessageAsJson(ctx, "Provide time as Iso8601 string", http.StatusNotAcceptable)
		return
	}

	time, timeConvErr := time.Parse(time.RFC3339Nano, strTime)

	if timeConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid time format . Provide time as Iso8601 string", http.StatusNotAcceptable)
		return
	}

	friendsCollec := config.MongoDB.Collection("friends")
	//To change time
	changeTimeData := bson.M{
		"$set": bson.M{
			"lastMessageTime": time,
		},
	}

	//Condition to change current user's friend's time
	conditionToChangeFromCurrentUser := bson.M{
		"currentUserId": currentUserId,
		"friendId":      userId,
	}
	//Condtion to change current user's time from friend's
	conditionToChangeFromFriend := bson.M{
		"currentUserId": userId,
		"friendId":      currentUserId,
	}

	//Filter to change two docs
	filterToChangeTwoDocs := bson.M{
		"$or": []bson.M{
			conditionToChangeFromCurrentUser,
			conditionToChangeFromFriend,
		},
	}

	timeUpdateRes, timeUpdateErr := friendsCollec.UpdateMany(context.TODO(), filterToChangeTwoDocs, changeTimeData)

	if timeUpdateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while updating time", http.StatusInternalServerError)
		return
	}

	if timeUpdateRes.ModifiedCount >= 1 {
		helpers.SendMessageAsJson(ctx, "Time updated successfully", http.StatusOK)
	} else {
		helpers.SendMessageAsJson(ctx, "Could not update time because of invalid user id", http.StatusInternalServerError)

	}

}
