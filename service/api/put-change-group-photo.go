package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"strconv"
)

func (rt *_router) setGroupPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	var requestJson = struct {
		NewGroupPhoto string `json:"newGroupPhoto"`
	}{}

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("usrId", token).Warnf("tried to change group photo of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the photo of another group", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	// Controllo che la chat che si vuole modificare sia un gruppo
	if isGroup, err := rt.db.IsAGroup(chatId); !isGroup {

		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Errorf("Chat <%d> not found", chatId)
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

	// Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	// Semplice controllo della stringa base64 per assicurarsi
	// che la stringa contenga solo caratteri usati dalla codifica base64
	re := regexp.MustCompile(`^([A-Za-z0-9+/=]+)$`)
	if !re.MatchString(requestJson.NewGroupPhoto) {
		http.Error(writer, "Bad request - The photo isn't codified correctly, or is not a photo", http.StatusBadRequest)
		context.Logger.Debug("The photo received in input is not in the base64 format")
		return
	}

	// Verifica che la stringa sia in formato base64 valido
	var photoData []byte
	photoData, err = base64.StdEncoding.DecodeString(requestJson.NewGroupPhoto)
	if err != nil {
		http.Error(writer, "Internal Server Error - Unable to decode the photo", http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Unable to decode the base64 string of the photo")
		return
	}

	context.Logger.Infof("user <%s> request to change group photo of group <%d>", token, chatId)
	// Aggiorno la propic del gruppo nel database
	if err := rt.db.SetGroupPhoto(chatId, photoData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Errorf("Group chat <%d> not found in database", chatId)
			http.Error(writer, "Group chat not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing group propic")
		http.Error(writer, "Unable to change group propic", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"chatPhoto": requestJson.NewGroupPhoto})
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
