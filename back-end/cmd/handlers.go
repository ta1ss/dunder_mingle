package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"social-network/back-end/forms"
	"social-network/back-end/models"
	"social-network/back-end/repository"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const postImagePath = "media/post_images/"
const groupImagePath = "media/group_images/"
const commentImagePath = "media/comment_images/"
const profileImagePath = "media/profile_images/"

// use utils.go JSON helpers(writeJson, readJson, errorJson)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "social-network active",
		Version: "1.0.0",
	}

	_ = app.writeJson(w, http.StatusOK, payload, nil)
}

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = app.readJson(w, r, &user)

	if !login(user.Email, user.Password) {
		fmt.Println("User authentication failed: ", user.Email, user.Password)
		_ = app.errorJson(w, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserByEmail(user.Email)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	user.UUID = setSessionID(user.ID)

	http.SetCookie(w, &http.Cookie{
		Name:    "SN-Session",
		Value:   user.UUID,
		Expires: time.Now().Add(24 * time.Hour),
	})

	fmt.Println("User authentication ok")
	_ = app.writeJson(w, http.StatusOK, user, nil)
}

func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	removeUUID(cookie.Value)
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	fmt.Println("User logged out, Cookie removed")
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}

// Profile Handlers

func (app *application) GetProfile(w http.ResponseWriter, r *http.Request) {
	var user models.User

	userId, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid user id"), http.StatusBadRequest)
		return
	}

	if userId == 0 {
		cookie, err := r.Cookie("SN-Session")
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
			return
		}
		user, err = repository.GetUserByCookie(cookie.Value)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
			return
		}
	} else {
		user, err = repository.GetUserById(userId)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, user, nil)
}

func (app *application) PutProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	var userUpdate models.User
	_ = app.readJson(w, r, &userUpdate)
	user.ProfileP = userUpdate.ProfileP

	err = repository.UpdateUserStatus(user.ID, userUpdate.ProfileP)
	if err != nil {
		log.Println(err)
		_ = app.errorJson(w, fmt.Errorf("failed to update user"), http.StatusInternalServerError)
		return
	}

	fmt.Println("User status updated")
	_ = app.writeJson(w, http.StatusOK, user, nil)
}

func (app *application) ProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid user id"), http.StatusBadRequest)
		return
	}

	if userId == 0 {
		cookie, err := r.Cookie("SN-Session")
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
			return
		}

		user, err := repository.GetUserByCookie(cookie.Value)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
			return
		}

		userId = user.ID
	}

	user, err := repository.GetUserById(userId)
	if err != nil {
		app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	file, err := os.Open(profileImagePath + user.Image)
	if err != nil {
		app.errorJson(w, fmt.Errorf("failed to open file"), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/jpeg")

	_, err = io.Copy(w, file)
	if err != nil {
		app.errorJson(w, fmt.Errorf("failed to write response"), http.StatusInternalServerError)
		return
	}
}

func (app *application) GroupImageHandler(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid group id"), http.StatusBadRequest)
		return
	}

	group, err := repository.GetGroup(groupId)
	if err != nil {
		app.errorJson(w, fmt.Errorf("group not found"), http.StatusUnauthorized)
		return
	}

	fmt.Println(groupImagePath + group.Img)
	file, err := os.Open(groupImagePath + group.Img)
	if err != nil {
		app.errorJson(w, fmt.Errorf("failed to open file"), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/jpeg")

	_, err = io.Copy(w, file)
	if err != nil {
		app.errorJson(w, fmt.Errorf("failed to write response"), http.StatusInternalServerError)
		return
	}
}

func (app *application) ProfileFollowersHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		log.Printf("getting userId error %v", err)
		return
	}

	followers, err := repository.GetUserFollow(userId)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to get followers"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, followers, nil)
}

func (app *application) GetFollowStatus(w http.ResponseWriter, r *http.Request) {
	var followStatus models.FollowStatus
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	currentUser := user.ID
	otherUser, _ := strconv.Atoi(chi.URLParam(r, "userID"))

	followStatus, err = repository.GetFollowStatus(currentUser, otherUser)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to get follow status"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, followStatus, nil)
}

func (app *application) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	var UserFollowRequest models.UserFollowRequest

	err := app.readJson(w, r, &UserFollowRequest)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to read json"), http.StatusInternalServerError)
		return
	}

	currentUser := UserFollowRequest.CurrentUser
	followedUser := UserFollowRequest.FollowedUser
	isRequested := UserFollowRequest.IsRequested
	isFollowed := UserFollowRequest.IsFollowed

	changeFollowStatus, err := repository.ChangeFollowStatus(currentUser, followedUser, isRequested, isFollowed)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to change follow status"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, changeFollowStatus, nil)
}

// Get All Users
func (app *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to get users"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, users, nil)
}

// Register Handler

func (app *application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// read response from the client
	_ = app.readJson(w, r, &user)

	err := forms.ValidateRegistrationForm(user)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("form not valid"), http.StatusUnauthorized)
		return
	}
	err = repository.EmailTaken(user.Email)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("email taken"), http.StatusUnauthorized)
		return
	}
	if user.Image == "" {
		user.Image = "default_profile.png"
	} else {
		imgName, err := saveBase64Image(user.Image, profileImagePath)
		if err != nil {
			fmt.Println(err)
		}
		user.Image = imgName
	}
	userID, err := registerUser(user)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("registration error"), http.StatusUnauthorized)
		return
	}
	user.UUID = setSessionID(userID)

	// set uuid as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "SN-Session",
		Value:   user.UUID,
		Expires: time.Now().Add(24 * time.Hour),
	})
	user.Password = ""
	user.PasswordConfirm = ""
	// send response back to the client
	fmt.Println("User authentication ok: ", user.Email)
	_ = app.writeJson(w, http.StatusOK, user, nil)
}

// User posts

func (app *application) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	var posts []models.Post

	if r.URL.Query().Has("id") { // Get selected user posts
		targetId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("getting userId error %v", err)
			return
		}
		posts, err = repository.GetUserPosts(targetId, user.ID)
		if err != nil {
			log.Println("GetUserPosts error: ", err)
		}
	} else { // Get all posts (news feed)
		posts, err = repository.GetAllPosts(user.ID)
		if err != nil {
			log.Println("GetAllPosts error: ", err)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, posts, nil)
}

func (app *application) PostUserPost(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	var post models.Post
	// read response from the client
	_ = app.readJson(w, r, &post)

	post.UserId = user.ID
	post.CreatedAt = time.Now()
	post.CreatedBy = user.FirstName + " " + user.LastName

	if post.Img != "" {
		img, err := saveBase64Image(post.Img, postImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("can't save image error: %v", err), http.StatusUnauthorized)
			return
		}
		post.Img = img
	}

	postId, err := repository.CreateUserPost(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot create post"), http.StatusInternalServerError)
		return
	}

	if post.Privacy == "Custom" && len(post.CustomPrivacy) > 0 {
		for _, targetId := range post.CustomPrivacy {
			err := repository.AddCustomPrivacy(postId, targetId)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("cannot add custom privacy to post"), http.StatusInternalServerError)
				return
			}
		}
	}
	_ = app.writeJson(w, http.StatusOK, post, nil)
}

func (app *application) DeleteUserPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	// read response from the client
	_ = app.readJson(w, r, &post)

	image, err := repository.GetUserPostImage(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot get post image"), http.StatusInternalServerError)
		return
	}

	err = repository.DeleteUserPost(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot delete post"), http.StatusInternalServerError)
		return
	}

	if image != "" {
		err = deleteImage(image, postImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("cannot delete post image"), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, post, nil)
}

func (app *application) GetUserPostAndComments(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(chi.URLParam(r, "postID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid post id"), http.StatusBadRequest)
		return
	}

	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	// get post where post id
	post, err := repository.GetUserPost(postId, user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("error getting user post"), http.StatusUnauthorized)
		return
	}

	// get comments where post id
	comments, err := repository.GetUserPostComments(postId)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("post comments not found"), http.StatusUnauthorized)
		return
	}

	data := struct {
		Post     models.Post      `json:"post"`
		Comments []models.Comment `json:"comments"`
	}{Post: post, Comments: comments}

	_ = app.writeJson(w, http.StatusOK, data, nil)
}

func (app *application) GetGroupPostAndComments(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(chi.URLParam(r, "postID"))
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("invalid post id"), http.StatusBadRequest)
		return
	}

	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	// get group post where post id
	post, err := repository.GetGroupPost(postId, user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("error getting group post"), http.StatusUnauthorized)
		return
	}

	// get group post comments where post id
	comments, err := repository.GetGroupPostComments(postId)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("group comments not found"), http.StatusUnauthorized)
		return
	}

	data := struct {
		Post     models.Post      `json:"post"`
		Comments []models.Comment `json:"comments"`
	}{Post: post, Comments: comments}

	_ = app.writeJson(w, http.StatusOK, data, nil)
}

func (app *application) PostComment(w http.ResponseWriter, r *http.Request) {
	commentType := chi.URLParam(r, "type")

	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	_ = app.readJson(w, r, &comment)

	comment.UserId = user.ID
	comment.UserImg = user.Image
	comment.CreatedAt = time.Now()
	comment.CreatedBy = user.FirstName + " " + user.LastName

	if comment.Img != "" {
		img, err := saveBase64Image(comment.Img, commentImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("can't save comment image error: %v", err), http.StatusUnauthorized)
			return
		}
		comment.Img = img
	}

	switch commentType {
	case "user":
		{
			comment.ID, err = repository.CreateUserComment(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("can't add comment to database: %v", err), http.StatusUnauthorized)
				return
			}
		}
	case "group":
		{
			comment.ID, err = repository.CreateGroupComment(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("can't add group comment to database: %v", err), http.StatusUnauthorized)
				return
			}
		}
	default:
		{
			_ = app.errorJson(w, fmt.Errorf("error creating a comment"), http.StatusBadRequest)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, comment, nil)
}

func (app *application) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentType := chi.URLParam(r, "type")

	var comment models.Comment
	var image string

	_ = app.readJson(w, r, &comment)

	switch commentType {
	case "user":
		{
			var err error
			image, err = repository.GetUserCommentImage(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("cannot get post image"), http.StatusInternalServerError)
				return
			}

			err = repository.DeleteUserComment(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("cannot delete post"), http.StatusInternalServerError)
				return
			}
		}
	case "group":
		{
			var err error
			image, err = repository.GetGroupCommentImage(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("cannot get group post image"), http.StatusInternalServerError)
				return
			}

			err = repository.DeleteGroupComment(comment)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("cannot delete group post"), http.StatusInternalServerError)
				return
			}
		}
	default:
		{
			_ = app.errorJson(w, fmt.Errorf("error creating a comment"), http.StatusBadRequest)
			return
		}
	}

	if image != "" {
		err := deleteImage(image, commentImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("cannot delete post image"), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, comment, nil)
}

// Group posts

func (app *application) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	if r.URL.Query().Has("id") { // Get group posts
		groupId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("getting groupId error %v", err)
			return
		}
		posts, err = repository.GetGroupPosts(groupId)
		if err != nil {
			log.Println("getPosts error: ", err)
		}
	} else {
		// Error?
	}
	_ = app.writeJson(w, http.StatusOK, posts, nil)
}

func (app *application) PostGroupPosts(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	var post models.Post
	// read response from the client
	_ = app.readJson(w, r, &post)

	post.UserId = user.ID
	post.CreatedAt = time.Now()
	post.CreatedBy = user.FirstName + " " + user.LastName

	if post.Img != "" {
		img, err := saveBase64Image(post.Img, postImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("can't save image error: %v", err), http.StatusUnauthorized)
			return
		}
		post.Img = img
	}
	err = repository.CreateGroupPost(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot create post"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, post, nil)
}

func (app *application) DeleteGroupPosts(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	// read response from the client
	_ = app.readJson(w, r, &post)

	image, err := repository.GetGroupPostImage(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot get group post image"), http.StatusInternalServerError)
		return
	}

	err = repository.DeleteGroupPost(post)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot delete group post"), http.StatusInternalServerError)
		return
	}

	if image != "" {
		err = deleteImage(image, groupImagePath)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("cannot delete post image"), http.StatusInternalServerError)
			return
		}
	}

	_ = app.writeJson(w, http.StatusOK, post, nil)
}

func (app *application) GetGroups(w http.ResponseWriter, r *http.Request) {
	var groups []models.Group
	groups, err := repository.GetAllGroups()
	if err != nil {
		log.Println("GetGroups error: ", err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, groups, nil)
}

func (app *application) PostGroups(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	_ = app.readJson(w, r, &group)
	group.CreatedBy = user.ID
	err = repository.GroupNameTaken(group.Title)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("group name taken"), http.StatusUnauthorized)
		return
	}
	if group.Img == "" {
		group.Img = "default_group.png"
	} else {
		imgName, err := saveBase64Image(group.Img, groupImagePath)
		if err != nil {
			log.Println(err)
		}
		group.Img = imgName
	}
	groupID, err := repository.CreateGroup(group)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("error creating group"), http.StatusUnauthorized)
		return
	}
	group.ID = groupID
	err = repository.AddUserToGroup(groupID, user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("error adding user to group"), http.StatusUnauthorized)
		return
	}
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) GetGroup(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	var group models.Group
	group, err = repository.GetGroup(groupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	err = repository.IsUserInGroup(groupId, user.ID)
	if err != nil {
		group.InGroup = false
	} else {
		group.InGroup = true
	}
	err = repository.HasRequested(groupId, user.ID)
	if err != nil {
		group.JoinRequested = false
	} else {
		group.JoinRequested = true
	}
	group.Img = fmt.Sprintf("http://localhost:8080/group/image/%d", groupId)
	group.InvitedUsersIds, err = repository.GetInvitedUsers(groupId)
	if err != nil {
		log.Println("Get group invited users error: ", err)
		return
	}
	group.InvitingUsersIds, err = repository.GetInvitingUsers(groupId)
	if err != nil {
		log.Println("Get group inviting users error: ", err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) GroupJoinHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	err = repository.AddGroupJoinRequest(groupId, user.ID)
	if err != nil {
		log.Printf("adding group join request error %v", err)
		return
	}
	var group models.Group
	group, err = repository.GetGroup(groupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	group.JoinRequested = true
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) RemoveGroupJoinHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	err = repository.RemoveGroupJoinRequest(groupId, user.ID)
	if err != nil {
		log.Printf("removing group join request error %v", err)
		return
	}
	var group models.Group
	group, err = repository.GetGroup(groupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	group.JoinRequested = false
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	err = repository.RemoveUserFromGroup(groupId, user.ID)
	if err != nil {
		log.Printf("removing group join request error %v", err)
		return
	}
	var group models.Group
	group, err = repository.GetGroup(groupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	group.JoinRequested = false
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) RemoveUserFromGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	err = repository.RemoveUserFromGroup(groupId, userId)
	if err != nil {
		log.Printf("removing group join request error %v", err)
		return
	}
	var group models.Group
	group, err = repository.GetGroup(groupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	group.JoinRequested = false
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

func (app *application) GroupInviteHandler(w http.ResponseWriter, r *http.Request) {
	var groupInvite models.GroupInvite
	_ = app.readJson(w, r, &groupInvite)
	err := repository.AddGroupInvites(groupInvite)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to add group invite"), http.StatusUnauthorized)
		return
	}
	group, err := repository.GetGroup(groupInvite.GroupId)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, group, nil)
}

// Followers

func (app *application) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// get cookie
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}

	// get user from the database
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	followers, err := repository.GetUserFollowers(user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to get followers"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, followers, nil)
}

func (app *application) GroupCreateEventHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	group, err := repository.GetGroup(groupId)
	if err != nil {
		log.Printf("getting group error %v", err)
		return
	}
	var event models.Event
	_ = app.readJson(w, r, &event)
	event.GroupId = groupId
	event.CreatedBy = user.ID
	event.CreatorName = user.FirstName + " " + user.LastName
	event.Img = group.Img
	event.GroupName = group.Title

	err = repository.CreateGroupEvent(event)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot create event"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, event, nil)

}

func (app *application) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		log.Printf("getting groupId error %v", err)
		return
	}
	var events []models.Event
	events, err = repository.GetGroupEvents(groupId)
	if err != nil {
		log.Println("Get Groups events error: ", err)
		return
	}
	// fmt.Println(events)
	_ = app.writeJson(w, http.StatusOK, events, nil)
}

func (app *application) PostGoingStatus(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil {
		log.Printf("getting eventId error %v", err)
		return
	}
	status, err := strconv.Atoi(r.URL.Query().Get("status"))
	if err != nil {
		log.Printf("getting status error %v", err)
		return
	}
	err = repository.AddEventGoingStatus(eventId, user.ID, status)
	if err != nil {
		log.Printf("adding going status error %v", err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}

func (app *application) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	var group models.Group
	// read response from the client
	_ = app.readJson(w, r, &group)

	err = repository.DeleteGroup(group)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cannot delete group"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, group, nil)

	err = repository.SetNotificationsSeen(user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to set notifications seen"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}

// Notifications
func (app *application) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	lastTimeStamp := r.URL.Query().Get("lastTimeStamp")
	notifications, err := repository.GetUserNotifications(user.ID, lastTimeStamp)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to fetch notification"), http.StatusInternalServerError)
		return
	}

	_ = app.writeJson(w, http.StatusOK, notifications, nil)
}

func (app *application) NotificationsSeenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	err = repository.SetNotificationsSeen(user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to set notifications seen"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}

func (app *application) NotificationActionHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Notification interface{} `json:"notification"`
		Action       string      `json:"action"`
	}
	_ = app.readJson(w, r, &requestBody)

	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}

	currentUser := user.ID
	notification := requestBody.Notification.(map[string]interface{})
	notificationId := notification["data"].(map[string]interface{})["id"].(float64)
	notificationType := notification["type"].(string)
	var groupJoinUser float64
	if optionalData, ok := notification["optionalData"]; ok && optionalData != nil {
		groupJoinUser = notification["optionalData"].(map[string]interface{})["id"].(float64)
	} else {
		groupJoinUser = 0
	}

	switch notificationType {
	case "follow_request":
		err := repository.HandleFollowRequest(currentUser, int(notificationId), requestBody.Action)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("failed to handle follow request"), http.StatusInternalServerError)
			return
		}
	case "group_join_request":
		if requestBody.Action == "accept" {
			err := repository.AddUserToGroup(int(notificationId), int(groupJoinUser))
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("failed to add user to group"), http.StatusInternalServerError)
				return
			}
		}
		err := repository.RemoveGroupJoinRequest(int(notificationId), int(groupJoinUser))
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("failed to remove group join request"), http.StatusInternalServerError)
			return
		}
	case "group_invitation":
		if requestBody.Action == "accept" {
			err := repository.AddUserToGroup(int(notificationId), currentUser)
			if err != nil {
				_ = app.errorJson(w, fmt.Errorf("failed to add user to group"), http.StatusInternalServerError)
				return
			}
		}

		err := repository.RemoveGroupInvite(int(notificationId), int(groupJoinUser), currentUser)
		if err != nil {
			_ = app.errorJson(w, fmt.Errorf("failed to remove group invite"), http.StatusInternalServerError)
			return
		}
	}

	_ = repository.UserActioned(currentUser, int(notificationId), int(groupJoinUser), notificationType)
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}

func (app *application) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil {
		log.Printf("getting eventId error %v", err)
		return
	}
	event, err := repository.GetEvent(eventId)
	if err != nil {
		log.Println("Get Groups event error: ", err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, event, nil)
}

func (app *application) ClearAllNotifications(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SN-Session")
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("cookie not found"), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUserByCookie(cookie.Value)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("user not found"), http.StatusUnauthorized)
		return
	}
	err = repository.ClearAllNotifications(user.ID)
	if err != nil {
		_ = app.errorJson(w, fmt.Errorf("failed to clear notifications"), http.StatusInternalServerError)
		return
	}
	_ = app.writeJson(w, http.StatusOK, nil, nil)
}
