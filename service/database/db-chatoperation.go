package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func (db *appdbimpl) GetUserChats(usrId string) ([]Chat, error) {

	/*
		In questa query eseguo il primo join per ottenere tutte le info delle chat di cui fa parte l'utente
		tramite il join nella tabella chats_table e quella dei partecipanti, associando il chatId
		successivamente faccio un join della tabella dei partecipanti su se stessa per trovare gli id dei
		partecipanti di una chat, che poi aggiungerò all'output finale tramite la funzione GROUP_CONCAT
		che restituisce una colonna di id concatenate dal carattere speciale ␟ (in questo caso)
	*/
	query := `
			SELECT C.chatId, C.chatName, C.chatPhoto, C.isGroup, GROUP_CONCAT(P2.usrId, '␟') AS participantsString
			FROM chat_participants_table AS P
			JOIN chats_table C ON P.chatId = C.chatId
			JOIN chat_participants_table P2 ON C.chatId = P2.chatId
			WHERE P.usrId = ?
			GROUP BY C.chatId, C.chatName, C.chatPhoto;`

	// Eseguo la query descritta primia
	rows, err := db.c.Query(query, usrId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logrus.WithError(closeErr).Error("rows.Close()")
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Controllo se la query ha restituito righe, se no, l'utente non fa parte di alcuna chat
	if !rows.Next() {
		return nil, ErrUserNoChat
	}

	// Inserisco i vari id trovati nell'array di output
	var userChats []Chat
	for rows.Next() {
		var chat Chat
		var chatPropicBytes []byte
		var participantsString string

		if err := rows.Scan(&chat.ChatId, &chat.ChatName, &chatPropicBytes, &chat.IsGroup, &participantsString); err != nil {
			return nil, err
		}

		// Splitto la stringa contenente i partecipanti
		chat.Participants = strings.Split(participantsString, "␟")

		/*
			Controllo se la chat è un gruppo o meno, se la chat è un gruppo
			salvo i valori di chatName e chatPhoto nella struct chat,
			sennò recupero le informazioni dei singoli utenti
		*/
		if chat.IsGroup {
			chat.ChatPhoto = base64.StdEncoding.EncodeToString(chatPropicBytes)
		} else {
			if len(chat.Participants) == 2 { // Ho una chat diretta tra due persone
				var otherParticipantId string
				for _, participant := range chat.Participants {
					if participant != usrId {
						otherParticipantId = participant
					}
				}

				user, err := db.GetUserInfo(otherParticipantId)
				if err != nil {
					return nil, err
				}
				chat.ChatName = user.UserName
				chat.ChatPhoto = user.UserPhoto
			} else { // C'è stato qualche problema, non possono esistere chat, che non sono gruppi con meno di 2 utenti
				return nil, fmt.Errorf("invalid participant number")
			}
		}

		// Aggiungo la chat allo slice di chats
		userChats = append(userChats, chat)
	}

	// Controllo se ci sono stati errori sulle righe
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userChats, nil
}

func (db *appdbimpl) InsertNewChat(creatorUsrId string, chatName string, chatPhotoData []byte, participants []string, isGroup bool, messageContent interface{}) (int, error) {

	tx, err := db.c.Begin()
	if err != nil {
		return -1, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	// Converto il valore bool in int
	isGroupVal := 0
	if isGroup {
		isGroupVal = 1
	}

	// Eseguo l'inserimento nel database
	var result sql.Result
	result, err = tx.Exec(`INSERT INTO chats_table (chatName, isGroup, chatPhoto) VALUES (?, ?, ?);`, chatName, isGroupVal, chatPhotoData)
	if err != nil {
		return -1, err
	}

	var newChatId64 int64
	newChatId64, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// Inserisco gli usrId nella lista da passare alla funzione Exec
	usrIds := make([]interface{}, len(participants)+1)
	usrIds[0] = creatorUsrId
	for i, participant := range participants {
		usrIds[i] = participant
	}

	// Ora devo creare le associazioni usrId <-> chatId nella chat_participants_table
	query := `INSERT OR IGNORE INTO chat_participants_table (chatId, usrId) SELECT ?, usrId FROM users_table WHERE usrId IN (` + strings.Repeat("?,", len(usrIds)-1) + `?);`
	if _, err := tx.Exec(query, append([]interface{}{int(newChatId64)}, usrIds...)...); err != nil {
		return -1, err
	}

	// Inserisco il messaggio iniziale nella tabella dei messaggi
	var contentType string
	switch value := messageContent.(type) {
	case []byte:
		contentType = "photo"
	case string:
		contentType = "text"
	default:
		return -1, fmt.Errorf("unsupported message content type: %T", value)
	}

	queryMessage := `INSERT INTO chat_messages_table (senderId, chatId, contentType, content, deliveryStatus, isForwarded) VALUES (?, ?, ?, ?, ?, ?) RETURNING msgId;`
	result, err = tx.Exec(queryMessage, creatorUsrId, int(newChatId64), contentType, messageContent, "sent", 0)
	if err != nil {
		return -1, err
	}

	// Recuper l'id del messaggio appena inserito
	var msgId int64
	msgId, err = result.LastInsertId()

	// Inserisco l'associazione tra messaggio e partecipante della chat dandogli il valoredi "not_received" e "read" per l'utente che ha inviato il messaggio
	query = `INSERT INTO message_status_table (msgId, receiverId, status)
			SELECT ?, usrId, 'not_received'
			FROM chat_participants_table
			WHERE chatId=? AND usrId != ?;`
	_, err = tx.Exec(query, msgId, int(newChatId64), creatorUsrId)
	if err != nil {
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}
	return int(newChatId64), nil
}

func (db *appdbimpl) DeleteChat(chatId int) error {

	// Inizializzo una transizione nel db, in quanto tutte queste operazioni al db sono considerate come una operazione atomica
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	// Rimuovo le info della chat da chats_table
	if _, err := tx.Exec("DELETE FROM chats_table WHERE chatId = ?;", chatId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetChatInfo(chatId int) (Chat, error) {
	var chat Chat
	var groupPropicByte []byte

	err := db.c.QueryRow(`SELECT isGroup, chatName, chatPhoto FROM chats_table WHERE chatId=?;`, chatId).Scan(&chat.IsGroup, &chat.ChatName, &groupPropicByte)
	if err != nil {
		return chat, err
	}
	chat.ChatId = chatId

	// Controllo se sia presente la foto nel database
	if len(groupPropicByte) > 0 {
		chat.ChatPhoto = base64.StdEncoding.EncodeToString(groupPropicByte)
	} else {
		chat.ChatPhoto = "" // se non è presente assegno la stringa vuota
	}

	// Ottengo gli id dei partecipanti della chat
	participants, participantsErr := db.GetChatPartecipants(chat.ChatId)
	if participantsErr != nil {
		return chat, participantsErr
	}
	chat.Participants = participants

	return chat, err
}

func (db *appdbimpl) RemoveUserFromChat(usrId string, chatId int) error {

	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	// Elimino la relazione tra chatId <-> usrId
	_, err = tx.Exec(`DELETE FROM chat_participants_table WHERE chatId= ? AND usrId= ?;`, chatId, usrId)
	if err != nil {
		return err
	}

	// Rimuovo tutti i messaggi mandati dall'utente in quella chat
	_, err = tx.Exec(`DELETE FROM chat_messages_table WHERE senderId= ? AND chatId= ?;`, usrId, chatId)
	if err != nil {
		return err
	}

	// Applico le modifiche al db
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) InsertUserInChat(usrId string, chatId int) error {

	_, err := db.c.Exec(`INSERT INTO chat_participants_table(chatId, usrId) VALUES (?, ?);`, chatId, usrId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetChatPartecipants(chatId int) ([]string, error) {

	// Recupero gli id degli user dalla chat_participants_table
	rows, err := db.c.Query(`SELECT usrId FROM chat_participants_table WHERE chatId = ?;`, chatId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Error("rows.Close()")
			}
		}
	}()

	// Inserisco i vari id trovati in un array
	var partecipants []string
	for rows.Next() {
		var usrId string

		if err := rows.Scan(&usrId); err != nil {
			return nil, err
		}
		partecipants = append(partecipants, usrId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return partecipants, nil
}

func (db *appdbimpl) GetChatParticipantsInfo(chatId int) ([]User, error) {

	// Uso di un join su una sotto tabella di chat_participants_table per ottenere tutti i dati sui partecipanti
	query := `SELECT users.*
	FROM  users_table AS users
	JOIN (
		SELECT usrId
		FROM chat_participants_table
		WHERE chatId = ?
	) AS participants
	ON users.usrId = participants.usrId;`

	rows, err := db.c.Query(query, chatId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Error("rows.Close()")
			}
		}
	}()

	var participants []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UsrId, &user.UserName, &user.UserPhoto)
		if err != nil {
			return nil, err
		}
		participants = append(participants, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return participants, nil
}

func (db *appdbimpl) CheckIfUserIsParticipant(chatId int, usrId string) (bool, error) {

	// Faccio una query per controllare se esiste una riga che ha l'associazione usrId <-> chatId, controllando se restituisce l'errore di NoRow
	err := db.c.QueryRow(`SELECT 1 FROM chat_participants_table WHERE chatId = ? AND usrId = ?;`, chatId, usrId).Scan(new(int))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}

func (db *appdbimpl) SetGroupName(chatId int, newName string) error {
	stmt, err := db.c.Prepare(`UPDATE chats_table SET chatName = ? WHERE chatId=?;`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Error("stmt.Close()")
			}
		}
	}()

	_, err = stmt.Exec(newName, chatId)
	if err != nil {
		return err
	}
	return err
}

func (db *appdbimpl) SetGroupPhoto(chatId int, newPhotoData []byte) error {

	stmt, err := db.c.Prepare(`UPDATE chats_table SET chatPhoto = ? WHERE chatId=?;`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Error("stmt.Close()")
			}
		}
	}()

	_, err = stmt.Exec(newPhotoData, chatId)
	if err != nil {
		return err
	}
	return err
}

func (db *appdbimpl) IsAGroup(chatId int) (bool, error) {

	var isGroup int
	err := db.c.QueryRow(`SELECT isGroup FROM chats_table WHERE chatId=?;`, chatId).Scan(&isGroup)
	if err != nil {
		return false, err
	}
	return isGroup == 1, nil
}
