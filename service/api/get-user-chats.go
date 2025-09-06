package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyConversations(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Tento di recuperare le chat di quell'user
	chats, err := rt.db.GetChatsOfUser(usrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("usrId", usrId).Debug("User doesn't have chat")
			chats = nil
		} else {
			context.Logger.WithError(err).WithField("usrId", usrId).Error("Error getting user chats")
			http.Error(writer, "Internal server error - Unable to retrive one or more conversation/s", http.StatusInternalServerError)
			return
		}
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
