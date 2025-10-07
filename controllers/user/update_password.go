package user

import (
	"chitchat/config"
	"chitchat/helpers"
	updatemodel "chitchat/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdatePasswordController(ctx *gin.Context) {
	//Parsing the data from json
	var updateModel updatemodel.ChangePasswordModel
	parseErr := ctx.ShouldBind(&updateModel)

	if parseErr != nil {
		fmt.Println("Parsing error in update/password : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Invalid json format", http.StatusNotAcceptable)
		return
	}
	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Getting the old password from userdetails table
	var oldPassword string
	config.GORM.
		Table("userdetails").
		Select("password").
		Where("id=?", currentUserId).
		Find(&oldPassword)

	/*Checking if the old password which is got from userdetails table
	is equal to the old password that user entered */

	isCorrect := config.CheckAndDecryptPassword(oldPassword, updateModel.CurrentPassword)

	if isCorrect {
		//If it is correct then updating the password
		//Encrypting the password that user entered
		encryptedPassword, encrypErr := config.EncryptPassword(updateModel.NewPassword)

		if encrypErr != nil {
			fmt.Println("Password encryption Error : ", encrypErr)
			return
		}

		//Updating the password
		db := config.GORM.
			Table("userdetails").
			Where("id=?", currentUserId).
			Update("password", encryptedPassword)

		if db.RowsAffected != 0 {
			helpers.SendMessageAsJson(ctx, "Password changed successfully", http.StatusOK)
			return
		}

		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)

	} else {
		helpers.SendMessageAsJson(ctx, "Current password is incorrect", http.StatusNotAcceptable)

	}

}
