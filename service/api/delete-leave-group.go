package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithField("usrId", token).Warnf("tried leave the group <%d> which isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't leave a group which are't member of", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}
}
