package database

import "fmt"

// Retrieve the conversations
func (db *appdbimpl) GetConversation(conversation_id int) ([]Message, error) {
	// Query data
	var sql = `SELECT MESSAGES.msg_id, ifnull(forwarded_from_msg_id, 0), conversation_id, MESSAGES.from_user_id,
				strftime ("%H:%M", sent_timestamp), seen, msg, iif(substr(msg, 1, 11)='data:image/', true, false) is_photo,
				ifnull(R1.reaction_user, "") reaction
				FROM MESSAGES
				LEFT JOIN
					(
						SELECT msg_id, from_user_id, group_concat(reaction || USERS.user_name) as reaction_user
						FROM REACTIONS
						INNER JOIN USERS ON USERS.user_id = REACTIONS.from_user_id
					) AS R1
					ON R1.msg_id=MESSAGES.msg_id
				LEFT JOIN USERS ON USERS.user_id = R1.from_user_id
				WHERE conversation_id = ` + fmt.Sprint(conversation_id) + ` ORDER BY sent_timestamp DESC`

	rows, err := db.c.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Msg_id, &msg.Forwarded_from_msg_id, &msg.Conversation_id, &msg.From_user_id,
			&msg.Sent_timestamp, &msg.Seen, &msg.Msg, &msg.Is_photo, &msg.Reaction)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	err = rows.Err()
	if err != nil {
		return messages, err
	}

	return messages, nil
}
