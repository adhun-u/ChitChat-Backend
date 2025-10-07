package user

import (
	"chitchat/config"
	"chitchat/helpers"
	usermodel "chitchat/models"
	"context"
	"strconv"
	"time"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------ FETCH REQUESTED USERS CONTROLLER -----------------
// For fetching current's users requests
func FetchRequestedUsers(ctx *gin.Context) {

	//Getting current user id from middleware
	currentUserId, exist := ctx.Get("userId")

	if !exist {
		fmt.Println("Userid does not exist ")
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
	collection := config.MongoDB.Collection("requestedUsers")

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Fetching the requsted users that requested current user to add
	var requestedUsers []usermodel.FetchRequestedUserModel
	cursor, dbErr := collection.Find(context, bson.M{"sentUserId": currentUserId}, options.Find().SetLimit(int64(limit)), options.Find().SetSkip(int64(offset)))

	if dbErr != nil {

		fmt.Println("Fetching db error in user/request (GET) : ", dbErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return

	}

	curErr := cursor.All(context, &requestedUsers)

	if curErr != nil {
		fmt.Println("Cursor error in user/request (GET) : ", curErr)
		return
	}

	if len(requestedUsers) != 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "requestedUsers": requestedUsers})

	} else {
		empty := []string{}
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "requestedUsers": empty})
	}
}
