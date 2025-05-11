package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"encoding/json"
	"errors"
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
	var user database.User
	usrId, err = rt.db.GetUsrIdByName(usernameLower)

	// Controllo se l'utente esiste ed è presente nel database se non è presente lo creo e provo ad inserirlo nel database
	// in caso di riuscita preparo la rosposta http e la invio
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("db error", err).Debug("utente non presente nel database")
			context.Logger.WithField("username", usernameLower).Debug("Creazione e aggiunta del nuovo utente al database")
			user, err = rt.db.InsertNewUser(requestJson.UserName)
			if err != nil {
				context.Logger.WithError(err).Error("Impossibile aggiungere nuovo user al database")
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			context.Logger.WithError(err).Error("Errore durante il recupero dei dati")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		user, err = rt.db.GetUserInfo(usrId)
		if err != nil {
			context.Logger.WithError(err).Error("unable to retrive user info")
			http.Error(writer, "Internal server error - unable to retrive user info after login", http.StatusInternalServerError)
			return
		}
	}

	rt.sendJsonResponse(writer, user, context)
	context.Logger.WithField("usrId", usrId).Info("login effettuato")
}

func (rt *_router) sendJsonResponse(writer http.ResponseWriter, user database.User, context reqcontext.RequestContext) {

	// Creo la risposta http contentente il token di autorizzazione e l'user
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"user": user})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the user")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return

	}

	// Scrivo nell'header della risposta il token bearer che in questo caso corrisponde all'usrId
	writer.Header().Set("Authorization", "Bearer "+user.UsrId)

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
