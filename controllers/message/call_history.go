package message

import (
	"chitchat/config"
	"chitchat/helpers"
	"chitchat/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ----- FETCH CALL HISTORY CONTROLLER -------
// For fetching call histories
func FetchCallHistories(ctx *gin.Context) {
	//Querying caller id
	callerId, callerIdIntConvErr := strconv.Atoi(ctx.Query("callerId"))
	//Querying callee id
	calleeId, calleeIdIntConvErr := strconv.Atoi(ctx.Query("calleeId"))

	if calleeIdIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide calleeId as integer", http.StatusNotAcceptable)
		return
	}

	if callerIdIntConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide callerId as integer", http.StatusNotAcceptable)
		return
	}

	var histories []models.CallHistoryModel

	callHistoriesCollec := config.MongoDB.Collection("callHistories")
	//Condition to find the doc
	conditionToFindDoc := bson.M{

		"callerId": callerId,
		"calleeId": calleeId,
	}

	fetchCallHistoryRes, fetchCallHistoryErr := callHistoriesCollec.Find(context.TODO(), conditionToFindDoc)

	if fetchCallHistoryErr != nil {
		helpers.SendMessageAsJson(ctx, "Something went wrong while fetching histories", http.StatusInternalServerError)
		fetchCallHistoryRes.Close(context.TODO())
		return
	}

	for fetchCallHistoryRes.Next(context.TODO()) {

		//Adding one by one
		var history models.CallHistoryModel

		cursorDecodeErr := fetchCallHistoryRes.Decode(&history)

		if cursorDecodeErr != nil {
			helpers.SendMessageAsJson(ctx, "Something went wrong while fetching histories", http.StatusInternalServerError)
			fetchCallHistoryRes.Close(context.TODO())
			return
		}

		histories = append(histories, history)
	}

	if len(histories) == 0 {
		var emptyList []string

		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "callHistories": emptyList})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Success", "callHistories": histories})
	}

}
