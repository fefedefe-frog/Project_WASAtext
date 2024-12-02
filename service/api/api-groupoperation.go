package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"Project_WASAtext/service/utilitytool"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(writer http.ResponseWriter, _ *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Infof("DELETE request to endpoint /chat/{chat_id}")

	//Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant{
		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithField("usrId", token).Warnf("tried leave the group <%d> which isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't leave a group which are't member of", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}


}

func (rt *_router) changeGroupName(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("PATCH request to endpoint /chat/{chat_id}")

	var requestJson= struct {
		NewGroupName string `json:"newGroupName"`
	}{}

	//Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant{

		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithField("usrId", token).Warnf("tried to change group name of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the name of another group", http.StatusForbidden)
			return
		}else if err != nil {
			context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
			http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
			return
		}
	}

	//Controllo che la chat che si vuole modificare sia un gruppo
	if isGroup, err := rt.db.IsAGroup(chatId); !isGroup{

		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithError(err).Errorf("Chat <%d> not found", chatId)
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}else if err != nil{
			context.Logger.WithError(err).Errorf("Unable to check if the chat <%d> is a group", chatId)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		context.Logger.WithField("chatId", chatId).Warn("User tried to change a name of a chat that isn't a group")
		http.Error(writer, "Forbidden - can't change the name of a non group chat", http.StatusForbidden)
		return
	}

	//Decodifica del JSON nella richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("user <%s> request to change group name of group <%d>", token, chatId)
	//Controllo se il nome è valido secondo i requisiti richiesti
	if err := utilitytool.NameIsValid(requestJson.NewGroupName); err != nil {
		switch {
		case errors.Is(err, utilitytool.ErrInvalidRegex):
			http.Error(writer, "Invalid name format, the name can't contain space at the start or end of the name", http.StatusBadRequest)
			rt.baseLogger.Debug("Invalid name format")

		case errors.Is(err, utilitytool.ErrNameShort):
			http.Error(writer, "the name must be at least 3 character long", http.StatusBadRequest)
			rt.baseLogger.Debug("name to short")

		case errors.Is(err, utilitytool.ErrNameLong):
			http.Error(writer, "the name can be max 16 character long", http.StatusBadRequest)
			rt.baseLogger.Debug("name to long")
		}
		return
	}

	//Aggiorno l'username nel database
	if err := rt.db.SetGroupName(chatId, requestJson.NewGroupName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Errorf("Group chat: <%s> not found in database", token)
			http.Error(writer, "Chat not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing group name")
		http.Error(writer, "Unable to rename", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	response := map[string]string{
		"chatName": requestJson.NewGroupName,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal Server Error - ", http.StatusInternalServerError)
		return
	}
	return
}

func (rt *_router) changeGroupPhoto(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string) {
	context.Logger.Info("PATCH request to endpoint /chat/{chat_id}/propic")

	var requestJson = struct {
		NewGroupPhoto string `json:"newGroupPhoto"`
	}{}

	//Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant{
		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithField("usrId", token).Warnf("tried to change group photo of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the photo of another group", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	//Controllo che la chat che si vuole modificare sia un gruppo
	if isGroup, err := rt.db.IsAGroup(chatId); !isGroup{

		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithError(err).Errorf("Chat <%d> not found", chatId)
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}else if err != nil{
			context.Logger.WithError(err).Errorf("Unable to check if the chat <%d> is a group", chatId)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		context.Logger.WithField("chatId", chatId).Warn("User tried to change a photo of a chat that isn't a group")
		http.Error(writer, "Forbidden - can't change the photo of a non group chat", http.StatusForbidden)
		return
	}

	//Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}

	context.Logger.Infof("user <%s> request to change group photo of group <%d>", token, chatId)
	//Aggiorno la propic del gruppo nel database
	if err := rt.db.SetGroupPhoto(chatId, requestJson.NewGroupPhoto); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Errorf("Group chat <%d> not found in database", chatId)
			http.Error(writer, "Group chat not found", http.StatusNotFound)
			return
		}
		rt.baseLogger.WithError(err).Error("Error changing group propic")
		http.Error(writer, "Unable to change group propic", http.StatusBadRequest)
		return
	}

	//Preparo la risposta
	response := map[string]string{
		"chatPhoto": requestJson.NewGroupPhoto,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

func (rt *_router) postAddUserToGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params, context reqcontext.RequestContext, token string){
	context.Logger.Info("POST request to endpoint /chat/{chat_id}/users")

	var requestJson = struct {
		UsrIdToAdd string `json:"usrIdToAdd"`
	}{}

	//Controllo che l'utente faccia effettivamente parte del gruppo
	chatId, err := strconv.Atoi(params.ByName("chat_id"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("invalid chat id")
		http.Error(writer, "invalid chat_id parameter", http.StatusBadRequest)
	}
	if isParticipant, err := rt.db.CheckIfUserIsParticipant(chatId, token); !isParticipant{
		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithField("usrId", token).Warnf("tried to change group photo of group <%d> which he isn't a member of", chatId)
			http.Error(writer, "Forbidden - can't change the photo of another group", http.StatusForbidden)
			return
		}
		context.Logger.WithError(err).Errorf("Error while checking if the user <%s> is member of the group <%d>", token, chatId)
		http.Error(writer, "Internal Server Error - can't check user paricipation of the group", http.StatusInternalServerError)
		return
	}

	//Decodifico la richiesta
	if err := json.NewDecoder(request.Body).Decode(&requestJson); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		rt.baseLogger.WithError(err).Error("Invalid JSON in requestBody")
		return
	}


	context.Logger.Infof("user <%s> request to add user <%s> to the group", token, requestJson.UsrIdToAdd)
	//Aggiungo l'utente alla chat
	err= rt.db.InsertUserInChat(requestJson.UsrIdToAdd, chatId)
	if err != nil {
		context.Logger.WithError(err).Error("Error adding user to group")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Debug("Successfully added user to group")

	//Preparo la risposta contentente la lista aggiornata di user id
	var chatParticipants []string
	chatParticipants, err = rt.db.GetChatPartecipants(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			context.Logger.WithField("chatId", chatId).Error("chat not found in the chat_participants_table")
			http.Error(writer, "Chat not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error getting chat participants")
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"users": chatParticipants})
	if marshalErr != nil{
		rt.baseLogger.WithError(marshalErr).Errorf("Failed to marshal the chat messages")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(responseJSON); err != nil {
		rt.baseLogger.WithError(err).Error("Json encoding error")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}