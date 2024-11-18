package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

//List all the user of the app
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, context reqcontext.RequestContext) {
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
	if err := json.NewEncoder(w).Encode(response); err != nil{
		http.Error(w, "Errore codifica JSON", http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("Errore nella codifica JSON")
		return
	}
}

//Retrive the info (name and propic) of an user, by its user id
func (rt *_router) getUserInfo(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	context.Logger.Info("richiesta all'endpoint /users")

	//Recupero l'user id dell'user interessato
	usrId := params.ByName("usr_id")

	//Tento di recuperare le informazioni di quell'user dal database
	user, err := rt.db.GetUserInfo(usrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.Infof("User id  not found: %s", usrId)
			http.NotFound(writer, request)
			return
		}
		rt.baseLogger.WithError(err).Errorf("Error getting user(%s) info", usrId)
		http.Error(writer, "Unable to retrive info", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	response := map[string]string{
		"usrId": user.UsrId,
		"userName": user.UserName,
		"userPhoto": user.UserPhoto,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

//Change username
func (rt *_router) patchChangeUserName(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	context.Logger.Info("Richiesta all'enpoint /users/{usr_id}")

	var requestJson= struct{
		NewUserName string `json:"newUserName"`
	}{}

	//Controllo che la richiesta arrivata dall'utente corrisponda alla modifica del suo nome, e non del nome di altri utenti
	usrToken := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	if usrToken != params.ByName("usr_id"){
		context.Logger.Warnf("user: %s tried to change username of users: %s", usrToken, params.ByName("usr_id"))
		http.Error(writer, "Unauthorized - can't change the name of another user", http.StatusUnauthorized)
		return
	}

	//Decodifica del JSON nella richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil{
		http.Error(writer, "Invalid JSON format" , http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}


	context.Logger.Infof("Richiesta di cambio nome da parte dell'user: %s || nuovo nome: %s", usrToken, requestJson.NewUserName)

	//Controllo se il nome è valido secondo i requisiti richiesti
	if _, err := utilitytool.UserNameIsValid(requestJson.NewUserName); err != nil{
		rt.baseLogger.WithField("not valid for", err).Info("nuovo nome non valido")
		http.Error(writer, fmt.Sprintf("Nuovo nome non valido: %s", err), http.StatusBadRequest)
		return
	}

	//Aggiorno l'username nel database
	if err := rt.db.SetUserName(usrToken, requestJson.NewUserName); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			rt.baseLogger.WithError(err).Errorf("User: %s not found in database", usrToken)
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
}

//Change user propic
func (rt *_router) patchChangeUserPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	context.Logger.Info("Richiesta all'enpoint /users/{usr_id}/propic")

	var requestJson= struct{
		NewUserPhoto string `json:"newUserPhoto"`
	}{}

	//Controllo che la richiesta arrivata dall'utente corrisponda alla modifica della sua propic, e non della propic di altri utenti
	usrToken := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	if usrToken != params.ByName("usr_id"){
		context.Logger.Warnf("user: %s tried to change propic of users: %s", usrToken, params.ByName("usr_id"))
		http.Error(writer, "Unauthorized - can't change the propic of another user", http.StatusUnauthorized)
		return
	}


	//Decodifico la richiesta
	err := json.NewDecoder(request.Body).Decode(&requestJson)
	if err != nil{
		http.Error(writer, "Invalid JSON format" , http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("Richiesta di cambio propic da parte dell'user: %s || nuova propic: %s...", usrToken, requestJson.NewUserPhoto[:10])
	//Aggiorno la propic nel database
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
}