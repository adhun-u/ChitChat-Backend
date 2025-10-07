package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ----------- FETCH ADDED USERS CONTROLLER -----------
// For fetching added users of a group
func FetchMembersController(ctx *gin.Context) {

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Querying group id
	groupId, convErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if convErr != nil {
		fmt.Println(convErr)
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusNotAcceptable)
		return
	}

	limit, limitIntConvErr := strconv.Atoi(ctx.Query("limit"))

	if limitIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide limit as integer", http.StatusNotAcceptable)
		return
	}

	page, pageIntConvErr := strconv.Atoi(ctx.Query("page"))

	if pageIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide page as integer", http.StatusNotAcceptable)
		return
	}

	offset := (page - 1) * limit

	if page == 1 {
		limit = limit - 1
	}

	var memberIds []models.GroupAddedUser
	var usersDetails []models.UserModel
	groupMembersCollec := config.MongoDB.Collection("groupMembers")

	//Condition
	condition := bson.M{
		"groupId": groupId.Hex(),
		"memberId": bson.M{
			"$ne": currentUserId,
		},
	}
	//Projection for getting added users
	projection := bson.M{
		"memberId": 1,
		"_id":      0,
	}

	if page == 1 {
		memberIds = append(memberIds, models.GroupAddedUser{UserId: int(currentUserId.(float64))})
	}

	//Fetching the other group members id
	fetchRes, fetchErr := groupMembersCollec.Find(
		context.TODO(),
		condition,
		options.Find().SetProjection(projection),
		options.Find().SetSkip(int64(offset)),
		options.Find().SetLimit(int64(limit)),
	)

	if fetchErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching members", http.StatusInternalServerError)
		return
	}

	defer fetchRes.Close(context.TODO())

	for fetchRes.Next(context.TODO()) {
		var memberId models.GroupAddedUser

		decodeErr := fetchRes.Decode(&memberId)
		if decodeErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while fetching members", http.StatusInternalServerError)
			return
		}

		memberIds = append(memberIds, memberId)

	}
	//For fetching each user's details from userdetails table
	for _, addedUser := range memberIds {
		var userDetails models.UserModel
		userDetailsFetchErr := config.
			GORM.
			Table("userdetails").
			Select("id", "username", "profilepic", "bio").
			Where("id=?", addedUser.UserId).
			Scan(&userDetails).Error

		if userDetailsFetchErr != nil {
			fmt.Println(userDetailsFetchErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}

		//Adding each user details to usersDetails
		usersDetails = append(usersDetails, userDetails)
	}
	if len(usersDetails) == 0 {
		empty := []string{}
		ctx.JSON(http.StatusOK, gin.H{"addedUsers": empty})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"addedUsers": usersDetails})
	}

}
