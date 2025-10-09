package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ----------- GET CURRENT USER ------------------
// To get current user details
func GetCurrentUserController(ctx *gin.Context) {
	//Getting userid in middleware
	userId, exist := ctx.Get("userId")

	fmt.Println("Got user id : ", userId)

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Getting all details from userdetails table
	var currentUserDetails usermodel.CurrentUserModel

	dbErr := config.GORM.
		Table("userdetails").
		Select("id", "username", "profilepic", "bio").
		Where("id=?", userId).
		First(&currentUserDetails).Error

	if dbErr != nil {
		fmt.Println("Querying error in /current : ", dbErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success", "currenUserDetails": currentUserDetails})

}
