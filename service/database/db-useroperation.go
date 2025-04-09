package database

import (
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

func (db *appdbimpl) InsertNewUser(userName string) (User, error) {
	var user User
	user.UserName = userName

	// Creo l'id dell'utente, rimescolando i caratteri del suo nome e aggiungendo delle cifre
	runes := []rune(strings.ToLower(userName)) // converto la stringa in array di caratteri

	// Mescola l'array di rune usando lo shuffle Fisher-Yates
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(runes) - 1; i > 0; i-- {
		j := r.Intn(i + 1)                      // Genera un indice casuale tra 0 e i
		runes[i], runes[j] = runes[j], runes[i] // Scambia i due elementi
	}

	// Converto l'array di rune di nuovo in una stringa, rimuovo eventuali spazi, aggiungo due numeri randomici e unisco le stringhe
	user.UsrId = fmt.Sprintf("%s%d%d", strings.ReplaceAll(string(runes), " ", ""), r.Intn(100), r.Intn(10))

	// Decodifica la stringa Base64 in byte
	defaultPropicBytes, errProp := base64.StdEncoding.DecodeString(defaultPropicBase64)
	if errProp != nil {
		return user, errProp
	}

	// Eseguo l'inserimento nel database
	_, err := db.c.Exec(`INSERT INTO users_table (usrId, userName, userPhoto) VALUES (?, ?, ?);`, user.UsrId, user.UserName, defaultPropicBytes)
	return user, err
}

func (db *appdbimpl) GetUsrIdByName(userName string) (string, error) {
	var usrId string

	err := db.c.QueryRow(`SELECT usrId FROM users_table WHERE userName = ?;`, userName).Scan(&usrId)
	if err != nil {
		return "", err
	}
	return usrId, nil
}

func (db *appdbimpl) SetUserName(usrId string, newName string) error {

	stmt, err := db.c.Prepare(`UPDATE users_table SET userName = ? WHERE usrId=?;`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(newName, usrId)
	if err != nil {
		return err
	}
	return err
}

func (db *appdbimpl) SetUserPhoto(usrId string, newPhotoData []byte) error {

	stmt, err := db.c.Prepare(`UPDATE users_table SET userPhoto = ? WHERE usrId=?;`)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("stmt.Close() error: %v", closeErr)
			}
		}
	}()

	_, err = stmt.Exec(newPhotoData, usrId)
	if err != nil {
		return err
	}
	return err
}

func (db *appdbimpl) GetUserInfo(usrId string) (User, error) {
	var user User
	var propicByte []byte

	err := db.c.QueryRow(`SELECT userName, userPhoto FROM users_table WHERE usrId=?;`, usrId).Scan(&user.UserName, &propicByte)
	if err != nil {
		return user, err
	}
	user.UsrId = usrId

	// Controllo se sia presente la foto nel database
	if len(propicByte) > 0 {
		user.UserPhoto = base64.StdEncoding.EncodeToString(propicByte)
	} else {
		user.UserPhoto = "" // se non è presente assegno la stringa vuota
	}

	return user, err
}

func (db *appdbimpl) GetUserNameById(usrId string) (string, error) {
	var username string

	err := db.c.QueryRow(`SELECT userName FROM users_table WHERE usrId=?;`, usrId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, err
}

func (db *appdbimpl) GetUsers(usrIdToIgnore string) ([]User, error) {

	rows, err := db.c.Query(`SELECT usrId, userName, userPhoto FROM users_table WHERE usrId != ?;`, usrIdToIgnore)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				logrus.WithError(closeErr).Errorf("rows.Close() error: %v", closeErr)
			}
		}
	}()

	// Itero su tutte le righe della tabella degli user
	var users []User
	for rows.Next() {
		var user User

		var propicBytes []byte
		if err := rows.Scan(&user.UsrId, &user.UserName, &propicBytes); err != nil {
			return nil, err
		}
		user.UserPhoto = base64.StdEncoding.EncodeToString(propicBytes)

		// Aggiungo l'utente all'array
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, err
}

func (db *appdbimpl) UsrIdExist(usrId string) (bool, error) {

	var exist int
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users_table WHERE usrId=?);`, usrId).Scan(&exist)
	if err != nil {
		// Non c'è bisogno di gestire `sql.ErrNoRows` esplicitamente,
		// in quanto QueryRow restituisce `false` per la variabile esistenza
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return exist == 1, nil
}
