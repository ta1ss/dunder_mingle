package repository

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
	"time"
)

func EmailTaken(email string) error {
	var taken int
	err := database.QueryRow("SELECT COUNT (*) FROM users WHERE email = $1", email).Scan(&taken)
	if err != nil {
		return err
	}
	if taken != 0 {
		return fmt.Errorf("e-mail taken")
	}
	return nil
}

func UpdateUserStatus(userID int, status int) error {
	_, err := database.Exec("UPDATE users SET profilePublic = ? WHERE id = ? ", status, userID)
	if err != nil {
		fmt.Println("Error updating user status: ", err)
		return err
	}

	return nil
}

func GetUserFollow(userID int) (models.UserFollow, error) {
	var userFollow models.UserFollow

	followers, err := GetUserFollowers(userID)
	if err != nil {
		log.Println("Error getting user followers: ", err)
		return userFollow, err
	}
	following, err := GetUserFollowing(userID)
	if err != nil {
		log.Println("Error getting user following: ", err)
		return userFollow, err
	}

	userFollow.Followers = followers
	userFollow.Following = following

	return userFollow, nil
}

func GetUserFollowers(userID int) ([]models.UserFollower, error) {
	var userFollowers []models.UserFollower

	followerRows, err := database.Query(`
		SELECT f.userId, f.followerId, u.firstName || ' ' || u.lastName AS followerName, u.img AS followerImg
		FROM followers f
		JOIN users u ON u.id = f.followerId
		WHERE f.userId = ? AND f.requested = 0
		ORDER BY u.lastName ASC
	`, userID)

	if err != nil {
		log.Println("Error getting user followers: ", err)
		return userFollowers, err
	}
	defer followerRows.Close()

	for followerRows.Next() {
		var userFollower models.UserFollower
		err := followerRows.Scan(&userFollower.UserID, &userFollower.FollowerID, &userFollower.FollowerName, &userFollower.UserImg)
		if err != nil {
			log.Println("Error scanning user followers: ", err)
			return userFollowers, err
		}
		userFollowers = append(userFollowers, userFollower)
	}

	return userFollowers, nil
}

func GetUserFollowing(userID int) ([]models.UserFollowing, error) {
	var userFollowing []models.UserFollowing

	followingRows, err := database.Query(`
	SELECT f.followerId as userId, f.userId as followingId, u.firstName || ' ' || u.lastName AS followingName, u.img AS followerImg
	FROM followers f
	JOIN users u ON u.id = f.userId
	WHERE f.followerId = ?
	ORDER BY u.lastName ASC
	`, userID)

	if err != nil {
		log.Println("Error getting user following: ", err)
		return userFollowing, err
	}
	defer followingRows.Close()

	for followingRows.Next() {
		var userFollow models.UserFollowing
		err := followingRows.Scan(&userFollow.UserID, &userFollow.FollowingID, &userFollow.FollowingName, &userFollow.UserImg)
		if err != nil {
			log.Println("Error scanning user following: ", err)
			return userFollowing, err
		}
		userFollowing = append(userFollowing, userFollow)
	}

	return userFollowing, nil
}

func GetUserById(id int) (models.User, error) {
	var user models.User
	row := database.QueryRow(`
		SELECT
			id, email, firstName, lastName, dateOfBirth, img, nickname, about, profilePublic, createdAt, updatedAt
		FROM
			users
		WHERE id = ?
	`, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Image,
		&user.Nickname,
		&user.About,
		&user.ProfileP,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		log.Println("getUserById error: ", err)
		return user, err
	}
	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Image, &user.Nickname, &user.About, &user.ProfileP, &user.CreatedAt, &user.UpdatedAt)
	user.Password = "contact admin to purchase password"
	if err != nil {
		log.Println("getUserByEmail error: ", err)
		return user, err
	}
	return user, nil
}

func GetUserByCookie(cookie string) (models.User, error) {
	var user models.User
	userEmail := getEmailByUUID(cookie)
	if userEmail == "" {
		return user, fmt.Errorf("user not found")
	}
	user, err := GetUserByEmail(userEmail)
	if err != nil {
		log.Println("getUserByCookie error: ", err)
		return user, err
	}
	user.Password = "contact admin to purchase password"
	return user, nil
}

func getEmailByUUID(uuid string) string {
	var email string
	var userID int
	// get the user id from the sessions table
	err := database.QueryRow("SELECT userId FROM sessions where uuid = ?", uuid).Scan(&userID)
	if err != nil {
		log.Println("Retrieving userID error: ", err)
		return ""
	}
	// get the email from the users table
	err = database.QueryRow("SELECT email FROM users where id = ?", userID).Scan(&email)
	if err != nil {
		log.Println("Retrieving email error: ", err)
		return ""
	}
	return email
}

func GetFollowStatus(currentUser, followedUser int) (models.FollowStatus, error) {
	var followStatus models.FollowStatus
	var requested int
	var isFollowing bool
	var isRequested bool

	row := database.QueryRow("SELECT requested FROM followers WHERE userId = ? AND followerId = ?", followedUser, currentUser).Scan(&requested)
	if row == sql.ErrNoRows {
		isFollowing = false
		isRequested = false
	} else if requested == 1 {
		isFollowing = false
		isRequested = true
	} else {
		isFollowing = true
		isRequested = false
	}

	followStatus.IsFollowing = isFollowing
	followStatus.IsRequested = isRequested

	return followStatus, nil
}

func ChangeFollowStatus(currentUser, followedUser int64, isRequested, isFollowed bool) (models.UserFollowRequest, error) {
	var userFollowRequest models.UserFollowRequest

	isRequestedInt := 0
	if isRequested {
		isRequestedInt = 1
	}

	if isFollowed {
		fmt.Println("adding new follower")
		_, err := database.Exec("INSERT INTO followers (userId, followerId) VALUES (?, ?)", followedUser, currentUser)
		if err != nil {
			log.Println("Error adding follower: ", err)
			return userFollowRequest, err
		}
		userFollowRequest.IsFollowed = true
	} else if !isFollowed && !isRequested {
		fmt.Println("removing follower")
		_, err := database.Exec("DELETE FROM followers WHERE userId = ? AND followerId = ?", followedUser, currentUser)
		if err != nil {
			log.Println("Error removing follower: ", err)
			return userFollowRequest, err
		}
		userFollowRequest.IsFollowed = false
	} else if isRequested {
		fmt.Println("adding follower with the request")
		_, err := database.Exec("INSERT INTO followers (userId, followerId, requested) VALUES (?, ?, ?)", followedUser, currentUser, isRequestedInt)
		if err != nil {
			log.Println("Error adding follower request: ", err)
			return userFollowRequest, err
		}
		_, err = database.Exec("INSERT INTO notifications (userId, sourceId, type, createdAt) VALUES (?, ?, ?, ?)", followedUser, currentUser, "follow_request", time.Now().Unix())
		if err != nil {
			log.Println("Error adding notification: ", err)
			return userFollowRequest, err
		}
		userFollowRequest.IsRequested = true
	}

	return userFollowRequest, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := database.Query(`
	SELECT id, email, firstName, lastName, dateOfBirth, img, nickname, about, profilePublic, createdAt, updatedAt
	FROM users
	ORDER BY lastName ASC
	`)
	if err != nil {
		log.Println("Error getting all users: ", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Image, &user.Nickname, &user.About, &user.ProfileP, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning user: ", err)
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserFullNameById(id int) (string, error) {
	var name string
	err := database.QueryRow(`
		SELECT
			firstName || ' ' || lastName
		FROM
			users
		WHERE id = ?
	`, id).Scan(&name)

	if err != nil {
		log.Println("getUserFullNameById error: ", err)
		return "", err
	}
	return name, nil
}
