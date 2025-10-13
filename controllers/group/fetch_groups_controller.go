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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --------- FETCH GROUP CONTROLLER ---------
func FetchGroupController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exists := ctx.Get("userId")

	if !exists {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
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

	groupCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")
	//Condition to find that how many groups does current user joined or created
	conditionToFindMember := bson.M{
		"memberId": currentUserId,
	}

	var groupIds []models.FetchGroupIdModel
	var groups []models.FetchGroupModel

	//Projection for groupMembers collection
	projectionForGroupMembers := bson.M{
		"_id":      0,
		"memberId": 0,
	}

	//Projection for groups collection
	projectionForGroup := bson.M{
		"imagePublicId": 0,
	}

	groupIdsFetchRes, groupIdsFetchErr := groupMembersCollec.
		Find(
			context.TODO(),
			conditionToFindMember,
			options.Find().SetProjection(projectionForGroupMembers),
			options.Find().SetSort(bson.D{{Key: "lastMessageTime", Value: -1}}),
			options.Find().SetLimit(int64(limit)),
			options.Find().SetSkip(int64(offset)),
		)

	if groupIdsFetchErr != nil {
		fmt.Println("Group ids fetch error : ", groupIdsFetchErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching groups", http.StatusInternalServerError)
		return
	}

	groupIdsCursorErr := groupIdsFetchRes.All(context.TODO(), &groupIds)

	if groupIdsCursorErr != nil {
		fmt.Println("Group ids cursor  error : ", groupIdsCursorErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching groups", http.StatusInternalServerError)
		return
	}

	for _, groupIdModel := range groupIds {

		var groupDetails models.FetchGroupModel
		groupId, objIdConvErr := primitive.ObjectIDFromHex(groupIdModel.GroupId)

		fmt.Println(groupId)

		if objIdConvErr != nil {
			fmt.Println("Object id conversion error : ", objIdConvErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while fetching group", http.StatusInternalServerError)
			return
		}
		//Condition to find group
		conditionToFindGroup := bson.M{
			"_id": groupId,
		}

		fetchGroupRes := groupCollec.FindOne(context.TODO(), conditionToFindGroup, options.FindOne().SetProjection(projectionForGroup))

		cursorErr := fetchGroupRes.Decode(&groupDetails)

		if cursorErr != mongo.ErrNoDocuments && cursorErr != nil {

			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}
		groups = append(groups, groupDetails)
	}

	if len(groups) == 0 {
		empty := []string{}
		fmt.Println("Empty groups")
		ctx.JSON(http.StatusOK, gin.H{"groups": empty})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"groups": groups})
	}

}
