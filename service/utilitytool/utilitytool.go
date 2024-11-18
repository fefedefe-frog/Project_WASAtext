package utilitytool

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

func UserNameIsValid (userName string) (bool, error){

	//Controllo che rispetti i regex richiesti
	if !regexp.MustCompile(`^\S.*\S$`).MatchString(userName){
		return false, errors.New("INVALID REGEX")
	}
	if utf8.RuneCountInString(userName)<3{	//Deve essere lungo almeno 3 caratteri
		return false, errors.New("NAME SHORT")
	}
	if utf8.RuneCountInString(userName)>16{	//Deve essere lungo massimo 16 caratteri
		return false, errors.New("NAME LONG")
	}

	return true, nil
}

func UsrIdIsValid (usrId string) (bool, error){

	//Controllo che rispetti i regex richiesti
	if !regexp.MustCompile(`\s+`).MatchString(usrId){
		return false, errors.New("INVALID ID REGEX")
	}
	if utf8.RuneCountInString(usrId)<3{	//Deve essere lungo almeno 3 caratteri
		return false, errors.New("ID SHORT")
	}
	if utf8.RuneCountInString(usrId)>32{	//Deve essere lungo massimo 32 caratteri
		return false, errors.New("ID LONG")
	}

	return true, nil
}