package group

import (
	"chitchat/helpers"
	"chitchat/models"
	"chitchat/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
)

var GroupCalls = make(map[string]bool)

// ------------- GROUP CALLS CONTROLLER ---------------------
// For handling group audio and video calls
func GroupCallsController(ctx *gin.Context) {
	var groupCallInfo models.GroupCallInfo
	//Parsing the necessary data to know which group wants to call
	parseErr := ctx.ShouldBind(&groupCallInfo)

	if parseErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid json format", http.StatusNotAcceptable)
		return
	}
	//Creating an access token
	accessToken := auth.NewAccessToken(os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
	grants := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         groupCallInfo.GroupName,
		CanPublish:   boolPtr(true),
		CanSubscribe: boolPtr(true),
	}
	//User data to show
	userMetaData := map[string]any{
		"profilePic": groupCallInfo.CurrentUserProfilePic,
		"username":   groupCallInfo.CurrentUserName,
	}

	//Converting the user meta data to json string to send as string
	jsonMetaData, _ := json.Marshal(userMetaData)

	accessToken.
		SetVideoGrant(grants).
		SetIdentity(fmt.Sprint(groupCallInfo.CurrentUserId)).
		SetMetadata(string(jsonMetaData)).
		SetValidFor(time.Minute)

	//Getting the jwt from the access token
	jwtToken, jwtTokenErr := accessToken.ToJWT()

	if jwtTokenErr != nil {

		fmt.Println("Jwt token error in 'Group call controller' ", jwtTokenErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong while getting token", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Room created", "token": jwtToken})

	//Sending notification to all group members
	if _, ok := GroupCalls[groupCallInfo.GroupId]; !ok {
		GroupCalls[groupCallInfo.GroupId] = true
		services.SendGroupCallNotification(groupCallInfo.GroupId, groupCallInfo.GroupName, groupCallInfo.GroupProfilePic, groupCallInfo.CallType)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
