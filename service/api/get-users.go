package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	users, err := rt.db.GetUsers()
	if err != nil {
		rt.baseLogger.Error("database error: ", err.Error())
	}

	//Imposto il tipo di contenuto come JSON
	w.Header().Set("content-type", "application/json")

	//Codifico i dati contenuti in users in formato JSON e li scrivo nella risposta HTTP
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Errore codifica JSON", http.StatusInternalServerError)
		rt.baseLogger.Error("Errore nella codifica JSON", err.Error())
	}
}
