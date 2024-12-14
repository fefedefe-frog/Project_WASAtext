package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)
//TODO controllo del contenuto che si vuole inserire come commento, sono accettate solo emoji
func (rt *_router) commentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}/messages/{msg_id}/comments")

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
	isParticipant, err =rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the user participant status")
		http.Error(writer, "Internal Server Error - Unable to check the user participant status", http.StatusInternalServerError)
		return
	}

	if !isParticipant {
		context.Logger.Warn("the user is not a participant of the chat of the message that was try to comment")
		http.Error(writer, "Forbidden - Can't comment a message of a group where aren't participant of", http.StatusForbidden)
		return
	}

	requestJson := struct {
		Content string `json:"content"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := rt.db.CommentMessage(msgId, token, requestJson.Content); err != nil {
		context.Logger.WithError(err).Error("Error while commenting message to chat")
		http.Error(writer, "Internal Server Error - Unable to comment message to chat", http.StatusInternalServerError)
		return
	}

	context.Logger.Debug("message commented successfully")
	writer.WriteHeader(http.StatusNoContent)
}