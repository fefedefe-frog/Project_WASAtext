package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setMessageStatusToRead(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Recupero il valore di chat_id e msg_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing param chat_id")
		http.Error(writer, "Bad request - Invalid param chat_id", http.StatusBadRequest)
		return
	}
	var msgId int
	msgId, err = strconv.Atoi(params.ByName("msg_id"))
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing param msg_id")
		http.Error(writer, "Bad request - Invalid param msg_id", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente che ha effettuato l'accesso faccia parte della chat del messaggio che vuole impostare come letto
	if exist, err := rt.db.CheckIfUserIsParticipant(chatId, usrId); err != nil {
		context.Logger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	} else if !exist {
		context.Logger.WithField("usrId", usrId).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user logged in is not a participant of the chat", http.StatusForbidden)
		return
	}

	// Aggiorno lo stato del messaggio
	if err := rt.db.UpdateMessageDeliveryStatusToRead(chatId, msgId, usrId); err != nil {
		context.Logger.WithError(err).Error("Error updating message status")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Invio la risposta senza corpo
	writer.WriteHeader(http.StatusAccepted)
}
