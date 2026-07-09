package utils

import (
	"errors"
	"log"
	"net/http"
)

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

func NewUserError(message string) error {
	return &UserError{Message: message}
}

func HandleError(w http.ResponseWriter, err error, defaultMsg string, defaultStatus int) {
	if err == nil {
		return
	}

	log.Printf("Internal Error: %v", err)

	var userErr *UserError
	if errors.As(err, &userErr) {
		http.Error(w, userErr.Message, http.StatusConflict)
		return
	}
	http.Error(w, defaultMsg, defaultStatus)
}
