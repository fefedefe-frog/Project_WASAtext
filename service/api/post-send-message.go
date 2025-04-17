package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (rt *_router) sendMessage(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Preparo la variabile che conterrà i valori della richiesta http
	// che mi aspetto di ricevere per questo metodo dell'endpoint
	requestJson := struct {
		ContentType string	`json:"contentType"`
		Content     string	`json:"content"`
		RespondTo   int		`json:"respondTo"`
	}{}

	// Recupero il chat id dai parametri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}

	// Controllo che l'utente faccia effettivamente parte del gruppo
	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).WithFields(logrus.Fields{"usrId": token, "chatId": chatId}).Errorf("Error while checking if the user is member of the group")
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.WithFields(logrus.Fields{"usrId": token, "groupId": chatId}).Warn("user tried send a message in a chat of which he isn't a member of")
		http.Error(writer, "Forbidden - can't send a message in a chat where you aren't member off", http.StatusForbidden)
		return
	}

	// Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.WithFields(logrus.Fields{"usrId": token, "chatId": chatId}).Info("user want to send a message in the chat")

	// Aggiungo il messaggio al db
	var newMessage database.Message
	newMessage.SenderId = token
	newMessage.RespondTo= requestJson.RespondTo
	newMessage.ContentType = requestJson.ContentType
	newMessage.Content = requestJson.Content
	newMessage.DeliveryStatus = "sent"
	newMessage.IsForwarded = false

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
