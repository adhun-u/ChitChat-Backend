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
	"go.mongodb.org/mongo-driver/mongo"
)

// ------------ FETCH USERS TO ADD MEMBER CONTROLLER -----------
// For fetching current user's friends for adding as group member
func FetchUserToAddMemberController(ctx *gin.Context) {

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}
	//Querying group id
	groupId := ctx.Query("groupId")

	if groupId == "" {
		helpers.SendMessageAsJson(ctx, "Provide groupId as string", http.StatusNotAcceptable)
		return
	}

	//Querying limit
	limit, limitIntConvErr := strconv.Atoi(ctx.Query("limit"))

	if limitIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide limit as integer", http.StatusNotAcceptable)
		return
	}

	//Querying page
	page, pageIntConvErr := strconv.Atoi(ctx.Query("page"))

	if pageIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide page as integer", http.StatusNotAcceptable)
		return
	}

	//Calculating offset
	offset := (page - 1) * limit

	var friends []models.UserModel

	friendsCollec := config.MongoDB.Collection("friends")

	pipeLine := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "currentUserId", Value: currentUserId},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "groupMembers"},
			{Key: "let", Value: bson.D{{Key: "fid", Value: "$friendId"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$groupId", groupId}},
					}},
				}}},
			}},
			{Key: "as", Value: "groupMember"},
		}}},

		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{
				{Key: "$not", Value: bson.D{
					{Key: "$in", Value: bson.A{"$friendId", "$groupMember.memberId"}},
				}},
			}},
		}}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "currentUserId", Value: 0},
			{Key: "groupMember", Value: 0},
		}}},

		bson.D{
			{Key: "$limit", Value: limit},
		},
		bson.D{
			{Key: "$skip", Value: offset},
		},
	}

	fetchFriendsRes, fetchFriendsErr := friendsCollec.Aggregate(context.TODO(), pipeLine)

	if fetchFriendsErr != nil {
		fmt.Println("fetching error : ", fetchFriendsErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching users", http.StatusInternalServerError)
		return
	}

	defer fetchFriendsRes.Close(context.TODO())

	//Fetching each user's details using their id
	for fetchFriendsRes.Next(context.TODO()) {

		var friendDetails models.FriendModel

		decodeErr := fetchFriendsRes.Decode(&friendDetails)

		fmt.Println(friendDetails)

		if decodeErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while fetching users", http.StatusInternalServerError)
			return
		}

		var userdetails models.UserModel

		userDetailsFetchErr := config.GORM.
			Table("userdetails").
			Select("id", "username", "profilepic", "bio").
			Where("id=?", friendDetails.FriendId).
			Scan(&userdetails).
			Error

		if userDetailsFetchErr != nil {
			fmt.Println("user details fetching error : ", userDetailsFetchErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while fetching users", http.StatusInternalServerError)
			return
		}
		friends = append(friends, userdetails)
	}

	if len(friends) == 0 {
		empty := []string{}

		ctx.JSON(http.StatusOK, gin.H{"addedUsers": empty})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"addedUsers": friends})
	}
}
