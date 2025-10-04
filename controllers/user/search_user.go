package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// -------------- SEARCH USER CONTROLLER --------------
func SearchUserController(ctx *gin.Context) {
	//Getting the current user id from middleware
	currentUserId, userIdExists := ctx.Get("userId")

	if !userIdExists {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}
	//Creating a context
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	//Querying username
	username := ctx.Query("username")
	//Querying to find how many data should bring
	limit, limitIntConvErr := strconv.Atoi(ctx.Query("limit"))
	//Querying page
	page, pageIntConvErr := strconv.Atoi(ctx.Query("page"))

	//Checking if the username is null
	if username == "" {
		empty := []string{}
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "users": empty})
		return
	}

	if limitIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide limit as integer", http.StatusNotAcceptable)
		return
	}
	if pageIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide page as integer", http.StatusNotAcceptable)
		return
	}
	var users []usermodel.SearchedUserModel

	//Offset to skip rows
	offset := (page - 1) * limit

	//Getting all users which are like the username
	findUsersErr := config.GORM.
		Table("userdetails").
		Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if findUsersErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching users", http.StatusInternalServerError)
		return
	}

	requestedUserCollec := config.MongoDB.Collection("requestedUsers")
	friendCollec := config.MongoDB.Collection("friends")

	//For checking if the current user requested any user
	for index, user := range users {
		//Counting if the current user requested any user
		reqDocCount, reqDocCollecErr := requestedUserCollec.
			CountDocuments(context, bson.M{
				"sentUserId":      user.Id,
				"requestedUserId": currentUserId,
			})

		if reqDocCollecErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while getting users details", http.StatusInternalServerError)
			return
		}

		//Counting to check if the user is current user's friend
		friendDocCount, friendDocErr := friendCollec.
			CountDocuments(context, bson.M{
				"$and": []bson.M{
					{"currentUserId": currentUserId},
					{"friendId": user.Id},
				},
			})

		if friendDocErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while getting users details", http.StatusInternalServerError)
			return
		}

		users[index].IsRequested = reqDocCount != 0
		users[index].IsAdded = friendDocCount != 0
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success", "users": users})
}
