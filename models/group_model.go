package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// For fetch group requests
type FetchGroupRequests struct {
	GroupId        string `json:"groupId" bson:"requestedGroupId"`
	GroupName      string `json:"groupName" bson:"requestedGroupName"`
	Username       string `json:"username"`
	UserId         int    `json:"userId" bson:"requestedUserId"`
	Userbio        string `json:"userBio"`
	UserProfilePic string `json:"userProfilepic"`
}

// For fetching group ids
type FetchGroupIdModel struct {
	GroupId string `json:"groupId" bson:"groupId"`
}

// For fetching all groups
type FetchGroupModel struct {
	GroupId            primitive.ObjectID `json:"groupId" bson:"_id"`
	GroupName          string             `json:"groupName" bson:"groupName"`
	GroupBio           string             `json:"groupBio" bson:"groupBio"`
	GroupImageUrl      string             `json:"groupImageUrl" bson:"groupImage"`
	GroupAdminUserId   int32              `json:"groupAdminUserId" bson:"groupAdminUserId"`
	GroupMembersCount  int32              `json:"groupMembersCount" bson:"groupMembersCount"`
	IsCurrentUserAdded bool               `json:"isCurrentUserAdded"`
	IsRequestSent      bool               `json:"isRequestSent"`
	CreatedAt          time.Time          `json:"createdAt" bson:"createdAt"`
	LastMessageTime    time.Time          `json:"lastMessageTime" bson:"lastMessageTime"`
}

// For getting members of a group
type GroupAddedUser struct {
	UserId int `json:"userId" bson:"memberId"`
}

// For sending request to add to a group
type RequestGroupModel struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
	AdminId   int    `json:"adminId"`
}

// For getting group call info
type GroupCallInfo struct {
	GroupId               string `json:"groupId"`
	GroupName             string `json:"groupName"`
	GroupProfilePic       string `json:"groupProfilePic"`
	CurrentUserId         int    `json:"currentUserId"`
	CurrentUserProfilePic string `json:"currentUserProfilePic"`
	CurrentUserName       string `json:"currentUserName"`
	CallType              string `json:"callType"`
}
