package database

// Retrieve the conversations
func (db *appdbimpl) GetMyConversations(from_user_id int) ([]Conversation, error) {
	// Query data
	rows, err := db.c.Query(`
						SELECT conversation_id, user_name as user_or_group_name, ifnull(user_photo, '') as user_or_group_photo, 0 is_group
						FROM CONVERSATIONS_USERS
						INNER JOIN USERS ON CONVERSATIONS_USERS.user_id=USERS.user_id
						WHERE CONVERSATIONS_USERS.conversation_id IN (
							SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS
							LEFT JOIN GROUPS ON CONVERSATIONS_USERS.conversation_id=GROUPS.conversation_id
							WHERE user_id = ? AND GROUPS.conversation_id IS NULL
						) AND CONVERSATIONS_USERS.user_id <> ?
						UNION ALL
						SELECT GROUPS.conversation_id, group_name as user_or_group_name, ifnull(group_photo, '') as user_or_group_photo, 1 is_group
						FROM GROUPS
						WHERE GROUPS.conversation_id IN (
							SELECT CONVERSATIONS_USERS.conversation_id FROM CONVERSATIONS_USERS WHERE user_id = ?
						);`, from_user_id, from_user_id, from_user_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.Conversation_id, &conv.User_or_group_name, &conv.User_or_group_photo, &conv.Is_group)
		if err != nil {
			return conversations, err
		}
		conversations = append(conversations, conv)
	}
	err = rows.Err()
	if err != nil {
		return conversations, err
	}

	return conversations, nil
}
