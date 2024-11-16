package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)


func (rt *_router) BearerAuth(h httpRouterHandler) httpRouterHandler{

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		//Estraggo l'header "Authentication" dalla richiesta
		authHeader := r.Header.Get("Authorization")

		//Controllo che l'header sia presente e che sia del tipo Bearer
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.WithField("authentication_header:", authHeader).Warn("Missing or invalid bearer Authentication header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//Estraggo il toked dall header(rimuovo "Bearer " dall'header)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		//Verifico se il token corrisponde ad un usrId di un utente già registrato
		exist, err := rt.db.UsrIdExist(token)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error checking user id existence")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !exist {
			//Se l'usrId non esiste e quindi l'utente non esiste e il token presente nel header non è valido restituisco una richiesta Unauthorized
			ctx.Logger.WithField("token", token).Warn("User id token not present in the db")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}


		//If the token is valid call the next handler in chain (usually, the handler function for the path)
		h(w, r, ps, ctx)
	}
}