package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/services"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// --------- UPDATE CURRENT USER DETAILS CONTROLLER ------------
// To update profilepic , fullname , bio
func UpdateCurrentUserDetailsController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	contxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if !exist {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	newname := ctx.PostForm("newname")
	newbio := ctx.PostForm("newbio")
	imageFile, _ := ctx.FormFile("profilePic")

	if newname != "" {
		//Updating current user name if the newname is not empty
		config.GORM.
			Table("userdetails").
			Where("id=?", currentUserId).
			Update("username", newname)
			//Also updating from requested users collection
		collection := config.MongoDB.Collection("requestedUserCollection")

		collection.
			FindOneAndUpdate(
				contxt,
				bson.M{"requestedUserId": currentUserId},
				bson.M{"requestedUsername": newname},
			)
	}

	if newbio != "" {
		//Updating current user name if the newbio is not empty
		config.GORM.
			Table("userdetails").
			Where("id=?", currentUserId).
			Update("bio", newbio)
			//Also updating from requested users collection
		collection := config.MongoDB.Collection("requestedUserCollection")

		collection.
			FindOneAndUpdate(
				contxt,
				bson.M{"requestedUserId": currentUserId},
				bson.M{"requestedUserbio": newbio},
			)
	}
	var newImageUrl string
	//Checking if file is not empty
	if imageFile != nil {
		//Uploading the file to Cloudinary to get url
		imageUrl, _ := services.UploadFileInFormData(imageFile)

		if imageUrl != "" {
			newImageUrl = imageUrl
			config.GORM.
				Table("userdetails").
				Where("id=?", currentUserId).
				Update("profilepic", imageUrl)

				//Also updating from requested users collection
			collection := config.MongoDB.Collection("requestedUserCollection")

			collection.
				FindOneAndUpdate(
					contxt,
					bson.M{"requestedUserId": currentUserId},
					bson.M{"requestedUserProfilePic": imageUrl},
				)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "username": newname, "bio": newbio, "imageUrl": newImageUrl})

}
