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

func (rt *_router) addToGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	var requestJson = struct {
		UsrIdToAdd string `json:"usrIdToAdd"`
	}{}

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Warn("user tried to modify group members of group which he isn't a member of")
			http.Error(writer, "Forbidden - can't add users to a group where you aren't part off", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	// Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("user <%s> request to add user <%s> to the group", token, requestJson.UsrIdToAdd)

	// Controllo che l'utente da aggiungere non faccia già parte del gruppo
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, requestJson.UsrIdToAdd); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithFields(logrus.Fields{"usrId": requestJson.UsrIdToAdd, "groupId": chatId}).Warn("The user you are trying to add is already part of the group")
			http.Error(writer, "Forbidden - can't add users to a group that he is already a member of", http.StatusForbidden)
			return
		}
		context.Logger.WithFields(logrus.Fields{"usrId": requestJson.UsrIdToAdd, "groupId": chatId}).Error("Error while checking if the user is a member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	// Aggiungo l'utente alla chat
	err = rt.db.InsertUserInChat(requestJson.UsrIdToAdd, chatId)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding user to group")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Debug("Successfully added user to group")

	// Tento di recuperare le info degli utenti
	var participants []database.User
	participants, err= rt.db.GetChatParticipantsInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("chatId", chatId).Error("chat not found in the chat_participants_table")
			http.Error(writer, "Chat not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat participants")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
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
