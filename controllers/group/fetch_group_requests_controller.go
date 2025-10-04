package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/models"
	"context"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --------- FETCH GROUP REQUESTS CONTROLLER -----------
// For fetching group requests
func FetchGroupRequestsController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	//Querying group id
	groupId := ctx.Query("groupId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
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
	//Fetching all requests from groupRequests collection using current user id
	var requests []models.FetchGroupRequests
	groupReqCollec := config.MongoDB.Collection("groupRequests")
	//Condition to fetch requests
	condition := bson.M{
		"requestedGroupId": groupId,
		"groupAdminId":     currentUserId,
	}

	//Projection for fetching required data only
	projection := bson.M{
		"requestedUserId":    1,
		"requestedGroupName": 1,
		"requestedGroupId":   1,
		"_id":                0,
	}

	//Fetching it
	cursor, fetchingErr := groupReqCollec.
		Find(
			context.TODO(),
			condition,
			options.Find().SetProjection(projection),
			options.Find().SetSkip(int64(offset)),
			options.Find().SetLimit(int64(limit)),
		)

	if fetchingErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Parsing all data
	cursorErr := cursor.All(context.TODO(), &requests)

	//Getting each user details using the id of user
	for index, requestedUser := range requests {

		var userDetails models.UserModel

		userDetailsFetchErr := config.GORM.
			Table("userdetails").
			Select("username", "profilepic", "bio").
			Where("id=?", requestedUser.UserId).
			Scan(&userDetails).
			Error

		if userDetailsFetchErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}
		requests[index].Username = userDetails.Username
		requests[index].UserProfilePic = userDetails.ProfilePic
		requests[index].Userbio = userDetails.Userbio
	}

	if cursorErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
	} else if len(requests) == 0 {
		fmt.Println("True")
		empty := []string{}

		ctx.JSON(http.StatusOK, gin.H{"requests": empty})
	} else {
		fmt.Println("false")
		ctx.JSON(http.StatusOK, gin.H{"requests": requests})
	}

}
