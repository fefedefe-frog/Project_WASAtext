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

	// Preparo la variabile che conterrà i valori della richiesta http
	// che mi aspetto di ricevere per questo metodo dell'endpoint
	requestJson := struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	}{}

	// Recupero il chat id dai parametri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	// Controllo che l'utente faccia effettivamente parte del gruppo
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

	// Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("user <%s> want to send a message to chat <%d>", token, chatId)

	// Aggiungo il messaggio al db
	var newMessage database.Message
	newMessage.SenderId = token
	newMessage.ContentType = requestJson.ContentType
	newMessage.Content = requestJson.Content
	newMessage.DeliveryStatus = "sent"
	newMessage.IsForwarded = false

	if err := rt.db.InsertMessage(newMessage, chatId); err != nil {
		context.Logger.WithError(err).Error("Error while inserting new message in the db")
		http.Error(writer, "Internal Server Error - Unable to send the message", http.StatusInternalServerError)
		return
	}

	// Invio la risposta senza corpo
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) deleteMessage(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("DELETE request to endpoint /chats/{chat_id}/messages/{msg_id}")

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

	if senderId != token {
		context.Logger.WithField("user", token).Warnf("Error while checking the sender id of a message")
		http.Error(writer, "Forbidden - You can't delete the message of another user", http.StatusForbidden)
		return
	}

	context.Logger.WithField("msgId", msgId).Debug("tryong to remove the message")
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

	context.Logger.WithField("msgId", msgId).Info("message removed successfully")
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}/messages/{msg_id}")

	msgId, err := strconv.Atoi(params.ByName("msg_id"))
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

	if err := rt.db.ForwardMessage(token, msgId, requestJson.ChatToForward); err != nil {
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

func (rt *_router) getMessageComments(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, _ string) {
	context.Logger.Info("GET request to endpoint /chats/{chat_id}/messages/{msg_id}/comments")

	// Recupero l'id del messaggio dai paramentri dell'endpoint
	msgId, err := strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid message id")
		http.Error(writer, "invalid msg_id parameter", http.StatusBadRequest)
		return
	}

	var messageComments []database.Comment
	messageComments, err = rt.db.GetMessageComments(msgId)
	if err != nil {
		if errors.Is(err, database.ErrMessageHaveNoComments) {
			http.Error(writer, "Not Found - message doesn't have any comment yet", http.StatusNotFound)
			return
		}else if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("Error while checking the message comments")
			http.Error(writer, "Not Found - message don't exist", http.StatusNotFound)
			return
		}
	}

	// Preparo la risposta contenente i commenti del messaggio
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"comments": messageComments})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) commentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}/messages/{msg_id}/comments")

	// Recupero l'id della chat del messaggio dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo che l'utente che vuole commentare il messaggio faccia parte della chat di quel messaggio
	var isParticipant bool
	isParticipant, err =rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the user participant status")
		http.Error(writer, "Internal Server Error - Unable to check the user participant status", http.StatusInternalServerError)
		return
	}

	if !isParticipant {
		context.Logger.Warn("the user is not a participant of the chat of the message that was try to comment")
		http.Error(writer, "Forbidden - Can't comment a message of a group where aren't participant of", http.StatusForbidden)
		return
	}

	requestJson := struct {
		Content string `json:"content"`
	}{}

	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := rt.db.CommentMessage(msgId, token, requestJson.Content); err != nil {
		context.Logger.WithError(err).Error("Error while commenting message to chat")
		http.Error(writer, "Internal Server Error - Unable to comment message to chat", http.StatusInternalServerError)
		return
	}

	context.Logger.Debug("message commented successfully")
	writer.WriteHeader(http.StatusNoContent)
}

func (rt *_router) uncommentMessage(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("DELETE request to endpoint /chats/{chat_id}/messages/{msg_id}/comments/{comment_id}")

	// Recupero l'id della chat e del commento dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var commentId int
	commentId, err = strconv.Atoi(params.ByName("comment_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid comment id")
		http.Error(writer, "invalid comment_id parameter", http.StatusBadRequest)
		return
	}

	// Controllo che l'utente che vuole rimuovere il commento sia colui che l'ha commentato in precedenza
	var isCommenter bool
	isCommenter, err =rt.db.CheckCommentAuthor(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the comment author")
		http.Error(writer, "Internal Server Error - Unable to check the comment author", http.StatusInternalServerError)
		return
	}

	if !isCommenter {
		context.Logger.Warn("the user is not the author of the comment that was try to uncomment")
		http.Error(writer, "Forbidden - Can't uncomment a message that isn't yours", http.StatusForbidden)
		return
	}

	if err := rt.db.UncommentMessage(commentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Error("The comment that is bein deleted doesn't exist")
			http.Error(writer, "Not Found - Comment not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error while uncommenting message to chat")
		http.Error(writer, "Internal Server Error - Unable to uncomment message to chat", http.StatusInternalServerError)
		return
	}

	context.Logger.Debug("message uncommented successfully")
	writer.WriteHeader(http.StatusNoContent)
}
