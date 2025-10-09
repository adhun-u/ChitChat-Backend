package user

import (
	"chitchat/config"
	"chitchat/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ----------- CHANGE DEVICE TOKEN CONTROLLER ----------
// For changing device when it refreshes
func ChangeFCMTokenController(ctx *gin.Context) {
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Querying fcm token
	fcmToken := ctx.Query("fcmToken")

	if fcmToken == "" {
		helpers.SendMessageAsJson(ctx, "Provide fcmToken", http.StatusNotAcceptable)
		return
	}

	//Changing the token
	updateErr := config.
		GORM.
		Table("userdetails").
		Update("deviceToken", fcmToken).
		Where("id=?", currentUserId).
		Error

	if updateErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while updating token", http.StatusInternalServerError)
		return
	} else {
		helpers.SendMessageAsJson(ctx, "Updated successfully", http.StatusOK)
	}

}
