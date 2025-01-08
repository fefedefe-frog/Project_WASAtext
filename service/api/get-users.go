package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUsers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, _ string) {

	users, err := rt.db.GetUsers()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Error("no users in database")
			http.Error(w, "No users in the database", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("database error")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Users []database.User `json:"users"`
	}{
		Users: users,
	}

	// Imposto il tipo di contenuto come JSON
	w.Header().Set("content-type", "application/json")

	// Codifico i dati contenuti in users in formato JSON e li scrivo nella risposta HTTP
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Errore codifica JSON", http.StatusInternalServerError)
		context.Logger.WithError(err).Error("Errore nella codifica JSON")
		return
	}
}
