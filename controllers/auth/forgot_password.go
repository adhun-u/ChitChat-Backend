package auth

import (
	"chitchat/config"
	"chitchat/helpers"
	authmodel "chitchat/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ------------- CHANGE PASSWORD CONTROLLER --------------------
func ForgotPasswordController(ctx *gin.Context) {
	var parsedData authmodel.ForgotPasswordModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&parsedData)

	if parseErr != nil {
		fmt.Println("Parse error in /changepassword: ", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Encrypting the password
	encryptedPassword, encryptionErr := config.EncryptPassword(parsedData.NewPassword)

	if encryptionErr != nil {
		fmt.Println("Encryption error /changepassword : ", encryptionErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Updating the column
	db := config.GORM.Table("userdetails").
		Where("email=?", parsedData.Email).
		Update("password", encryptedPassword)

	if db.Error != nil {
		fmt.Println("Updation error in update/password : ", db.Error)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Checking if the row affected
	if int(db.RowsAffected) != 0 {
		helpers.SendMessageAsJson(ctx, "Password has been changed", http.StatusOK)
		return
	} else {
		helpers.SendMessageAsJson(ctx, "Email is not registered", http.StatusNotAcceptable)
		return
	}

}
