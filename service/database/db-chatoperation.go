package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func (db *appdbimpl) GetChatsOfUser(usrId string) ([]Chat, error) {

	/*
		In questa query eseguo il primo join per ottenere tutte le info delle chat di cui fa parte l'utente
		tramite il join nella tabella chats_table e quella dei partecipanti, associando il chatId
		successivamente faccio un join della tabella dei partecipanti su se stessa per trovare tutti i
		partecipanti di una chat, di cui recupero tutte le informazioni
	*/
	query := `
			SELECT 
			    C.chatId,
			    C.chatName,
			    C.chatPhoto,
			    C.isGroup,
			    U.usrId,
			    U.userName,
			    U.userPhoto	
			FROM chat_participants_table AS P
			JOIN chats_table C ON P.chatId = C.chatId
			JOIN chat_participants_table P2 ON C.chatId = P2.chatId
			JOIN users_table U ON P2.userId = U.userId
			WHERE P.usrId = ?
			GROUP BY C.chatId`

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

	// Inserisco le informazioni ricavate dalla lettura delle varie righe in una mappa
	chatMap := make(map[int]*Chat)
	for rows.Next() {
		var participant User
		var chatId int
		var chatName string
		var chatPhoto []byte
		var isGroupInt int

		err = rows.Scan(&chatId, &chatName, &chatPhoto, &isGroupInt, &participant.UsrId, &participant.UserName, &participant.UserPhoto)
		if err != nil {
			return nil, err
		}

		// Controllo se ho già recuperato le informazioni della chat da altre righe,
		// se non è così allora aggiungo la chat alla mappa passandogli tutte le informazioni
		if _, exists := chatMap[chatId]; !exists {
			chatMap[chatId] = &Chat{
				ChatId:       chatId,
				ChatName:     chatName,
				ChatPhoto:    chatPhoto,
				IsGroup:      isGroupInt != 0,
				Participants: []User{},
			}
		}

		// Aggiungo le informazioni del partecipante alla chat
		chatMap[chatId].Participants = append(chatMap[chatId].Participants, participant)
	}

	// Converto la mappa in array di Chat
	var userChats []Chat
	for _, chat := range chatMap {
		userChats = append(userChats, *chat)
	}

	// Controllo se ci sono stati errori sulle righe
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userChats, nil
}

func (db *appdbimpl) InsertNewChat(participants []string, chat Chat, messageTextContent string, messagePhotoContent []byte) (int, error) {

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
	if chat.IsGroup {
		isGroupVal = 1
	}

	// Eseguo l'inserimento nel database
	var result sql.Result
	result, err = tx.Exec(`INSERT INTO chats_table (chatName, isGroup, chatPhoto) VALUES (?, ?, ?);`, chat.ChatName, isGroupVal, chat.ChatPhoto)
	if err != nil {
		return -1, err
	}

	var newChatId64 int64
	newChatId64, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	// Inserisco gli usrId nella lista da passare alla funzione Exec
	usrIds := make([]interface{}, len(participants))
	for i, participant := range participants {
		usrIds[i] = participant
	}

	// Ora devo creare le associazioni usrId <-> chatId nella chat_participants_table
	query := `INSERT OR IGNORE INTO chat_participants_table (chatId, usrId) SELECT ?, usrId FROM users_table WHERE usrId IN (` + strings.Repeat("?,", len(usrIds)-1) + `?);`
	if _, err := tx.Exec(query, append([]interface{}{int(newChatId64)}, usrIds...)...); err != nil {
		return -1, err
	}

	queryMessage := `INSERT INTO chat_messages_table (senderId, chatId, textContent, photoContent, deliveryStatus, isForwarded) VALUES (?, ?, ?, ?, ?, ?) RETURNING msgId;`
	result, err = tx.Exec(queryMessage, participants[len(participants)-1], int(newChatId64), messageTextContent, messagePhotoContent, "sent", 0)
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
	_, err = tx.Exec(query, msgId, int(newChatId64), participants[len(participants)-1])
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

	query := `
			SELECT 
			    C.chatName,
			    C.chatPhoto,
			    C.isGroup,
			    U.usrId,
			    U.userName,
			    U.userPhoto	
			FROM chats_table AS C
			JOIN chat_participants_table P ON C.chatId = P.chatId
			JOIN users_table U ON P.userId = U.userId
			WHERE C.chatId = ?`

	// Eseguo la query descritta primia
	rows, err := db.c.Query(query, chatId)
	if err != nil {
		return Chat{}, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logrus.WithError(closeErr).Error("rows.Close()")
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Inserisco le informazioni ricavate dalla lettura delle varie righe in una mappa
	var chat Chat
	for rows.Next() {
		var participant User
		var isGroupInt int

		err = rows.Scan(&chat.ChatName, &chat.ChatPhoto, &isGroupInt, &participant.UsrId, &participant.UserName, &participant.UserPhoto)
		if err != nil {
			return Chat{}, err
		}

		chat.ChatId = chatId
		chat.IsGroup = isGroupInt != 0

		// Aggiungo le informazioni del partecipante alla chat
		chat.Participants = append(chat.Participants, participant)
	}
	if err := rows.Err(); err != nil {
		return Chat{}, err
	}

	return chat, nil
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

	// Rimuovo tutti i messaggi mandati dall'utente in quella chat se la chat non è già stata eliminata
	_, err = tx.Exec(`DELETE FROM chat_messages_table WHERE senderId= ? AND chatId= ? AND EXISTS (SELECT 1 FROM chats_table WHERE chatId= ?);`, usrId, chatId, chatId)
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

func (db *appdbimpl) FindChatFromParticipants(participants []string, isGroup bool) (int, error) {

	isGroupInt := 0
	if isGroup {
		isGroupInt = 1
	}
	participantNum := len(participants)
	placeholders := strings.TrimSuffix(strings.Repeat("?,", participantNum), ",")

	// Preparo gli args variabili contententi solo usrId dei partecipanti
	args := make([]interface{}, 0, participantNum+3)
	args = append(args, isGroupInt, participantNum)
	for i, participant := range participants {
		args[i] = participant
	}
	args = append(args, participantNum)

	query := fmt.Sprintf(`
		SELECT PD.chatId
		FROM chat_participants_table AS PD
		JOIN chat_table C ON PD.chatId = C.chatId WHERE C.isGroup = ?
		GROUP BY PD.chatId
		HAVING COUNT(DISTINCT usrId) = ?
			AND SUM(CASE WHEN usrId IN (%s) THEN 1 ELSE 0 END) = ?`, placeholders)

	var chatId int
	err := db.c.QueryRow(query, args...).Scan(&chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return -1, err
	}
	return chatId, nil
}
