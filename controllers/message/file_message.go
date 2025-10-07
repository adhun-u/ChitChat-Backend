package message

import (
	"chitchat/helpers"
	"chitchat/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ------- UPLOAD FILE CONTROLLER --------------
// For uploading file in message
func UploadFileController(ctx *gin.Context) {
	fileType := ctx.PostForm("type")
	//If the file is image
	switch fileType {
	case "image":
		//Getting image from form file
		image, imageErr := ctx.FormFile("file")

		if imageErr != nil {
			fmt.Println("image uploading error : ", imageErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if image != nil {
			//Uploading image
			imageUrl, publicId := services.UploadFileInFormData(image)

			if imageUrl != "" {
				ctx.JSON(http.StatusOK, gin.H{"fileUrl": imageUrl, "publicId": publicId, "type": "image"})
				return
			}
		} else {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}
	case "audio":
		//Getting audio from form file
		audio, audioErr := ctx.FormFile("file")

		if audioErr != nil {
			fmt.Println("Audio error : ", audioErr)
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if audio != nil {

			//Uploading audio
			audioUrl, publicId := services.UploadFileInFormData(audio)

			if audioUrl != "" {
				ctx.JSON(http.StatusOK, gin.H{"fileUrl": audioUrl, "publicId": publicId, "type": "audio"})
				return
			}

		} else {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		}
	case "voice":

		//Getting voice from form file
		voice, voiceErr := ctx.FormFile("file")

		if voiceErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if voice != nil {
			voiceUrl, publicId := services.UploadFileInFormData(voice)

			if voiceUrl != "" {
				ctx.JSON(http.StatusOK, gin.H{"fileUrl": voiceUrl, "publicId": publicId, "type": "voice"})
			}
		} else {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		}

	}

}

// ------------------- DELETE FILE CONTROLLER -----------
// For deleting a file from cloudinary
func DeleteFileController(ctx *gin.Context) {

	//Querying public id from the request to identify the file
	publicId := ctx.Query("publicId")

	delErr := services.DeleteFile(publicId)

	if delErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	helpers.SendMessageAsJson(ctx, "Deleted file successfully", http.StatusOK)

}
