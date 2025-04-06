package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getMessageComments(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Recupero l'id della chat del messaggio dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo che l'utente che vuole commentare il messaggio faccia parte della chat di quel messaggio
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the user participant status")
		http.Error(writer, "Internal Server Error - Unable to check the user participant status", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.Warn("the user is not a participant of the chat of the message that was try to comment")
		http.Error(writer, "Forbidden - Can't get the comments of a message from a chat where aren't participant of", http.StatusForbidden)
		return
	}

	var messageComments []database.Comment
	messageComments, err = rt.db.GetMessageComments(msgId)
	if err != nil {
		context.Logger.WithError(err).Warn("Error while checking the message comments")
		http.Error(writer, "Internal server errror", http.StatusInternalServerError)
		return
	}

	// Preparo la risposta contenente i commenti del messaggio
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"comments": messageComments})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the message comments")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
