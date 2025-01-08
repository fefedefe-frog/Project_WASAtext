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

func (rt *_router) getUserInfo(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, _ string) {

	// Recupero l'user id dell'user interessato
	usrId := params.ByName("usr_id")
	if err := utilitytool.UsrIdIsValid(usrId); err != nil {
		context.Logger.Infof("Invalid usrId : %v", err)
		http.Error(writer, "Invalid usrId", http.StatusBadRequest)
		return
	}

	// Tento di recuperare le informazioni di quell'user dal database
	user, err := rt.db.GetUserInfo(usrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.Infof("User id  not found: %s", usrId)
			http.NotFound(writer, request)
			return
		}
		context.Logger.WithError(err).Errorf("Error getting %s info", usrId)
		http.Error(writer, "Unable to retrive info", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	var responseUserJSON []byte
	responseUserJSON, err = json.Marshal(user)
	if err != nil {
		context.Logger.WithError(err).Errorf("Failed to marshal the user")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseUserJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
