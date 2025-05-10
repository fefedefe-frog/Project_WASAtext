package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func (rt *_router) sendMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Recupero il chat id dai parametri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}

	// Controllo che l'utente faccia effettivamente parte del gruppo
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, usrId)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": usrId, "chatId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": usrId, "groupId": chatId}).Warn("user tried send a message in a chat of which he isn't a member of")
		http.Error(writer, "Forbidden - can't send a message in a chat where you aren't member off", http.StatusForbidden)
		return
	}

	// Limito la memoria per il parsing del form
	err = request.ParseMultipartForm(32 << 20)
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing multipart form")
		http.Error(writer, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	textContent := request.FormValue("textContent")

	var photoContent []byte
	// Carico l'immagine contenuta nella richiesta http
	photoContentFile, _, err := request.FormFile("content")
	if err != nil {
		context.Logger.WithError(err).Error("Error getting file content")
		http.Error(writer, "Internal server error - Error getting file content", http.StatusInternalServerError)
		return
	}
	defer func() {
		err := photoContentFile.Close()
		if err != nil {
			context.Logger.WithError(err).Error("Error closing file")
		}
	}()

	// Leggo il contenuto dell'immagine
	photoContent, err = io.ReadAll(photoContentFile)
	if err != nil {
		context.Logger.WithError(err).Error("Error reading file")
		http.Error(writer, "Internal server error - Error reading file", http.StatusInternalServerError)
		return
	}

	if textContent == "" && len(photoContent) == 0 {
		context.Logger.Warn("No text or photo content found in the request for starting a new chat")
		http.Error(writer, "Bad request - The message need to have at least one type of content between text or an image", http.StatusBadRequest)
		return
	}

	// Recupero e converto in intero l'id del messaggio a cui sto rispondendo
	respondToStr := request.FormValue("respondTo")
	respondTo, err := strconv.Atoi(respondToStr)
	if err != nil {
		context.Logger.WithError(err).Error("Error converting respondTo to int")
		http.Error(writer, "Internal server error - Error converting string to int", http.StatusInternalServerError)
		return
	}

	context.Logger.WithFields(logrus.Fields{"usrId": usrId, "chatId": chatId}).Info("user want to send a message in the chat")

	// Aggiungo il messaggio al db
	newMessage := database.Message{
		MsgId:          -1,
		SenderId:       usrId,
		RespondTo:      respondTo,
		TextContent:    textContent,
		PhotoContent:   photoContent,
		DeliveryStatus: "sent",
		IsForwarded:    false,
	}

	newMessage.MsgId, err = rt.db.InsertMessage(newMessage, chatId)
	if err != nil {
		context.Logger.WithError(err).Error("Error while inserting new message in the db")
		http.Error(writer, "Internal Server Error - Unable to send the message", http.StatusInternalServerError)
		return
	}

	// Se tutti i controlli vanno a buon fine procedo a preparare la risposta e inviarla
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"message": newMessage})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
