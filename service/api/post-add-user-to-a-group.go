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

func (rt *_router) addToGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	var requestJson = struct {
		UsrIdToAdd string `json:"usrId"`
	}{}

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "Bad request - Invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, usrId)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Warn("user tried add a particpiant to a group which he isn't a member of")
		http.Error(writer, "Forbidden - can't add users to a group where you aren't part off", http.StatusForbidden)
		return
	}

	// Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.WithField("usrId", usrId).Info("user request to add user/s to the group")

	// Controllo che l'utente da aggiungere non faccia già parte del gruppo
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, requestJson.UsrIdToAdd)
	if err != nil {
		context.Logger.WithFields(logrus.Fields{"usrId": requestJson.UsrIdToAdd, "groupId": chatId}).Error("Error while checking if the user is a member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": requestJson.UsrIdToAdd, "groupId": chatId}).Warn("The user you are trying to add is already part of the group")
		http.Error(writer, "Can't add users to a group that he is already a member of", http.StatusNoContent)
		return
	}

	// Aggiungo l'utente alla chat
	err = rt.db.InsertUserInChat(requestJson.UsrIdToAdd, chatId)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding user to group")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	context.Logger.Debug("Successfully added user to group")

	// Tento di recuperare le info degli utenti
	var participants []database.User
	participants, err = rt.db.GetChatParticipantsInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Error("Unable to retrive one or more user info from the db")
			http.Error(writer, "Not found - Not found one or more user info from the db", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat participants info")
		http.Error(writer, "Internal Server Error - Unable to retrive the info", http.StatusInternalServerError)
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
