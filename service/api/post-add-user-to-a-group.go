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

func (rt *_router) addUserToGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chat/{chat_id}/users")

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
			context.Logger.WithField("usrId", token).Warnf("tried to change group photo of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the photo of another group", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
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
	// Aggiungo l'utente alla chat
	err = rt.db.InsertUserInChat(requestJson.UsrIdToAdd, chatId)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding user to group")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Debug("Successfully added user to group")

	// Preparo la risposta contentente la lista aggiornata di user id
	var chatParticipants []string
	chatParticipants, err = rt.db.GetChatPartecipants(chatId)
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

	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"users": chatParticipants})
	if marshalErr != nil {
		rt.baseLogger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseJSON); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
