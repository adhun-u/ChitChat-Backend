package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---------- FETCH Friends WITH LAST MESSAGE CONTROLLER -----------------
// To fetch friends of current user with last message
func FetchFriendsWithLastMessageController(ctx *gin.Context) {
	//Getting the current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
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
	friendsCollec := config.MongoDB.Collection("friends")

	var friendIds []usermodel.FriendModel

	//Condition for fetching current user's friends
	conditionToFetchFriends := bson.M{
		"currentUserId": currentUserId,
	}

	//Projection
	projection := bson.M{
		"currentUserId": 0,
		"_id":           0,
	}

	//Fetching current user's friends
	friendsRes, friendsErr := friendsCollec.
		Find(
			context.TODO(),
			conditionToFetchFriends,
			options.Find().SetProjection(projection),
			options.Find().SetSort(bson.D{{Key: "lastMessageTime", Value: -1}}),
			options.Find().SetLimit(int64(limit)),
			options.Find().SetSkip(int64(offset)),
		)

	if friendsErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while friends", http.StatusInternalServerError)
		return
	}

	cursorErr := friendsRes.All(context.TODO(), &friendIds)

	if cursorErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching friends", http.StatusInternalServerError)
		return
	}

	defer friendsRes.Close(context.TODO())

	var friendsDetails []usermodel.GetAddedUserDetails
	for index, eachAddedUser := range friendIds {
		//Fetching each user's details from userdetails using the user's id
		var friendDetails usermodel.GetAddedUserDetails
		//Getting each user's details using the user id
		config.GORM.
			Table("userdetails").
			Select("id,username,profilepic", "bio").
			Where("id=?", eachAddedUser.FriendId).
			Scan(&friendDetails)
			//Getting pending last message and pending message count
		tempMessageCollection := config.MongoDB.Collection("tempMessages")
		//A condition to find any pending message of current user with this user
		countCondition := bson.M{
			"$and": []bson.M{
				{"senderId": eachAddedUser.FriendId},
				{"receiverId": currentUserId},
			},
		}
		//Counting how many documents are satisfied with this condition
		count, _ := tempMessageCollection.CountDocuments(context.TODO(), countCondition)

		if count != 0 {
			//Then fetching last message the user with current user
			//Condition for fetching last message
			lastMessageCondition := bson.M{
				"$and": []bson.M{
					{"senderId": eachAddedUser.FriendId},
					{"receiverId": currentUserId},
					{"$and": []bson.M{
						{"$ne": bson.M{
							"type": "audioCall",
						}}, {
							"$ne": bson.M{
								"type": "videoCall",
							},
						},
					}},
				},
			}
			//Selecting single field
			projection := bson.M{
				"_id":  0,
				"time": 0,
			}
			//Fetching
			tempMessageCollection.
				FindOne(
					context.TODO(),
					lastMessageCondition,
					options.FindOne().SetProjection(projection),
				).Decode(&friendDetails.LastPendingMessage)
		}
		friendDetails.PendingMessagesCount = int(count)
		friendDetails.LastPendingMessage.Time = friendIds[index].LastMessageTime
		//Adding to the friendsDetails
		friendsDetails = append(friendsDetails, friendDetails)
	}
	//Checking if the getAddedUsers is empty
	if len(friendsDetails) == 0 {
		empty := []string{}
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message":    "Success",
				"addedUsers": empty,
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":    "Success",
			"addedUsers": friendsDetails,
		},
	)

}
