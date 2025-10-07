package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ----------- ACCEPT GROUP REQUEST CONTROLLER ------------
// For accepting group requests
func AcceptGroupRequestController(ctx *gin.Context) {
	//Querying the id of the user who has to be accepted
	acceptedUserId, intConvErr := strconv.Atoi(ctx.Query("userId"))
	//Querying group id
	groupId, objIdConErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))
	//Querying group name
	groupName := ctx.Query("groupName")
	//Querying group image
	groupImage := ctx.Query("groupImage")
	//Querying time
	acceptedTimeStr := ctx.Query("time")

	if acceptedTimeStr == "" {
		helpers.SendMessageAsJson(ctx, "Provide time as iso8601 string", http.StatusNotAcceptable)
		return
	}

	parsedTime, timeParseErr := time.Parse(time.RFC3339Nano, acceptedTimeStr)

	if timeParseErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid time . Provide time as iso8601 string", http.StatusNotAcceptable)
		return
	}

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide userId as integer", http.StatusNotAcceptable)
		return
	}

	if objIdConErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusNotAcceptable)
		return
	}

	groupsCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")
	groupRequeCollec := config.MongoDB.Collection("groupRequests")

	//Inserting the accepted user to groupMembers collection
	//Data to insert
	insertData := bson.M{
		"groupId":  groupId.Hex(),
		"memberId": acceptedUserId,
	}

	//condition to find the group to increment group members count
	conditionToFindGroup := bson.M{
		"_id": groupId,
	}

	//For changing some fields
	incrementMembersCount := bson.M{
		"$inc": bson.M{
			"groupMembersCount": 1,
		},
		"$set": bson.M{
			"lastMessageTime": parsedTime,
		},
	}
	_, insertErr := groupMembersCollec.InsertOne(context.TODO(), insertData)

	if insertErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while accepting request", http.StatusInternalServerError)
		return
	}

	updateRes, updateErr := groupsCollec.UpdateOne(context.TODO(), conditionToFindGroup, incrementMembersCount)

	if updateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while accepting request", http.StatusInternalServerError)
		return
	}

	if updateRes.ModifiedCount == 0 {
		helpers.SendMessageAsJson(ctx, "Invalid group", http.StatusNotFound)
		return
	}
	//Then deleting the accepted user from groupRequests collection
	//Data to delete
	deleteData := bson.M{
		"requestedUserId": acceptedUserId,
	}

	deleteResult := groupRequeCollec.FindOneAndDelete(context.TODO(), deleteData)

	if deleteResult.Err() != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while accepting requests", http.StatusInternalServerError)
		return
	}

	helpers.SendMessageAsJson(ctx, "Accepted request", http.StatusOK)
	var acceptedDeviceToken string
	//Getting accepted user device token to send notification
	fetchErr := config.
		GORM.Table("userdetails").
		Select("deviceToken").
		Where("id=?", acceptedUserId).
		Scan(&acceptedDeviceToken).Error

	if fetchErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Sending notification
	services.SendNotification(acceptedDeviceToken, groupName, "Accepted your request", groupImage)

	//Subscribing accepted user to get group messages as notification
	var deviceToken string

	findErr := config.
		GORM.
		Table("userdetails").
		Select("deviceToken").
		Where("id=?", acceptedUserId).
		Scan(&deviceToken).Error

	if findErr != nil {
		fmt.Println("Device token finding error while subscribing")
		return
	}

	//Subscribing the accepted user to a this group topic to get calls and messages as notification
	services.SubscribeUserToGroupMessageTopic(deviceToken, ctx.Query("groupId"))
	services.SubscribeToGroupCallTopic(deviceToken, ctx.Query("groupId"))
}
