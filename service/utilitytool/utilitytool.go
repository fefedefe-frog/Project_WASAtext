package utilitytool

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrInvalidRegex = errors.New("INVALID REGEX")
	ErrNameShort    = errors.New("NAME SHORT")
	ErrNameLong     = errors.New("NAME LONG")
	ErrIdShort      = errors.New("ID SHORT")
	ErrIdLong       = errors.New("ID LONG")
)

func NameIsValid(userName string) error {

	// Controllo che rispetti i regex richiesti
	if !regexp.MustCompile(`^\S.*\S$`).MatchString(userName) {
		return ErrInvalidRegex
	}
	if utf8.RuneCountInString(userName) < 3 { //Deve essere lungo almeno 3 caratteri
		return ErrNameShort
	}
	if utf8.RuneCountInString(userName) > 16 { // Deve essere lungo massimo 16 caratteri
		return ErrNameLong
	}

	return nil
}

func UsrIdIsValid(usrId string) error {

	// Controllo che rispetti i regex richiesti
	if !regexp.MustCompile(`\s*`).MatchString(usrId) {
		return ErrInvalidRegex
	}
	if utf8.RuneCountInString(usrId) < 3 { // Deve essere lungo almeno 3 caratteri
		return ErrIdShort
	}
	if utf8.RuneCountInString(usrId) > 32 { // Deve essere lungo massimo 32 caratteri
		return ErrIdLong
	}

	return nil
}
