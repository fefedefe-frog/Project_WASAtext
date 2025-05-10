package api

import (
	"Project_WASAtext/service/api/reqcontext"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func (rt *_router) setUserPhoto(writer http.ResponseWriter, request *http.Request, _ httprouter.Params, context reqcontext.RequestContext, usrId string) {

	// Limito la memoria per il parsing del form
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		context.Logger.WithError(err).Error("Error parsing multipart form")
		http.Error(writer, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Recupero il valore della nuova foto
	photoFile, _, err := request.FormFile("newUserPhoto")
	if err != nil {
		context.Logger.WithError(err).Error("Error getting file content")
		http.Error(writer, "Bad request - Error getting file content", http.StatusBadRequest)
		return
	}
	defer func() {
		err := photoFile.Close()
		if err != nil {
			context.Logger.WithError(err).Error("Error closing file")
		}
	}()

	// Leggo il contenuto dell'immagine
	var newUserPhoto []byte
	newUserPhoto, err = io.ReadAll(photoFile)
	if err != nil {
		context.Logger.WithError(err).Error("Error reading file")
		http.Error(writer, "Internal server error - Error reading file", http.StatusInternalServerError)
		return
	}

	// Aggiorno la propic nel database
	if err := rt.db.SetUserPhoto(usrId, newUserPhoto); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Logger.WithError(err).Errorf("User: %s not found in database", usrId)
			http.Error(writer, "Not found - User not found", http.StatusNotFound)
			return
		}
		context.Logger.WithError(err).Error("Error changing user propic")
		http.Error(writer, "Bad request - Unable to change propic", http.StatusBadRequest)
		return
	}

	// Preparo la risposta
	responseJSON, marshalErr := json.Marshal(map[string]interface{}{"userPhoto": base64.StdEncoding.EncodeToString(newUserPhoto)})
	if marshalErr != nil {
		context.Logger.WithError(marshalErr).Errorf("Failed to marshal the new user photo")
		http.Error(writer, "Internal server error - failed json conversion", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write(responseJSON); err != nil {
		context.Logger.WithError(err).Error("Errore preaparazione risposta html")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
