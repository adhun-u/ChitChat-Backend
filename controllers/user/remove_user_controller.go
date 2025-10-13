package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ----------- REMOVE USER CONTROLLER -------------
// For removing a user from friends list of current user
func RemoveUserController(ctx *gin.Context) {

	//Getting current user id from middlewares
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Querying user's id who is going to be removed from friends
	removeUserId, intConvErr := strconv.Atoi(ctx.Query("removeUserId"))

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide removeUserId as integer", http.StatusNotAcceptable)
		return
	}

	friendsCollec := config.MongoDB.Collection("friends")

	//Document to delete among current user's friends
	deleteFromCurrentUser := bson.M{
		"currentUserId": currentUserId,
		"friendId":      removeUserId,
	}

	//Document to delete among the user's friends who is going to removed
	deleteFromRemovedUsr := bson.M{
		"currentUserId": removeUserId,
		"friendId":      currentUserId,
	}

	//Filter to remove these two docs
	filter := bson.M{
		"$or": bson.A{
			deleteFromCurrentUser,
			deleteFromRemovedUsr,
		},
	}

	//Removing the user from the collection
	removeRes, removeErr := friendsCollec.DeleteMany(context.TODO(), filter)

	if removeErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while removing", http.StatusInternalServerError)
		return
	}

	if removeRes.DeletedCount == 2 {
		helpers.SendMessageAsJson(ctx, "Removed successfully", http.StatusOK)
		//Getting current user device token as well as removed user's
		var deviceTokens []string
		var userIds = []int{
			int(currentUserId.(float64)), removeUserId,
		}
		findErr := config.
			GORM.
			Table("userdetails").
			Select("deviceToken").
			Where("id IN ?", userIds).
			Scan(&deviceTokens).Error

		if findErr != nil {
			fmt.Println("Device token finding error in remove user controller : ", findErr)
			return
		}

		//Unsubscribing the each to not get message notifications

		services.UnSubscribeFromUserMessageTopic(deviceTokens[0], removeUserId, int(currentUserId.(float64)))

		services.UnSubscribeFromUserMessageTopic(deviceTokens[1], int(currentUserId.(float64)), removeUserId)
	} else {
		helpers.SendMessageAsJson(ctx, "Invalid user", http.StatusNotFound)
	}

}
