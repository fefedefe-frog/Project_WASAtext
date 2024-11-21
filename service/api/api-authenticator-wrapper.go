package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)



// httpRouterHandlerAuthenticated is the signature for functions that needs to be authenticated and need a token in
// addition to those required by the httprouter package.
type httpRouterHandlerAuthenticated func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext, string)



func (rt *_router) BearerAuth(fn httpRouterHandlerAuthenticated) httpRouterHandler{

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		//Estraggo l'header "Authentication" dalla richiesta
		authHeader := r.Header.Get("Authorization")

		//Controllo che l'header sia presente e che sia del tipo Bearer
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.WithField("authentication_header:", authHeader).Warn("Missing or invalid bearer Authentication header")
			http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
			return
		}

		//Estraggo il toked dall header(rimuovo "Bearer " dall'header)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if err := utilitytool.UsrIdIsValid(token); err != nil{
			ctx.Logger.WithError(err).Warn("Invalid token")
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}


		//Verifico se il token corrisponde ad un usrId di un utente già registrato
		exist, err := rt.db.UsrIdExist(token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.Logger.WithError(err).Warn("Token does not exist")
				http.Error(w, "Unauthorized - token not valid or deprecated", http.StatusInternalServerError)
				return
			}
			ctx.Logger.WithError(err).Warn("Error checking token existence")
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
		fn(w, r, ps, ctx, token)
	}
}