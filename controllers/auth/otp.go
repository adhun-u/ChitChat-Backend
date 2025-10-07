package auth

import (
	"chitchat/helpers"
	authmodel "chitchat/models"
	"chitchat/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var otpEntry = make(map[string]*authmodel.OtpEntry)

// ----------- SEND OTP CONTROLLER ----------------------------
func SendOTPController(ctx *gin.Context) {
	var otpModel authmodel.SendOTPModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&otpModel)

	if parseErr != nil {
		fmt.Println("Parsing error in otp/send : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Generating otp
	otp := services.GenerateOtp()
	otpEntry[otpModel.Email] = &authmodel.OtpEntry{
		Email:      otpModel.Email,
		Otp:        otp,
		ExpiryTime: time.Now().Add(1 * time.Minute),
	}
	//Creating new gomailer
	mailer := gomail.NewMessage()
	//Sending otp
	err := services.SendOtp(mailer, otp, otpModel.Email)

	if err != nil {
		fmt.Println("Mail error in otp/send : ", err)
		helpers.SendMessageAsJson(ctx, "Invalid email", http.StatusNotAcceptable)
		return
	}
	helpers.SendMessageAsJson(ctx, "OTP has been sent successfully", http.StatusOK)

}

// --------- VARIFY OTP CONTROLLER ------------------------------
func VerifyOTPController(ctx *gin.Context) {
	var verificationModel authmodel.OTPVerificationModel
	//Parsing the data from json
	parseErr := ctx.ShouldBind(&verificationModel)

	if parseErr != nil {
		fmt.Println("Parse error in otp/verify : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	isVerified := services.VerifyOtp(otpEntry[verificationModel.Email], verificationModel.EnteredOtp, verificationModel.Email)

	if !isVerified {
		helpers.SendMessageAsJson(ctx, "Invalid OTP", http.StatusNotAcceptable)
		return
	} else {
		helpers.SendMessageAsJson(ctx, "Verified successfully", http.StatusOK)
		return
	}
}
