package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var response = "success"

func (app *Application) Constants(apiRouter *chi.Mux) {
	apiRouter.Route("/constants", func(r chi.Router) {
		r.Get("/get-province", app.GetThaiProvince)
		r.Get("/get-district", app.GetThaiDistrict)
		r.Get("/get-sub-district", app.GetThaiSubDistrict)
	})
}

// @Summary Get Thai Province
// @Description Get all Thai Province.
// @ID get-province
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.Province} "Province"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Province not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-province [get]
func (app *Application) GetThaiProvince(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.Constant.GetThaiProvince()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get Thai District
// @Description Get all Thai District.
// @ID get-district
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.District} "District"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "District not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-district [get]
func (app *Application) GetThaiDistrict(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.Constant.GetThaiDistrict()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get Thai SubDistrict
// @Description Get all Thai SubDistrict.
// @ID get-sub-district
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.SubDistrict} "SubDistrict"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "SubDistrict not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-sub-district [get]
func (app *Application) GetThaiSubDistrict(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.Constant.GetThaiSubDistrict()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}
