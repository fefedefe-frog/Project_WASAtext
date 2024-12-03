package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) startNewChat(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	requestJson= struct {
		ChatName string `json:"chatName"`
		ChatPhoto string `json:"chatPhoto"`
		IsGroup bool `json:"isGroup"`
		
	}{}
	rt.db.InsertNewChat()
}

func (rt *_router) getUserChats(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'endpoint /chats")

	//Tento di recuperare le chat di quell'user
	chats, err := rt.db.GetUserChats(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "NO_USER_CHATS" {
			rt.baseLogger.Info("User doesn't have chat")
			http.Error(writer, "The user doesn't have any chat", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Errorf("Error getting user %s chats", token)
		http.Error(writer, "Unable to retrive info", http.StatusInternalServerError)
		return
	}

	//Preparo la risposta contentente tutte le info delle chats
	responseChatsJSON, marshalErr := json.Marshal(map[string]interface{}{"chats": chats})
	if marshalErr != nil {
		rt.baseLogger.WithError(marshalErr).Errorf("Failed to marshal the chats")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	//Scrivo la risposta json
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseChatsJSON); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

func (rt *_router) getChatInfo(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'endpoint /chat/{chat_id}")

	//Recupero il valore di chat_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			if errors.Is(numErr.Err, strconv.ErrSyntax) {
				rt.baseLogger.WithError(numErr.Err).Error("the param chat_id is not a valid number")
			} else if errors.Is(numErr.Err, strconv.ErrRange) {
				rt.baseLogger.WithError(numErr.Err).Error("the param chat_id range is out of range")
			} else {
				rt.baseLogger.WithError(err).Error("Error parsing param chat_id")
			}
			http.Error(writer, "invalid param chat_id", http.StatusBadRequest)
			return
		}
	}

	//Controllo se l'utente che ha effettuato l'accesso faccia parte della chat di cui vuole ricavare le informazioni
	if exist, err := rt.db.CheckIfUserIsParticipant(chatId, token); err != nil {
		rt.baseLogger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	} else if !exist {
		context.Logger.WithField("usrId", token).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user logged in is not a participant of the chat", http.StatusForbidden)
		return
	}

	chat, dbErr := rt.db.GetChatInfo(chatId)
	if dbErr != nil {
		if errors.Is(dbErr, sql.ErrNoRows) {
			rt.baseLogger.WithError(dbErr).Warn("The chat doesn't exist, in the db")
			http.Error(writer, "Not found - chat not exist", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(dbErr).Error("Error getting chat info")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}

	//Preparo la risposta
	responseChatJSON, marshalErr := json.Marshal(chat)
	if marshalErr != nil {
		rt.baseLogger.WithError(marshalErr).Errorf("Failed to marshal the chat")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")

	if _, err := writer.Write(responseChatJSON); err != nil {
		rt.baseLogger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

func (rt *_router) getChatMessages(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("Richiesta all'endpoint /chat/{chat_id}/messages")

	//Recupero il valore di chat_id dai parametri dell'enpoint e controllo che sia un numero valido
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			if errors.Is(numErr.Err, strconv.ErrSyntax) {
				rt.baseLogger.WithError(numErr.Err).Error("the param chat_id is not a valid number")
			} else if errors.Is(numErr.Err, strconv.ErrRange) {
				rt.baseLogger.WithError(numErr.Err).Error("the param chat_id range is out of range")
			} else {
				rt.baseLogger.WithError(err).Error("Error parsing param chat_id")
			}
			http.Error(writer, "invalid param chat_id", http.StatusBadRequest)
			return
		}
	}

	//Controllo se l'utente che ha effettuato l'accesso faccia parte della chat di cui vuole ricavare i messaggi
	if exist, err := rt.db.CheckIfUserIsParticipant(chatId, token); err != nil {
		rt.baseLogger.WithError(err).Error("Error checking if user is participant")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	} else if !exist {
		context.Logger.WithField("usrId", token).Warn("User is not a participant")
		http.Error(writer, "Forbidden - The user logged in is not a participant of the chat", http.StatusForbidden)
		return
	}

	//Recupero i messaggi della chat dal database
	chatMessages, err := rt.db.GetChatMessages(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Warn("The chat have no messages, shouldn't be possible")
			http.Error(writer, "The chat is empty", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error getting chat messages")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Se tutti i controlli vanno a buon fine procedo a preparare la risposta e inviarla
	responseMessagesJSON, marshalErr := json.Marshal(map[string]interface{}{"messages": chatMessages})
	if marshalErr != nil {
		rt.baseLogger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	//Scrivo la risposta
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseMessagesJSON); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}
