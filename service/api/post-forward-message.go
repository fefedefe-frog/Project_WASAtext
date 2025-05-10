package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (rt *_router) forwardMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid message id parameter")
		http.Error(writer, "invalid parameter", http.StatusBadRequest)
		return
	}
	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
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

	// Controllo se l'utente che vuole inoltrare il messaggio fa parte della chat
	// del messaggio da inoltrare, e della chat dove verrà inoltrato il messaggio
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, usrId)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Warn("user tried to forward a message of a chat which he isn't a member of")
		http.Error(writer, "Forbidden - can't forward a message of a chat where you aren't part off", http.StatusForbidden)
		return
	}
	var isParticipatForward bool
	isParticipatForward, err = rt.db.CheckIfUserIsParticipant(requestJson.ChatToForward, usrId)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": usrId, "groupId": requestJson.ChatToForward}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipatForward {
		context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": requestJson.ChatToForward}).Warn("user tried to forward a message to a chat which he isn't a member of")
		http.Error(writer, "Forbidden - can't forward a message to a chat where you aren't part off", http.StatusForbidden)
		return
	}

	if err := rt.db.ForwardMessage(usrId, msgId, requestJson.ChatToForward); err != nil {
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
