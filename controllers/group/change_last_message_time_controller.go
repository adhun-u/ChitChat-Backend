package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -------- CHANGE LAST MESSAGE TIME CONTROLLER ------------
// For changing last message time to sort as first
func ChangeGroupLastMessageTimeController(ctx *gin.Context) {

	//Parsing group id str to mongo db object id
	groupId, objIdConvErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if objIdConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusNotAcceptable)
		return
	}

	//Parsing time
	parsedTime, parseTimeErr := time.Parse(time.RFC3339Nano, ctx.Query("time"))

	if parseTimeErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid time format . Provide time as iso8601 string", http.StatusNotAcceptable)
		return
	}

	groupsCollec := config.MongoDB.Collection("groups")

	//Document for finding group
	docToFindGroup := bson.M{
		"_id": groupId,
	}
	//Document for changing time
	docToChangeTime := bson.M{
		"$set": bson.M{
			"lastMessageTime": parsedTime,
		},
	}

	timeUpdateRes, timeUpdateErr := groupsCollec.UpdateOne(context.TODO(), docToFindGroup, docToChangeTime)

	if timeUpdateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while updating time", http.StatusInternalServerError)
		return
	}

	if timeUpdateRes.ModifiedCount != 0 {
		helpers.SendMessageAsJson(ctx, "Updated successfully", http.StatusOK)

	} else {
		helpers.SendMessageAsJson(ctx, "Could not update time", http.StatusInternalServerError)

	}
}
