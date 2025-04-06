package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) uncommentMessage(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {

	// Recupero l'id della chat e del commento dai paramentri dell'endpoint
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
		return
	}
	var commentId int
	commentId, err = strconv.Atoi(params.ByName("comment_id"))
	if err != nil {
		context.Logger.WithError(err).Error("invalid comment id")
		http.Error(writer, "invalid comment_id parameter", http.StatusBadRequest)
		return
	}

	var isParticipant bool
	isParticipant, err = rt.db.CheckIfUserIsParticipant(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the user participant status")
		http.Error(writer, "Internal Server Error - Unable to check the user participant status", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		context.Logger.Warn("the user is not a participant of the chat of the message that was try to comment")
		http.Error(writer, "Forbidden - Can't uncomment a message of a chat where aren't a participant", http.StatusForbidden)
		return
	}

	// Controllo che l'utente che vuole rimuovere il commento sia colui che l'ha commentato in precedenza
	var isCommenter bool
	isCommenter, err = rt.db.CheckCommentAuthor(chatId, token)
	if err != nil {
		context.Logger.WithError(err).Error("Error while checking the comment author")
		http.Error(writer, "Internal Server Error - Unable to check the comment author", http.StatusInternalServerError)
		return
	}

	if !isCommenter {
		context.Logger.Warn("the user is not the author of the comment that was try to uncomment")
		http.Error(writer, "Forbidden - Can't uncomment a message that isn't yours", http.StatusForbidden)
		return
	}

	// Controllo che l'utente che vuole commentare il messaggio faccia parte della chat di quel messaggio

	if err := rt.db.UncommentMessage(commentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Error("The comment that is bein deleted doesn't exist")
			http.Error(writer, "Not Found - Comment not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error while uncommenting message to chat")
		http.Error(writer, "Internal Server Error - Unable to uncomment message to chat", http.StatusInternalServerError)
		return
	}

	context.Logger.Debug("message uncommented successfully")
	writer.WriteHeader(http.StatusNoContent)
}
