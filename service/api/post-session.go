package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"unicode/utf8"
)



func (rt *_router) postSession(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	context.Logger.Info("Richiesta dall'endpoint /session")

	var requestJson= struct {
		UserName string `json:"userName"`
	}{}

	//Decodifica il corpo della richiesta JSON
	err := json.NewDecoder(request.Body).Decode(&requestJson)

	if err != nil {
		http.Error(writer, "Invalid JSON format" , http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Info(fmt.Sprintf("Tentativo di login da user: %s", requestJson.UserName))
	//Controllo che rispetti i regex richiesti e la lunghezza minima e massima
	if !regexp.MustCompile(`^\S.*\S$`).MatchString(requestJson.UserName){
		http.Error(writer, "Invalid name format, the name can't contain space at the start or end of the name", http.StatusBadRequest)
		rt.baseLogger.Debug("Invalid name format")
		return
	}
	if utf8.RuneCountInString(requestJson.UserName)<3{	//Deve essere lungo almeno 3 caratteri
		http.Error(writer, fmt.Sprintf("%s", "the name must be at least 3 character long"), http.StatusBadRequest)
		rt.baseLogger.Debug("login name to short")
		return
	}
	if utf8.RuneCountInString(requestJson.UserName)>16{	//Deve essere lungo massimo 16 caratteri
		http.Error(writer, fmt.Sprintf("%s", "the name can be max 16 character long"), http.StatusBadRequest)
		rt.baseLogger.Debug("login name to long")
		return
	}

	usrId, err := rt.db.GetUsrIdByName(requestJson.UserName)

	//Controllo se l'utente esiste ed è presente nel database se non è presente lo creo e provo ad inserirlo nel database
	//in caso di riuscita preparo la rosposta http e la invio
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			rt.baseLogger.WithField("db error", err).Debug("utente non presente nel database")

			rt.baseLogger.Debug(fmt.Sprintf("Creazione e aggiunta del nuovo utente '%s' al database", requestJson.UserName))
			user, err:= rt.db.InsertNewUser(requestJson.UserName)

			if err != nil {
				rt.baseLogger.WithError(err).Error("Impossibile aggiungere nuovo user al database")
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}

			rt.sendJsonResponse(writer, user.UsrId)
			return
		}else{
			rt.baseLogger.WithError(err).Error("Errore durante il recupero dei dati")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	rt.sendJsonResponse(writer, usrId)
	return
}

func (rt *_router) sendJsonResponse(writer http.ResponseWriter, usrId string) {

	//Creo la risposta http contentente il token di autorizzazione e l'usrId (in questo caso entrambi sono la stessa cosa
	response := map[string]string{
		"usrId": usrId,
	}

	//Scrivo nell'header della risposta il token bearer che in questo caso corrisponde all'usrId
	writer.Header().Set("Authorization", "Bearer "+usrId)

	//Scrivo nella risposta l'usrId
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}