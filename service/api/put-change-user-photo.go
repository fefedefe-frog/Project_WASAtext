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
)

func (rt *_router) setMyPhoto(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

	var requestJson = struct {
		NewUserPhoto string `json:"newUserPhoto"`
	}{}

	// Decodifico la richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	// Semplice controllo della stringa base64 per assicurarsi
	// che la stringa contenga solo caratteri usati dalla codifica base64
	re := regexp.MustCompile(`^([A-Za-z0-9+/=]+)$`)
	if !re.MatchString(requestJson.NewUserPhoto) {
		http.Error(writer, "Bad request - The photo isn't codified correctly, or is not a photo", http.StatusBadRequest)
		context.Logger.Debug("The photo received in input is not in the base64 format")
		return
	}

	// Verifica che la stringa sia in formato base64 valido
	var photoData []byte
	photoData, err = base64.StdEncoding.DecodeString(requestJson.NewUserPhoto)
	if err != nil {
		http.Error(writer, "Internal Server Error - Unable to decode the photo", http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Unable to decode the base64 string of the photo")
		return
	}

	// Aggiorno la propic nel database
	context.Logger.Infof("Richiesta di cambio propic da parte dell'user: %s || nuova propic: %s...", token, requestJson.NewUserPhoto[:10])
	if err := rt.db.SetUserPhoto(token, photoData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Errorf("User: %s not found in database", token)
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error changing user propic")
		http.Error(writer, "Unable to change propic", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"chatName": requestJson.NewUserPhoto})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the new user photo")
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
