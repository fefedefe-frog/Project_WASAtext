package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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