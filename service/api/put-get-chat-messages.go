package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getConversation(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	var requestJson = struct {
		MsgId int `json:"msgId"`
	}{}

	// Recupero il valore di chat_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "Bad request - Invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente che ha effettuato l'accesso faccia parte della chat di cui vuole ricavare i messaggi
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, usrId)
	if err != nil {
		context.Logger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithField("usrId", usrId).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user is not a participant of the chat", http.StatusForbidden)
		return
	}

	// Decodifica il corpo della richiesta JSON
	err = json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	// Recupero i messaggi della chat dal database
	var chatMessages []database.Message
	chatMessages, err = rt.db.GetChatMessages(chatId, usrId, requestJson.MsgId)
	if err != nil {
		if errors.Is(err, database.ErrUpdateMessageStatus) {

			context.Logger.WithError(err).Warn("Error updating the message status of the messages retrived")
		} else {
			context.Logger.WithError(err).Error("Error getting chat messages")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
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
	if _, err := writer.Write(responseMessagesJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
