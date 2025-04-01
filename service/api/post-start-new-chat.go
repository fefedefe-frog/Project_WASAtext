package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) startNewChat(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

	requestJson := struct {
		ChatInfo struct {
			ChatName     string   `json:"chatName"`
			ChatPhoto    string   `json:"chatPhoto"`
			IsGroup      bool     `json:"isGroup"`
			Participants []string `json:"participants"`
		} `json:"chatInfo"`

		FirstMessage struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"firstMessage"`
	}{}

	chatInfo := &requestJson.ChatInfo
	firstMessage := &requestJson.FirstMessage

	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		context.Logger.WithError(err).Error("error while decoding body request")
		http.Error(writer, "Bad Request - invalid json format", http.StatusBadRequest)
		return
	}

	if chatInfo.IsGroup {
		if chatInfo.ChatPhoto == "" {
			chatInfo.ChatPhoto = database.DefaultGroupPhotoBase64
		}
		if chatInfo.ChatName == "" {
			chatInfo.ChatName = "Gruppo"
		}
	}

	// Decodifica la stringa Base64 in byte
	var groupPhotoData []byte
	groupPhotoData, err = base64.StdEncoding.DecodeString(chatInfo.ChatPhoto)
	if err != nil {
		http.Error(writer, "Internal Server Error - Unable to decode the photo", http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Unable to decode the base64 string of the photo")
		return
	}

	// Assegno il valore di messageContent a seconda del contenuto del messaggio
	var messageContent interface{}
	if firstMessage.ContentType == "photo" {
		var convErr error
		messageContent, convErr = base64.StdEncoding.DecodeString(firstMessage.Content)
		if convErr != nil {
			http.Error(writer, "Internal Server Error - Unable to decode the photo", http.StatusInternalServerError)
			context.Logger.WithError(convErr).Error("Unable to decode the base64 string of the photo")
			return
		}
	} else {
		messageContent = firstMessage.Content
	}

	var newChatId int
	newChatId, err = rt.db.InsertNewChat(token, chatInfo.ChatName, groupPhotoData, chatInfo.Participants, chatInfo.IsGroup, messageContent)
	if err != nil {
		context.Logger.WithError(err).Error("unable to insert a new chat in the db")
		http.Error(writer, "Internal server error - unable to create the chat", http.StatusInternalServerError)
		return
	}

	context.Logger.WithField("chatId", newChatId).Info("chat created successfully")

	var chat database.Chat
	chat, err = rt.db.GetChatInfo(newChatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Warn("The chat doesn't exist, in the db")
			http.Error(writer, "Not found - chat not exist", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat info")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Controllo se la chat è un gruppo o meno, se non è un gruppo procedo
	// a recuperare le informazioni dell'altro utente partecipante alla chat
	if !chat.IsGroup {

		otherUsrId := chat.Participants[0]
		if otherUsrId == token {
			otherUsrId = chat.Participants[1]
		}

		user, err := rt.db.GetUserInfo(otherUsrId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				context.Logger.WithError(err).Warn("The other participant of the chat doesn't exist, in the db")
				http.Error(writer, "Not found - other participant not exist", http.StatusNotFound)
				return
			}
			context.Logger.WithError(err).Warn("Error getting other participant info")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		chat.ChatName = user.UserName
		chat.ChatPhoto = user.UserPhoto
	}

	// preparo la risposta, contenente le info della chat
	responseJSON, marshalErr := json.Marshal(chat)
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat")
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
