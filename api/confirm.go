package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Confirm(apiRouter *chi.Mux) {
	apiRouter.Route("/confirm", func(r chi.Router) {
		r.Post("/receive", app.ConfirmReceive)
		r.Post("/send", app.ConfirmSend)
	})
}

// @Summary Confirm receive
// @Description Confirm receive
// @ID confirm-receive
// @Tags Confirm
// @Accept json
// @Produce json
// @Param searchWord query string true "input TTFID or trackingNo"
// @Success 200 {object} Response{result=string} "confirm receive result"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /confirm/receive [post]
func (app *Application) ConfirmReceive(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	searchWord := queryValues.Get("searchWord")
	res, err := app.Service.Confirm.ConfirmReceive(searchWord)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Confirm send
// @Description Confirm send
// @ID confirm-send
// @Tags Confirm
// @Accept json
// @Produce json
// @Param searchWord query string true "input TTFID or trackingNo"
// @Success 200 {object} Response{result=string} "confirm send result"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /confirm/send [post]
func (app *Application) ConfirmSend(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	searchWord := queryValues.Get("searchWord")
	res, err := app.Service.Confirm.ConfirmSend(searchWord)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}
