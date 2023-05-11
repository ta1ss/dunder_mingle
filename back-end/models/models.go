package models

import "time"

type User struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"passwordConfirm"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DateOfBirth     time.Time `json:"dateOfBirth"`
	Image           string    `json:"image"`
	Nickname        string    `json:"nickname"`
	About           string    `json:"about"`
	ProfileP        int       `json:"profileP"`
	Online          int       `json:"online"`
	MessageRead     int       `json:"messageRead"`
	LastMessage     time.Time `json:"lastMessage"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	UUID            string    `json:"uuid"`
}

// Followers
type UserFollower struct {
	UserID       int    `json:"userId"`
	FollowerID   int    `json:"followerId"`
	FollowerName string `json:"followerName"`
	UserImg      string `json:"userImg"`
}

type UserFollowing struct {
	UserID        int    `json:"userId"`
	FollowingID   int    `json:"followingId"`
	FollowingName string `json:"followingName"`
	UserImg       string `json:"userImg"`
}

type UserFollow struct {
	Followers []UserFollower  `json:"followers"`
	Following []UserFollowing `json:"following"`
}

type FollowStatus struct {
	IsFollowing bool `json:"isFollowing"`
	IsRequested bool `json:"isRequested"`
}

type UserFollowRequest struct {
	CurrentUser  int64 `json:"currentUser"`
	FollowedUser int64 `json:"followedUser"`
	IsRequested  bool  `json:"isRequested"`
	IsFollowed   bool  `json:"isFollowed"`
}

// Notifications
type Notification struct {
	ID           int         `json:"id"`
	UserId       int         `json:"userId"`
	Type         string      `json:"type"`
	SourceId     int         `json:"source"`
	Data         interface{} `json:"data"`
	Seen         int         `json:"seen"`
	CreatedAt    time.Time   `json:"createdAt"`
	Actioned     int         `json:"actioned"`
	OptionalData interface{} `json:"optionalData"`
}

type Post struct {
	ID            int       `json:"id"`
	UserId        int       `json:"userId"`
	GroupId       int       `json:"groupId"`
	UserImg       string    `json:"userImg"`
	CreatedBy     string    `json:"createdBy"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Img           string    `json:"img"`
	Comments      int       `json:"comments"`
	Privacy       string    `json:"privacy"`
	CustomPrivacy []int     `json:"customPrivacy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Comment struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userId"`
	PostId    int       `json:"postId"`
	GroupId   int       `json:"groupId"`
	UserImg   string    `json:"userImg"`
	CreatedBy string    `json:"createdBy"`
	Body      string    `json:"body"`
	Img       string    `json:"img"`
	CreatedAt time.Time `json:"createdAt"`
}

type Group struct {
	ID               int       `json:"id"`
	CreatedBy        int       `json:"createdBy"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Img              string    `json:"img"`
	CreatedAt        time.Time `json:"createdAt"`
	InGroup          bool      `json:"inGroup"`
	JoinRequested    bool      `json:"joinRequested"`
	GroupMembers     []User    `json:"groupMembers"`
	InvitedUsersIds  []int     `json:"invitedUsersIds"`
	InvitingUsersIds []int     `json:"invitingUsersIds"`
}

type GroupInvite struct {
	UserId          int   `json:"userId"`
	GroupId         int   `json:"groupId"`
	InvitedUsersIds []int `json:"invitedUsersIds"`
}

type Event struct {
	ID            int       `json:"id"`
	GroupId       int       `json:"groupId"`
	CreatedBy     int       `json:"createdBy"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	EventStart    time.Time `json:"eventStart"`
	EventEnd      time.Time `json:"eventEnd"`
	CreatedAt     time.Time `josn:"createdAt"`
	CreatorName   string    `json:"creatorName"`
	GoingUsers    []User    `json:"goingUsers"`
	NotGoingUsers []User    `json:"notGoingUsers"`
	Img           string    `json:"img"`
	GroupName     string    `json:"groupName"`
}

type Message struct {
	ID          int       `json:"id"`
	GroupId     int       `json:"groupId"`
	SenderId    int       `json:"senderId"`
	UserName    string    `json:"userName"`
	TargetId    int       `json:"targetId"`
	TargetName  string    `json:"targetName"`
	Body        string    `json:"body"`
	MessageRead int       `json:"messageRead"`
	CreatedAt   time.Time `josn:"createdAt"`
}
