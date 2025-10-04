package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchSentRequestUsers(ctx *gin.Context) {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Getting the current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("Current user id does not exist : ")
		helpers.SendMessageAsJson(ctx, "Provide valid token", http.StatusNotAcceptable)
		return
	}

	//Querying limit
	limit, limitIntConvErr := strconv.Atoi(ctx.Query("limit"))

	if limitIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide limit as integer", http.StatusNotAcceptable)
		return
	}
	//Querying page
	page, pageIntConvErr := strconv.Atoi(ctx.Query("page"))

	if pageIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide page as integer", http.StatusNotAcceptable)
		return
	}

	offset := (page - 1) * limit

	//Fetching all users that current user sent request
	sentUsersCollection := config.MongoDB.Collection("sentUsers")

	var sentUsers []usermodel.FetchSentRequestUsersModel
	cursor, findErr := sentUsersCollection.Find(
		context,
		bson.M{"requestedUserId": currentUserId},
		options.Find().SetLimit(int64(limit)),
		options.Find().SetSkip(int64(offset)),
	)

	if findErr != nil {
		fmt.Println("finding error in user/sent : ", findErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	cursor.All(context, &sentUsers)
	if len(sentUsers) != 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "sentUsers": sentUsers})
	} else {
		empty := []string{}
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "sentUsers": empty})
	}

}
