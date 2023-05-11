package repository

import (
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
)

func GetGroupChatMessages(groupId int) ([]models.Message, error) {
	var messages []models.Message
	rows, err := database.Query(`
	SELECT
		group_messages.groupId, group_messages.userId, group_messages.body, group_messages.createdAt, users.firstName || ' ' || users.lastName AS fullName
	FROM
		group_messages
	INNER JOIN
		users ON group_messages.userId = users.id
	WHERE
		group_messages.groupId = ?
	ORDER BY
		group_messages.createdAt DESC;
	`, groupId)
	if err != nil {
		log.Println("Get Group Messages error: ", err)
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.GroupId,
			&message.SenderId,
			&message.Body,
			&message.CreatedAt,
			&message.UserName)
		if err != nil {
			log.Println("Get group Messages error: ", err)
			return messages, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func CreateGroupMessage(message models.Message) error {
	_, err := database.Exec(`
	INSERT INTO
		group_messages (groupId, userId, body)
	VALUES (?, ?, ?)`, message.GroupId, message.SenderId, message.Body)
	if err != nil {
		log.Println("Create group message: ", err)
		return err
	}

	return nil
}
