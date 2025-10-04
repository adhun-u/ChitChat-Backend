package group

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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------------- REMOVE MEMBER CONTROLLER ----------
// For removing a member from a group
func RemoveMemberController(ctx *gin.Context) {

	//Querying group id to identify where to delete
	groupId, objIdConvErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if objIdConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusNotAcceptable)
		return
	}

	//Querying user id to identify whom to delete
	userId, intConvErr := strconv.Atoi(ctx.Query("userId"))

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide userId as integer", http.StatusNotAcceptable)
		return
	}

	groupCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")
	//Condition to find group
	conditionToFindGroup := bson.M{
		"_id": groupId,
	}
	//For decrementing group members count
	decrementMembersCount := bson.M{
		"$inc": bson.M{
			"groupMembersCount": -1,
		},
	}
	//Condition to find member to delete
	conditionToDelete := bson.M{
		"memberId": userId,
	}

	delRes, delErr := groupMembersCollec.DeleteOne(context.TODO(), conditionToDelete)

	if delErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while deleting", http.StatusInternalServerError)
		return
	}

	if delRes.DeletedCount == 0 {
		helpers.SendMessageAsJson(ctx, "Member does not exist", http.StatusNotFound)
		return
	}

	updateRes, updateErr := groupCollec.UpdateOne(context.TODO(), conditionToFindGroup, decrementMembersCount)

	if updateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while removing", http.StatusInternalServerError)
		return
	}

	if updateRes.ModifiedCount != 0 {
		helpers.SendMessageAsJson(ctx, "Removed successfully", http.StatusOK)
		//Getting the removed user device token
		var removedUserDeviceToken string

		findErr := config.
			GORM.
			Table("userdetails").
			Select("deviceToken").
			Where("id=?", userId).
			Scan(&removedUserDeviceToken).Error

		if findErr != nil {
			fmt.Println("Device token finding error in 'Remove memeber controller'", findErr)
			return
		}

		//Unsubscribing the user from the group topic to don't get the group call and message notifications
		services.UnsubscribeFromGroupMessageTopic(removedUserDeviceToken, ctx.Query("groupId"))
		services.UnsubscribeFromGroupCallTopic(removedUserDeviceToken, ctx.Query("groupId"))
	}

}
