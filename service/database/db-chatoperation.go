package database

import (
	"encoding/base64"
	"strings"
)

func (db *appdbimpl) GetUserChats(usrId string) ([]Chat, error) {
	var userChats []Chat

	//Recupero gli id delle chat dalla chat_partecipants_table
	rows, err := db.c.Query(`SELECT chatId FROM chat_partecipants_table WHERE usrId = ?`, usrId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	//Inserisco i vari id trovati in un array
	var userChatsId []int
	for rows.Next(){
		var chatId int


		if err := rows.Scan(&chatId); err != nil{
			return nil, err
		}
		userChatsId = append(userChatsId, chatId)
	}

	//Controllo se ci sono stati errori sulle righe
	if err := rows.Err();err != nil{
		return userChats, err
	}

	query := "SELECT chatId, chatName, chatType, chatPhoto FROM chats_table WHERE chatId IN (" + strings.Repeat("?", len(userChatsId)-1) +"?)"

	chatRows, err := db.c.Query(query, toInterfaceSlice(userChatsId)...)
	if err != nil{
		return nil, err
	}
	defer chatRows.Close()

	//Aggiungo le chat dell'utente allo slice userChats
	for chatRows.Next(){
		var chat Chat
		var chatPropicBytes []byte

		if err := chatRows.Scan(&chat.ChatId, &chat.ChatName, &chat.ChatType, &chatPropicBytes); err != nil{
			return nil, err
		}
		chat.ChatPhoto= base64.StdEncoding.EncodeToString(chatPropicBytes)

		//Aggiungo la chat allo slice di chats
		userChats = append(userChats, chat)
	}

	if err := chatRows.Err();err != nil{
		return nil, err
	}

	return userChats, nil
}


func (db *appdbimpl) InsertNewChat(cN string, cP string, cT string, ps []string) (string, error) {
	//TODO implement me
	panic("implement me")
}


func (db *appdbimpl) DeleteChat(chatId int) error {
	//TODO implement me rimuove ogni associazione tra chat e usrId in chat_partecipants_table, rimuove ogni messaggio associato a quella chat da messa, rimuove la chat
	panic("implement me")
}


func (db *appdbimpl) GetChatInfo(chatId int) (Chat, error) {
	//TODO implement me
	panic("implement me")
}


//From a chatId received in input the function search in the messages table all the messages with that chatId
//then save it in a variable and procede to elaborate the data for returning a slice of messages
func (db *appdbimpl) GetChatMessages(chatId int) ([]Message, error) {
	var messages []Message

	rows, err := db.c.Query(`SELECT msgId, senderId, contentType, content, timestamp FROM chat_messages_table WHERE chatId = ?`, chatId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	//Itero su tutte le righe della tabella degli user
	for rows.Next() {
		var message Message

		var contentRaw []byte
		err := rows.Scan(&message.MsgId, &message.SenderId, &message.ContentType, &contentRaw, &message.Timestamp)
		if err != nil{
			return nil, err
		}

		//controllo se il contenuto è una foto
		switch message.ContentType{
		case "photo":
			message.Content = base64.StdEncoding.EncodeToString(contentRaw)
		case "text":
			message.Content = string(contentRaw)
		default:
			message.Content= string(contentRaw)
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
	//TODO implement me
	panic("implement me")
}


//The function retrive all the partecipants ID of a chat by its chatId gived in input
func (db *appdbimpl) GetChatPartecipants(chatId int) ([]string, error){

	//Recupero gli id delle chat dalla chat_partecipants_table
	rows, err := db.c.Query(`SELECT usrId FROM chat_partecipants_table WHERE chatId = ?`, chatId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	//Inserisco i vari id trovati in un array
	var partecipants []string
	for rows.Next(){
		var usrId string


		if err := rows.Scan(&usrId); err != nil{
			return nil, err
		}
		partecipants = append(partecipants, usrId)
	}

	if err := rows.Err();err != nil{
		return nil, err
	}

	return partecipants, nil
}


func (db *appdbimpl) SetGroupName(chatId int, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (db *appdbimpl) SetGroupPhoto(chatId int, newPhoto string) error {
	//TODO implement me
	panic("implement me")
}





func toInterfaceSlice(ids []int) []interface{} {
	out := make([]interface{}, len(ids))
	for i, id := range ids{
		out[i] = id
	}
	return out
}