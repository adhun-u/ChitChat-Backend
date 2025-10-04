package urls

import (
	"chitchat/controllers/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAuthUrls(route *gin.Engine) {

	authRoute := route.Group("auth/")

	{
		//For registering user with email
		authRoute.POST("register/email", auth.RegistrationWithEmailController)
		//For registering user with google
		authRoute.POST("register/google", auth.RegisterationWithGoogleController)
		//For loging with email
		authRoute.POST("login/email", auth.LoginWithEmail)
		//For loging with google
		authRoute.POST("login/google", auth.LoginWithGoogle)
		//For checking username and email
		authRoute.POST("check", auth.CheckEmailController)
		//For sending otp
		authRoute.POST("otp/send", auth.SendOTPController)
		//For verifying otp
		authRoute.POST("otp/verify", auth.VerifyOTPController)
		//For changing password from forgot page
		authRoute.PATCH("change/password", auth.ForgotPasswordController)
	}

}
