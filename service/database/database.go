package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	UsrId     string `json:"usrId"`
	UserName  string `json:"userName"`
	UserPhoto string `json:"userPhoto"`
}

type Chat struct {
	ChatId       int      `json:"chatId"`
	ChatName     string   `json:"chatName"`
	ChatPhoto    string   `json:"chatPhoto"`
	ChatType     string   `json:"chatType"`
	Participants []string `json:"participants"`
}

type Message struct {
	MsgId          int       `json:"msgId"`
	SenderId       string    `json:"senderId"`
	ContentType    string    `json:"contentType"`
	Content        string    `json:"content"`
	DeliveryStatus int       `json:"deliveryStatus"`
	Timestamp      string    `json:"timestamp"`
	Comments       []Comment `json:"comments"`
}

type Comment struct {
	CommenterId string `json:"commenterId"`
	Content     string `json:"content"`
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//User operations
	InsertNewUser(userName string) (User, error)      //Aggiunge un nuovo utente al database
	GetUsrIdByName(userName string) (string, error)   //Ottiene l'usrId dal nome dell'user
	SetUserName(usrId string, newName string) error   //Imposta un nuovo username di un user, dal suo usrId
	SetUserPhoto(usrId string, newPhoto string) error //Imposta una nuova foto di un user, dal suo usrId
	GetUserInfo(usrId string) (User, error)           //Ottiene le informazioni di un user, capire come gestire le info user
	UsrIdExist(usrId string) (bool, error)            //Controlla se un user esiste dal suo usrId

	GetUsers() ([]User, error) //Ottiene un array di tutti gli users, capire come gestire il ritorno  di tutti gli utenti

	//Chat operations
	GetUserChats(usrId string) ([]Chat, error)         //Ottiene le chat di un utente
	InsertNewChat(usersId []string) (Chat, error)      //Aggiunge una nuova chat al database
	DeleteChat(chatId int) error                       //Elimina una chat dal database
	GetChatInfo(chatId int) (Chat, error)              //Ottiene le info di una chat
	GetChatMessages(chatId int) ([]Message, error)     //Ottiene una lista di messaggi della chat
	RemoveUserFromChat(usrId string, chatId int) error //Rimuove l'associazione tra un utente ed una chat
	GetChatPartecipants(chatId int) ([]string, error)  //Ottiene una lista di utenti di una chat

	//Group operations
	SetGroupName(chatId int, newName string) error   //Imposta un nuovo nome al gruppo
	SetGroupPhoto(chatId int, newPhoto string) error //Imposta una nuova foto per il gruppo

	//Message operations
	InsertMessage(message Message, chatId int) error     //Aggiunge un messaggio e le sue relative informazioni
	RemoveMessage(msgId int, chatId int) error           //Rimuove un messaggio e le sue relative informazioni
	ForwardMessage(msgId int, chatIdToForwatd int) error //Inoltra un messaggio esistente da una chat ad un'altra
	GetMessageComments(msgId int) ([]Comment, error)     //Ottiene i commenti di un messaggio
	CommentMessage(msgId int, comment Comment) error     //Inserisce il commento di un messaggio
	UncommentMessage(msgId int, commentId int) error     //Rimuove il commento di un messaggio

	//DB operations
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
    			usrId TEXT PRIMARY KEY, 
    			userName TEXT, 
    			userPhoto BLOB
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
    			chatId INTEGER PRIMARY KEY,
    			chatName TEXT, 
    			chatType TEXT, 
    			chatPhoto BLOB
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	//Check for chat_partecipants_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat_partecipants_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chat_partecipants_table(
    			chatId INTEGER, 
    			usrId TEXT, 
    			FOREIGN KEY (chatId) REFERENCES chats_table(chatId), 
    			FOREIGN KEY (usrId) REFERENCES users_table(usrId), 
    			PRIMARY KEY (chatId, usrId)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	//Check for chat_messages_table
	tableName = ""
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat_messages_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chat_messages_table(
    			msgId INTEGER PRIMARY KEY, 
    			senderId TEXT, 
    			chatId INTEGER, 
    			contentType TEXT, 
    			content BLOB,
				deliveryStatus TEXT,
    			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, 
    			FOREIGN KEY (senderId) REFERENCES users_table(usrId), 
    			FOREIGN KEY (chatId) REFERENCES chats_table(chatId)
        );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
