package api

import (
	"boilerplate-backend-go/dto"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Application) TTF(apiRouter *chi.Mux) {
	apiRouter.Route("/ttf", func(r chi.Router) {
		r.Post("/create/draft", app.CreateDraft)
		r.Post("/create/req", app.CreateTTFReq)
		r.Post("/create/ret", app.CreateTTFRet)
		r.Post("/edit/draft", app.EditDraft)
		r.Get("/get-all", app.GetAll)
		r.Get("/get-all/by", app.GetAllByOpr)
		r.Get("/get-by", app.GetBy)
	})
}

// @Summary create draft
// @Description crete draft document ttf.
// @ID ttf-create-draft
// @Tags TTF
// @Accept json
// @Produce json
// @Param create-draft body dto.InputCreateDoc true "request body for create draft in JSON format"
// @Success 200 {object} Response{result=[]string} "Return from create draft"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/create/draft [post]
func (app *Application) CreateDraft(w http.ResponseWriter, r *http.Request) {
	req := dto.InputCreateDoc{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}
	ttfNumber, err := app.Service.TTF.CreateDraft(req, "draft")
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, ttfNumber, http.StatusOK)
}

// @Summary edit draft
// @Description edit draft document ttf.
// @ID ttf-edit-draft
// @Tags TTF
// @Accept json
// @Produce json
// @Param edit-draft body dto.InputEditDoc true "request body for edit draft in JSON format"
// @Success 200 {object} Response{result=[]string} "Return from create draft"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/edit/draft [post]
func (app *Application) EditDraft(w http.ResponseWriter, r *http.Request) {
	req := dto.InputEditDoc{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	ttfNumber, err := app.Service.TTF.EditDraft(req, "draft")
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, ttfNumber, http.StatusOK)
}

// @Summary create request ttf
// @Description create rquest document ttf.
// @ID ttf-create-request
// @Tags TTF
// @Accept json
// @Produce json
// @Param create-draft body dto.InputCreateDoc true "request body for create draft in JSON format"
// @Param create-draft body dto.InputCreateDoc true "request body for create draft in JSON format"
// @Success 200 {object} Response{result=[]string} "Return from create draft"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/create/req [post]
func (app *Application) CreateTTFReq(w http.ResponseWriter, r *http.Request) {
	req := dto.InputCreateDoc{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	ttfNumber, err := app.Service.TTF.CreateTTF(req, "req")
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, ttfNumber, http.StatusOK)
}

// @Summary create return ttf
// @Description create return document ttf.
// @ID ttf-create-return
// @Tags TTF
// @Accept json
// @Produce json
// @Param create-draft body dto.InputCreateDoc true "request body for create draft in JSON format"
// @Param create-draft body dto.InputCreateDoc true "request body for create draft in JSON format"
// @Success 200 {object} Response{result=[]string} "Return from create draft"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/create/ret [post]
func (app *Application) CreateTTFRet(w http.ResponseWriter, r *http.Request) {
	req := dto.InputCreateDoc{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}

	ttfNumber, err := app.Service.TTF.CreateTTF(req, "ret")
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, ttfNumber, http.StatusOK)
}

// @Summary Get all
// @Description Get all.
// @ID ttf-get-all
// @Tags TTF
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]dto.GetAllTTF} "Get all document have been created successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/get-all [get]
func (app *Application) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.TTF.GetAll()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get all by operation's status
// @Description Get all by operation's status.
// @ID ttf-get-all-by-opr
// @Tags TTF
// @Accept json
// @Produce json
// @Param status query string true "Operation's status"
// @Success 200 {object} Response{result=[]dto.GetAllTTF} "Get all document by operation's status have been created successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/get-all/by [get]
func (app *Application) GetAllByOpr(w http.ResponseWriter, r *http.Request) {
	// Step 1: Retrieve the "status" parameter from the query string
	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		http.Error(w, "Operation status is required", http.StatusBadRequest)
		return
	}

	// Step 2: Convert the "status" parameter to an integer
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		http.Error(w, "Invalid status parameter. Must be an integer.", http.StatusBadRequest)
		return
	}

	// Step 3: Call the service with the integer "status"
	res, err := app.Service.TTF.GetAllByOpr(status)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get by document's ID
// @Description Get by document's ID.
// @ID ttf-get-by
// @Tags TTF
// @Accept json
// @Produce json
// @Param number query string true "Document Number"
// @Param id query string true "Document ID"
// @Success 200 {object} Response{result=[]dto.GetTTFBy} "Get all document by document's ID status have been created successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Api not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /ttf/get-by [get]
func (app *Application) GetBy(w http.ResponseWriter, r *http.Request) {
	// Step 1: Retrieve the "number" and "id" parameters from the query string
	docNumber := r.URL.Query().Get("number")
	docID := r.URL.Query().Get("id")

	// Step 2: Validate that both docNumber and docID are provided
	if docNumber == "" || docID == "" {
		http.Error(w, "Document number and ID are required", http.StatusBadRequest)
		return
	}

	res, err := app.Service.TTF.GetBy(docNumber, docID)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}
