package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --------------- EXIT GROUP CONTROLLER -------------------
// For leaving from a group
func ExitGroupController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")
	//Querying group id
	groupId, objConvErr := primitive.ObjectIDFromHex(ctx.Query("groupId"))

	if objConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide valid groupId", http.StatusInternalServerError)
		return
	}

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusInternalServerError)
		return
	}

	//Querying current group members count
	currentGroupMembersCount, memberCountIntConvErr := strconv.Atoi(ctx.Query("membersCount"))

	if memberCountIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide membersCount as integer", http.StatusNotAcceptable)
		return
	}

	groupCollec := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")
	//Condition to decrement groupMembers count
	conditionToDecrement := bson.M{
		"$inc": bson.M{
			"groupMembersCount": -1,
		},
	}
	//Condition to find the group
	conditionToFind := bson.M{
		"_id": groupId,
	}

	//User to remove from groupMembers collection
	dataToRemove := bson.M{
		"memberId": currentUserId,
	}

	//Removing current user from the group
	updateRes, updateErr := groupCollec.UpdateOne(context.TODO(), conditionToFind, conditionToDecrement)

	if updateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while existing", http.StatusInternalServerError)
		return
	}
	if updateRes.ModifiedCount != 0 {

		//Deleting current user from groupMembers collection
		delRes, delErr := groupMembersCollec.DeleteOne(context.TODO(), dataToRemove)
		if delErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while existing", http.StatusInternalServerError)
			return
		}

		if delRes.DeletedCount != 0 {

			helpers.SendMessageAsJson(ctx, "Exited successfully", http.StatusOK)
			//Unsubscribing current user from the group topic to not getting the notification
			var deviceToken string

			findErr := config.
				GORM.
				Table("userdetails").
				Select("deviceToken").
				Where("id=?", currentUserId).
				Scan(&deviceToken).Error

			if findErr != nil {
				fmt.Println("Device token finding error while unsubscribing : ", findErr)
				return
			}

			services.UnsubscribeFromGroupMessageTopic(deviceToken, ctx.Query("groupId"))
			services.UnsubscribeFromGroupCallTopic(deviceToken, ctx.Query("groupId"))

			if currentGroupMembersCount-1 == 0 {

				docToDelete := bson.M{
					"_id": groupId,
				}
				//Then deleting the group from groups collection
				_, groupDelErr := groupCollec.DeleteOne(context.TODO(), docToDelete)

				if groupDelErr != nil {
					fmt.Println("Group deletion error : ", groupDelErr)
					return
				}

			}

		}

	}

}
