package urls

import (
	"chitchat/controllers/user"
	"chitchat/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserUrls(route *gin.Engine) {

	userRoute := route.Group("user/")

	{
		//For getting current user
		userRoute.GET("current", middlewares.MiddleWare, user.GetCurrentUserController)
		//For updating current user's details
		userRoute.PATCH("currentuser", middlewares.MiddleWare, user.UpdateCurrentUserDetailsController)
		//For changing current user's password from "change password page"
		userRoute.PATCH("password", middlewares.MiddleWare, user.UpdatePasswordController)
		//For searching a user
		userRoute.GET("search", middlewares.MiddleWare, user.SearchUserController)
		//For requesting a user to add
		userRoute.POST("request", middlewares.MiddleWare, user.RequestUserController)
		//For fetching requests of current user's
		userRoute.GET("request", middlewares.MiddleWare, user.FetchRequestedUsers)
		//For accepting requests
		userRoute.POST("request/accept", middlewares.MiddleWare, user.AcceptRequest)
		//For declining a request
		userRoute.DELETE("request/decline", middlewares.MiddleWare, user.DeclineRequestController)
		//For fetching the users that current sent request to add
		userRoute.GET("sent", middlewares.MiddleWare, user.FetchSentRequestUsers)
		//For fetching added users
		userRoute.GET("friendsWithLastMessage", middlewares.MiddleWare, user.FetchFriendsWithLastMessageController)
		//For withdrawing a request
		userRoute.DELETE("withdraw", middlewares.MiddleWare, user.WithdrawRequestController)
		//For removing specific user from addedUsers
		userRoute.DELETE("remove", middlewares.MiddleWare, user.RemoveUserController)
		//For changing last message time
		userRoute.PATCH("lastMessageTime", middlewares.MiddleWare, user.ChangeLastMessageController)
	}
}
