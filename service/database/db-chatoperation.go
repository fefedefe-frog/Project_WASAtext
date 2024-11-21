package database

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func (db *appdbimpl) GetUserChats(usrId string) ([]Chat, error) {
	var userChats []Chat

	//Recupero gli id delle chat dalla chat_participants_table
	rows, err := db.c.Query(`SELECT chatId FROM chat_participants_table WHERE usrId = ?`, usrId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
			}
		}
	}()

	//Inserisco i vari id trovati in un array
	var userChatsId []int
	for rows.Next() {
		var chatId int

		if err := rows.Scan(&chatId); err != nil {
			return nil, err
		}
		userChatsId = append(userChatsId, chatId)
	}

	//Controllo se ci sono stati errori sulle righe
	if err := rows.Err(); err != nil {
		return userChats, err
	}

	query := "SELECT chatName, chatType, chatPhoto FROM chats_table WHERE chatId IN (" + strings.Repeat("?", len(userChatsId)-1) + "?)"

	chatRows, err := db.c.Query(query, toInterfaceSlice(userChatsId)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := chatRows.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("chatRows.Close() error: %v", closeErr)
			}
		}
	}()

	//Aggiungo le chat dell'utente allo slice userChats
	for chatRows.Next() {
		var chat Chat
		var chatPropicBytes []byte

		if err := chatRows.Scan(&chat.ChatId, &chat.ChatName, &chat.IsGroup, &chatPropicBytes); err != nil{
			return nil, err
		}
		chat.ChatPhoto = base64.StdEncoding.EncodeToString(chatPropicBytes)

		participants, participantsErr := db.GetChatPartecipants(chat.ChatId)

		if participantsErr != nil{
			return nil, participantsErr
		}

		chat.Participants = participants

		//Aggiungo la chat allo slice di chats
		userChats = append(userChats, chat)
	}

	if err := chatRows.Err(); err != nil {
		return nil, err
	}
	return userChats, nil
}

func (db *appdbimpl) InsertNewChat(participants []string, chatName string, chatPhoto string, isGroup bool) (Chat, error) {
	var chat Chat
	var groupPhotoBytes []byte

	//Se la chat è un gruppo controllo se sono stati dati il nome e la propic
	if isGroup{
		if  chatName == ""{	//Assegno un nome di default
			chat.ChatName= "Gruppo"
		}
		if chatPhoto == ""{ //Assegno una propic di default
			chat.ChatPhoto= defaultGroupPhotoBase64
		}
	}else{
		chat.ChatName= ""
		chat.ChatPhoto= ""
	}

	var errProp error
	// Decodifica la stringa Base64 in byte
	groupPhotoBytes, errProp = base64.StdEncoding.DecodeString(chat.ChatPhoto)
	if errProp != nil {
		return chat, errProp
	}

	//Conto quante chat sono presenti nel database per poi sommare 1 al valore ottenuto e assengnarlo come chatId della nuova chat
	var chatsCount int
	if err := db.c.QueryRow("SELECT COUNT(chatId) FROM chats_table").Scan(&chatsCount); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			chatsCount= 0
		}else{
			return chat, err
		}
	}
	chat.ChatId= chatsCount + 1

	tx, err := db.c.Begin()
	if  err != nil{
		return chat, err
	}
	defer func(err error) {
		if err != nil{
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}(err)

	//Eseguo l'inserimento nel database
	_, execErr := tx.Exec(`INSERT INTO chats_table (chatId, chatName, isGroup, chatPhoto) VALUES (?, ?, ?, ?)`, chat.ChatId, chat.ChatName, chat.IsGroup, groupPhotoBytes)
	if execErr != nil{
		return chat, execErr
	}

	//Ora devo creare le associazioni usrId <-> chatId nella chat_participants_table
	stmt, prepErr := tx.Prepare("INSERT INTO chat_participants_table (chatId, usrId) VALUES (?, ?)")
	if prepErr != nil{
		return chat, prepErr
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	for _, usrId := range participants {
		_, err := stmt.Exec(chat.ChatId, usrId)
		if err != nil {
			return chat, err // Interrompe l'inserimento se c'è un errore
		}
	}


	if txErr := tx.Commit();  txErr != nil{
		err= txErr
		return chat, err
	}
	return chat, nil
}

func (db *appdbimpl) DeleteChat(chatId int) error {
	//Inizializzo una transizione nel db, in quanto tutte queste operazioni al db sono considerate come una operazione atomica
	tx, err := db.c.Begin()
	if err != nil{
		return err
	}
	defer func(err error) {
		if err != nil{
			if rbErr := tx.Rollback(); rbErr != nil {
				// Se il rollback fallisce, logghiamo l'errore di rollback
				logrus.WithError(rbErr).Error("Errore durante il rollback")
			}
		}
	}(err)

	//Rimuovo ogni associazione usrId <-> chatId da chat_participants_table
	_, delPartErr := tx.Exec("DELETE FROM chat_participants_table WHERE chatId = ?", chatId)
	if delPartErr != nil{
		return delPartErr
	}

	//TODO rimuove ogni messaggio associato a quella chat da chat_messages_table


	//Rimuovo le info della chat da chats_table

	if _, err := tx.Exec("DELETE FROM chats_table WHERE chatId = ?", chatId); err != nil{
		return err
	}

	if err := tx.Commit();  err != nil{
		return err
	}
	return nil
}

func (db *appdbimpl) GetChatInfo(chatId int) (Chat, error) {
	var chat Chat
	var groupPropicByte []byte

	err := db.c.QueryRow(`SELECT isGroup, chatName, chatPhoto FROM chats_table WHERE chatId=?`, chatId).Scan(&chat.IsGroup, &chat.ChatName, &groupPropicByte)
	if err != nil{
		return chat, err
	}
	chat.ChatId= chatId

	//Controllo se sia presente la foto nel database
	if len(groupPropicByte) > 0{
		chat.ChatPhoto = base64.StdEncoding.EncodeToString(groupPropicByte)
	} else {
		chat.ChatPhoto = ""  //se non è presente assegno la stringa vuota
	}

	return chat, err
}

/*
	GetChatMessages Message, from a chatId received in input the function search in the messages table all the messages
	with that chatId then save it in a variable and procede to elaborate the data for returning a slice of messages
*/
func (db *appdbimpl) GetChatMessages(chatId int) ([]Message, error) {
	var messages []Message

	// Cerco tutte le righe che contengono il chatId corrispondente a quello interessato
	rows, err := db.c.Query(`SELECT msgId, senderId, contentType, content, timestamp FROM chat_messages_table WHERE chatId = ?`, chatId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
			}
		}
	}()

	//Itero su tutte le righe ottenute precedentemente
	for rows.Next() {
		var message Message

		var contentRaw []byte
		err := rows.Scan(&message.MsgId, &message.SenderId, &message.ContentType, &contentRaw, &message.Timestamp)
		if err != nil {
			return nil, err
		}

		//controllo se il contenuto è una foto
		switch message.ContentType{
		case "photo":
			message.Content = base64.StdEncoding.EncodeToString(contentRaw)
		case "text":
			message.Content = string(contentRaw)
		default:
			message.Content = string(contentRaw)
		}

		//Aggiungo l'utente all'array
		messages = append(messages, message)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return messages, err
}

func (db *appdbimpl) RemoveUserFromChat(usrId string, chatId int) error {
	//TODO usare db.c.Begin()
	_, err := db.c.Exec("DELETE FROM chat_participants_table WHERE chatId = ? AND usrId = ?", chatId, usrId)
	if err != nil{
		return err
	}

	return nil

	//TODO rimuove tutti i messaggi inviati dall'utente lui inviati
}

// GetChatPartecipants The function retrive all the partecipants ID of a chat by its chatId gived in input
func (db *appdbimpl) GetChatPartecipants(chatId int) ([]string, error) {

	//Recupero gli id delle chat dalla chat_participants_table
	rows, err := db.c.Query(`SELECT usrId FROM chat_participants_table WHERE chatId = ?`, chatId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
			}
		}
	}()

	//Inserisco i vari id trovati in un array
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

func (db *appdbimpl) SetGroupName(chatId int, newName string) error {
	stmt, err := db.c.Prepare(`UPDATE chats_table SET chatName = ? WHERE chatId=?`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(newName, chatId)
	if err != nil{
		return err
	}
	return err
}

func (db *appdbimpl) SetGroupPhoto(chatId int, newPhoto string) error {
	// Verifica che la stringa sia in formato base64 valido
	_, err := base64.StdEncoding.DecodeString(newPhoto)
	if err != nil{
		return err
	}

	//Semplice controllo della stringa base64 per assicurarsi
	//che la stringa contenga solo caratteri usati dalla codifica base64
	re := regexp.MustCompile(`^([A-Za-z0-9+/=]+)$`)
	if !re.MatchString(newPhoto) {
		return errors.New("la stringa base64 non rappresenta un'immagine valida")
	}

	// Decodifica la stringa base64
	data, errPropic := base64.StdEncoding.DecodeString(newPhoto)
	if errPropic != nil {
		return errPropic
	}


	stmt, err := db.c.Prepare(`UPDATE chats_table SET chatPhoto = ? WHERE chatId=?`)
	if err != nil {
		return err //log.Fatal("errore nella preparazione della query:", err)	//usare logger giusto, ritorno errore che viene gestito quando la funzione è chiamata
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil{
			if err == nil{
				err = closeErr
			}else{
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(data, chatId)
	if err != nil {
		return err //log.Fatal("errore nell'esecuzione della  query:", err)
	}
	return err
}

func toInterfaceSlice(ids []int) []interface{} {
	out := make([]interface{}, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}
