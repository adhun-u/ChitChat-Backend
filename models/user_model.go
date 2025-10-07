package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// For accepting friend request
type AcceptRequestModel struct {
	UserId int `json:"userId"`
}

// For requesting user to message
type AddRequestUserModel struct {
	RequestedUserId         int    `json:"requestedUserId"`
	RequestedUsername       string `json:"requestedUsername"`
	RequestedUserProfilePic string `json:"requestedUserProfilePic"`
	RequestedUserbio        string `json:"requestedUserbio"`
	RequestedDate           string `json:"requestedDate"`
}

// For getting current user's friend
type FriendModel struct {
	CurrentUserId   int       `bson:"currentUserId"`
	FriendId        int       `bson:"friendId"`
	LastMessageTime time.Time `bson:"lastMessageTime"`
}

// To give added users details
type GetAddedUserDetails struct {
	UserId               int    `json:"userId" gorm:"column:id"`
	Username             string `json:"username" gorm:"column:username"`
	ProfilePic           string `json:"profilePic" gorm:"column:profilepic"`
	Userbio              string `json:"bio" gorm:"column:bio"`
	PendingMessagesCount int    `json:"pendingMessageCount"`
	LastPendingMessage
}

type LastPendingMessage struct {
	LastPendingMessage string    `json:"lastPendingMessage" bson:"message"`
	Time               time.Time `json:"time" bson:"time"`
	Type               string    `json:"type" bson:"type"`
	ImageText          string    `json:"imageText" bson:"imageText"`
}

// For giving requested users details
type FetchRequestedUserModel struct {
	Id                      primitive.ObjectID `json:"id" bson:"_id"`
	RequstedUserId          int                `json:"requestedUserId" bson:"requestedUserId"`
	RequestedUsername       string             `json:"requestedUsername" bson:"requestedUsername"`
	RequestedUserProfilePic string             `json:"profilePic" bson:"requestedUserProfilePic"`
	RequestedUserbio        string             `json:"bio" bson:"requestedUserbio"`
	RequestedDate           string             `json:"requestedDate" bson:"requestedDate"`
}

// For retrieving users that current user sent
type FetchSentRequestUsersModel struct {
	SentUserId         int    `json:"sentUserId" bson:"sentUserId"`
	SentUsername       string `json:"sentUsername" bson:"sentUsername"`
	SentUserProfilePic string `json:"sentUserProfilePic" bson:"sentUserProfilePic"`
	SentUserbio        string `json:"sentUserbio" bson:"sentUserbio"`
	SentDate           string `json:"sentDate" bson:"sentDate"`
}

// To get current user
type CurrentUserModel struct {
	Id         int    `json:"userId" gorm:"column:id"`
	Username   string `json:"username" gorm:"column:username"`
	Authtype   string `json:"authtype" gorm:"column:authtype"`
	Email      string `json:"email" gorm:"column:email"`
	ProfilePic string `json:"profilePic" gorm:"column:profilepic"`
	Userbio    string `json:"bio" gorm:"column:bio"`
}

// To retrieve or store user details
type UserModel struct {
	UserId      int    `json:"userId" gorm:"column:id"`
	Username    string `json:"username" gorm:"column:username"`
	ProfilePic  string `json:"profilePic" gorm:"column:profilepic"`
	Userbio     string `json:"bio" gorm:"column:bio"`
	DeviceToken string `json:"-" gorm:"column:deviceToken"`
}

// For giving searched user details
type SearchedUserModel struct {
	Id          int    `json:"id" gorm:"column:id"`
	Username    string `json:"username" gorm:"column:username"`
	ProfilePic  string `json:"profilePic" gorm:"column:profilepic"`
	UserBio     string `json:"bio" gorm:"column:bio"`
	IsRequested bool   `json:"isRequested"`
	IsAdded     bool   `json:"isAdded"`
}
