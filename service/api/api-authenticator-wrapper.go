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

func (rt *_router) BearerAuth(fn httpRouterHandlerAuthenticated) httpRouterHandler {

	return func(writer http.ResponseWriter, request *http.Request, parameters httprouter.Params, context reqcontext.RequestContext) {

		// Estraggo l'header "Authentication" dalla richiesta
		authHeader := request.Header.Get("Authorization")

		// Controllo che l'header sia presente e che sia del tipo Bearer
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			context.Logger.WithField("authentication_header:", authHeader).Warn("Missing or invalid bearer Authentication header")
			http.Error(writer, "Unauthorized - missing token", http.StatusUnauthorized)
			return
		}

		// Estraggo il toked dall header(rimuovo "Bearer " dall'header)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if err := utilitytool.UsrIdIsValid(token); err != nil {
			context.Logger.WithError(err).Warnf("Invalid token: <%s>", token)
			http.Error(writer, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		// Verifico se il token corrisponde ad un usrId di un utente già registrato
		exist, err := rt.db.UsrIdExist(token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				context.Logger.WithError(err).Warnf("Token: <%s> does not exist", token)
				http.Error(writer, "Unauthorized - token not valid or deprecated", http.StatusInternalServerError)
				return
			}
			context.Logger.WithError(err).Warn("Error checking token existence")
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !exist {
			// Se l'usrId non esiste e quindi l'utente non esiste e il token presente nel header non è valido restituisco una richiesta Unauthorized
			context.Logger.WithField("token", token).Warn("User id token not present in the db")
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid call the next handler in chain (usually, the handler function for the path)
		fn(writer, request, parameters, context, token)
	}
}
