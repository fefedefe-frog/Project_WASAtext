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

func (rt *_router) leaveChat(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Recupero il chatId dai parametri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
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
		context.Logger.WithFields(logrus.Fields{"usrId": token, "groupChatId": chatId}).Warn("tried leave the group which isn't a member of")
		http.Error(writer, "Forbidden - can't leave a group which are't member of", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveUserFromChat(token, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Error("The chat or the user doesn't exist")
			http.Error(writer, "Not found - The chat or the user doesn't exist", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Error("The chat or the user doesn't exist")
		http.Error(writer, "Internal server error - Unable to leave the chat", http.StatusForbidden)
		return
	}

	context.Logger.WithFields(logrus.Fields{"usrId": token, "chatId": chatId}).Debug("leaved the chat")
	writer.WriteHeader(http.StatusNoContent)
}
