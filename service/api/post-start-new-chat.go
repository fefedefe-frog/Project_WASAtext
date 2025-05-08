package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) startNewChat(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

	// Limito la memoria per il parsing del form
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing multipart form")
		http.Error(writer, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Recupero le info della chat dai dati del form
	rawChatInfo := request.FormValue("chatInfo")
	var chatInfo struct {
		ChatName     string   `json:"chatName"`
		ChatPhoto    string   `json:"chatPhoto"`
		IsGroup      bool     `json:"isGroup"`
		Participants []string `json:"participants"`
	}
	if err := json.Unmarshal([]byte(rawChatInfo), &chatInfo); err != nil {
		context.Logger.WithError(err).Error("Error unmarshalling chatInfo")
		http.Error(writer, "errore recuper informazioni dal form", http.StatusBadRequest)
		return
	}

	// Recupero i dati del messaggio da inviare per iniziare la chat
	contentType := request.FormValue("contentType")
	var messageContent []byte
	if contentType == "photo" {
		// Carico l'immagine contenuta nella richiesta http
		file, _, err := request.FormFile("content")
		if err != nil {
			context.Logger.WithError(err).Error("Error getting file content")
			http.Error(writer, "Internal server error - Error getting file content", http.StatusInternalServerError)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				context.Logger.WithError(err).Error("Error closing file")
			}
		}()

		// Leggo il contenuto dell'immagine
		messageContent, err = io.ReadAll(file)
		if err != nil {
			context.Logger.WithError(err).Error("Error reading file")
			http.Error(writer, "Internal server error - Error reading file", http.StatusInternalServerError)
			return
		}

	} else {
		messageContent = []byte(request.FormValue("content"))
	}

	// Preparo i dati della chat
	if chatInfo.IsGroup {
		if chatInfo.ChatPhoto == "" {
			chatInfo.ChatPhoto = database.DefaultGroupPhotoBase64
		}
		if chatInfo.ChatName == "" {
			chatInfo.ChatName = "Gruppo"
		}
	} else {
		chatInfo.ChatPhoto = ""
		chatInfo.ChatName = ""
	}

	// Decodifica la stringa Base64 in byte
	var groupPhotoData []byte
	groupPhotoData, err = base64.StdEncoding.DecodeString(chatInfo.ChatPhoto)
	if err != nil {
		http.Error(writer, "Internal Server Error - Unable to decode the photo", http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Unable to decode the base64 string of the photo")
		return
	}

	var newChatId int
	chatInfo.Participants = append(chatInfo.Participants, token)
	newChatId, err = rt.db.InsertNewChat(token, chatInfo.ChatName, groupPhotoData, chatInfo.Participants, chatInfo.IsGroup, contentType, messageContent)
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
		http.Error(writer, "Internal server error - Error while recovering the info of the chat just created", http.StatusInternalServerError)
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
			http.Error(writer, "Internal server error - Error while recovering the info of the chat just created", http.StatusInternalServerError)
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
