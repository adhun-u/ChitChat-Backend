package services

import (
	models "chitchat/models"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

/*

 THIS FILE IS MAINLY USED FOR GENERATING OTP , SENDING OTP AND VERIFYING

*/

// To generate 4 digit otp
func GenerateOtp() string {
	var otp string = ""
	for i := 1; i <= 4; i++ {
		//Generating random number
		otp += strconv.Itoa(rand.Intn(9))
	}
	return otp
}

// To send 4 digit otp
func SendOtp(mail *gomail.Message, otp string, email string) error {
	//Setting from address
	mail.SetHeader("From", os.Getenv("COMPANY_EMAIL"))
	//Setting to address
	mail.SetHeader("To", email)
	//Setting subject
	mail.SetHeader("Subject", fmt.Sprintf("Your OTP is %s will expire within a minute", otp))
	//Setting body
	mail.SetBody("text/plain", "Don't disclose your otp to anyone for your security")
	//Creating a new dialer for sending the otp
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("COMPANY_EMAIL"), os.Getenv("EMAIL_APP_PASS"))
	//Sending otp
	err := dialer.DialAndSend(mail)
	return err
}

// To verify otp
func VerifyOtp(otpEntry *models.OtpEntry, enteredOTP string, enteredEmail string) bool {
	//Checking if the entered email is equal to sent email
	if otpEntry.Email != enteredEmail {
		otpEntry = &models.OtpEntry{}
		return false
	}
	//Checking if the entered otp is equal to sent otp
	if otpEntry.Otp != enteredOTP {
		otpEntry = &models.OtpEntry{}
		return false
	}
	//Current time
	var currentTime = time.Now()
	//Checking if the current time is after otp sent time
	if currentTime.After(otpEntry.ExpiryTime) {
		otpEntry = &models.OtpEntry{}
		return false
	}
	return true
}
