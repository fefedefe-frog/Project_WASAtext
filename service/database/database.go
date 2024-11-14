package database

/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//User operations
	GetUsrIdByName(userName string) (string, error)   	//Ottiene l'usrId dal nome dell'user
	SetUserName(usrId string, newName string) error   	//Imposta un nuovo username di un user, dal suo usrId
	SetUserPhoto(usrId string, newPhoto string) error	//Imposta una nuova foto di un user, dal suo usrId
	GetUserInfo(usrId string) (string, error)         	//Ottiene le informazioni di un user, capire come gestire le info user

	GetUsers() ([]string, error)              			//Ottiene un array di tutti gli users
	GetUserChats(usrId string) ([]int, error)			//Ottiene le chat di un user dal suo usrId


	//Chat operations



	//Group operations



	//Message operations


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
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users_table';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "" +
			"CREATE TABLE users_table (" +
			"usrId TEXT PRIMARY KEY, " +
			"user_name TEXT, " +
			"user_photo BLOB" +
			");"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Check for chats_table
	tableName = ""
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='chats_table';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "" +
			"CREATE TABLE chats_table (" +
			"chatId INTEGER PRIMARY KEY, " +
			"usrId TEXT, " +
			"chat_type TEXT, " +
			"chat_photo BLOB" +
			");"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	//Check for chat_partecipants_table
	tableName = ""
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='chat_partecipants_table';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "" +
			"CREATE TABLE chat_partecipants_table (" +
			"chatId INTEGER, " +
			"usrId TEXT," +
			"FOREIGN KEY (chat_id) REFERENCES chats_table(chat_id), " +
			"FOREIGN KEY (usrId) REFERENCES users_table(usrId), " +
			"PRIMARY KEY (chat_id, usrId)" +
			");"
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	//Check for chat_messages_table
	tableName = ""
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='chat_messages_table';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := "" +
			"CREATE TABLE chat_messages_table (" +
			"msgId INTEGER PRIMARY KEY, " +
			"usrId TEXT, " +
			"chatId INTEGER, " +
			"msg_type TEXT, " +
			"msg_content BLOB, " +
			"timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"FOREIGN KEY (usrId) REFERENCES users_table(usrId), " +
			"FOREIGN KEY (chat_id) REFERENCES chats_table(chatId)" +
			");"
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
