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

func (rt *_router) sendMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats/{chat_id}")

	requestJson := struct {
		ContentType 	string `json:"contentType"`
		Content     	string `json:"content"`
	}{}

	//Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant{
		if errors.Is(err, sql.ErrNoRows){
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
	newMessage.SenderId= token
	newMessage.ContentType = requestJson.ContentType
	newMessage.Content = requestJson.Content
	newMessage.DeliveryStatus= "sent"
	newMessage.IsForwarded= false

	if err := rt.db.InsertMessage(newMessage, chatId); err != nil {
		rt.baseLogger.WithError(err).Error("Error while inserting new message in the db")
		http.Error(writer, "Internal Server Error - Unable to send the message", http.StatusInternalServerError)
		return
	}

	//Invio la risposta senza corpo
	writer.WriteHeader(http.StatusNoContent)
	return
}

func (rt *_router) deleteMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) forwardMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) getMessageComments(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) commentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}

func (rt *_router) uncommentMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
}