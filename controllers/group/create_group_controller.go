package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --------- CREATE GROUP CONTROLLER -----------
// For creating a group
func CreateGroupController(ctx *gin.Context) {

	//Getting group name as post form
	groupName := ctx.PostForm("groupName")
	//Getting group image as form file
	groupImage, _ := ctx.FormFile("groupImage")
	//Getting group bio as post form
	groupBio := ctx.PostForm("groupBio")
	//Getting the time when this group is created as post form
	createdAt, timeParsErr := time.Parse(time.RFC3339Nano, ctx.PostForm("createdAt"))

	if timeParsErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid time format . Provide time as iso8601 string", http.StatusNotAcceptable)
		return
	}
	//Parsing group adming id as integer
	groupAdminUserId, idErr := strconv.Atoi(ctx.PostForm("groupAdminUserId"))

	if idErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide groupAdminUserId as integer in form data", http.StatusNotAcceptable)
		return
	}
	//Uploading the image to cloudinary to get a url
	imageUrl, publicId := services.UploadFileInFormData(groupImage)
	groupCollection := config.MongoDB.Collection("groups")
	groupMembersCollec := config.MongoDB.Collection("groupMembers")

	//Data to insert in groups collection
	groupData := bson.M{
		"groupName":         groupName,
		"groupBio":          groupBio,
		"groupImage":        imageUrl,
		"groupAdminUserId":  groupAdminUserId,
		"imagePublicId":     publicId,
		"createdAt":         createdAt,
		"lastMessageTime":   createdAt,
		"groupMembersCount": 1,
	}

	//Adding to database
	insertRes, insertErr := groupCollection.InsertOne(context.TODO(), groupData)

	if insertErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while creating group", http.StatusInternalServerError)
	} else {
		//Parsing the group id from the result
		insertedId := insertRes.InsertedID.(primitive.ObjectID)
		groupId := insertedId.Hex()

		//Data to insert in groupMembers collection (also adding the user who created this group)
		groupMemberData := bson.M{
			"groupId":  groupId,
			"memberId": groupAdminUserId,
		}

		_, groupMemberInseErr := groupMembersCollec.InsertOne(context.TODO(), groupMemberData)

		if groupMemberInseErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while adding admin", http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Group created successfully", "groupId": insertRes.InsertedID.(primitive.ObjectID).Hex()})
		//Subscribing group admin to a fcm topic to get group messages and call notifications
		var deviceToken string

		findErr := config.
			GORM.
			Table("userdetails").
			Select("deviceToken").
			Where("id=?", groupAdminUserId).
			Scan(&deviceToken).Error

		if findErr != nil {
			fmt.Println("Device token finding error while creating a group : ", findErr)
			return
		}

		services.SubscribeUserToGroupMessageTopic(deviceToken, groupId)
		services.SubscribeToGroupCallTopic(deviceToken, groupId)
	}

}
