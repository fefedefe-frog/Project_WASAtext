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

	// Recupero gli id delle chat associate all'usrId dalla chat_participants_table
	rows, err := db.c.Query(`SELECT chatId FROM chat_participants_table WHERE usrId = ?;`, usrId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
			}
		}
	}()

	// Inserisco i vari id trovati in un array
	var userChatsId []int
	for rows.Next() {
		var chatId int

		if err := rows.Scan(&chatId); err != nil {
			return nil, err
		}
		userChatsId = append(userChatsId, chatId)
	}

	// Controllo se ci sono stati errori sulle righe
	if err := rows.Err(); err != nil {
		return userChats, err
	}

	if len(userChatsId) == 0 {
		return userChats, ErrUserNoChat
	}

	// Recupero tutte le informazioni delle chat passando la lista di chatId ottenute in precedenza
	query := "SELECT chatName, isGroup, chatPhoto FROM chats_table WHERE chatId IN (" + strings.Repeat("?", len(userChatsId)-1) + "?);"
	var chatRows *sql.Rows
	chatRows, err = db.c.Query(query, toInterfaceSlice(userChatsId)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := chatRows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("chatRows.Close() error: %v", closeErr)
			}
		}
	}()

	// Aggiungo le chat dell'utente allo slice userChats
	for chatRows.Next() {
		var chat Chat
		var chatPropicBytes []byte

		if err := chatRows.Scan(&chat.ChatId, &chat.ChatName, &chat.IsGroup, &chatPropicBytes); err != nil {
			return nil, err
		}

		// Ottengo gli id dei partecipanti della chat
		participants, participantsErr := db.GetChatPartecipants(chat.ChatId)
		if participantsErr != nil {
			return nil, participantsErr
		}
		chat.Participants = participants

		/*
			Controllo se la chat è un gruppo o meno, se la chat è un gruppo
			salvo i valori di chatName e chatPhoto nella struct chat,
			sennò recupero le informazioni dei singoli utenti
		*/
		if chat.IsGroup {
			chat.ChatPhoto = base64.StdEncoding.EncodeToString(chatPropicBytes)
		} else {
			var secondParticipantId string
			if len(participants) == 2 {
				for _, participant := range participants {
					if participant != usrId {
						secondParticipantId = participant
					}
				}
			} else {
				return nil, ErrChatParticipantNumber
			}
			user, err := db.GetUserInfo(secondParticipantId)
			if err != nil {
				return nil, err
			}
			chat.ChatName = user.UserName
			chat.ChatPhoto = user.UserPhoto
		}

		// Aggiungo la chat allo slice di chats
		userChats = append(userChats, chat)
	}

	if err := chatRows.Err(); err != nil {
		return nil, err
	}
	return userChats, nil
}

func (db *appdbimpl) InsertNewChat(participants []string, chatName string, chatPhoto string, isGroup bool) (int, error) {
	var groupPhotoBytes []byte

	// Se la chat è un gruppo controllo se sono stati dati il nome e la propic
	if isGroup {
		if chatName == "" { // Assegno un nome di default
			chatName = "Gruppo"
		}
		if chatPhoto == "" { // Assegno una propic di default
			chatPhoto = defaultGroupPhotoBase64
		}
	} else {
		chatName = ""
		chatPhoto = ""
	}

	var errProp error
	// Decodifica la stringa Base64 in byte
	groupPhotoBytes, errProp = base64.StdEncoding.DecodeString(chatPhoto)
	if errProp != nil {
		return -1, errProp
	}

	// Conto quante chat sono presenti nel database per poi sommare 1 al valore ottenuto e assengnarlo come chatId della nuova chat
	var chatsCount int
	if err := db.c.QueryRow("SELECT COUNT(chatId) FROM chats_table;").Scan(&chatsCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chatsCount = 0
		} else {
			return -1, err
		}
	}
	newChatId := chatsCount + 1

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
	_, err = tx.Exec(`INSERT INTO chats_table (chatId, chatName, isGroup, chatPhoto) VALUES (?, ?, ?, ?);`, newChatId, chatName, isGroupVal, groupPhotoBytes)
	if err != nil {
		return -1, err
	}

	// Ora devo creare le associazioni usrId <-> chatId nella chat_participants_table
	var stmt *sql.Stmt
	stmt, err = tx.Prepare("INSERT INTO chat_participants_table (chatId, usrId) VALUES (?, ?);")
	if err != nil {
		return -1, err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	for _, usrId := range participants {
		exist, err := db.UsrIdExist(usrId)
		if exist {
			if _, err := stmt.Exec(newChatId, usrId); err != nil {
				return -1, err // Interrompe l'inserimento se c'è un errore
			}
		} else {
			if err != nil {
				// se non riesco a controllare se l'utente esiste lo segnaolo, se la chat era tra due persone annullo la sua creazione
				logrus.WithError(err).WithField("usrId", usrId).Error("unable to add user to the group")
				if !isGroup {
					return -1, err
				}
			}
			// se non esiste nessun utente nel db, lo segnalo semplicemente
			logrus.WithField("usrId", usrId).Info("user does not exist")
		}
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}
	return newChatId, nil
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

	// Recupero gli id delle chat dalla chat_participants_table
	rows, err := db.c.Query(`SELECT usrId FROM chat_participants_table WHERE chatId = ?;`, chatId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
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
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(newName, chatId)
	if err != nil {
		return err
	}
	return err
}

func (db *appdbimpl) SetGroupPhoto(chatId int, newPhoto string) error {
	// Verifica che la stringa sia in formato base64 valido
	_, err := base64.StdEncoding.DecodeString(newPhoto)
	if err != nil {
		return err
	}

	// Semplice controllo della stringa base64 per assicurarsi
	// che la stringa contenga solo caratteri usati dalla codifica base64
	re := regexp.MustCompile(`^([A-Za-z0-9+/=]+)$`)
	if !re.MatchString(newPhoto) {
		return errors.New("la stringa base64 non rappresenta un'immagine valida")
	}

	// Decodifica la stringa base64
	data, errPropic := base64.StdEncoding.DecodeString(newPhoto)
	if errPropic != nil {
		return errPropic
	}

	var stmt *sql.Stmt
	stmt, err = db.c.Prepare(`UPDATE chats_table SET chatPhoto = ? WHERE chatId=?;`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(data, chatId)
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

func toInterfaceSlice(ids []int) []interface{} {
	out := make([]interface{}, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}
