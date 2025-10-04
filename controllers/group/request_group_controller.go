package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/models"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ----------- REQUEST GROUP CONTROLLER ----------------
// For sending request to admin when a user sent request
func RequestGroupController(ctx *gin.Context) {

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var groupDetails models.RequestGroupModel
	//Parsing the group details from request
	parseErr := ctx.ShouldBind(&groupDetails)

	if parseErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Getting current user details
	var currentUser models.CurrentUserModel
	fetchUserDetailsErr := config.GORM.
		Table("userdetails").
		Select("username", "profilepic").
		Where("id=?", currentUserId).
		Scan(&currentUser).Error

	if fetchUserDetailsErr != nil {
		fmt.Println("Fetching current user details error : ", fetchUserDetailsErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Then inserting to the groupRequests collection
	groupReqCollec := config.MongoDB.Collection("groupRequests")
	//Data to insert
	requestedUserDetails := bson.M{
		"requestedUserId":    currentUserId,
		"requestedGroupName": groupDetails.GroupName,
		"requestedGroupId":   groupDetails.GroupId,
		"groupAdminId":       groupDetails.AdminId,
	}
	//Inserting the data
	_, insertErr := groupReqCollec.InsertOne(context.TODO(), requestedUserDetails)

	if insertErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Getting device token of admin
	var adminDeviceToken string
	deviceTokenFetchErr := config.GORM.Table("userdetails").
		Select("deviceToken").
		Where("id=?", groupDetails.AdminId).
		Scan(&adminDeviceToken).Error

	//Sending a notification to admin that the current user is requested to join
	if deviceTokenFetchErr == nil {
		services.
			SendNotification(
				adminDeviceToken,
				currentUser.Username,
				fmt.Sprintf("wants to join %s", groupDetails.GroupName),
				currentUser.ProfilePic,
			)
	} else {
		fmt.Println("Fetching device token error : ", deviceTokenFetchErr)
	}

	ctx.JSON(http.StatusOK, gin.H{"groupId": groupDetails.GroupId})
}
