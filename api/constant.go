 package api

//for dropdown

import (
	"boilerplate-backend-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) Constants(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/constants")
		api.GET("/get-province", app.GetThaiProvince)
		api.GET("/get-district", app.GetThaiDistrict)
		api.GET("/get-sub-district", app.GetThaiSubDistrict)
		// api.GET("/get-postcode", app.GetPostCode) // รอพี่ไบรท์ดึงข้อมูล
		api.GET("/get-warehouse", app.GetWarehouse)
		api.GET("/get-product", app.GetProduct)
		//api.GET("/get-customer", app.GetCustomer) // รอพี่ไบรท์ดึงข้อมูล
		api.GET("/search-product", app.SearchProduct)

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
func (app *Application) GetThaiProvince(c *gin.Context) {
	result, err := app.Service.Constant.GetThaiProvince(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Thai Province successfully ]", result, http.StatusOK)
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
func (app *Application) GetThaiDistrict(c *gin.Context) {
	result, err := app.Service.Constant.GetThaiDistrict(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Thai District successfully ]", result, http.StatusOK)

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
func (app *Application) GetThaiSubDistrict(c *gin.Context) {
	result, err := app.Service.Constant.GetThaiSubDistrict(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Thai SubDistrict successfully ]", result, http.StatusOK)

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
// func (app *Application) GetPostCode(c *gin.Context) {
// 	result, err := app.Service.Constant.GetPostCode(c.Request.Context())
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}
// 		handleResponse(c, true, "[ Get PostCode successfully ]", result, http.StatusOK)

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
func (app *Application) GetWarehouse(c *gin.Context) {
	result, err := app.Service.Constant.GetWarehouse(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Warehouse successfully ]", result, http.StatusOK)

}

// @Summary Get ProductAll with Pagination
// @Description Get paginated products
// @ID get-productAll-paginated
// @Tags Constants
// @Accept json
// @Produce json
// @Param page query int true "Page number" default(1)
// @Param limit query int true "Limit per page" default(4)
// @Success 200 {object} Response{result=[]entity.ROM_V_ProductAll, total=int} "Paginated Product List"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-product [get]
func (app *Application) GetProduct(c *gin.Context) {

	page, limit := utils.ParsePagination(c.Request)

	result, err := app.Service.Constant.GetProduct(c.Request.Context(), page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Product successfully ]", result, http.StatusOK)
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
// func (app *Application) GetCustomer(c *gin.Context) {
// 	result, err := app.Service.Constant.GetThaiSubDistrict(c.Request.Context())
// 	if err != nil {
// 		handleError(c, err)
// 		return
// 	}
// 		handleResponse(c, true, "[ Get Customer successfully ]", result, http.StatusOK)

// }

// @Summary Search Product by Keyword
// @Description Search for products by name or SKU with a limit
// @ID search-product
// @Tags Constants
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param limit query int true "Limit per page" default(10)
// @Success 200 {object} Response{result=[]entity.ROM_V_ProductAll} "Search Results"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/search-product [get]
func (app *Application) SearchProduct(c *gin.Context) {
	keyword := c.Query("keyword")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	result, err := app.Service.Constant.SearchProduct(c.Request.Context(), keyword, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Search Product successfully ]", result, http.StatusOK)
}
