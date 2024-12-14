package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getMessageComments(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, _ string) {
	context.Logger.Info("GET request to endpoint /chats/{chat_id}/messages/{msg_id}/comments")

	// Recupero l'id del messaggio dai paramentri dell'endpoint
	msgId, err := strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid message id")
		http.Error(writer, "invalid msg_id parameter", http.StatusBadRequest)
		return
	}

	var messageComments []database.Comment
	messageComments, err = rt.db.GetMessageComments(msgId)
	if err != nil {
		if errors.Is(err, database.ErrMessageHaveNoComments) {
			http.Error(writer, "Not Found - message doesn't have any comment yet", http.StatusNotFound)
			return
		}else if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("Error while checking the message comments")
			http.Error(writer, "Not Found - message don't exist", http.StatusNotFound)
			return
		}
	}

	// Preparo la risposta contenente i commenti del messaggio
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"comments": messageComments})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}