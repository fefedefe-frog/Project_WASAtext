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

func (rt *_router) getUserInfo(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, _ string) {

	// Recupero l'user id dell'user interessato
	usrId := params.ByName("usr_id")
	if err := utilitytool.UsrIdIsValid(usrId); err != nil {
		context.Logger.WithError(err).WithField("usrId", usrId).Info("Invalid usrId")
		http.Error(writer, "Bad Request - Invalid usrId", http.StatusBadRequest)
		return
	}

	// Tento di recuperare le informazioni di quell'user dal database
	user, err := rt.db.GetUserInfo(usrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("usrId", usrId).Info("User id not found")
			http.Error(writer, "Not found - The user doesn't exist", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).WithField("usrId", usrId).Error("Error getting user info")
		http.Error(writer, "Bad Request - Unable to retrive info", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	var responseUserJSON []byte
	responseUserJSON, err = json.Marshal(user)
	if err != nil {
		context.Logger.WithError(err).Errorf("Failed to marshal the user")
		http.Error(writer, "Internal server error - Failed json conversion", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseUserJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta http")
		http.Error(writer, "Internal server error - Error while preparing the http response", http.StatusInternalServerError)
		return
	}
}
