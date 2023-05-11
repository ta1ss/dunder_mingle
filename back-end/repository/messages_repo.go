package repository

import (
	"database/sql"
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
	"time"
)

func GetMessages(user1, user2 int) ([]models.Message, error) {
	var messages []models.Message
	rows, err := database.Query(`
	SELECT
		user_messages.*,
		user1.firstName || ' ' || user1.lastName AS userFullName,
		user2.firstName || ' ' || user2.lastName AS targetFullName
	FROM 
		user_messages
	INNER JOIN 
		users AS user1 ON user_messages.senderId = user1.id
	INNER JOIN 
		users AS user2 ON user_messages.targetId = user2.id
	WHERE
		(user_messages.senderId = ? AND user_messages.targetId = ?)
	OR
		(user_messages.senderId = ? AND user_messages.targetId = ?)
	ORDER BY
		user_messages.createdAt DESC;
	`, user1, user2, user2, user1)
	if err != nil {
		log.Println("GetMessages error: ", err)
		return messages, err
	}

	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.SenderId,
			&message.TargetId,
			&message.Body,
			&message.MessageRead,
			&message.CreatedAt,
			&message.UserName,
			&message.TargetName)
		if err != nil {
			log.Println("GetMessages error: ", err)
			return messages, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func CreateUserMessage(message models.Message) error {
	var read int
	if message.SenderId == message.TargetId {
		read = 1
	}
	_, err := database.Exec(`
	INSERT INTO
		user_messages (senderId, targetId, body, messageRead)
	VALUES (?, ?, ?, ?)`, message.SenderId, message.TargetId, message.Body, read)
	if err != nil {
		log.Println("Create user message: ", err)
		return err
	}

	return nil
}

func GetLastMessageTime(user1, user2 int) (int, int, time.Time, error) {
	var senderId int
	var messageRead int
	var createdAt time.Time

	row := database.QueryRow(`
		SELECT 
			senderId, messageRead, createdAt
		FROM
			user_messages
		WHERE
			(user_messages.senderId = ? AND user_messages.targetId = ?)
		OR
			(user_messages.senderId = ? AND user_messages.targetId = ?)
		ORDER BY
			createdAt DESC
		LIMIT 1
	`, user1, user2, user2, user1)

	err := row.Scan(&senderId, &messageRead, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return senderId, 1, time.Unix(0, 0), nil
		}
		return 0, 1, time.Unix(0, 0), err
	}

	return senderId, messageRead, createdAt, nil
}

func MessagesRead(senderId, targetId int) error {
	stmt := `
	UPDATE 
		user_messages
	SET 
		messageRead = 1
	WHERE
		targetId = ?
	AND
		senderId = ?
	`
	_, err := database.Exec(stmt, targetId, senderId)

	if err != nil {
		return err
	}

	return nil
}