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

func (rt *_router) getChatParticipants(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("usrId", token).Warnf("tried to change group photo of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the photo of another group", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	// Tento di recuperare le info degli utenti
	var participants []database.User
	participants, err= rt.db.GetChatParticipantsInfo(chatId)
	if err != nil {
		context.Logger.WithError(err).WithField("chatId", chatId).Errorf("Error getting chat participants")
		http.Error(writer, "Unable to retrive participants of the chat", http.StatusInternalServerError)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"participants": participants})
	if marshalErr != nil {
		context.Logger.WithError(err).Errorf("Failed to marshal the user")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
