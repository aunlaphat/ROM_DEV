package api

import (
	entity "boilerplate-backend-go/Entity"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var response = "success"

func (app *Application) Constants(apiRouter *chi.Mux) {
	apiRouter.Route("/constants", func(r chi.Router) {
		r.Get("/get-type-req", app.GetTypeReq)
		r.Get("/get-type-foc", app.GetTypeFoc)
		r.Get("/get-type-logistic", app.GetTypeLogistic)
		r.Get("/get-province", app.GetThaiProvince)
		r.Get("/get-district", app.GetThaiDistrict)
		r.Get("/get-sub-district", app.GetThaiSubDistrict)
		r.Get("/get-warehouse-location", app.GetWarehouseLoaction)
		r.Post("/get-sku", app.GetSku)
	})
}

// @Summary Get Type req
// @Description Get all Type req.
// @ID get-type-req
// @Tags Constants
// @Accept json
// @Produce json
// @Param permissionID query string true "input Permission ID"
// @Success 200 {object} Response{result=[]entity.TypeReqPermission} "Type Req"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Type Req not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-type-req [get]
func (app *Application) GetTypeReq(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	permissionID := queryValues.Get("permissionID")
	res, err := app.Service.Constant.GetTypeReq(permissionID)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get Type foc
// @Description Get all Type foc.
// @ID get-type-foc
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.TypeFoc} "Type FOC"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Type Foc not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-type-foc [get]
func (app *Application) GetTypeFoc(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.Constant.GetTypeFoc()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary Get Type logistic
// @Description Get all Type logistic.
// @ID get-type-logistic
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.TypeLogistic} "Type Logistic"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Type Foc not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-type-logistic [get]
func (app *Application) GetTypeLogistic(w http.ResponseWriter, r *http.Request) {
	res, err := app.Service.Constant.GetTypeLogistic()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
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

// @Summary get warehouse & location by user's id
// @Description get warehouse & location by user's id.
// @ID get-warehouse-location
// @Tags Constants
// @Accept json
// @Produce json
// @Param userID query string true "User's ID"
// @Success 200 {object} Response{result=[]service.SiteWarehouseLocations} "Warehouse Location by user'ID"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "warehoues location not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-warehouse-location [get]
func (app *Application) GetWarehouseLoaction(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "user'id is required", http.StatusBadRequest)
		return
	}
	res, err := app.Service.Constant.GetWarehouseLocation(userID)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary get sku by warehouse, location and site
// @Description get sku by warehouse, location and site.
// @ID get-sku
// @Tags Constants
// @Accept json
// @Produce json
// @Param get-sku body entity.WarehouseLocation true "Get Sku"
// @Success 200 {object} Response{result=[]entity.GetSKU} "Get sku by warehouse, location and site"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "sku not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-sku [post]
func (app *Application) GetSku(w http.ResponseWriter, r *http.Request) {
	req := entity.WarehouseLocation{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		HandleError(w, err)
		return
	}
	res, err := app.Service.Constant.GetSku(req)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}
