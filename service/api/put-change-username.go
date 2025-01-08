package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) putChangeUserName(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

	var requestJson = struct {
		NewUserName string `json:"newUserName"`
	}{}

	// Decodifica del JSON nella richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("Richiesta di cambio nome da parte dell'user: %s || nuovo nome: %s", token, requestJson.NewUserName)

	// Controllo se il nome è valido secondo i requisiti richiesti
	if err := utilitytool.NameIsValid(requestJson.NewUserName); err != nil {
		switch {
		case errors.Is(err, utilitytool.ErrInvalidRegex):
			http.Error(writer, "Invalid name format, the name can't contain space at the start or end of the name", http.StatusBadRequest)
			context.Logger.Debug("Invalid name format")

		case errors.Is(err, utilitytool.ErrNameShort):
			http.Error(writer, "the name must be at least 3 character long", http.StatusBadRequest)
			context.Logger.Debug("login name to short")

		case errors.Is(err, utilitytool.ErrNameLong):
			http.Error(writer, "the name can be max 16 character long", http.StatusBadRequest)
			context.Logger.Debug("login name to long")
		}
		return
	}

	// Aggiorno l'username nel database
	if err := rt.db.SetUserName(token, requestJson.NewUserName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Errorf("User: %s not found in database", token)
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error changing user name")
		http.Error(writer, "Unable to rename", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"chatName": requestJson.NewUserName})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the new username")
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
