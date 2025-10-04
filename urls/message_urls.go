package urls

import (
	"chitchat/controllers/message"
	"chitchat/middlewares"
	"chitchat/socket"

	"github.com/gin-gonic/gin"
)

func RegisterMessageUrls(route *gin.Engine) {

	messageRoute := route.Group("/message")

	{
		//To fetch the temporary message
		messageRoute.GET("/tempMessages", middlewares.MiddleWare, message.FetchTempMessagesController)
		//To connect the socket
		messageRoute.GET("/ws", socket.ConnectMessageSocket)
		//To delete temporary messages
		messageRoute.DELETE("/tempMessages", middlewares.MiddleWare, message.DeleteTempMessageController)
		//To delete single message
		messageRoute.DELETE("/oneMessage", middlewares.MiddleWare, message.DeleteSingleMessage)
		//To fetch seen indication
		messageRoute.GET("/seenIndication", middlewares.MiddleWare, message.FetchSeenStatusController)
		//To delete seen indication
		messageRoute.DELETE("/seenIndication", middlewares.MiddleWare, message.DeleteSeenStatusController)
		//To save seen info
		messageRoute.POST("/seenInfo", middlewares.MiddleWare, message.SaveSeenInfoController)
		//To upload file (audio,video,image) as message
		messageRoute.POST("/file", middlewares.MiddleWare, message.UploadFileController)
		//To delete file
		messageRoute.DELETE("/file", middlewares.MiddleWare, message.DeleteFileController)
		//To connect call communication websocket
		messageRoute.GET("/call/ws", socket.ConnectCallSocket)
	}

}
