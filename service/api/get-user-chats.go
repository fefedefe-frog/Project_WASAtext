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

func (rt *_router) getUserChats(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("GET request to endpoint /chats")

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

	rt.baseLogger.Debugf("user ha {%d} chats", len(chats))
	for i, chat := range chats {
		rt.baseLogger.Debugf("chat{%d}:\n\t chatId -> {%d}\n\t chatName -> {%s}\n\t chatPhoto -> {%s}\n\t participants -> {%d}\n", i, chat.ChatId, chat.ChatName, chat.ChatPhoto[:10], len(chat.Participants))
	}


	// Scrivo la risposta json
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(map[string]interface{}{"chats": chats}); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
