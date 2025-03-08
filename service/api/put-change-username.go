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

func (rt *_router) setMyUserName(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {

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

	// Controllo se l'username scelto non sia già in uso
	if idFinded, err := rt.db.GetUsrIdByName(requestJson.NewUserName); err != nil {
		// Banalmente controllo se, data la chiamata alla funzione precedente, se l'username non è presente nel db
		// allora dovrei ricevere l'errore ErrNoRows, ovvero non è stato trovato nessun utente con quel nome
		if !errors.Is(err, sql.ErrNoRows) {
			// Se ho un errore, e non è l' errore no rows, c'è stato un problema
			http.Error(writer, "Internal Server Error - Unable to check if the username is available", http.StatusInternalServerError)
			context.Logger.WithError(err).Error("unable to check if the username selected is available")
			return
		}
	} else if idFinded != "" { // La ricerca tramite usename ha dato risultato quindi l'username è già usato da altri
		http.Error(writer, "Username already exists", http.StatusBadRequest)
		context.Logger.Debug("username selected is not available")
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
