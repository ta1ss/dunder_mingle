package repository

import (
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
)

func GetGroupPostComments(id int) ([]models.Comment, error) {
	var comments []models.Comment
	rows, err := database.Query(`
	SELECT
		group_comments.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img
	FROM 
		group_comments
	INNER JOIN 
		users ON group_comments.userId = users.id
	WHERE
		group_comments.postId = ?
	ORDER BY
		group_comments.createdAt ASC;
	`, id)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return comments, err
	}

	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Body,
			&comment.Img,
			&comment.CreatedAt,
			&comment.CreatedBy,
			&comment.UserImg)
		if err != nil {
			log.Println("GetAllPosts error: ", err)
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateGroupComment(comment models.Comment) (int, error) {
	result, err := database.Exec(`
	INSERT INTO
		group_comments (postId, userId, body, img)
	VALUES (?, ?, ?, ?)`, comment.PostId, comment.UserId, comment.Body, comment.Img)
	if err != nil {
		log.Println("Create group post comment: ", err)
		return 0, err
	}

	commentId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentId), nil
}

func DeleteGroupComment(comment models.Comment) error {
	_, err := database.Exec(`
	DELETE FROM
		group_comments
	WHERE
		id = ?
	AND
		userId = ?
	`, comment.ID, comment.UserId)
	if err != nil {
		log.Println("Can't delete group comment: ", err)
		return err
	}

	return nil
}

func GetGroupCommentImage(comment models.Comment) (string, error) {
	err := database.QueryRow(`
	SELECT
		img
	FROM
		group_comments
	WHERE
		id = ?
	`, comment.ID).Scan(&comment.Img)
	if err != nil {
		log.Println("Can't get group comment image: ", err)
		return "", err
	}

	return comment.Img, nil
}