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
)

// List all the user of the app
func (rt *_router) getUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, _ string) {
	context.Logger.Info("Richiesta all'endpoint /users")

	users, err := rt.db.GetUsers()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("no users in database")
			http.Error(w, "No users in the database", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("database error")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Users []database.User `json:"users"`
	}{
		Users: users,
	}

	//Imposto il tipo di contenuto come JSON
	w.Header().Set("content-type", "application/json")

	//Codifico i dati contenuti in users in formato JSON e li scrivo nella risposta HTTP
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Errore codifica JSON", http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("Errore nella codifica JSON")
		return
	}
	return
}

// Retrive the info (name and propic) of a user, by its user id
func (rt *_router) getUserInfo(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, _ string) {
	context.Logger.Info("richiesta all'endpoint /users")

	//Recupero l'user id dell'user interessato
	usrId := params.ByName("usr_id")
	if err := utilitytool.UsrIdIsValid(usrId); err != nil{
		rt.baseLogger.Infof("Invalid usrId : %v", err)
		http.Error(writer, "Invalid usrId", http.StatusBadRequest)
		return
	}

	//Tento di recuperare le informazioni di quell'user dal database
	user, err := rt.db.GetUserInfo(usrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.Infof("User id  not found: %s", usrId)
			http.NotFound(writer, request)
			return
		}
		rt.baseLogger.WithError(err).Errorf("Error getting %s info", usrId)
		http.Error(writer, "Unable to retrive info", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	responseUserJSON, err := json.Marshal(user)
	if err != nil{
		rt.baseLogger.WithError(err).Errorf("Failed to marshal the user")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(responseUserJSON);err != nil {
		rt.baseLogger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

// Change username
func (rt *_router) patchChangeUserName(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'enpoint /users/{usr_id}")

	var requestJson = struct {
		NewUserName string `json:"newUserName"`
	}{}

	//Controllo che la richiesta arrivata dall'utente corrisponda alla modifica del suo nome, e non del nome di altri utenti
	if token != params.ByName("usr_id") {
		context.Logger.WithField("usrId", token).Warnf("tried to change username of users: %s", params.ByName("usr_id"))
		http.Error(writer, "Forbidden - can't change the name of another user", http.StatusForbidden)
		return
	}

	//Decodifica del JSON nella richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("Richiesta di cambio nome da parte dell'user: %s || nuovo nome: %s", token, requestJson.NewUserName)

	//Controllo se il nome è valido secondo i requisiti richiesti
	if err := utilitytool.NameIsValid(requestJson.NewUserName); err != nil {
		switch {
		case errors.Is(err, utilitytool.ErrInvalidRegex):
			http.Error(writer, "Invalid name format, the name can't contain space at the start or end of the name", http.StatusBadRequest)
			rt.baseLogger.Debug("Invalid name format")

		case errors.Is(err, utilitytool.ErrNameShort):
			http.Error(writer, "the name must be at least 3 character long", http.StatusBadRequest)
			rt.baseLogger.Debug("login name to short")

		case errors.Is(err, utilitytool.ErrNameLong):
			http.Error(writer, "the name can be max 16 character long", http.StatusBadRequest)
			rt.baseLogger.Debug("login name to long")
		}
		return
	}

	//Aggiorno l'username nel database
	if err := rt.db.SetUserName(token, requestJson.NewUserName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Errorf("User: %s not found in database", token)
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing user name")
		http.Error(writer, "Unable to rename", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	response := map[string]string{
		"userName": requestJson.NewUserName,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

// Change user propic
func (rt *_router) patchChangeUserPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'enpoint /users/{usr_id}/propic")

	var requestJson = struct {
		NewUserPhoto string `json:"newUserPhoto"`
	}{}

	if err := utilitytool.UsrIdIsValid(params.ByName("usr_id")); err != nil{
		rt.baseLogger.Infof("Invalid usrId : %v", err)
		http.Error(writer, "Invalid usrId", http.StatusBadRequest)
		return
	}

	//Controllo che la richiesta arrivata dall'utente corrisponda alla modifica della sua propic, e non della propic di altri utenti
	if token != params.ByName("usr_id"){
		context.Logger.WithField("usrId", token).Warnf("tried to change propic of users: %s", params.ByName("usr_id"))
		http.Error(writer, "Forbidden - can't change the propic of another user", http.StatusForbidden)
		return
	}

	//Decodifico la richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("Richiesta di cambio propic da parte dell'user: %s || nuova propic: %s...", token, requestJson.NewUserPhoto[:10])
	//Aggiorno la propic nel database
	if err := rt.db.SetUserPhoto(token, requestJson.NewUserPhoto); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Errorf("User: %s not found in database", token)
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing user propic")
		http.Error(writer, "Unable to change propic", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	response := map[string]string{
		"userPhoto": requestJson.NewUserPhoto,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}