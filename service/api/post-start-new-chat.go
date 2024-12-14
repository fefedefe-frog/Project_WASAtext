package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) startNewChat(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("POST request to endpoint /chats")

	requestJson := struct {
		ChatName     string   `json:"chatName"`
		ChatPhoto    string   `json:"chatPhoto"`
		IsGroup      bool     `json:"isGroup"`
		Participants []string `json:"participants"`
	}{}

	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		context.Logger.WithError(err).Error("error while decoding body request")
		http.Error(writer, "Bad Request - invalid json format", http.StatusBadRequest)
		return
	}

	// aggiungo l'utente che effettua la richiesta alla lista dei partecipanti
	requestJson.Participants = append(requestJson.Participants, token)
	var newChatId int
	newChatId, err = rt.db.InsertNewChat(requestJson.Participants, requestJson.ChatName, requestJson.ChatPhoto, requestJson.IsGroup)
	if err != nil {
		context.Logger.WithError(err).Error("unable to insert a new chat in the db")
		http.Error(writer, "Internal server error - unable to create the chat", http.StatusInternalServerError)
		return
	}

	context.Logger.WithField("chatId", newChatId).Info("chat created successfully")

	// preparo la risposta
	var newChatInfo database.Chat
	newChatInfo, err = rt.db.GetChatInfo(newChatId)
	if err != nil {
		context.Logger.WithError(err).Error("unable to get the chat info")
		http.Error(writer, "Internal server error - unable to retrieve the info of the chat just created", http.StatusInternalServerError)
		return
	}

	var responseJSON []byte
	responseJSON, err = json.Marshal(newChatInfo)
	if err != nil {
		context.Logger.WithError(err).Error("unable to marshal the chat info")
		http.Error(writer, "Internal server error - failed json conversion while preparing the response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error - failed to send the response", http.StatusInternalServerError)
		return
	}
}