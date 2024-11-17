package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)


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

	//Imposto il tipo di contenuto come JSON
	w.Header().Set("content-type", "application/json")

	//Codifico i dati contenuti in users in formato JSON e li scrivo nella risposta HTTP
	if err := json.NewEncoder(w).Encode(users); err != nil{
		http.Error(w, "Errore codifica JSON", http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("Errore nella codifica JSON")
		return
	}
}
