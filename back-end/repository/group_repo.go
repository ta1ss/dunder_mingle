package repository

import (
	"fmt"
	"log"
	"social-network/back-end/database"
	"social-network/back-end/models"
	"time"
)

func GroupNameTaken(title string) error {
	var taken int
	err := database.QueryRow("SELECT COUNT (*) FROM groups WHERE title = $1", title).Scan(&taken)
	if err != nil {
		return err
	}
	if taken != 0 {
		return fmt.Errorf("group name taken")
	}
	return nil
}

func CreateGroup(group models.Group) (groupID int, err error) {
	result, err := database.Exec("INSERT INTO groups (`createdBy`,`title`,`description`,`img`) VALUES (?,?,?,?)", group.CreatedBy, group.Title, group.Description, group.Img)
	if err != nil {
		return 0, err
	}
	userIDint64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userIDint64), nil
}

func GetAllGroups() ([]models.Group, error) {
	var groups []models.Group
	rows, err := database.Query(`SELECT * FROM groups ORDER BY groups.createdAt DESC;`)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID,
			&group.CreatedBy,
			&group.Title,
			&group.Description,
			&group.Img,
			&group.CreatedAt)
		if err != nil {
			log.Println("GetAllGroups error: ", err)
			return groups, err
		}
		users, err := GetAllGroupMembers(group.ID)
		if err != nil {
			log.Println("GetAllGroups error: ", err)
			return groups, err
		}
		group.GroupMembers = users
		groups = append(groups, group)
	}
	return groups, nil
}

// why is calling all the time?
func GetAllGroupMembers(groupId int) ([]models.User, error) {
	var users []models.User
	var userIds []int
	rows, err := database.Query(`SELECT userId FROM user_groups WHERE groupId = ?`, groupId)
	if err != nil {
		log.Println("GetAllGroupMembers error: ", err)
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			log.Println("GetAllGroupMembers error: ", err)
			return users, err
		}
		userIds = append(userIds, userId)
		// users = append(users, user)
	}
	// fmt.Println(userIds)
	for i := 0; i < len(userIds); i++ {
		user, err := GetUserById(userIds[i])
		if err != nil {
			return users, nil
		}
		users = append(users, user)

	}
	return users, nil
}

func GetGroup(groupID int) (models.Group, error) {
	var group models.Group
	rows, err := database.Query(`SELECT * FROM groups WHERE id = ?`, groupID)
	if err != nil {
		log.Println("GetGroup error: ", err)
		return group, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&group.ID,
			&group.CreatedBy,
			&group.Title,
			&group.Description,
			&group.Img,
			&group.CreatedAt)
		if err != nil {
			log.Println("GetAllPosts error: ", err)
			return group, err
		}
	}
	users, err := GetAllGroupMembers(group.ID)
	if err != nil {
		log.Println("GetAllGroups error: ", err)
		return group, err
	}
	group.GroupMembers = users
	return group, nil
}

func IsUserInGroup(groupId int, userId int) error {
	var inGroup int
	err := database.QueryRow("SELECT COUNT(*) FROM user_groups WHERE groupId = $1 AND userId = $2", groupId, userId).Scan(&inGroup)
	if err != nil {
		return err
	}
	if inGroup == 0 {
		return fmt.Errorf("not in group")
	}
	return nil
}

func AddGroupJoinRequest(groupId int, userId int) error {
	_, err := database.Exec("INSERT INTO group_join_requests (`groupId`,`userId`) VALUES (?,?)", groupId, userId)
	if err != nil {
		return err
	}

	group, _ := GetGroup(groupId)
	groupCreator := group.CreatedBy

	_, err = database.Exec("INSERT INTO notifications (userId, sourceId, type, createdAt, groupJoinerId) VALUES (?, ?, ?, ?, ?)", groupCreator, groupId, "group_join_request", time.Now().Unix(), userId)
	if err != nil {
		log.Println("group_join_request notification add error: ", err)
		return err
	}

	return nil
}

func RemoveGroupJoinRequest(groupId int, userId int) error {
	_, err := database.Exec("DELETE FROM group_join_requests WHERE groupId = ? AND userId = ?", groupId, userId)
	if err != nil {
		return err
	}
	_, err = database.Exec("DELETE FROM notifications WHERE sourceId = ? AND groupJoinerId = ? AND type = ?", groupId, userId, "group_join_request")
	if err != nil {
		log.Println("group_join_request notification delete error: ", err)
		return err
	}
	return nil
}

func HasRequested(groupId, userId int) error {
	var requested int
	err := database.QueryRow("SELECT COUNT(*) FROM group_join_requests WHERE groupId = $1 AND userId = $2", groupId, userId).Scan(&requested)
	if err != nil {
		return err
	}
	if requested == 0 {
		return fmt.Errorf("already requested")
	}
	return nil
}

func AddUserToGroup(groupId int, userId int) error {
	_, err := database.Exec(`INSERT INTO user_groups (groupId, userId)
	SELECT * FROM (SELECT ?, ?) AS tmp
	WHERE NOT EXISTS (
		SELECT 1 FROM user_groups 
		WHERE groupId = ? AND userId = ?
	) LIMIT 1;`, groupId, userId, groupId, userId)
	if err != nil {
		return err
	}
	return nil
}

func CreateGroupPost(post models.Post) error {
	_, err := database.Exec(`
	INSERT INTO
		group_posts (groupId, userId, title, body, img)
	VALUES (?, ?, ?, ?, ?)`, post.GroupId, post.UserId, post.Title, post.Body, post.Img)
	if err != nil {
		log.Println("Create group post error: ", err)
		return err
	}
	return nil
}

func CreateGroupEvent(event models.Event) error {
	eventID, err := database.Exec(`
	INSERT INTO
		group_events (groupId, createdBy, title, description, dateStart, dateEnd, creatorName, img, groupName)
	VALUES (?, ?, ?, ?, ?, ?, ?,?,?);
	SELECT last_insert_rowid(); `, event.GroupId, event.CreatedBy, event.Title, event.Description, event.EventStart.Format("2006-01-02 15:04:05"), event.EventEnd.Format("2006-01-02 15:04:05"), event.CreatorName, event.Img, event.GroupName)
	if err != nil {
		log.Println("Create group event error: ", err)
		return err
	}
	eventIDInt, err := eventID.LastInsertId()
	if err != nil {
		return err
	}
	members, _ := GetAllGroupMembers(event.GroupId)
	// fmt.Println("Adding notification to members: ")
	for i := 0; i < len(members); i++ {
		// fmt.Println(members[i].ID)
		if members[i].ID != event.CreatedBy {
			_, err := database.Exec("INSERT INTO notifications (userId, sourceId, type, createdAt, actioned) VALUES (?, ?, ?, ?, ?)", members[i].ID, int(eventIDInt), "group_event", time.Now().Unix(), "1")
			if err != nil {
				log.Println("Create group event notificaiton error: ", err)
				return err
			}
		}
	}
	return nil
}

func GetGroupEvents(groupID int) ([]models.Event, error) {
	var events []models.Event
	rows, err := database.Query(`SELECT * FROM group_events WHERE groupId = ? ORDER BY group_events.createdAt DESC;`, groupID)
	if err != nil {
		log.Println("GetGroupEvents error: ", err)
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.GroupId,
			&event.CreatedBy,
			&event.Title,
			&event.Description,
			&event.EventStart,
			&event.EventEnd,
			&event.CreatedAt,
			&event.CreatorName,
			&event.Img,
			&event.GroupName,
		)
		if err != nil {
			log.Println("GetGroupEvents error: ", err)
			return events, err
		}
		var goingUsers []models.User
		var notGoingUsers []models.User
		rows, err := database.Query(`SELECT userId FROM user_events WHERE eventId = ? AND status = 1`, event.ID)
		if err != nil {
			log.Println("get usetId by events error: ", err)
			return events, err
		}
		defer rows.Close()
		for rows.Next() {
			var goingUser int
			err := rows.Scan(&goingUser)
			if err != nil {
				log.Println("get usetId by events error: ", err)
				return events, err
			}
			user, err := GetUserById(goingUser)
			if err != nil {
				return events, nil
			}
			goingUsers = append(goingUsers, user)
		}
		rows, err = database.Query(`SELECT userId FROM user_events WHERE eventId = ? AND status = 0`, event.ID)
		if err != nil {
			log.Println("get usetId by events error: ", err)
			return events, err
		}
		defer rows.Close()
		for rows.Next() {
			var notGoingUser int
			err := rows.Scan(&notGoingUser)
			if err != nil {
				log.Println("get usetId by events error: ", err)
				return events, err
			}
			user, err := GetUserById(notGoingUser)
			if err != nil {
				return events, nil
			}
			notGoingUsers = append(notGoingUsers, user)
		}
		event.GoingUsers = goingUsers
		event.NotGoingUsers = notGoingUsers
		events = append(events, event)
	}
	return events, nil
}

func AddEventGoingStatus(eventId int, userId int, status int) error {
	var hasStatus int
	err := database.QueryRow("SELECT COUNT(*) FROM user_events WHERE eventId = $1 AND userId = $2", eventId, userId).Scan(&hasStatus)
	if err != nil {
		return err
	}
	if hasStatus == 0 {
		_, err := database.Exec(`
			INSERT INTO
				user_events (eventId, userId, status)
			VALUES (?, ?, ?)`, eventId, userId, status)
		if err != nil {
			log.Println("Event going status error: ", err)
			return err
		}
	} else {
		_, err := database.Exec(`
			UPDATE
				user_events 
				SET  status = CASE
				WHEN status = 0 THEN 1
				WHEN status = 1 THEN 0
				ELSE status
				END
			WHERE status IN (0, 1) AND eventId = ? AND userId = ?`, eventId, userId)
		if err != nil {
			log.Println("Event going status error: ", err)
			return err
		}
	}
	return nil
}

func DeleteGroup(group models.Group) error {
	_, err := database.Exec(`DELETE FROM groups WHERE id = ?`, group.ID)
	if err != nil {
		log.Println("Can't delete post: ", err)
		return err
	}

	return nil
}

func AddGroupInvites(groupInvite models.GroupInvite) error {
	for i := 0; i < len(groupInvite.InvitedUsersIds); i++ {
		_, err := database.Exec("INSERT INTO group_invitations (`groupId`,`userId`, invitedUserId) VALUES (?,?,?)", groupInvite.GroupId, groupInvite.UserId, groupInvite.InvitedUsersIds[i])
		if err != nil {
			return err
		}
		_, err = database.Exec("INSERT INTO notifications (userId, sourceId, type, createdAt, groupJoinerId) VALUES (?, ?, ?, ?, ?)", groupInvite.InvitedUsersIds[i], groupInvite.GroupId, "group_invitation", time.Now().Unix(), groupInvite.UserId)
		if err != nil {
			log.Println("Create group invitation notificaiton error: ", err)
			return err
		}
	}
	return nil
}

func RemoveGroupInvite(groupId, userId, invitedUserId int) error {
	_, err := database.Exec("DELETE FROM group_invitations WHERE groupId = ? AND userId = ? AND invitedUserId = ?", groupId, userId, invitedUserId)
	if err != nil {
		return err
	}
	var otherInvitesCount int
	err = database.QueryRow("SELECT COUNT (*) FROM group_invitations WHERE groupId = ? AND invitedUserId = ?", groupId, invitedUserId).Scan(&otherInvitesCount)
	if err != nil {
		return err
	}
	if otherInvitesCount != 0 {
		_, err := database.Exec("DELETE FROM group_invitations WHERE groupId IN (?) AND invitedUserId IN (?)", groupId, invitedUserId)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetInvitedUsers(groupId int) ([]int, error) {
	var invitedUsersIds []int
	rows, err := database.Query(`SELECT invitedUserId FROM group_invitations WHERE groupId = ?`, groupId)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return invitedUsersIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			log.Println("GetAllGroups error: ", err)
			return invitedUsersIds, err
		}
		invitedUsersIds = append(invitedUsersIds, userId)
	}
	return invitedUsersIds, nil
}

func GetInvitingUsers(groupId int) ([]int, error) {
	var invitingUsersIds []int
	rows, err := database.Query(`SELECT userId FROM group_invitations WHERE groupId = ?`, groupId)
	if err != nil {
		log.Println("GetUserPosts error: ", err)
		return invitingUsersIds, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			log.Println("GetAllGroups error: ", err)
			return invitingUsersIds, err
		}
		invitingUsersIds = append(invitingUsersIds, userId)
	}
	return invitingUsersIds, nil
}

func GetEvent(eventID int) (models.Event, error) {
	var event models.Event
	err := database.QueryRow(`SELECT * FROM group_events WHERE id = ?`, eventID).Scan(
		&event.ID,
		&event.GroupId,
		&event.CreatedBy,
		&event.Title,
		&event.Description,
		&event.EventStart,
		&event.EventEnd,
		&event.CreatedAt,
		&event.CreatorName,
		&event.Img,
		&event.GroupName,
	)
	if err != nil {
		log.Println("GetGroupEvents error: ", err)
		return event, err
	}
	var goingUsers []models.User
	rows, err := database.Query(`SELECT userId FROM user_events WHERE eventId = ? AND status = 1`, eventID)
	if err != nil {
		log.Println("get usetId by events error: ", err)
		return event, err
	}
	defer rows.Close()
	for rows.Next() {
		var goingUser int
		err := rows.Scan(&goingUser)
		if err != nil {
			log.Println("get usetId by events error: ", err)
			return event, err
		}
		user, err := GetUserById(goingUser)
		if err != nil {
			return event, nil
		}
		goingUsers = append(goingUsers, user)
	}
	event.GoingUsers = goingUsers
	var notGoingUsers []models.User
	rows, err = database.Query(`SELECT userId FROM user_events WHERE eventId = ? AND status = 0`, eventID)
	if err != nil {
		log.Println("get usetId by events error: ", err)
		return event, err
	}
	defer rows.Close()
	for rows.Next() {
		var notGoingUser int
		err := rows.Scan(&notGoingUser)
		if err != nil {
			log.Println("get usetId by events error: ", err)
			return event, err
		}
		user, err := GetUserById(notGoingUser)
		if err != nil {
			return event, nil
		}
		notGoingUsers = append(notGoingUsers, user)
	}
	event.NotGoingUsers = notGoingUsers
	return event, nil
}

func RemoveUserFromGroup(groupId int, userId int) error {
	_, err := database.Exec("DELETE FROM user_groups WHERE groupId = ? AND userId = ?", groupId, userId)
	if err != nil {
		return err
	}
	return nil
}
