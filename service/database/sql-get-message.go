package database

import (
	"fmt"
)

// Retrieve a message
func (db *appdbimpl) GetMessage(msg_id int) (Message, error) {
	var msg Message
	var err error
	//var row sql.row

	// Query data
	var sql = `SELECT msg_id, ifnull(forwarded_from_msg_id, 0), conversation_id, from_user_id,
				sent_timestamp, seen, ifnull(reaction, ''), msg FROM MESSAGES
				WHERE msg_id = ` + fmt.Sprint(msg_id)

	row := db.c.QueryRow(sql)
	if err != nil {
		return msg, err
	}

	err = row.Scan(&msg.Msg_id, &msg.Forwarded_from_msg_id, &msg.Conversation_id,
				&msg.From_user_id, &msg.Sent_timestamp, &msg.Seen, &msg.Reaction, &msg.Msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
