package repository

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
	"time"
)

func GetUserNotifications(userID int, lastTimeStamp string) ([]models.Notification, error) {
	var notifications []models.Notification
	var rows *sql.Rows
	var err error

	t, _ := time.Parse(time.RFC3339, lastTimeStamp)
	unixTimestamp := t.Unix()
	rows, err = database.Query("SELECT id, userId, sourceId, type, seen, createdAt, actioned FROM notifications WHERE userId = ? AND (createdAt > ? OR ? = 'null')", userID, unixTimestamp, unixTimestamp)

	if err != nil {
		log.Println("Error retrieving notifications: ", err)
		return notifications, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.UserId, &notification.SourceId, &notification.Type, &notification.Seen, &notification.CreatedAt, &notification.Actioned)
		if err != nil {
			log.Println("Error scanning notifications: ", err)
			return notifications, err
		}

		switch notification.Type {
		case "follow_request":
			user, err := GetUserById(notification.SourceId)
			if err != nil {
				log.Println("Error getting user: ", err)
				return notifications, err
			}
			notification.Data = user

		case "group_invitation":
			data, optionalData, err := getGroupInvitation(userID, notification.SourceId, notification.ID, "group_invitation")
			if err != nil {
				log.Println("Error getting group invitation: ", err)
				return notifications, err
			}
			notification.Data = data
			notification.OptionalData = optionalData

		case "group_join_request":
			data, optionalData, err := getGroupInvitation(userID, notification.SourceId, notification.ID, "group_join_request")
			if err != nil {
				log.Println("Error getting group join request: ", err)
				return notifications, err
			}
			notification.Data = data
			notification.OptionalData = optionalData

		case "group_event":
			data, err := GetEvent(notification.SourceId)
			if err != nil {
				log.Println("Error getting event: ", err)
				return notifications, err
			}
			notification.Data = data
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func getGroupInvitation(userID, groupID, notificationID int, notificationType string) (models.Group, models.User, error) {
	group, err := GetGroup(groupID)
	if err != nil {
		log.Println("Error getting group: ", err)
		return group, models.User{}, err
	}

	var groupJoinerId int
	if notificationType == "group_join_request" {
		err = database.QueryRow("SELECT groupJoinerId FROM notifications WHERE userId = ? AND sourceId = ? AND type = ? AND id = ?", userID, groupID, "group_join_request", notificationID).Scan(&groupJoinerId)
	} else if notificationType == "group_invitation" {
		err = database.QueryRow("SELECT groupJoinerId FROM notifications WHERE userId = ? AND sourceId = ? AND type = ? AND id = ?", userID, groupID, "group_invitation", notificationID).Scan(&groupJoinerId)
	} else {
		return group, models.User{}, fmt.Errorf("invalid notification type")
	}

	if err != nil {
		log.Println("Error getting groupJoinerId: ", err)
		return group, models.User{}, err
	}

	groupJoiner, err := GetUserById(groupJoinerId)
	if err != nil {
		log.Println("Error getting groupJoiner: ", err)
		return group, models.User{}, err
	}

	return group, groupJoiner, nil
}

func SetNotificationsSeen(userID int) error {
	_, err := database.Exec("UPDATE notifications SET seen = 1 WHERE userId = ?", userID)
	if err != nil {
		log.Println("Error setting notifications seen: ", err)
		return err
	}

	return nil
}

func HandleFollowRequest(currentUser, id int, action string) error {
	if action == "accept" {
		_, err := database.Exec("UPDATE followers SET requested = 0 WHERE userId = ? and followerId = ?", currentUser, id)
		if err != nil {
			log.Println("Error accepting follow request: ", err)
			return err
		}
	} else {
		_, err := database.Exec("DELETE FROM followers WHERE userId = ? and followerId = ?", currentUser, id)
		if err != nil {
			log.Println("Error declining follow request: ", err)
		}
	}

	return nil
}

func UserActioned(currentUser, notificationId, groupJoinUser int, notificationType string) error {
	_, err := database.Exec("UPDATE notifications SET actioned = 1 WHERE userId = ? and sourceId = ? and type = ? and (groupJoinerId = ? OR ? != 'group_join_request')", currentUser, notificationId, notificationType, groupJoinUser, notificationType)
	if err != nil {
		log.Println("Error setting user actioned: ", err)
		return err
	}
	return nil
}

func ClearAllNotifications(currentUser int) error {
	_, err := database.Exec("DELETE FROM notifications WHERE userId = ?", currentUser)
	if err != nil {
		log.Println("Error clearing notifications: ", err)
		return err
	}
	return nil
}
