package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyConversations(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

	// Tento di recuperare le chat di quell'user
	chats, err := rt.db.GetUserChats(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, database.ErrUserNoChat) {
			context.Logger.Info("User doesn't have chat")
			http.Error(writer, "The user doesn't have any chat", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Errorf("Error getting user %s chats", token)
		http.Error(writer, "Unable to retrive info", http.StatusInternalServerError)
		return
	}

	// Preparo la risposta json
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"chats": chats})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat info")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta json
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
