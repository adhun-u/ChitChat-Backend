package auth

import (
	"chitchat/config"
	"chitchat/helpers"
	authmodel "chitchat/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// -------- LOGIN WITH EMAIL CONTROLLER -----------------------
// To login with email and password
func LoginWithEmail(ctx *gin.Context) {
	var loginCredentials authmodel.LoginWithEmailModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&loginCredentials)

	if parseErr != nil {
		fmt.Println("Parsing error in login/email : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Creating a struct for saving the data
	type UserDetails struct {
		UserId   int64  `gorm:"column:id"`
		Password string `gorm:"column:password"`
	}

	fmt.Println("Email : ", loginCredentials.Email)
	var userdata UserDetails
	//Getting password from table
	dbErr := config.GORM.
		Table("userdetails").
		Select("password,id").
		Where("email=?", loginCredentials.Email).
		Pluck("password", &userdata).Error

	fmt.Println("Got details : ", userdata.UserId)

	if dbErr != nil {
		fmt.Println("Login error in login/email : ", dbErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Checking if the password is empty
	if userdata.Password == "" {
		helpers.SendMessageAsJson(ctx, "Invalid email and password", http.StatusNotAcceptable)
		return
	}
	//Comparing the encrypted password and entered password
	isCorrect := config.CheckAndDecryptPassword(userdata.Password, loginCredentials.Password)

	if !isCorrect {
		helpers.SendMessageAsJson(ctx, "Incorrect password", http.StatusNotAcceptable)
		return
	} else {
		//Updating the FCM token to database
		config.GORM.
			Table("userdetails").
			Where("id=?", userdata.UserId).
			Update("deviceToken", loginCredentials.DeviceToken)
		//Generating token
		token := config.GenerateJWTtoken(int(userdata.UserId), "email")

		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{"message": "Logged successfully", "token": token})
			return
		} else {
			helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		}
	}

}

// ------------ LOGIN WITH GOOGLE CONTROLLER -------------------
// To login with google
func LoginWithGoogle(ctx *gin.Context) {
	var loginCredentials authmodel.LoginWithGoogleModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&loginCredentials)

	if parseErr != nil {
		fmt.Println("Parsing error in login/google", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Creating a struct for saving the data
	type UserDetails struct {
		UserId int64  `gorm:"column:id"`
		Email  string `gorm:"column:email"`
	}
	var userDetails UserDetails
	//Checking if the row exists
	config.GORM.
		Table("userdetails").
		Select("id,email").
		Where("email=? AND authtype='google' AND username=?", loginCredentials.Email, loginCredentials.Username).
		First(&userDetails)

	//Checking whether the email and are null or not
	if userDetails.Email != "" && userDetails.UserId != 0 {
		//Updating the FCM token to database
		config.GORM.
			Table("userdetails").
			Where("id=?", userDetails.UserId).
			Update("deviceToken", loginCredentials.DeviceToken)
		//Generating token
		token := config.GenerateJWTtoken(int(userDetails.UserId), "google")
		ctx.JSON(http.StatusOK, gin.H{"message": "Logged successfully", "token": token})
		return
	} else {
		helpers.SendMessageAsJson(ctx, "Email is not registered", http.StatusNotAcceptable)
	}
}
