package models

import "sync"

// For sending message
type MessageModel struct {
	SenderId              int    `json:"senderId" bson:"senderId"`
	ReceiverId            int    `json:"receiverId" bson:"receiverId"`
	SenderName            string `json:"senderName"`
	ParentMessageSenderId int    `json:"parentMessageSenderId" bson:"parentMessageSenderId"`
	ParentMessageType     string `json:"parentMessageType" bson:"parentMessageType"`
	ParentText            string `json:"parentText" bson:"parentText"`
	ParentAudioDuration   string `json:"parentAudioDuration" bson:"parentAudioDuration"`
	ParentVoiceDuration   string `json:"parentVoiceDuration" bson:"parentVoiceDuration"`
	SenderProfilePic      string `json:"senderProfilePic"`
	ChatId                string `json:"chatId" bson:"chatId"`
	Time                  string `json:"time" bson:"time"`
	Message               string `json:"message" bson:"message"`
	Type                  string `json:"type" bson:"type"`
	VoiceUrl              string `json:"voiceUrl" bson:"voiceUrl"`
	VoiceDuration         string `json:"voiceDuration" bson:"voiceDuration"`
	ImageUrl              string `json:"imageUrl" bson:"imageUrl"`
	ImageText             string `json:"imageText" bson:"imageText"`
	AudioUrl              string `json:"audioUrl" bson:"audioUrl"`
	AudioDuration         string `json:"audioDuration" bson:"audioDuration"`
	AudioTitle            string `json:"audioTitle" bson:"audioTitle"`
	Exit                  bool   `json:"exit"`
	IsSeen                bool   `json:"isSeen" bson:"isSeen"`
	IsRead                bool   `json:"isRead" bson:"isRead"`
	RepliedMessage        bool   `json:"repliedMessage" bson:"repliedMessage"`
}

// For indicating if user is typing or audio recording
type InteractionIndicatorModel struct {
	ReceiverId int    `json:"receiverId"`
	Indication string `json:"indication"`
}

// For entering chat connection
type ChatConnectionModel struct {
	IsInChat   bool         `json:"isInChat"`
	ReceiverId int          `json:"receiverId"`
	Mutext     sync.RWMutex `json:"_"`
}

// For creating unique key for each connection
type ChatKey struct {
	SenderId   int
	ReceiverId int
}

// For fetching read status
type ReadStatusModel struct {
	SenderId   int  `bson:"senderId"`
	ReceiverId int  `bson:"receiverId"`
	IsSeen     bool `bson:"isSeen"`
}

// For creating peer to peer connection for audio calling
type CallSignal struct {
	CallerName string       `json:"callerName"`
	CallerId   int          `json:"callerId"`
	CalleeId   int          `json:"calleeId"`
	CallType   string       `json:"callType"`
	Type       string       `json:"type"`
	Data       string       `json:"data"`
	Mutex      sync.RWMutex `json:"_"`
}

//For making candidates
type Candidate struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type CandidateKey struct {
	CallerId int
	CalleeId int
}

//For fetching call histories
type CallHistoryModel struct {
	ChatId   string `json:"chatId" bson:"chatId"`
	CallerId int    `json:"callerId" bson:"callerId"`
	CalleeId int    `json:"calleeId" bson:"callerId"`
	CallType string `json:"callType" bson:"callType"`
	CallTime string `json:"callTime" bson:"callTime"`
}
