package urls

import (
	"chitchat/controllers/group"
	"chitchat/middlewares"
	"chitchat/socket"

	"github.com/gin-gonic/gin"
)

func RegisterGroupUrls(route *gin.Engine) {

	groupRoute := route.Group("/group")

	{
		//For creating group
		groupRoute.POST("/create", middlewares.MiddleWare, group.CreateGroupController)
		//For getting current user joined or created groups
		groupRoute.GET("/get", middlewares.MiddleWare, group.FetchGroupController)
		//For searching groups
		groupRoute.GET("/search", middlewares.MiddleWare, group.SearchGroupController)
		//For requesting to add in group
		groupRoute.POST("/request", middlewares.MiddleWare, group.RequestGroupController)
		//For fetching requests of a group
		groupRoute.GET("/requests", middlewares.MiddleWare, group.FetchGroupRequestsController)
		//For fetching added users of a group
		groupRoute.GET("/addedUsers", middlewares.MiddleWare, group.FetchMembersController)
		//For fetching users to add to a group
		groupRoute.GET("/users", middlewares.MiddleWare, group.FetchUserToAddMemberController)
		//For accepting group request
		groupRoute.POST("/acceptRequest", middlewares.MiddleWare, group.AcceptGroupRequestController)
		//For removing group request
		groupRoute.DELETE("/declineRequest", middlewares.MiddleWare, group.DeclineGroupRequestController)
		//For editing group info
		groupRoute.PATCH("/edit", middlewares.MiddleWare, group.EditGroupController)
		//For adding a person
		groupRoute.POST("/add", middlewares.MiddleWare, group.AddMemberController)
		//For deleting a user from a group
		groupRoute.DELETE("/delete/user", middlewares.MiddleWare, group.RemoveMemberController)
		//For handling group websocket
		groupRoute.GET("/ws", socket.ConnectGroupSocket)
		//For existing from a group
		groupRoute.DELETE("/exit", middlewares.MiddleWare, group.ExitGroupController)
		//For sending group notification
		groupRoute.POST("/notification", middlewares.MiddleWare, group.SendGroupMessageNotificationController)
		//For creating room for audio or video call for a group
		groupRoute.POST("/call", middlewares.MiddleWare, group.GroupCallsController)
		//For changing last message time
		groupRoute.PATCH("/changeTime", middlewares.MiddleWare, group.ChangeGroupLastMessageTimeController)
	}
}
