package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------- DECLINE GROUP REQUEST CONTROLLER ---------
func DeclineGroupRequestController(ctx *gin.Context) {

	//Querying the user id who is gonna be deleted
	userId, intConvErr := strconv.Atoi(ctx.Query("userId"))
	//Querying group id to identify where to delete
	groupId, objIdConvErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if objIdConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusOK)
		return
	}

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	groupReqCollec := config.MongoDB.Collection("groupRequests")
	//Condition for deleting the user from groupRequests collection
	conditionToDelete := bson.M{
		"$and": []bson.M{
			{
				"requestedUserId": userId,
			}, {
				"requestedGroupId": groupId,
			},
		},
	}

	delRes, delErr := groupReqCollec.DeleteOne(context.TODO(), conditionToDelete)

	if delErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if delRes.DeletedCount != 0 {
		helpers.SendMessageAsJson(ctx, "Declined successfully", http.StatusOK)
	}

}
