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

// ---------- ADD MEMBER CONTROLLER ------------
// For adding new person to a group
func AddMemberController(ctx *gin.Context) {
	//Querying group id
	groupId, objConvErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if objConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusNotAcceptable)
		return
	}

	//Querying user id of the user who is gonna be added
	userId, intConvErr := strconv.Atoi(ctx.Query("userId"))

	if intConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide userId as integer", http.StatusNotAcceptable)
		return
	}

	groupCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")

	//Data to insert in groupMembers collection
	dataToInsert := bson.M{
		"groupId":  groupId.Hex(),
		"memberId": userId,
	}
	//Condition to find the group
	conditionToFindGroup := bson.M{
		"_id": groupId,
	}

	//Condition to increment group members count
	incrementData := bson.M{
		"$inc": bson.M{
			"groupMembersCount": 1,
		},
	}

	_, insertGroupMemberErr := groupMembersCollec.InsertOne(context.TODO(), dataToInsert)

	if insertGroupMemberErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while adding", http.StatusInternalServerError)
		return
	}

	updateResult, updateErr := groupCollec.UpdateOne(context.TODO(), conditionToFindGroup, incrementData)

	if updateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while adding", http.StatusInternalServerError)
		return
	}

	if updateResult.ModifiedCount != 0 {
		helpers.SendMessageAsJson(ctx, "Added successfully", http.StatusOK)

		//Subscribing the added user to get group messages and group call as notifications
		//Getting device token of the user
		var deviceToken string
		findErr := config.GORM.
			Table("userdetails").
			Select("deviceToken").
			Where("id=?", userId).
			Scan(&deviceToken).Error

		if findErr != nil {
			fmt.Println("Finding user device token error while adding : ", findErr)
			return
		}
		fmt.Println("Device token : ", deviceToken)
		fmt.Println(ctx.Query("groupId"))
		services.SubscribeUserToGroupMessageTopic(deviceToken, ctx.Query("groupId"))
		services.SubscribeToGroupCallTopic(deviceToken, ctx.Query("groupId"))
	}

}
