package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (rt *_router) deleteMessage(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Recupero l'id della chat e del messaggio dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid message id")
		http.Error(writer, "invalid msg_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo che il messaggio esista, e sia stato mandato dall'utente che lo vuole eliminare
	var senderId string
	senderId, err = rt.db.GetSenderIdByMsgId(msgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("Message not found in the database")
			http.Error(writer, "Not Found - Message not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error while checking the sender id of a message")
		http.Error(writer, "Internal Server Error - Unable to delete the message", http.StatusInternalServerError)
		return
	}

	if senderId != usrId {
		context.Logger.WithField("user", usrId).Warnf("User tryied to delete a message of which he isn't the author")
		http.Error(writer, "Forbidden - You can't delete the message of another user", http.StatusForbidden)
		return
	}

	context.Logger.WithField("msgId", msgId).Debug("trying to remove the message")
	err = rt.db.RemoveMessage(msgId, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).WithFields(logrus.Fields{"msgId": msgId, "chatId": chatId}).Warn("Unable to find and remove a message with that msgId and chatId the")
			http.Error(writer, "Not Found - Unable to find the message by the msgId and chatId gived", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).WithFields(logrus.Fields{"msgId": msgId, "chatId": chatId}).Warn("Unable to find and remove a message with that msgId and chatId the")
		http.Error(writer, "Internal Server Error - Unable to delete the message", http.StatusInternalServerError)
		return
	}

	context.Logger.WithField("msgId", msgId).Debug("message removed successfully")
	writer.WriteHeader(http.StatusAccepted)
}
