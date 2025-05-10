package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Recupero il chatId dai paramentri del'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "Bad request - Invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, usrId)
	if err != nil {

		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Warn("user tried to change a photo of group which he isn't a member of")
		http.Error(writer, "Forbidden - can't change the photo of a group where you aren't part off", http.StatusForbidden)
		return
	}

	// Controllo che la chat che si vuole modificare sia un gruppo
	if isGroup, err := rt.db.IsAGroup(chatId); !isGroup {

		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).WithField("groupChatId", chatId).Error("Chat not found")
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			context.Logger.WithError(err).Errorf("Unable to check if the chat <%d> is a group", chatId)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		context.Logger.WithField("chatId", chatId).Warn("User tried to change a photo of a chat that isn't a group")
		http.Error(writer, "Forbidden - can't change the photo of a non group chat", http.StatusForbidden)
		return
	}

	// Limito la memoria per il parsing del form
	err = request.ParseMultipartForm(32 << 20)
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing multipart form")
		http.Error(writer, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Recupero il valore della nuova foto
	photoFile, _, err := request.FormFile("newGroupPhoto")
	if err != nil {
		context.Logger.WithError(err).Error("Error getting file content")
		http.Error(writer, "Bad request - Error getting file content", http.StatusBadRequest)
		return
	}
	defer func() {
		err := photoFile.Close()
		if err != nil {
			context.Logger.WithError(err).Error("Error closing file")
		}
	}()

	// Leggo il contenuto dell'immagine
	var newGroupPhoto []byte
	newGroupPhoto, err = io.ReadAll(photoFile)
	if err != nil {
		context.Logger.WithError(err).Error("Error reading file")
		http.Error(writer, "Internal server error - Error reading file", http.StatusInternalServerError)
		return
	}

	context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Info("request to change group photo")
	// Aggiorno la propic del gruppo nel database
	if err := rt.db.SetGroupPhoto(chatId, newGroupPhoto); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Errorf("Group chat <%d> not found in database", chatId)
			http.Error(writer, "Group chat not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error changing group propic")
		http.Error(writer, "Unable to change group propic", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"chatPhoto": newGroupPhoto})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the new group chat photo")
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
