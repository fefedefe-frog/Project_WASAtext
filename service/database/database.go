package database

import (
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUpdateMessageStatus = errors.New("unable to update message status")
)

type User struct {
	UsrId     string                 `json:"usrId"`     // Unique user id
	UserName  string                 `json:"userName"`  // name of the user
	UserPhoto utilitytool.BytesPhoto `json:"userPhoto"` // Propic of the user
}

type Chat struct {
	ChatId       int                    `json:"chatId"`       // Unique chat id
	IsGroup      bool                   `json:"isGroup"`      // Indicate if the chat is a group or a 1to1
	ChatName     string                 `json:"chatName"`     // Group name, or other username
	ChatPhoto    utilitytool.BytesPhoto `json:"chatPhoto"`    // Group photo, or other user propic
	Participants []User                 `json:"participants"` // Array that contain all the user participating in the chat
}

type Message struct {
	MsgId          int                    `json:"msgId"`          // Unique message id
	SenderId       string                 `json:"senderId"`       // User id of the user that send the message
	RespondTo      int                    `json:"respondTo"`      // Msg Id of the message at which is responding
	TextContent    string                 `json:"contentType"`    // Define the type of the content if text OR photo
	PhotoContent   utilitytool.BytesPhoto `json:"content"`        // Content of the message which can be text or a photo as a string in base64
	DeliveryStatus string                 `json:"deliveryStatus"` // Indicate the status of the message
	Timestamp      string                 `json:"timestamp"`      // Date when the message is sent
	Comments       []Comment              `json:"comments"`       // Array of Comment that store the reaction of other users
	IsForwarded    bool                   `json:"isForwarded"`    // Bool value that say if the message is forwarded or not
}

type Comment struct {
	CommentId   int    `json:"commentId"`   // Comment id
	CommenterId string `json:"commenterId"` // User id of the commenter
	Content     string `json:"content"`     // Content that can be an emoji
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// User operations

	// InsertNewUser User, add a new user in the database
	InsertNewUser(userName string) (User, error)
	// GetUsrIdByName string, retrive the User.UsrId of the user by passing its username
	GetUsrIdByName(userName string) (string, error)
	// SetUserName error, update the name of a user, by passing its usrId and its newName
	SetUserName(usrId string, newName string) error
	// SetUserPhoto update the user propic, the function receive the user id and the new propic as a string in base64 format
	SetUserPhoto(usrId string, newPhoto []byte) error
	// GetUserInfo User, search in the database a user with the usrId passed in the fuction, and retrive it info
	GetUserInfo(usrId string) (User, error)
	// GetUserNameById string, retrive the name of an user by its id
	GetUserNameById(usrId string) (string, error)
	// GetUsers User, get an array of all the users present in the db
	GetUsers(usrIdToIgnore string) ([]User, error)
	// UsrIdExist bool, check if a user in present in the db by its user id
	UsrIdExist(usrId string) (bool, error)

	// Chat operations

	// GetChatsOfUser Chat, retrive the chats of a user by passing its usrId
	GetChatsOfUser(usrId string) ([]Chat, error)
	// InsertNewChat int, add a new chat in the db, the function receive an array of users that are in the chat
	InsertNewChat(participants []string, chat Chat, messageTextContent string, messagePhotoContent []byte) (int, error)
	// DeleteChat error, remove a chat from the db, also remove all the message of that chat from the db
	DeleteChat(chatId int) error
	// GetChatInfo Chat, retrive all the info of a chat from the db
	GetChatInfo(chatId int) (Chat, error)
	// RemoveUserFromChat error, remove the associations betwheen a user and a chat
	RemoveUserFromChat(usrId string, chatId int) error
	// InsertUserInChat error, make the relation usrId <-> chatId
	InsertUserInChat(usrId string, chatId int) error
	// GetChatPartecipants string, retrive all the user id of all the users in a chat
	GetChatPartecipants(chatId int) ([]string, error)
	// GetChatParticipantsInfo []User, retrive all info of the participants of a chat
	GetChatParticipantsInfo(chatId int) ([]User, error)
	// CheckIfUserIsParticipant bool, check if exist the relation between the chatId and usrId given in input
	CheckIfUserIsParticipant(chatId int, userId string) (bool, error)
	// SetGroupName error, update the name of a group chat
	SetGroupName(chatId int, newName string) error
	// SetGroupPhoto error, update the group photo
	SetGroupPhoto(chatId int, newPhotoData []byte) error
	// IsAGroup bool, check if the chat is a group chat or not
	IsAGroup(chatId int) (bool, error)
	// FindChatFromParticipants int, try to retrive the chatId of a chat from it's participant list
	FindChatFromParticipants(participants []string, isGroup bool) (int, error)

	// Message operations

	// GetChatMessages Message, retrive the messages starting from a specified msgId of a chat
	GetChatMessages(chatId int, usrId string, msgId int) ([]Message, error)
	// InsertMessage error, insert a message in the database
	InsertMessage(message Message, chatId int) (int, error)
	// RemoveMessage error, remove a message from the database
	RemoveMessage(msgId int, chatId int) error
	// ForwardMessage error, forward an existing message with the msgId gived in input to another chat
	ForwardMessage(forwarderId string, msgId int, chatIdToForwatd int) error
	// UpdateMessageDeliveryStatusToRead error, update the delivery status of a message (And all its previus) to read, for the specified user
	UpdateMessageDeliveryStatusToRead(msgId int, chatId int, usrId string) error
	// GetMessageComments []Comment, get all the comment of a message by its msgId
	GetMessageComments(msgId int) ([]Comment, error)
	// CommentMessage error, comment a message
	CommentMessage(msgId int, commenterId string, content string) error
	// UncommentMessage error, remove a comment from a message by the commenterId of the message
	UncommentMessage(commentId int) error
	// GetMessageById Message, retrive a message from the database by its msgId
	GetMessageById(msgId int) (Message, error)
	// GetSenderIdByMsgId string, retrive the user id of the sender of the message by its msgId
	GetSenderIdByMsgId(msgId int) (string, error)
	// CheckCommentAuthor bool, check if the user is the author of a comment
	CheckCommentAuthor(commentId int, usrId string) (bool, error)

	// Ping error, verify the connection of the db
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if users_table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users_table(
    			usrId TEXT PRIMARY KEY NOT NULL, 
    			userName TEXT NOT NULL, 
    			userPhoto BLOB NOT NULL DEFAULT '',
    			UNIQUE (userName)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for chats_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chats_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chats_table(
    			chatId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    			isGroup INTEGER NOT NULL CHECK (isGroup IN (0, 1)), 
    			chatName TEXT, 
    			chatPhoto BLOB
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for chat_participants_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat_participants_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chat_participants_table(
    			chatId INTEGER, 
    			usrId TEXT, 
    			FOREIGN KEY (chatId) REFERENCES chats_table(chatId) ON DELETE CASCADE, 
    			FOREIGN KEY (usrId) REFERENCES users_table(usrId) ON DELETE CASCADE, 
    			PRIMARY KEY (chatId, usrId),
    			CONSTRAINT unique_participant UNIQUE  (chatId, usrId)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for messages_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages_table(
    			msgId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
    			senderId TEXT NOT NULL,
    			respondTo INTEGER DEFAULT NULL, 
    			chatId INTEGER NOT NULL, 
    			textContent TEXT NOT NULL DEFAULT '', 
    			photoContent BLOB NOT NULL,
				deliveryStatus TEXT NOT NULL CHECK (deliveryStatus IN ('sent', 'received', 'read')),
    			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    			isForwarded INTEGER NOT NULL CHECK (isForwarded IN (0, 1)),
    			FOREIGN KEY (senderId) REFERENCES users_table(usrId) ON DELETE CASCADE, 
    			FOREIGN KEY (chatId) REFERENCES chats_table(chatId) On DELETE CASCADE,
    			FOREIGN KEY (respondTo) REFERENCES messages_table(msgId) ON DELETE SET NULL
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for message_comments_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='message_comments_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE message_comments_table(
    			commentId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    			msgId INTEGER, 
    			commenterId TEXT, 
    			content TEXT NOT NULL,
    			FOREIGN KEY (commenterId) REFERENCES users_table(usrId) ON DELETE CASCADE, 
    			FOREIGN KEY (msgId) REFERENCES messages_table(msgId) ON DELETE CASCADE,
    			CONSTRAINT unique_comment UNIQUE (msgId, commenterId)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for message_status_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='message_status_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE message_status_table(
    			msgId INTEGER,
    			receiverId TEXT,
    			status TEXT NOT NULL CHECK (status IN ('not_received', 'received', 'read')),
    			PRIMARY KEY (msgId, receiverId),
    			FOREIGN KEY (msgId) REFERENCES messages_table(msgId) ON DELETE CASCADE,
    			FOREIGN KEY (receiverId) REFERENCES users_table(usrId) ON DELETE CASCADE,
    			CONSTRAINT unique_status UNIQUE (msgId, receiverId)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Trigger for auto deleting a chat with only one member
	_, err = db.Exec(`CREATE TRIGGER IF NOT EXISTS delete_chats_with_one_member
						AFTER DELETE ON chat_participants_table
						BEGIN
							DELETE FROM chats_table
							WHERE chatId = OLD.chatId
							AND (SELECT COUNT(*) FROM chat_participants_table WHERE chatId = OLD.chatId) < 2;
						END;`)
	if err != nil {
		return nil, fmt.Errorf("error initializing db trigger: %w", err)
	}

	// Trigger for auto updating a message delivery status
	_, err = db.Exec(`CREATE TRIGGER IF NOT EXISTS update_message_status
						AFTER UPDATE OF status ON message_status_table
						BEGIN
							UPDATE messages_table
						    SET deliveryStatus = 'read'
						    WHERE msgId = NEW.msgId
						    AND (
						        SELECT COUNT(*)
						        FROM message_status_table
						        WHERE msgId = NEW.msgId AND status != 'read'
						    ) = 0;
							
							UPDATE messages_table
						    SET deliveryStatus = 'received'
						    WHERE msgId = NEW.msgId
						    AND (
						        SELECT COUNT(*)
						        FROM message_status_table
						        WHERE msgId = NEW.msgId AND status != 'received'
						    ) = 0;
						END;`)
	if err != nil {
		return nil, fmt.Errorf("error initializing db trigger: %w", err)
	}

	// Index initialization for faster db interaction used in some function
	_, err = db.Exec(`-- Indici per le Primary key
							CREATE INDEX IF NOT EXISTS idx_users_usrId ON users_table(usrId);
       						CREATE INDEX IF NOT EXISTS idx_chat_chatId ON chats_table(chatId);
       						CREATE INDEX IF NOT EXISTS idx_message_msgId ON messages_table(msgId);
       						-- Indici per le interazioni della participants_table
       						CREATE INDEX idx_participants_usrId ON chat_participants_table(usrId);
							CREATE INDEX idx_participants_chatId ON chat_participants_table(chatId);
							CREATE INDEX idx_participants_chatANDusr ON chat_participants_table(chatId, usrId);
							-- Indici per le interazioni della messages_table
							CREATE INDEX IF NOT EXISTS idx_chat_messages_chatId_timestamp ON messages_table(chatId, timestamp);
							-- Indici per le interazioni della status_table
							CREATE INDEX IF NOT EXISTS idx_message_status_msgId ON message_status_table(msgId);
							CREATE INDEX IF NOT EXISTS idx_message_status_receiverId ON message_status_table(receiverId);
							CREATE INDEX IF NOT EXISTS idx_message_status_status_receiverId_chatId ON message_status_table(receiverId, msgId);`)
	if err != nil {
		return nil, fmt.Errorf("error initializing db index: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
