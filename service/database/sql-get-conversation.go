package database

import "fmt"

// Retrieve the conversations
func (db *appdbimpl) GetConversation(conversation_id int) ([]Message, error) {
	// Query data
	var sql = `SELECT msg_id, ifnull(forwarded_from_msg_id, 0), conversation_id, from_user_id,
				strftime ("%H:%M", sent_timestamp), seen, ifnull(reaction, ''), msg FROM MESSAGES
				WHERE conversation_id = ` + fmt.Sprint(conversation_id)

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
			&msg.Sent_timestamp, &msg.Seen, &msg.Reaction, &msg.Msg)
		if err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}
