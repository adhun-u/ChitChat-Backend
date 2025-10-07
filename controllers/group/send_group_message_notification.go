package group

import (
	"chitchat/helpers"
	"chitchat/models"
	"chitchat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// -------------- SEND GROUP MESSAGE NOTIFICATION CONTROLLER -------------
// For sending a group message as notification to members
func SendGroupMessageNotificationController(ctx *gin.Context) {
	//Parsing necessary data to send notification
	var notification models.GroupMessageNotificationModel

	parseErr := ctx.ShouldBind(&notification)

	if parseErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Sending the message as notification to all group members who are not in chat
	services.SendGroupMessageNotification(notification.GroupId, notification.Title, notification.Body, notification.GroupProfilePic, notification.MessageType)

	helpers.SendMessageAsJson(ctx, "Notification sent successfully", http.StatusOK)

}
