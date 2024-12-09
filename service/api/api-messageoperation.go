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

func (rt *_router) sendMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}")

	//Preparo la variabile che conterrà i valori della richiesta http
	//che mi aspetto di ricevere per questo metodo dell'endpoint
	requestJson := struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	}{}

	//Recupero il chat id dai parametri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	//Controllo che l'utente faccia effettivamente parte del gruppo
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.Warnf("user <%s> tried to send a message to the chat <%d> which he isn't a member of", token, chatId)
			http.Error(writer, "Forbidden - can't send a message of a chat which aren't member of", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the chat <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the chat", http.StatusInternalServerError)
		return
	}

	//Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("user <%s> want to send a message to chat <%d>", token, chatId)
	//Aggiungo il messaggio al db
	var newMessage database.Message
	newMessage.SenderId = token
	newMessage.ContentType = requestJson.ContentType
	newMessage.Content = requestJson.Content
	newMessage.DeliveryStatus = "sent"
	newMessage.IsForwarded = false

	if err := rt.db.InsertMessage(newMessage, chatId); err != nil {
		rt.baseLogger.WithError(err).Error("Error while inserting new message in the db")
		http.Error(writer, "Internal Server Error - Unable to send the message", http.StatusInternalServerError)
		return
	}

	//Invio la risposta senza corpo
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) deleteMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("DELETE request to endpoint /chats/{chat_id}/messages/{msg_id}")

	//Recupero l'id della chat e del messaggio dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid message id")
		http.Error(writer, "invalid msg_id parameter", http.StatusBadRequest)
	}

	//Controllo che il messaggio esista, e sia stato mandato dall'utente che lo vuole eliminare
	var senderId string
	senderId, err = rt.db.GetSenderIdByMsgId(msgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Warn("Message not found in the database")
			http.Error(writer, "Not Found - Message not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error while checking the sender id of a message")
		http.Error(writer, "Internal Server Error - Unable to delete the message", http.StatusInternalServerError)
		return
	}

	if senderId != token {
		rt.baseLogger.WithField("user", token).Warnf("Error while checking the sender id of a message")
		http.Error(writer, "Forbidden - You can't delete the message of another user", http.StatusForbidden)
		return
	}

	rt.baseLogger.WithField("msgId", msgId).Debug("tryong to remove the message")
	err = rt.db.RemoveMessage(msgId, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).WithFields(logrus.Fields{"msgId": msgId, "chatId": chatId}).Warn("Unable to find and remove a message with that msgId and chatId the")
			http.Error(writer, "Not Found - Unable to find the message by the msgId and chatId gived", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).WithFields(logrus.Fields{"msgId": msgId, "chatId": chatId}).Warn("Unable to find and remove a message with that msgId and chatId the")
		http.Error(writer, "Internal Server Error - Unable to delete the message", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.WithField("msgId", msgId).Info("message removed successfully")
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}/messages/{msg_id}")

	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id parameter")
		http.Error(writer, "invalid parameter", http.StatusBadRequest)
		return
	}

	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid message id parameter")
		http.Error(writer, "invalid parameter", http.StatusBadRequest)
		return
	}

	//Recupero l'id della chat a cui ba inoltrato il messaggio
	requestJson := struct {
		ChatToForward int `json:"chatToForward"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := rt.db.ForwardMessage(token, msgId, requestJson.ChatToForward); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Warn("Message not found in the database")
			http.Error(writer, "Not Found - Message not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error while forwarding message to chat")
		http.Error(writer, "Internal Server Error - Unable to forward message to chat", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Debug("message forwarded successfully")
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getMessageComments(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) commentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) uncommentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}
