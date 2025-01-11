/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
*/

package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	User_id    int    `json:"user_id"`
	User_name  string `json:"user_name"`
	User_photo string `json:"user_photo"`
}

type Group struct {
	Conversation_id int    `json:"conversation_id"`
	Group_name      string `json:"group_name"`
	Group_photo     string `json:"group_photo"`
}

type Message struct {
	Msg_id                int    `json:"msg_id"`
	Forwarded_from_msg_id int    `json:"forwarded_from_msg_id"`
	Conversation_id       int    `json:"conversation_id"`
	From_user_id          int    `json:"from_user_id"`
	Sent_timestamp        string `json:"sent_timestamp"`
	Seen                  int    `json:"seen"`
	Reaction              string `json:"reaction"`
	Msg                   string `json:"msg"`
}

type Conversation struct {
	Conversation_id     int    `json:"conversation_id"`
	User_or_group_name  string `json:"user_or_group_name"`
	User_or_group_photo string `json:"user_or_group_photo"`
	Is_group            int    `json:"is_group"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetIdFromName(string) (int, error)
	GetUserIdAndPhoto(user_name string) (int, string, error)
	GetAllUsers() ([]User, error)
	CreateUser(user_id int, user_name string) error
	SetMyUserName(id int, name string) error
	SetMyPhoto(id int, photo string) error
	GetNextUserId() (int, error)
	GetNextConversationId() (int, error)
	GetNextMessageId() (int, error)
	GetMessage(msg_id int) (Message, error)
	GetMessageText(msg_id int) (string, error)
	ChangeName(id int, new_name string) error
	GetMyConversations(from_user_id int) ([]Conversation, error)
	GetConversation(conversation_id int) ([]Message, error)
	IsUserInConversation(user_id int, conversation_id int) (bool, error)
	GetConversationId(user_id1 int, user_id2 int) (int, error)
	CreateConversation(conversation_id int, user_id1 int, user_id2 int) error
	SendMessage(conversation_id int, from_user_id int, message string, forwarded_from_msg_id int) error
	DeleteMessage(msg_id int) (sql.Result, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// Returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Check if table exists. If not, the database is empty, so we need to create the structure
	var tableName string
	var err error

	// Check and create table USERS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE USERS (
					user_id	INTEGER NOT NULL,
					user_name TEXT NOT NULL UNIQUE,
					user_photo TEXT,
					PRIMARY KEY(user_id));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table MESSAGES
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='MESSAGES';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE MESSAGES (
					msg_id INTEGER NOT NULL,
					conversation_id INTEGER NOT NULL,
					from_user_id INTEGER NOT NULL,
					sent_timestamp TEXT NOT NULL DEFAULT '1970-01-01 00:00:00',
					msg TEXT NOT NULL,
					forwarded_from_msg_id INTEGER,
					seen INTEGER NOT NULL DEFAULT 0,
					reaction TEXT,
					PRIMARY KEY(msg_id));`)

		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table CONVERSATIONS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='CONVERSATIONS_USERS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE CONVERSATIONS_USERS (
					conversation_id INTEGER NOT NULL,
					user_id INTEGER NOT NULL,
					PRIMARY KEY("conversation_id","user_id"));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	// Check and create table GROUPS
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='GROUPS';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`CREATE TABLE GROUPS (
					conversation_id INTEGER NOT NULL,
					group_name INTEGER NOT NULL,
					group_photo INTEGER,
					PRIMARY KEY("conversation_id"));`)
		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
