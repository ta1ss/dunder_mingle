package repository

import (
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
)

func GetUserPost(id, userId int) (models.Post, error) {
	var post models.Post
	err := database.QueryRow(`
	SELECT
		user_posts.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img
	FROM 
		user_posts
	INNER JOIN 
		users ON user_posts.userId = users.id
	LEFT JOIN
		followers ON user_posts.id = followers.userId AND followers.requested = 0
	LEFT JOIN
		user_custom_posts ON user_posts.id = user_custom_posts.postId
	WHERE
		user_posts.id = ?
	AND
		((user_posts.userId = ?)
	OR
		(users.profilePublic = 1 AND user_posts.privacy = 'Public')
	OR
		(users.profilePublic = 0 AND user_posts.privacy = 'Public' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Private' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Custom' AND user_custom_posts.postId = user_posts.id AND user_custom_posts.targetId = ?))
	`, id, userId, userId, userId, userId).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Body,
		&post.Img,
		&post.Privacy,
		&post.CreatedAt,
		&post.CreatedBy,
		&post.UserImg)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return post, err
	}

	return post, nil
}

func GetAllPosts(id int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := database.Query(`
	SELECT DISTINCT
		user_posts.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img,
		COUNT(DISTINCT user_comments.id) AS comments
	FROM 
		user_posts
	INNER JOIN 
		users ON user_posts.userId = users.id
	LEFT JOIN
		followers ON user_posts.id = followers.userId AND followers.requested = 0
	LEFT JOIN
		user_custom_posts ON user_posts.id = user_custom_posts.postId
	LEFT JOIN 
		user_comments ON user_posts.id = user_comments.postId
	WHERE
		(user_posts.userId = ?)
	OR
		(users.profilePublic = 1 AND user_posts.privacy = 'Public')
	OR
		(users.profilePublic = 0 AND user_posts.privacy = 'Public' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Private' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Custom' AND user_custom_posts.postId = user_posts.id AND user_custom_posts.targetId = ?)
	GROUP BY
		user_posts.id
	ORDER BY
		user_posts.createdAt DESC;
	`, id, id, id, id)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Body,
			&post.Img,
			&post.Privacy,
			&post.CreatedAt,
			&post.CreatedBy,
			&post.UserImg,
			&post.Comments)
		if err != nil {
			log.Println("GetAllPosts error: ", err)
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func GetUserPosts(targetId, id int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := database.Query(`
	SELECT DISTINCT
		user_posts.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img,
		COUNT(DISTINCT user_comments.id) AS comments
	FROM 
		user_posts
	INNER JOIN 
		users ON user_posts.userId = users.id
	LEFT JOIN
		followers ON user_posts.id = followers.userId AND followers.requested = 0
	LEFT JOIN
		user_custom_posts ON user_posts.id = user_custom_posts.postId
	LEFT JOIN 
		user_comments ON user_posts.id = user_comments.postId
	WHERE 
		user_posts.userId = ?
	AND
		((? = ?)
	OR
		(users.profilePublic = 1 AND user_posts.privacy = 'Public')
	OR
		(users.profilePublic = 0 AND user_posts.privacy = 'Public' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Private' AND followers.userId = user_posts.userId AND followers.followerId = ?)
	OR
		(user_posts.privacy = 'Custom' AND user_custom_posts.postId = user_posts.id AND user_custom_posts.targetId = ?))
	GROUP BY
		user_posts.id
	ORDER BY
		user_posts.createdAt DESC;
	`, targetId, targetId, id, id, id, id)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Body,
			&post.Img,
			&post.Privacy,
			&post.CreatedAt,
			&post.CreatedBy,
			&post.UserImg,
			&post.Comments)
		if err != nil {
			log.Println("GetUserPosts error: ", err)
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func GetGroupPosts(id int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := database.Query(`
	SELECT 
		group_posts.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img,
		COUNT(DISTINCT group_comments.id) AS comments
	FROM 
		group_posts
	INNER JOIN 
		users ON group_posts.userId = users.id
	LEFT JOIN 
		group_comments ON group_posts.id = group_comments.postId
	WHERE 
		groupId = ?
	GROUP BY
		group_posts.id
	ORDER BY 
		createdAt DESC
	`, id)
	if err != nil {
		log.Println("GetGroupPosts error: ", err)
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.GroupId,
			&post.UserId,
			&post.Title,
			&post.Body,
			&post.Img,
			&post.CreatedAt,
			&post.CreatedBy,
			&post.UserImg,
			&post.Comments)
		if err != nil {
			log.Println("GetGroupPosts error: ", err)
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func GetGroupPost(id, userId int) (models.Post, error) {
	var post models.Post
	err := database.QueryRow(`
	SELECT
		group_posts.*,
		users.firstName || ' ' || users.lastName AS fullName,
		users.img
	FROM 
		group_posts
	INNER JOIN 
		users ON group_posts.userId = users.id
	LEFT JOIN
		user_groups ON group_posts.groupId = user_groups.groupId
	WHERE
		group_posts.id = ?
	AND
		user_groups.userId = ?
	`, id, userId).Scan(
		&post.ID,
		&post.GroupId,
		&post.UserId,
		&post.Title,
		&post.Body,
		&post.Img,
		&post.CreatedAt,
		&post.CreatedBy,
		&post.UserImg)
	if err != nil {
		log.Println("GetGroupPost error: ", err)
		return post, err
	}

	return post, nil
}

func CreateUserPost(post models.Post) (int, error) {
	result, err := database.Exec(`
	INSERT INTO
		user_posts (userId, title, body, img, privacy)
	VALUES (?, ?, ?, ?, ?)`, post.UserId, post.Title, post.Body, post.Img, post.Privacy)
	if err != nil {
		log.Println("Create user post error: ", err)
		return 0, err
	}

	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postId), nil
}

func GetUserPostImage(post models.Post) (string, error) {
	err := database.QueryRow(`
	SELECT
		img
	FROM
		user_posts
	WHERE
		id = ?
	`, post.ID).Scan(&post.Img)
	if err != nil {
		log.Println("Can't get post image: ", err)
		return "", err
	}

	return post.Img, nil
}

func DeleteUserPost(post models.Post) error {
	_, err := database.Exec(`
	DELETE FROM
		user_posts 
	WHERE
		id = ?
	AND
		userId = ?
	`, post.ID, post.UserId)
	if err != nil {
		log.Println("Can't delete post: ", err)
		return err
	}

	return nil
}

func GetGroupPostImage(post models.Post) (string, error) {
	err := database.QueryRow(`
	SELECT
		img
	FROM
		group_posts
	WHERE
		id = ?
	`, post.ID).Scan(&post.Img)
	if err != nil {
		log.Println("Can't get post image: ", err)
		return "", err
	}

	return post.Img, nil
}

func DeleteGroupPost(post models.Post) error {
	_, err := database.Exec(`
	DELETE FROM
		group_posts 
	WHERE
		id = ?
	AND
		userId = ?
	`, post.ID, post.UserId)
	if err != nil {
		log.Println("Can't delete post: ", err)
		return err
	}

	return nil
}

func AddCustomPrivacy(postId, targetId int) error {
	_, err := database.Exec(`
	INSERT INTO
		user_custom_posts (postId, targetId)
	VALUES (?, ?)`, postId, targetId)
	if err != nil {
		log.Println("Create user post error: ", err)
		return err
	}

	return nil
}
