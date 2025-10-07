package models

import "time"

//To register user with google
type RegisterationWithGoogleModel struct {
	Id          int64  `gorm:"primaryKey"`
	Username    string `json:"username" gorm:"column:username"`
	Email       string `json:"email" gorm:"column:email"`
	ProfilePic  string `json:"profilepic" gorm:"column:profilepic"`
	Authtype    string `gorm:"column:authtype"`
	Bio         string `gorm:"column:bio"`
	DeviceToken string `json:"deviceToken" gorm:"column:deviceToken"`
}

//To register user with email and password
type RegistrationWithEmailModel struct {
	Id          int64  `gorm:"primaryKey"`
	Email       string `json:"email" gorm:"column:email"`
	Username    string `json:"username" gorm:"column:username"`
	Password    string `json:"password" gorm:"column:password"`
	ProfilePic  string `gorm:"column:profilepic"`
	Userbio     string `gorm:"column:bio"`
	Authtype    string `gorm:"column:authtype"`
	DeviceToken string `json:"deviceToken" gorm:"column:deviceToken"`
}

//Login with google model
type LoginWithGoogleModel struct {
	Email       string `json:"email" gorm:"column:email"`
	Username    string `json:"username" gorm:"column:username"`
	DeviceToken string `json:"deviceToken" gorm:"column:deviceToken"`
}

//Login user with email model
type LoginWithEmailModel struct {
	Email       string `json:"email" gorm:"column:email"`
	Password    string `json:"password" gorm:"column:password"`
	DeviceToken string `json:"deviceToken" gorm:"column:deviceToken"`
}

//To check if the username or email already registered
type CheckUsernameAndEmailModel struct {
	Email string `json:"email" gorm:"column:email"`
}

//Forgot password model
type ForgotPasswordModel struct {
	Email       string `json:"email" gorm:"column:email"`
	NewPassword string `json:"newPassword"`
}

//Send otp model
type SendOTPModel struct {
	Email string `json:"email"`
}

//For saving otp data
type OtpEntry struct {
	Otp        string
	Email      string
	ExpiryTime time.Time
}

//For otp verification
type OTPVerificationModel struct {
	Email      string `json:"email"`
	EnteredOtp string `json:"otp"`
}
