package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getChatInfo(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Recupero il valore di chat_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "Bad request - Invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente che ha effettuato l'accesso faccia parte della chat di cui vuole ricavare le informazioni
	if exist, err := rt.db.CheckIfUserIsParticipant(chatId, token); err != nil {
		context.Logger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error - Error while checking if user is a member of the chat", http.StatusInternalServerError)
		return
	} else if !exist {
		context.Logger.WithField("usrId", token).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user logged in is not a participant of the chat", http.StatusForbidden)
		return
	}

	var chat database.Chat
	chat, err = rt.db.GetChatInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("The chat doesn't exist, in the db")
			http.Error(writer, "Not found - chat not exist", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat info")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Controllo se la chat è un gruppo o meno, se non è un gruppo procedo
	// a recuperare le informazioni dell'altro utente partecipante alla chat
	if !chat.IsGroup {

		otherUsrId := chat.Participants[0]
		if otherUsrId == token {
			otherUsrId = chat.Participants[1]
		}

		user, err := rt.db.GetUserInfo(otherUsrId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				context.Logger.WithError(err).Warn("The other participant of the chat doesn't exist, in the db")
				http.Error(writer, "Not found - other participant not exist", http.StatusNotFound)
				return
			}
			context.Logger.WithError(err).Warn("Error getting other participant info")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		chat.ChatName = user.UserName
		chat.ChatPhoto = user.UserPhoto
	}

	// Preparo la risposta
	responseChatJSON, marshalErr := json.Marshal(chat)
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseChatJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
