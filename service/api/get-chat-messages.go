package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getChatMessages(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'endpoint /chat/{chat_id}/messages")

	// Recupero il valore di chat_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			if errors.Is(numErr.Err, strconv.ErrSyntax) {
				context.Logger.WithError(numErr.Err).Error("the param chat_id is not a valid number")
			} else if errors.Is(numErr.Err, strconv.ErrRange) {
				context.Logger.WithError(numErr.Err).Error("the param chat_id range is out of range")
			} else {
				context.Logger.WithError(err).Error("Error parsing param chat_id")
			}
			http.Error(writer, "invalid param chat_id", http.StatusBadRequest)
			return
		}
	}

	// Controllo se l'utente che ha effettuato l'accesso faccia parte della chat di cui vuole ricavare i messaggi
	if exist, err := rt.db.CheckIfUserIsParticipant(chatId, token); err != nil {
		context.Logger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	} else if !exist {
		context.Logger.WithField("usrId", token).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user logged in is not a participant of the chat", http.StatusForbidden)
		return
	}

	// Recupero i messaggi della chat dal database
	chatMessages, err := rt.db.GetChatMessages(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("The chat have no messages, shouldn't be possible")
			http.Error(writer, "The chat is empty", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat messages")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Se tutti i controlli vanno a buon fine procedo a preparare la risposta e inviarla
	responseMessagesJSON, marshalErr := json.Marshal(map[string]interface{}{"messages": chatMessages})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseMessagesJSON); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
