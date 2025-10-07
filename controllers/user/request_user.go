package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/models"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ------------- REQUEST USER CONTROLLER ------------------
// For sending a request to be friend
func RequestUserController(ctx *gin.Context) {
	var requestModel models.AddRequestUserModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&requestModel)

	if parseErr != nil {
		fmt.Println("parse error in user/add : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Invalid json format", http.StatusNotAcceptable)
		return
	}
	//Getting the current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("Id does not exist")
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Fetching current user details using above id
	var currentUserDetails models.UserModel

	//Fetch data
	dbErr := config.GORM.
		Table("userdetails").
		Select("username", "profilepic", "bio").
		Where("id=?", currentUserId).
		Find(&currentUserDetails).Error

	if dbErr != nil {
		fmt.Println("Fetching error in user/request : ", dbErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	requestedUserCollection := config.MongoDB.Collection("requestedUsers")
	//Creating a context
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//Cancelling the context to free the memory
	defer cancel()

	//Inserting into requestedUsers collection
	_, requestResultErr := requestedUserCollection.
		InsertOne(context, bson.M{
			"sentUserId":              requestModel.RequestedUserId,
			"requestedUserId":         currentUserId,
			"requestedUsername":       currentUserDetails.Username,
			"requestedUserProfilePic": currentUserDetails.ProfilePic,
			"requestedUserbio":        currentUserDetails.Userbio,
			"requestedDate":           requestModel.RequestedDate,
		})

	if requestResultErr != nil {
		fmt.Println("Insertion error in user/add : ", requestResultErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	sentUsersCollection := config.MongoDB.Collection("sentUsers")

	//Inserting the data to sentUsers collection
	_, sentUserResultErr := sentUsersCollection.InsertOne(
		context, bson.M{
			"sentUserId":         requestModel.RequestedUserId,
			"sentUsername":       requestModel.RequestedUsername,
			"sentUserbio":        requestModel.RequestedUserbio,
			"sentUserProfilePic": requestModel.RequestedUserProfilePic,
			"requestedUserId":    currentUserId,
			"sentDate":           requestModel.RequestedDate,
		})

	if sentUserResultErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Fetching FCM token from userdetails
	var requestedUserDeviceToken string
	config.GORM.
		Table("userdetails").
		Select("deviceToken").
		Where("id=?", requestModel.RequestedUserId).
		Scan(&requestedUserDeviceToken)
	//Sending notification to requested user
	services.SendNotification(
		requestedUserDeviceToken,
		currentUserDetails.Username,
		"Sent you a request",
		currentUserDetails.ProfilePic,
	)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Request is sent",
		"isSent":     true,
		"sentUserId": requestModel.RequestedUserId,
	})

}
