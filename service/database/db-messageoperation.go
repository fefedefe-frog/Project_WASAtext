package database

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) GetChatMessages(chatId int, usrId string, msgId int) ([]Message, error) {
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

	// Procedo ad aggiornare lo stato di tutti i messaggi che vengono ricevuti a partire dall'ultimo messaggio già ricevuto
	query :=
		`	UPDATE message_status_table
			SET status = 'received'
			WHERE msgId IN (
				SELECT ms.msgId
				FROM message_status_table ms
						 JOIN messages_table m ON ms.msgId = m.msgId
				WHERE ms.receiverId = ?
				  AND ms.status = 'not_received'
				  AND m.timestamp > COALESCE((
					SELECT MIN(m2.timestamp)
					FROM message_status_table ms2
							 JOIN messages_table m2 ON ms2.msgId = m2.msgId
					WHERE ms2.receiverId = ?
					  AND m2.chatId = ?
					  AND ms2.status = 'read'
				), '0001-01-01T23:59:59'))
			AND receiverId= ?;`

	_, err = tx.Exec(query, usrId, usrId, chatId, usrId)
	if err != nil {
		return messages, ErrUpdateMessageStatus
	}

	// Cerco tutte le righe che contengono il chatId corrispondente a quello interessato
	var rows *sql.Rows
	rows, err = tx.Query(`SELECT msgId, senderId, respondTo, textContent, photoContent, deliveryStatus, timestamp, isForwarded FROM messages_table WHERE chatId = ? AND msgId > ? ORDER BY timestamp, msgId;`, chatId, msgId)
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
		var isForwardedInt int
		var respondTo sql.NullInt64
		err := rows.Scan(&message.MsgId, &message.SenderId, &respondTo, &message.TextContent, &message.PhotoContent, &message.DeliveryStatus, &message.Timestamp, &isForwardedInt)
		if err != nil {
			return nil, err
		}

		if respondTo.Valid {
			message.RespondTo = int(respondTo.Int64)
		} else {
			message.RespondTo = -1
		}

		message.IsForwarded = isForwardedInt != 0

		// Aggiungo il messaggio all'array
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

func (db *appdbimpl) GetChatLastMessage(chatId int, usrId string) (Message, error) {
	var message Message
	tx, err := db.c.Begin()
	if err != nil {
		return Message{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	// Procedo ad aggiornare lo stato di tutti i messaggi che vengono ricevuti a partire dall'ultimo messaggio già ricevuto
	query :=
		`	UPDATE message_status_table
			SET status = 'received'
			WHERE msgId IN (
				SELECT ms.msgId
				FROM message_status_table ms
						 JOIN messages_table m ON ms.msgId = m.msgId
				WHERE ms.receiverId = ?
				  AND ms.status = 'not_received'
				  AND m.timestamp > COALESCE((
					SELECT MIN(m2.timestamp)
					FROM message_status_table ms2
							 JOIN messages_table m2 ON ms2.msgId = m2.msgId
					WHERE ms2.receiverId = ?
					  AND m2.chatId = ?
					  AND ms2.status = 'read'
				), '0001-01-01T23:59:59'))
			AND receiverId= ?;`

	_, err = tx.Exec(query, usrId, usrId, chatId, usrId)
	if err != nil {
		return Message{}, err
	}

	var isForwardedInt int
	var respondTo sql.NullInt64
	query = `SELECT msgId, senderId, respondTo, textContent, photoContent, deliveryStatus, timestamp, isForwarded FROM messages_table WHERE chatId = ? ORDER BY timestamp DESC, msgId LIMIT 1;`
	err = tx.QueryRow(query, chatId).Scan(&message.MsgId, &message.SenderId, &respondTo, &message.TextContent, &message.PhotoContent, &message.DeliveryStatus, &message.Timestamp, &isForwardedInt)
	if err != nil {
		return Message{}, err
	}

	if respondTo.Valid {
		message.RespondTo = int(respondTo.Int64)
	} else {
		message.RespondTo = -1
	}

	message.IsForwarded = isForwardedInt != 0

	if err := tx.Commit(); err != nil {
		return Message{}, err
	}
	return message, nil
}

func (db *appdbimpl) InsertMessage(message Message, chatId int) (int, string, error) {

	tx, err := db.c.Begin()
	if err != nil {
		return -1, "", err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}()

	var isForwarded = 0
	if message.IsForwarded {
		isForwarded = 1
	}

	var result sql.Result
	if message.RespondTo != -1 {
		query := `INSERT INTO messages_table (senderId, respondTo, chatId, textContent, photoContent, deliveryStatus, isForwarded) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING msgId;`
		result, err = tx.Exec(query, message.SenderId, message.RespondTo, chatId, message.TextContent, message.PhotoContent, message.DeliveryStatus, isForwarded)

	} else {
		query := `INSERT INTO messages_table (senderId, chatId, textContent, photoContent, deliveryStatus, isForwarded) VALUES (?, ?, ?, ?, ?, ?) RETURNING msgId;`
		result, err = tx.Exec(query, message.SenderId, chatId, message.TextContent, message.PhotoContent, message.DeliveryStatus, isForwarded)

	}
	if err != nil {
		return -1, "", err
	}

	// Recuper l'id del messaggio appena inserito
	var msgId int64
	msgId, err = result.LastInsertId()

	var timestamp string
	err = tx.QueryRow(`SELECT timestamp FROM messages_table WHERE msgId=?`, int(msgId)).Scan(&timestamp)
	if err != nil {
		return -1, "", err
	}

	// Inserisco l'associazione tra messaggio e partecipante della chat dandogli il valoredi "not_received" e "read" per l'utente che ha inviato il messaggio
	query := `INSERT INTO message_status_table (msgId, receiverId, status)
			SELECT ?, usrId, 'not_received'
			FROM chat_participants_table
			WHERE chatId=? AND usrId != ?;`
	_, err = tx.Exec(query, msgId, chatId, message.SenderId)
	if err != nil {
		return -1, "", err
	}

	if err := tx.Commit(); err != nil {
		return -1, "", err
	}
	return int(msgId), timestamp, nil
}

func (db *appdbimpl) RemoveMessage(msgId int, chatId int) error {

	_, err := db.c.Exec(`DELETE FROM messages_table WHERE msgId= ? AND chatId= ?;`, msgId, chatId)
	if err != nil {
		return err
	}

	// Controllo se la chat ha altri messaggi sennò la elimino
	_, err = db.c.Exec(`DELETE FROM chats_table WHERE chatId= ? AND (SELECT COUNT(*) FROM messages_table WHERE chatId= ?) = 0;`, chatId, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) ForwardMessage(forwarderId string, msgId int, chatIdToForwatd int) error {

	// Recupero il messaggio da inoltrare
	var message Message
	err := db.c.QueryRow(`SELECT textContent, photoContent FROM messages_table WHERE msgId=?;`, msgId).Scan(&message.TextContent, &message.PhotoContent)
	if err != nil {
		return err
	}

	// Imposto il nuovo senderId del messaggio, e aggiorno il valore di isForwarded a true
	message.SenderId = forwarderId
	message.DeliveryStatus = "sent"
	message.RespondTo = -1
	message.IsForwarded = true

	if _, _, err := db.InsertMessage(message, chatIdToForwatd); err != nil {
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
	query :=
		`	UPDATE message_status_table
			SET status = 'read'
			WHERE msgId IN (
				SELECT ms.msgId
				FROM message_status_table ms
						 JOIN messages_table m ON ms.msgId = m.msgId
				WHERE ms.receiverId = ?
				  AND ms.status != 'read'
				  AND m.timestamp > COALESCE((
					SELECT MIN(m2.timestamp)
					FROM message_status_table ms2
							 JOIN messages_table m2 ON ms2.msgId = m2.msgId
					WHERE ms2.receiverId = ?
					  AND m2.chatId = ?
					  AND ms2.status = 'read'
				), '0001-01-01T23:59:59')
				  AND ms.status != 'read'
			)
			AND receiverId= ?;`
	_, err = tx.Exec(query, usrId, usrId, chatId, usrId)
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
	var isForwarded int

	query := `SELECT senderId, respondTo, textContent, photoContent, deliveryStatus, timestamp, isForwarded FROM messages_table WHERE msgId = ?`
	err := db.c.QueryRow(query, msgId).Scan(&message.SenderId, &message.RespondTo, &message.TextContent, &message.PhotoContent, &message.DeliveryStatus, &message.Timestamp, &isForwarded)
	if err != nil {
		return message, err
	}

	// Converto il valore di isForwarded
	message.IsForwarded = isForwarded != 0
	return message, nil
}

func (db *appdbimpl) GetSenderIdByMsgId(msgId int) (string, error) {

	var senderId string
	err := db.c.QueryRow(`SELECT senderId FROM messages_table WHERE msgId= ?;`, msgId).Scan(&senderId)
	if err != nil {
		return "", err
	}
	return senderId, nil
}

func (db *appdbimpl) GetMessageComments(msgId int) ([]Comment, error) {

	var comments []Comment

	rows, err := db.c.Query(`SELECT commentId, commenterId, content FROM message_comments_table WHERE msgId=?;`, msgId)
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Error("rows.Close()")
			}
		}
	}()
	if err != nil {
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
