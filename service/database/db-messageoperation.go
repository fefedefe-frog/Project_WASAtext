package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) GetChatMessages(chatId int, usrId string) ([]Message, error) {
	var messages []Message

	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	// Procedo ad aggiornare lo stato di tutti i messaggi ceh vengono ricevuti
	_, err = tx.Exec(`UPDATE message_status_table SET status= 'received' WHERE msgId IN (SELECT msgId FROM chat_messages_table WHERE chatId=?) AND receiverId=?`, chatId, usrId)
	if err != nil {
		return messages, ErrUpdateMessageStatus
	}

	// Cerco tutte le righe che contengono il chatId corrispondente a quello interessato
	var rows *sql.Rows
	rows, err = tx.Query(`SELECT msgId, senderId, contentType, content, deliveryStatus, timestamp FROM chat_messages_table WHERE chatId = ? ORDER BY timestamp DESC, msgId;`, chatId)
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

	// Itero su tutte le righe ottenute precedentemente
	for rows.Next() {
		var message Message

		var contentRaw []byte
		err := rows.Scan(&message.MsgId, &message.SenderId, &message.ContentType, &contentRaw, &message.DeliveryStatus, &message.Timestamp)
		if err != nil {
			return nil, err
		}

		// controllo se il contenuto è una foto
		if message.ContentType == "photo" {
			message.Content = base64.StdEncoding.EncodeToString(contentRaw)
		} else {
			message.Content = string(contentRaw)
		}

		// Aggiungo l'utente all'array
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return messages, err
}

// TODO modificare gestione stringa qui

func (db *appdbimpl) InsertMessage(message Message, chatId int) error {

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

	// Controllo il tipo di contenuto che ha il messaggio
	var messageContent interface{}
	if message.ContentType == "photo" {
		messageContent, err = base64.StdEncoding.DecodeString(message.Content)
		if err != nil {
			return err
		}
	} else {
		messageContent = message.Content
	}

	var isForwarded = 0
	if message.IsForwarded {
		isForwarded = 1
	}

	// Inserisco il messaggio nella tabella dei messaggi
	message.DeliveryStatus = "sent"
	query := `INSERT INTO chat_messages_table (senderId, chatId, contentType, content, deliveryStatus, isForwarded) VALUES (?, ?, ?, ?, ?, ?) RETURNING msgId;`
	var result sql.Result
	result, err = tx.Exec(query, message.SenderId, chatId, message.ContentType, messageContent, message.DeliveryStatus, isForwarded)
	if err != nil {
		return err
	}

	// Recuper l'id del messaggio appena inserito
	var msgId int64
	msgId, err = result.LastInsertId()

	// Inserisco l'associazione tra messaggio e partecipante della chat dandogli il valoredi "not_received" e "read" per l'utente che ha inviato il messaggio
	query = `INSERT INTO message_status_table (msgId, receiverId, status)
			SELECT ?, usrId, 'not_received'
			FROM chat_participants_table
			WHERE chatId=? AND usrId != ?;`
	_, err = tx.Exec(query, msgId, chatId, message.SenderId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(msgId int, chatId int) error {

	_, err := db.c.Exec(`DELETE FROM chat_messages_table WHERE msgId= ? AND chatId= ?;`, msgId, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) ForwardMessage(forwarderId string, msgId int, chatIdToForwatd int) error {

	// Recupero il messaggio da inoltrare
	var message Message
	var contentBytes []byte
	err := db.c.QueryRow(`SELECT contentType, content FROM chat_messages_table WHERE msgId=?;`, msgId).Scan(&message.ContentType, contentBytes)
	if err != nil {
		return err
	}

	// Converto il contenuto del messaggio
	if message.ContentType == "photo" {
		message.Content = base64.StdEncoding.EncodeToString(contentBytes)
	} else {
		message.Content = string(contentBytes)
	}

	// Imposto il nuovo senderId del messaggio, e aggiorno il valore di isForwarded a true
	message.SenderId = forwarderId
	message.IsForwarded = true

	if err := db.InsertMessage(message, chatIdToForwatd); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UpdateMessageDeliveryStatusToRead(msgId int, chatId int, usrId string) error {

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

	// Aggiorno il messaggio che interessa
	_, err = tx.Exec(`UPDATE message_status_table SET status= 'read' WHERE msgId=? AND receiverId=?;`, msgId, usrId)
	if err != nil {
		return err
	}

	// Aggiorno tutti i messaggi precedenti legati alla chat di quel messaggio
	// in una subquery recupero tutti i msgId legati alla chat passata come chatId,
	_, err = tx.Exec(`UPDATE message_status_table
							SET status = 'read'
							WHERE msgId IN (
								SELECT ms.msgId
								FROM message_status_table ms
									JOIN chat_messages_table cm ON ms.msgId = cm.msgId
								WHERE ms.receiverId = ?
								  AND cm.chatId = ?
								  AND cm.timestamp < (
									SELECT MIN(cm2.timestamp)
									FROM message_status_table ms2
											 JOIN chat_messages_table cm2 ON ms2.msgId = cm2.msgId
									WHERE ms2.receiverId = ?
									  AND cm2.chatId = ?
									  AND ms2.status = 'read'
								)
								  AND ms.status != 'read'
							)
							  AND receiverId = ?;`, usrId, chatId, usrId, chatId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetMessageById(msgId int) (Message, error) {
	var message Message
	var rawContent []byte
	var isForwarded int
	query := `SELECT senderId, contentType, content, deliveryStatus, timestamp, isForwarded FROM chat_messages_table WHERE msgId= ?;`
	err := db.c.QueryRow(query, msgId).Scan(&message.SenderId, &message.ContentType, &rawContent, &message.DeliveryStatus, &message.Timestamp, &message.Comments, &isForwarded)
	if err != nil {
		return message, err
	}

	// Controllo se il contenuto è una foto e la elaboro
	switch message.ContentType {
	case "photo":
		message.Content = base64.StdEncoding.EncodeToString(rawContent)
	case "text":
		message.Content = string(rawContent)
	default:
		message.Content = string(rawContent)
	}

	// Converto il valore di isForwarded
	message.IsForwarded = isForwarded != 0

	return message, nil
}

func (db *appdbimpl) GetSenderIdByMsgId(msgId int) (string, error) {

	var senderId string
	err := db.c.QueryRow(`SELECT senderId FROM chat_messages_table WHERE msgId= ?;`, msgId).Scan(&senderId)
	if err != nil {
		return "", err
	}
	return senderId, nil
}

func (db *appdbimpl) GetMessageComments(msgId int) ([]Comment, error) {

	var comments []Comment

	rows, err := db.c.Query(`SELECT commentId, commenterId, content FROM message_comments_table WHERE msgId=?;`, msgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageHaveNoComments
		}
		return nil, err
	}

	// Itero su tutti i commenti del messaggio
	for rows.Next() {
		var comment Comment

		err := rows.Scan(&comment.CommentId, &comment.CommenterId, &comment.Content)
		if err != nil {
			return nil, err
		}

		// Aggiungo l'utente all'array
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}

func (db *appdbimpl) CommentMessage(msgId int, commenterId string, content string) error {

	_, err := db.c.Exec(`INSERT INTO message_comments_table (msgId, commenterId, content) VALUES (?, ?, ?)`, msgId, commenterId, content)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UncommentMessage(commentId int) error {

	_, err := db.c.Exec(`DELETE FROM message_comments_table WHERE commentId= ?;`, commentId)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) CheckCommentAuthor(commentId int, usrId string) (bool, error) {

	// Faccio una query per controllare se esiste una riga che ha l'associazione usrId <-> chatId, controllando se restituisce l'errore di NoRow
	err := db.c.QueryRow(`SELECT 1 FROM message_comments_table WHERE commentId = ? AND commenterId = ?;`, commentId, usrId).Scan(new(int))
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}
