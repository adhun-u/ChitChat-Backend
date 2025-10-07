package socket

import (
	"chitchat/helpers"
	"chitchat/models"
	"chitchat/services"

	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader for upgrading normal http connection to websocket connection
var callSocketUpgrader = &websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// Creating a variable for adding callers to send remote data
var callers = make(map[int]*websocket.Conn)

// Mutex for locking and unlocking memory
var callSockMutext sync.Mutex

// Offers
var offers = make(map[int]string)

// Answers
var answers = make(map[int]string)

// Candidates
var candidates = make(map[models.CandidateKey][]models.Candidate)

// --------------- CALL COMMUNICATION SOCKET CONTROLLER ---------------
// For handling audio or video call
func ConnectCallSocket(ctx *gin.Context) {
	//Querying current user id
	currentUserId, intConErr := strconv.Atoi(ctx.Query("currentUserId"))
	//Querying caller profilepic
	profilePic := ctx.Query("currentUserProfilePic")
	if intConErr != nil {
		helpers.SendMessageAsJson(ctx, "Provide caller id as an integer", http.StatusNotAcceptable)
		return
	}
	//Querying opposite user id
	oppositeUserId, userIdintConvErr := strconv.Atoi(ctx.Query("oppositeUserId"))

	if userIdintConvErr != nil {
		helpers.SendMessageAsJson(ctx, "Invalid opposite user id provide oppositeUserId as integer", http.StatusNotAcceptable)
		return
	}
	//Upgrading connection
	conn, connErr := callSocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if connErr != nil {
		helpers.SendMessageAsJson(ctx, "Websocket connection error while connection", http.StatusInternalServerError)
		return
	}

	defer func() {
		//Clearing all resouce related to current user when current user leaves
		callSockMutext.Lock()
		delete(callers, currentUserId)
		delete(answers, currentUserId)
		delete(candidates, models.CandidateKey{CallerId: currentUserId, CalleeId: oppositeUserId})
		callSockMutext.Unlock()
	}()
	callSockMutext.Lock()
	//Creating a user with the connection
	callers[currentUserId] = conn
	callSockMutext.Unlock()

	for {
		fmt.Println("Call socket connected")
		var signal models.CallSignal

		//Reading all data as json
		readErr := conn.ReadJSON(&signal)
		if readErr != nil {
			fmt.Println("Call socket read error : ", readErr)
			return
		}

		fmt.Println("Type : ", signal.Type)
		//If caller created an offer to call callee , then sending the callee that caller is calling
		if signal.Type == "CREATE-OFFER" {
			//Saving the offer to get the offer when callee enters
			signal.Mutex.Lock()
			offers[currentUserId] = signal.Data
			signal.Mutex.Unlock()
			//Sending the notification to callee
			sendCallNotification(signal.CallerName, signal.CallerId, signal.CalleeId, signal.CallType, profilePic)
		}
		//For getting the offer that caller created when callee asks
		signal.Mutex.RLock()
		if signal.Type == "GET-OFFER" {
			if calleeConn, ok := callers[currentUserId]; ok {
				calleeConn.WriteJSON(gin.H{"type": "GET-OFFER", "data": offers[signal.CallerId]})
			}
		}
		signal.Mutex.RUnlock()
		// For posting the answer for caller to get the answer that callee posted
		if signal.Type == "POST-ANSWER" {
			//Saving the answer to send caller to that callee posted answer
			signal.Mutex.Lock()
			answers[currentUserId] = signal.Data
			signal.Mutex.Unlock()
			//Sending the indication that callee answered the call
			signal.Mutex.RLock()
			if callerConn, ok := callers[signal.CallerId]; ok {
				callerConn.WriteJSON(gin.H{"type": "POST-ANSWER", "data": answers[signal.CalleeId]})
				callerConn.WriteJSON(gin.H{"type": "GET-CANDIDATE", "candidates": candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}]})
			}
			if calleeConn, ok := callers[signal.CalleeId]; ok {
				calleeConn.WriteJSON(gin.H{"type": "GET-CANDIDATE", "candidates": candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}]})
			}
			signal.Mutex.RUnlock()
		}
		// For saving the candidates
		if signal.Type == "POST-CANDIDATE" {
			fmt.Println("ADDED : ", signal.Data)
			signal.Mutex.Lock()
			candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}] = append(candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}], models.Candidate{
				Type: signal.Type,
				Data: signal.Data,
			})
			signal.Mutex.Unlock()
			signal.Mutex.RLock()
			if currentUserId == signal.CallerId {
				if calleeConn, ok := callers[signal.CalleeId]; ok {
					calleeConn.WriteJSON(gin.H{"type": "GET-CANDIDATE", "candidates": candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}]})
				}
			} else {
				if callerConn, ok := callers[signal.CallerId]; ok {
					callerConn.WriteJSON(gin.H{"type": "GET-CANDIDATE", "candidates": candidates[models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId}]})
				}
			}
			signal.Mutex.RUnlock()
		}
		//For ending the call and reminding the caller or callee that call has ended
		if signal.Type == "CALL-END" || signal.Type == "CALL-CONNECTING" || signal.Type == "RINGING" || signal.Type == "CALL-CONNECTED" || signal.Type == "DECLINE" {
			signal.Mutex.RLock()
			if currentUserId == signal.CallerId {
				if calleeConn, ok := callers[signal.CalleeId]; ok {
					calleeConn.WriteJSON(gin.H{"type": signal.Type})
				}
				if signal.Type == "CALL-END" {
					signal.Mutex.Lock()
					//Clearing offers that current user created for audio or video call
					delete(offers, currentUserId)
					//Clearing answer that callee added
					delete(answers, signal.CalleeId)
					signal.Mutex.Unlock()
				}
			} else {
				if callerConn, ok := callers[signal.CallerId]; ok {
					callerConn.WriteJSON(gin.H{"type": signal.Type})
				}
				if signal.Type == "CALL-END" {
					signal.Mutex.Lock()
					//Clearing offers that caller created for audio or video call
					delete(offers, signal.CallerId)
					//Clearing answer that current user added
					delete(answers, currentUserId)
					//Clearing candidates info that caller and callee added
					delete(candidates, models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId})
					signal.Mutex.Unlock()
				}
				if signal.Type == "DECLINE" {
					signal.Mutex.Lock()
					//Clearing offer that caller created if callee declined the call
					delete(offers, signal.CallerId)
					//Clearing candidates info that caller and callee added
					delete(candidates, models.CandidateKey{CallerId: signal.CallerId, CalleeId: signal.CalleeId})
					signal.Mutex.Unlock()
				}
			}
			signal.Mutex.RUnlock()
		}
	}

}

// For sending calling notification
func sendCallNotification(callerName string, callerId int, calleeId int, callType string, profilePic string) {
	notificationId := fmt.Sprintf("%d%d", callerId, calleeId)
	services.SendCallNotification(
		callType,
		callerName,
		strconv.Itoa(callerId),
		strconv.Itoa(calleeId),
		profilePic,
		notificationId,
	)
}
