package auth

import (
	"chitchat/config"
	"chitchat/helpers"
	authmodel "chitchat/models"
	"chitchat/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// -------- REGISTER USER WITH EMAIL CONTROLLER -----------------
// To register user with email
func RegistrationWithEmailController(ctx *gin.Context) {

	var parsedData authmodel.RegistrationWithEmailModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&parsedData)

	if parseErr != nil {
		fmt.Println("Invalid json : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Invalid json", http.StatusNotAcceptable)
		return
	}
	//Encrypting the password
	encryptedPassword, encryptionErr := config.EncryptPassword(parsedData.Password)

	//Checking if there is an error during encryption
	if encryptionErr != nil {
		fmt.Println("Password encryption Error : ", encryptionErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Creating an interface for inserting data
	userDetails := authmodel.RegistrationWithEmailModel{
		Email:       parsedData.Email,
		Username:    parsedData.Username,
		Password:    encryptedPassword,
		Authtype:    "email",
		ProfilePic:  "",
		Userbio:     "",
		DeviceToken: parsedData.DeviceToken,
	}
	//Inserting the data
	result := config.GORM.Table("userdetails").Create(&userDetails)
	//Checking whether row affected or not
	if result.Error != nil {
		fmt.Println("Insertion error in register/email : ", result.Error)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)

	} else if int(result.RowsAffected) != 0 {
		//Inserting the FCM token to database
		//Generating token
		token := config.GenerateJWTtoken(int(userDetails.Id), "email")
		ctx.JSON(http.StatusOK, gin.H{"message": "Account has been created", "token": token})
	} else {
		helpers.SendMessageAsJson(ctx, "Could not create account", http.StatusInternalServerError)
	}

}

// ---------- REGISTER USER WITH GOOGLE CONTROLLER ------------------
// To register user with google
func RegisterationWithGoogleController(ctx *gin.Context) {

	var parsedData authmodel.RegisterationWithGoogleModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&parsedData)

	if parseErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusNotAcceptable)
		return
	}

	var emailCount int64
	//Retrieving the data from userdetails to check whether the google account already exists or not
	config.GORM.
		Table("userdetails").
		Where("email=?", parsedData.Email).Count(&emailCount)

	if int(emailCount) != 0 {
		helpers.SendMessageAsJson(ctx, "Email is already in use", http.StatusAlreadyReported)
		return
	}
	//Uploading the image
	imageUrl := services.UploadFileWithURL(parsedData.ProfilePic)

	//Inserting the data
	userDetails := authmodel.RegisterationWithGoogleModel{
		Username:    parsedData.Username,
		Email:       parsedData.Email,
		Authtype:    "google",
		ProfilePic:  imageUrl,
		Bio:         "",
		DeviceToken: parsedData.DeviceToken,
	}
	result := config.GORM.Table("userdetails").Create(&userDetails)

	//Checking whether row affected or not
	if (result.RowsAffected) != 0 {
		//Generating token
		token := config.GenerateJWTtoken(int(userDetails.Id), "google")
		ctx.JSON(http.StatusOK, gin.H{"message": "Account has been created", "token": token})
	} else {
		helpers.SendMessageAsJson(ctx, "Could not create account", http.StatusInternalServerError)
	}

}
