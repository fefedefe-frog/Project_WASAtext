package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) putChangeUserPhoto(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("PUT request to endpoint /profile/propic")

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

	// Aggiorno la propic nel database
	context.Logger.Infof("Richiesta di cambio propic da parte dell'user: %s || nuova propic: %s...", token, requestJson.NewUserPhoto[:10])
	if err := rt.db.SetUserPhoto(token, requestJson.NewUserPhoto); err != nil {
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
	response := map[string]string{
		"userPhoto": requestJson.NewUserPhoto,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		context.Logger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}
