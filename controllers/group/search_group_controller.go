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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --------- SEARCH GROUP CONTROLLER ---------
// For searching groups
func SearchGroupController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}
	//Querying group name
	searchedGroupName := ctx.Query("groupName")
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

	offset := (page - 1) * limit
	groupCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")
	groupReqCollec := config.MongoDB.Collection("groupRequests")
	var groups []models.FetchGroupModel
	//Condition to fetch groups which is similar to searched group name
	condition := bson.M{
		"groupName": bson.M{
			"$regex":   searchedGroupName,
			"$options": "i",
		},
	}

	//Projection for fetching required only data
	projection := bson.M{
		"groupName":        1,
		"groupImage":       1,
		"groupAdminUserId": 1,
		"groupBio":         1,
		"_id":              1,
	}
	//Fetching it
	cursor, searchingErr := groupCollec.
		Find(context.TODO(), condition, options.Find().SetProjection(projection), options.Find().SetSkip(int64(offset)), options.Find().SetLimit(int64(limit)))

	if searchingErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Parsing the data
	cursorErr := cursor.All(context.TODO(), &groups)

	if cursorErr != nil {
		fmt.Println("Cursor error in search group : ", cursorErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//For checking if current user added is added in the group or requested
	for index, group := range groups {
		//Condition for finding whether the current user is added or not
		findCondition := bson.M{
			"groupId":  group.GroupId.Hex(),
			"memberId": currentUserId,
		}
		//Condition for finding whether the current user is requested or not
		isRequestedCondition := bson.M{
			"requestedUserId": currentUserId,
		}
		findingErr := groupMembersCollec.FindOne(context.TODO(), findCondition).Err()
		isRequestFindErr := groupReqCollec.FindOne(context.TODO(), isRequestedCondition).Err()

		groups[index].IsCurrentUserAdded = findingErr != mongo.ErrNoDocuments
		groups[index].IsRequestSent = isRequestFindErr != mongo.ErrNoDocuments
	}

	if len(groups) == 0 {
		empty := []string{}

		ctx.JSON(http.StatusOK, gin.H{"searchResult": empty})

	} else {
		ctx.JSON(http.StatusOK, gin.H{"searchResult": groups})
	}

}
