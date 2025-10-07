package group

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------- EDIT GROUP CONTROLLER -----------
// For changing group image , group name or group bio
func EditGroupController(ctx *gin.Context) {

	//Group id
	groupId, convErr := primitive.ObjectIDFromHex(ctx.PostForm("groupId"))

	if convErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide group id as string", http.StatusNotAcceptable)
		return
	}

	//Group name
	groupName := ctx.PostForm("groupName")

	//Group image
	groupImage, _ := ctx.FormFile("groupImage")

	//Group bio
	groupBio := ctx.PostForm("groupBio")

	groupCollec := config.MongoDB.Collection("groups")

	//Condition to find the group
	conditionToFindGroup := bson.M{
		"_id": groupId,
	}

	//Changing name of the group if the name is not empty
	if groupName != "" {

		//For updating name
		updateNameDoc := bson.M{
			"$set": bson.M{
				"groupName": groupName,
			},
		}

		updateErr := groupCollec.FindOneAndUpdate(context.TODO(), conditionToFindGroup, updateNameDoc).Err()

		if updateErr != nil {
			fmt.Println("Name updation error : ", updateErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while updating group details", http.StatusInternalServerError)
			return
		}
	}

	//Changing bio of the group if the bio is not empty
	if groupBio != "" {

		//For updating bio
		updateBioDoc := bson.M{
			"$set": bson.M{
				"groupBio": groupBio,
			},
		}
		updateErr := groupCollec.FindOneAndUpdate(context.TODO(), conditionToFindGroup, updateBioDoc).Err()

		if updateErr != nil {
			fmt.Println("Bio updation error : ", updateErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while updating group details", http.StatusInternalServerError)
		}
	}

	//Changing image of the group if image is not null
	var groupImageUrl string
	if groupImage != nil {

		//Getting previous public id of the group image to delete the source from cloudinary
		var prevGroupData bson.M
		//Projection for getting publicId only
		projection := bson.M{
			"imagePublicId": 1,
			"_id":           0,
		}
		findingErr := groupCollec.FindOne(context.TODO(), conditionToFindGroup, options.FindOne().SetProjection(projection)).Decode(&prevGroupData)

		if findingErr != nil {
			fmt.Println("Prev public id finding error : ", findingErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while updating group details", http.StatusInternalServerError)
		}

		//Getting the prev id from document
		prevPublicId, ok := prevGroupData["imagePublicId"].(string)

		if ok {

			deleErr := services.DeleteFile(prevPublicId)
			if deleErr != nil {
				fmt.Println("Old image deletion Err : ", deleErr)
				helpers.SendMessageAsJson(ctx, "Something went wrong while updating group details", http.StatusInternalServerError)
				return
			}

		}
		prevGroupData = nil
		//Uploading new image to cloudinary
		imageUrl, publicId := services.UploadFileInFormData(groupImage)

		//For updating image and publicId
		updateimageDoc := bson.M{
			"$set": bson.M{
				"groupImage":    imageUrl,
				"imagePublicId": publicId,
			},
		}

		updateErr := groupCollec.FindOneAndUpdate(context.TODO(), conditionToFindGroup, updateimageDoc).Err()

		if updateErr != nil {
			fmt.Println("Group image updation error : ", updateErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong while updating group details", http.StatusInternalServerError)
			return
		}

		groupImageUrl = imageUrl

	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated successfully", "groupName": groupName, "groupBio": groupBio, "groupImageUrl": groupImageUrl})
}
