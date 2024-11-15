package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func (db *appdbimpl) InsertNewUser(userName string) (User, error){
	var user User
	user.UserName= userName

	//Creo l'id dell'utente, rimescolando i caratteri del suo nome e aggiungendo delle cifre
	runes := []rune(userName)		//converto la stringa in array di caratteri

	//Mescola l'array di rune usando lo shuffle Fisher-Yates
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(runes) - 1; i > 0; i-- {
		j := r.Intn(i + 1) 						// Genera un indice casuale tra 0 e i
		runes[i], runes[j] = runes[j], runes[i] // Scambia i due elementi
	}

	//Converto l'array di rune di nuovo in una stringa, aggiungo due numeri randomici e unisco le due stringhe
	user.UsrId= fmt.Sprintf("%s%d%d", string(runes), r.Intn(100), r.Intn(10))
	user.UserPhoto= ""

	//Eseguo l'inserimento nel database
	_, err := db.c.Exec("INSERT INTO users_table (usrId, userName, userPhoto) VALUES (?, ?, ?)",user.UsrId,user.UserName,user.UserPhoto)
	return user, err
}

func (db *appdbimpl) GetUsrIdByName(userName string) (string, error){
	var usrId string

	err := db.c.QueryRow("SELECT usrId FROM users_table WHERE userName = ?", userName).Scan(&usrId)

	return usrId, err
}

func (db *appdbimpl) SetUserName(usrId string, newName string) error{

	stmt, err := db.c.Prepare("UPDATE users_table SET userName = ? WHERE usrId=?")
	if err != nil {
		log.Fatal("errore nella preparazione della query:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newName, usrId)
	if err != nil {
		log.Fatal("errore nell'esecuzione della  query:", err)
	}
	return err
}

func (db *appdbimpl) SetUserPhoto(usrId string, newPhoto string) error{

	stmt, err := db.c.Prepare("UPDATE users_table SET userPhoto = ? WHERE usrId=?")
	if err != nil {
		log.Fatal("errore nella preparazione della query:", err)	//usare logger giusto
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPhoto, usrId)
	if err != nil {
		log.Fatal("errore nell'esecuzione della  query:", err)
	}
	return err
}

func (db *appdbimpl) GetUserInfo(usrId string) (User, error){
	var user User

	err := db.c.QueryRow("SELECT userName, userPhoto FROM users_table WHERE usrId=?", usrId).Scan(&user.UserName, &user.UserPhoto)
	user.UsrId = usrId
	return user, err
}

func (db *appdbimpl) GetUsers() ([]User, error){
	var users []User

	rows, err := db.c.Query("SELECT usrId, userName, userPhoto FROM users_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Itero su tutte le righe della tabella degli user
	for rows.Next() {
		var user User

		err := rows.Scan(&user.UsrId, &user.UserName, &user.UserPhoto)
		if err != nil {
			return nil, err
		}

		//Aggiungo l'utente all'array
		users = append(users, user)
	}

	return users, err
}