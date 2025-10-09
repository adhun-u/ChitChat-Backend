package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ----------- ACCEPT REQUEST CONTROLLER ----------------
// For accepting the request that current user got
func AcceptRequest(ctx *gin.Context) {
	var acceptRequestModel usermodel.AcceptRequestModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&acceptRequestModel)

	if parseErr != nil {
		fmt.Println("Parse error in user/request/accept : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Invalid json format", http.StatusNotAcceptable)
		return
	}
	//Getting the current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("Current user id does not exist")
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
	}

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

	var currentUser usermodel.UserModel
	var requestedUser usermodel.UserModel

	//Retrieving the current user details
	currentFetchErr := config.GORM.
		Table("userdetails").
		Select("id", "username", "profilepic", "bio", "deviceToken").
		Where("id=?", currentUserId).
		First(&currentUser).Error
	//Retrieving the requested user details
	oppositeUserFetchErr := config.GORM.
		Table("userdetails").
		Select("id", "username", "profilepic", "bio", "deviceToken").
		Where("id=?", acceptRequestModel.UserId).
		First(&requestedUser).Error

	if oppositeUserFetchErr != nil || currentFetchErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Println("Current user id and token : ", currentUser.UserId, currentUser.DeviceToken)
	fmt.Println("Opposite user id and token : ", requestedUser.UserId, requestedUser.DeviceToken)
	friendsCollec := config.MongoDB.Collection("friends")
	requestUsersCollec := config.MongoDB.Collection("requestedUsers")

	currentUserDoc := bson.M{
		"currentUserId":   currentUserId,
		"friendId":        acceptRequestModel.UserId,
		"lastMessageTime": time,
	}

	acceptedUserDoc := bson.M{
		"currentUserId":   acceptRequestModel.UserId,
		"friendId":        currentUserId,
		"lastMessageTime": time,
	}

	_, insertFriendErr := friendsCollec.InsertMany(
		context.TODO(),
		[]any{
			currentUserDoc,
			acceptedUserDoc,
		},
	)

	if insertFriendErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while accepting request", http.StatusInternalServerError)
		return
	}

	//Getting the sentUsers collection for deletion
	sentUsersCollection := config.MongoDB.Collection("sentUsers")
	//Condition to delete document from requestedUsers and sentUsers collection
	deleteFilter := bson.M{"$or": []bson.M{
		{"requestedUserId": requestedUser.UserId},
		{"sentUserId": currentUser.UserId}}}
	//Deleting the doc from requestedUsers collection
	_, reqDeleResultErr := requestUsersCollec.DeleteMany(context.TODO(), deleteFilter)

	if reqDeleResultErr != nil {
		fmt.Println("Request deleting error in user/request/accept : ", reqDeleResultErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while operating some operations", http.StatusInternalServerError)
		return
	}

	//Deleting the doc from sentUsers collection
	_, sentDeletErr := sentUsersCollection.DeleteMany(context.TODO(), deleteFilter)

	if sentDeletErr != nil {
		fmt.Println("sent deleting error in user/request/accep : ", sentDeletErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while operating some operations", http.StatusInternalServerError)
		return
	}
	//Sending notification to the user who has been accepted
	services.SendNotification(requestedUser.DeviceToken, currentUser.Username, "Accepted your request", currentUser.ProfilePic)
	ctx.JSON(http.StatusOK, gin.H{"message": "Accepted", "requestedUserId": requestedUser.UserId})

	//Subscribing accepted user to call and message topic to get notifications
	services.SubscribeToUserMessageTopic(requestedUser.DeviceToken, currentUser.UserId, requestedUser.UserId)
	services.SubscribeToUserCallTopic(requestedUser.DeviceToken, currentUser.UserId, requestedUser.UserId)
	//Subscribing current user to accepted user's topic to get call and message notifications
	services.SubscribeToUserMessageTopic(currentUser.DeviceToken, requestedUser.UserId, currentUser.UserId)
	services.SubscribeToUserCallTopic(currentUser.DeviceToken, requestedUser.UserId, currentUser.UserId)
}
