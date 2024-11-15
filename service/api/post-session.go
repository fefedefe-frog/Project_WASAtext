package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type requestJson struct {
	UserName string `json:"userName"`
}

func (rt *_router) postSession(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {

	var session requestJson

	//Decodifica il corpo della richiesta JSON
	err := json.NewDecoder(request.Body).Decode(&session)
	rt.baseLogger.Info(fmt.Sprintf("Richiesta dall'endpoint /session arrivata utente: %s", session.UserName))

	if err != nil {
		http.Error(writer, "Invalid JSON format" , http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	usrId, err := rt.db.GetUsrIdByName(session.UserName)

	//Controllo se l'utente esiste ed è presente nel database se non è presente lo creo e provo ad inserirlo nel database
	//in caso di riuscita preparo la rosposta http e la invio
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			rt.baseLogger.WithField("db error", err).Debug("utente non presente nel database")

			rt.baseLogger.Debug(fmt.Sprintf("Creazione e aggiunta del nuovo utente '%s' al database", session.UserName))
			user, err:= rt.db.InsertNewUser(session.UserName)

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

	//TODO check if is ok
}

func (rt *_router) sendJsonResponse(writer http.ResponseWriter, usrId string) {

	//Creo la risposta http contentente il token di autorizzazione e l'usrId (in questo caso entrambi sono la stessa cosa
	response := map[string]string{
		"token": usrId,
		"usrId": usrId,
	}
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}