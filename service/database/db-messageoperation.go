package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) InsertMessage(message Message, chatId int) error {

	//Conto quanti messaggi sono presenti nel db per poi sommare 1 al valore ottenuto e assegnarlo come msgId del nuovo messaggio
	var messageCounts int
	if err := db.c.QueryRow("SELECT COUNT(msgId) FROM chat_messages_table").Scan(&messageCounts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			messageCounts = 0
		} else {
			return err
		}
	}
	message.MsgId = messageCounts + 1

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

	query := `INSERT INTO chat_messages_table (msgId, senderId, chatId, contentType, content, deliveryStatus, timestamp, isForwarded) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	//Controllo il tipo di contenuto che ha il messaggio
	var messageContent interface{}
	if message.ContentType == "photo" {
		messageContent, err = base64.StdEncoding.DecodeString(message.Content)
		if err != nil {
			return err
		}
	} else {
		messageContent = message.Content
	}

	if _, err := tx.Exec(query, message.MsgId, message.SenderId, chatId, message.ContentType, messageContent, message.DeliveryStatus, message.Timestamp, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveMessage(msgId int, chatId int) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) ForwardMessage(forwarderId string, msgId int, chatIdToForwatd int) error {

	//Recupero il messaggio da inoltrare
	var message Message
	var contentBytes []byte
	query := `SELECT contentType, content FROM chat_messages_table WHERE msgId=?`
	err := db.c.QueryRow(query, msgId).Scan(&message.ContentType, contentBytes)
	if err != nil {
		return err
	}

	//Converto il contenuto del messaggio
	if message.ContentType == "photo" {
		message.Content = base64.StdEncoding.EncodeToString(contentBytes)
	} else {
		message.Content = string(contentBytes)
	}

	//Imposto il nuovo senderId del messaggio, e aggiorno il valore di isForwarded a true
	message.SenderId = forwarderId
	message.IsForwarded = true

	if err := db.InsertMessage(message, chatIdToForwatd); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetMessageComments(msgId int) ([]Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) CommentMessage(msgId int, comment Comment) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) UncommentMessage(msgId int, commentId int) error {
	//TODO implement me
	panic("implement me")
}
