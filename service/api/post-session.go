package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (rt *_router) doLogin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext) {

	var requestJson = struct {
		UserName string `json:"userName"`
	}{}

	// Decodifica il corpo della richiesta JSON
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		context.Logger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	// Converto il contenuto in lower case
	usernameLower := strings.ToLower(requestJson.UserName)

	context.Logger.Infof("Tentativo di login da user: '%s'", usernameLower)
	// Controllo che rispetti i regex richiesti e la lunghezza minima e massima
	if err := utilitytool.NameIsValid(usernameLower); err != nil {
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

	var usrId string
	usrId, err = rt.db.GetUsrIdByName(usernameLower)

	// Controllo se l'utente esiste ed è presente nel database se non è presente lo creo e provo ad inserirlo nel database
	// in caso di riuscita preparo la rosposta http e la invio
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("db error", err).Debug("utente non presente nel database")

			context.Logger.Debug(fmt.Sprintf("Creazione e aggiunta del nuovo utente '%s' al database", usernameLower))
			user, err := rt.db.InsertNewUser(requestJson.UserName)

			if err != nil {
				context.Logger.WithError(err).Error("Impossibile aggiungere nuovo user al database")
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}
			rt.sendJsonResponse(writer, user.UsrId, context)
			return
		} else {
			context.Logger.WithError(err).Error("Errore durante il recupero dei dati")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	context.Logger.WithField("usrId", usrId).Info("login effettuato")
	rt.sendJsonResponse(writer, usrId, context)
}

func (rt *_router) sendJsonResponse(writer http.ResponseWriter, usrId string, context reqcontext.RequestContext) {

	// Creo la risposta http contentente il token di autorizzazione e l'usrId (in questo caso entrambi sono la stessa cosa
	response := map[string]string{
		"usrId": usrId,
	}

	// Scrivo nell'header della risposta il token bearer che in questo caso corrisponde all'usrId
	writer.Header().Set("Authorization", "Bearer "+usrId)

	// Scrivo nella risposta l'usrId
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
