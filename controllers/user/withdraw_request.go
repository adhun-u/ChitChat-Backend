package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ------- WITHDRAW REQUEST CONTROLLER --------------
// To withdraw a request after requested
func WithdrawRequestController(ctx *gin.Context) {

	//Getting a user's id who is going to be withdrawn
	userIdStr := ctx.Query("userId")
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}
	//Converting string userid into int user id
	userId, convErr := strconv.Atoi(userIdStr)

	if convErr != nil {
		helpers.SendMessageAsJson(ctx, "Providee userId as integer", http.StatusNotAcceptable)
		return
	}

	sentUserCollec := config.MongoDB.Collection("sentUsers")

	//Deleting the user from current user's sentusers list
	//Condition
	delSentUser := bson.M{
		"$and": []bson.M{
			{"requestedUserId": currentUserId},
			{"sentUserId": userId},
		},
	}
	//Deleting
	_, sentUserDeletionErr := sentUserCollec.DeleteOne(context.TODO(), delSentUser)

	if sentUserDeletionErr != nil {
		helpers.SendMessageAsJson(ctx, "Couldn't withdraw", http.StatusInternalServerError)
		return
	}
	requestedUserCollec := config.MongoDB.Collection("requestedUsers")

	//Deleting the user from current user's sentusers list
	//Condition
	delRequestedUser := bson.M{
		"$and": []bson.M{
			{"requestedUserId": currentUserId},
			{"sentUserId": userId},
		},
	}
	//Deleting
	_, requestedUserDelErr := requestedUserCollec.DeleteOne(context.TODO(), delRequestedUser)

	if requestedUserDelErr != nil {
		helpers.SendMessageAsJson(ctx, "Couldn't withdraw", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request withdrawn", "withdrawnUserId": userId})

}
