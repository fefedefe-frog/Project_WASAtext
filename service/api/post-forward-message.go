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

func (rt *_router) forwardMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}/messages/{msg_id}")

	msgId, err := strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid message id parameter")
		http.Error(writer, "invalid parameter", http.StatusBadRequest)
		return
	}

	// Recupero l'id della chat a cui ba inoltrato il messaggio
	requestJson := struct {
		ChatToForward int `json:"chatToForward"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := rt.db.ForwardMessage(token, msgId, requestJson.ChatToForward); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("Message not found in the database")
			http.Error(writer, "Not Found - Message not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error while forwarding message to chat")
		http.Error(writer, "Internal Server Error - Unable to forward message to chat", http.StatusInternalServerError)
		return
	}

	context.Logger.Debug("message forwarded successfully")
	writer.WriteHeader(http.StatusNoContent)
}
