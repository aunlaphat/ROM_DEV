package api

//for dropdown

import (
	"boilerplate-backend-go/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Constants(apiRouter *chi.Mux) {
	apiRouter.Route("/constants", func(r chi.Router) {
		r.Get("/get-province", app.GetThaiProvince)
		r.Get("/get-district", app.GetThaiDistrict)
		r.Get("/get-sub-district", app.GetThaiSubDistrict)
		// r.Get("/get-postcode", app.GetPostCode) // รอพี่ไบรท์ดึงข้อมูล
		r.Get("/get-warehouse", app.GetWarehouse)
		r.Get("/get-productAll", app.GetProductAll)
		r.Get("/get-productAlls", app.GetProductAlls)
		//r.Get("/get-customer", app.GetCustomer) // รอพี่ไบรท์ดึงข้อมูล
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
	result, err := app.Service.Constant.GetThaiProvince()
	if err != nil {
		handleError(w, err)
		return
	}
	handleResponse(w, true, "⭐ Get Thai Province successfully ⭐", result, http.StatusOK)
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
	result, err := app.Service.Constant.GetThaiDistrict()
	if err != nil {
		handleError(w, err)
		return
	}
	handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

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
	result, err := app.Service.Constant.GetThaiSubDistrict()
	if err != nil {
		handleError(w, err)
		return
	}
	handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

}

// // @Summary Get Thai GetPostCode
// // @Description Get all Thai GetPostCode.
// // @ID get-post-code
// // @Tags Constants
// // @Accept json
// // @Produce json
// // @Success 200 {object} Response{result=[]entity.GetPostCode} "GetPostCode"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 404 {object} Response "SubDistrict not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /constants/get-sub-district [get]
// func (app *Application) GetPostCode(w http.ResponseWriter, r *http.Request) {
// 	result, err := app.Service.Constant.GetPostCode()
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 		handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

// }

// @Summary Get Warehouse
// @Description Get Warehouse
// @ID get-warehouse
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.Warehouse} "Warehouse"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "SubDistrict not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-warehouse [get]
func (app *Application) GetWarehouse(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.Constant.GetWarehouse()
	if err != nil {
		handleError(w, err)
		return
	}
	handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

}

// @Summary Get ProductAll
// @Description Get all product
// @ID get-productAll
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.ROM_V_ProductAll} "Product"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "SubDistrict not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-productAll [get]
func (app *Application) GetProductAll(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.Constant.GetProductAll()
	if err != nil {
		handleError(w, err)
		return
	}
	handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

}

// @Summary Get ProductAll with Pagination
// @Description Get paginated products
// @ID get-productAll-paginated
// @Tags Constants
// @Accept json
// @Produce json
// @Param page query int true "Page number" default(1)
// @Param limit query int true "Limit per page" default(10)
// @Success 200 {object} Response{result=[]entity.ROM_V_ProductAll, total=int} "Paginated Product List"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-productAlls [get]
func (app *Application) GetProductAlls(w http.ResponseWriter, r *http.Request) {

	page, limit := utils.ParsePagination(r)

	result, err := app.Service.Constant.GetProductAllWithPagination(page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)
}

// // @Summary Get Customer
// // @Description Get inform customer
// // @ID get-customer
// // @Tags Constants
// // @Accept json
// // @Produce json
// // @Success 200 {object} Response{result=[]entity.ROM_V_Customer} "Customer"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 404 {object} Response "SubDistrict not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /constants/get-customer [get]
// func (app *Application) GetCustomer(w http.ResponseWriter, r *http.Request) {
// 	result, err := app.Service.Constant.GetThaiSubDistrict()
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 		handleResponse(w, true, "⭐ successfully ⭐", result, http.StatusOK)

// }
