package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (rt *_router) getChatParticipants(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "Bad request - Invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Warn("user tried to retrive the participants of a group which he isn't a member of")
		http.Error(writer, "Forbidden - can't retrive the info of a group where you aren't part off", http.StatusForbidden)
		return
	}

	// Tento di recuperare le info degli utenti
	var participants []database.User
	participants, err = rt.db.GetChatParticipantsInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Error("Error getting chat participants")
			http.Error(writer, "Not found - The participants of the chat doesn't exist", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).WithField("chatId", chatId).Error("Error getting chat participants")
		http.Error(writer, "Internal server error - Unable to retrive participants of the chat", http.StatusInternalServerError)
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
