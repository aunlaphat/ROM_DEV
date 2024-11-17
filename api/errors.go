package api

import (
	"boilerplate-backend-go/errors"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case errors.AppError:
		http.Error(w, e.Message, e.Code)
	case error:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
