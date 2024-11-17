package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (rt *_router) patchChangeUserPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	context.Logger.Info("Richiesta all'enpoint /users/{usr_id}/propic")

	var requestJson= struct{
		Operation string `json:"op"`;
		Path string `json:"path"`;
		NewUserPhoto string `json:"value"`
	}{}

	//Controllo che la richiesta arrivata dall'utente corrisponda alla modifica della sua propic, e non della propic di altri utenti
	usrToken := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	if usrToken != params.ByName("usr_id"){
		context.Logger.Warnf("user: %s tried to change propic of users: %s", usrToken, params.ByName("usr_id"))
		http.Error(writer, "Unauthorized - can't change the propic of another user", http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil{
		http.Error(writer, "Invalid JSON format" , http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("Richiesta di cambio propic da parte dell'user: %s || nuova propic: %s...", usrToken, requestJson.NewUserPhoto[:10])
	if err := rt.db.SetUserPhoto(usrToken, requestJson.NewUserPhoto); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			rt.baseLogger.WithError(err).Errorf("User: %s not found in database", usrToken)
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing user propic")
		http.Error(writer, "Unable to change propic", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"userName": requestJson.NewUserPhoto,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}