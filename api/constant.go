package api

//for dropdown
import (
	Status "boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) Constants(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/constants")
	api.GET("/search-province", app.SearchProvince)
	api.GET("/get-district", app.GetDistrict)
	api.GET("/get-sub-district", app.GetSubDistrict)
	api.GET("/get-postal-code", app.GetPostalCode)

	api.GET("/get-warehouse", app.GetWarehouse)
	api.GET("/get-product", app.GetProduct)
	api.GET("/search-customer", app.SearchCustomer)
	api.GET("/get-customer-id", app.GetCustomerID)
	api.GET("/get-customer-info", app.GetCustomerInfoByCustomerID)
	api.GET("/get-invoice-names", app.GetInvoiceNamesByCustomerID)
	api.GET("/search-product", app.SearchProduct)

}

// SearchProvince godoc
// @Summary 	Search Province by keyword
// @Description Retrieve the list of provinces by keyword
// @ID 		search-province
// @Tags    Constants
// @Accept  json
// @Produce json
// @Param   keyword query string true "Province search keyword"
// @Success 200 {object} Response{result=[]entity.Province} "List of provinces"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router  /constants/search-province [get]
func (app *Application) SearchProvince(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")

	keyword = strings.TrimSpace(keyword) // ลบช่องว่างหน้าหลังข้อความกันการค้นหาผิดเพราะค่าว่าง
	if keyword == "" {
		app.Logger.Warn("[ keyword is required ]")
		handleError(c, Status.BadRequestError("[ keyword is required ]"))
		return
	}

	provinces, err := app.Service.Constant.SearchProvince(c.Request.Context(), keyword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.Logger.Warn("[ No matching data found ]", zap.String("keyword", keyword))
			handleError(c, Status.NotFoundError("[ No matching data found for keyword: %s ]", keyword))
			return
		}
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(provinces) == 0 {
		app.Logger.Info("[ No data found ]", zap.String("keyword", keyword))
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Provinces retrieved successfully ]", provinces, http.StatusOK)
}

// GetDistrict godoc
// @Summary Get District by ProvinceCode
// @Description Retrieve a list of districts by province code
// @ID get-district-by-province
// @Tags Constants
// @Accept json
// @Produce json
// @Param provinceCode query string true "Province code"
// @Success 200 {object} Response{result=[]entity.District} "List of districts"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-district [get]
func (app *Application) GetDistrict(c *gin.Context) {
	provinceCode := c.DefaultQuery("provinceCode", "")

	districts, err := app.Service.Constant.GetDistrict(c.Request.Context(), provinceCode)
	if err != nil {
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(districts) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Districts retrieved successfully ]", districts, http.StatusOK)
}

// GetSubDistrict godoc
// @Summary Get Subdistrict by DistrictCode
// @Description Retrieve a list of subdistricts by district code
// @ID get-subdistrict-by-district
// @Tags Constants
// @Accept json
// @Produce json
// @Param districtCode query string true "District code"
// @Success 200 {object} Response{result=[]entity.SubDistrict} "List of subdistricts"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-sub-district [get]
func (app *Application) GetSubDistrict(c *gin.Context) {
	districtCode := c.DefaultQuery("districtCode", "")

	subdistricts, err := app.Service.Constant.GetSubDistrict(c.Request.Context(), districtCode)
	if err != nil {
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(subdistricts) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Subdistricts retrieved successfully ]", subdistricts, http.StatusOK)
}

// GetPostalCode godoc
// @Summary Get Postal Code by SubdistrictCode
// @Description Retrieve postal code by subdistrict code
// @ID get-postalcode-by-subdistrict
// @Tags Constants
// @Accept json
// @Produce json
// @Param subdistrictCode query string true "Subdistrict code"
// @Success 200 {string} Response{result=[]entity.PostalCode} "Postal code"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-postal-code [get]
func (app *Application) GetPostalCode(c *gin.Context) {
	subdistrictCode := c.DefaultQuery("subdistrictCode", "")

	postalCode, err := app.Service.Constant.GetPostalCode(c.Request.Context(), subdistrictCode)
	if err != nil {
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(postalCode) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Postal code retrieved successfully ]", postalCode, http.StatusOK)
}

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
		app.Logger.Error("[ Error ]", zap.Error(err))
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
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Product successfully ]", result, http.StatusOK)
}

// @Summary Search Customer by CustomerID or InvoiceName
// @Description Search for customers by CustomerID or InvoiceName with pagination support (limit and offset)
// @ID search-customer
// @Tags Constants
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param searchType query string true "Search by 'CustomerID' or 'InvoiceName'"
// @Param offset query int false "Offset for pagination (default is 0)" default(0)
// @Param limit query int false "Limit for number of customers to return (default is 4)" default(4)
// @Success 200 {object} Response{result=[]entity.InvoiceInformation} "Search Results"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "No matching customer found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/search-customer [get]
func (app *Application) SearchCustomer(c *gin.Context) {
	keyword := c.Query("keyword")
	searchType := c.DefaultQuery("searchType", "") // ค่าเป็น "CustomerID" หรือ "InvoiceName"
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))

	keyword = strings.TrimSpace(keyword) // ลบช่องว่างหน้าหลังข้อความกันการค้นหาผิดเพราะค่าว่าง
	if keyword == "" {
		app.Logger.Warn("[ keyword is required ]")
		handleError(c, Status.BadRequestError("[ keyword is required ]"))
		return
	}

	result, err := app.Service.Constant.SearchCustomer(c.Request.Context(), keyword, searchType, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.Logger.Warn("[ No data found ]", zap.String("keyword", keyword))
			handleError(c, Status.NotFoundError("[ No data found for keyword: %s ]", keyword))
			return
		}
		app.Logger.Error("[ Error ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]", zap.String("keyword", keyword))
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Get Customer successfully ]", result, http.StatusOK)
}

// @Summary Get Customer IDs
// @Description Retrieve all customer IDs
// @ID get-customer-ids
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]entity.InvoiceInformation} "List of customer IDs"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-customer-id [get]
func (app *Application) GetCustomerID(c *gin.Context) {
    result, err := app.Service.Constant.GetCustomerID(c.Request.Context())
    if err != nil {
        app.Logger.Error("[ Error ]", zap.Error(err))
        handleError(c, err)
        return
    }

    handleResponse(c, true, "[ Get Customer IDs successfully ]", result, http.StatusOK)
}

// @Summary Get Customer Info by CustomerID
// @Description Retrieve customer information by customer ID
// @ID get-customer-info-by-id
// @Tags Constants
// @Accept json
// @Produce json
// @Param customerID query string true "Customer ID"
// @Param offset query int false "Offset for pagination (default is 0)" default(0)
// @Param limit query int false "Limit for number of customers to return (default is 4)" default(4)
// @Success 200 {object} Response{result=[]entity.InvoiceInformation} "Customer information"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-customer-info [get]
func (app *Application) GetCustomerInfoByCustomerID(c *gin.Context) {
    customerID := c.Query("customerID")
    if customerID == "" {
        app.Logger.Warn("[ customerID is required ]")
        handleError(c, Status.BadRequestError("[ customerID is required ]"))
        return
    }

    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))

    result, err := app.Service.Constant.GetCustomerInfoByCustomerID(c.Request.Context(), customerID, offset, limit)
    if err != nil {
        app.Logger.Error("[ Error ]", zap.Error(err))
        handleError(c, err)
        return
    }

    handleResponse(c, true, "[ Get Customer Info successfully ]", result, http.StatusOK)
}

// @Summary Get Invoice Names by CustomerID
// @Description Retrieve invoice names by customer ID
// @ID get-invoice-names-by-customer-id
// @Tags Constants
// @Accept json
// @Produce json
// @Param customerID query string true "Customer ID"
// @Param offset query int false "Offset for pagination (default is 0)" default(0)
// @Param limit query int false "Limit for number of invoices to return (default is 4)" default(4)
// @Success 200 {object} Response{result=[]entity.InvoiceInformation} "Invoice names"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/get-invoice-names [get]
func (app *Application) GetInvoiceNamesByCustomerID(c *gin.Context) {
    customerID := c.Query("customerID")
    if customerID == "" {
        app.Logger.Warn("[ customerID is required ]")
        handleError(c, Status.BadRequestError("[ customerID is required ]"))
        return
    }

    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))

    result, err := app.Service.Constant.GetInvoiceNamesByCustomerID(c.Request.Context(), customerID, offset, limit)
    if err != nil {
        app.Logger.Error("[ Error ]", zap.Error(err))
        handleError(c, err)
        return
    }

    handleResponse(c, true, "[ Get Invoice Names successfully ]", result, http.StatusOK)
}

// @Summary Search Product by Keyword
// @Description Search for products by SKU or NAMEALIAS with pagination support (limit and offset)
// @ID search-product
// @Tags Constants
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Param searchType query string true "Search by 'SKU' or 'NAMEALIAS'"
// @Param offset query int false "Offset for pagination (default is 0)" default(0)
// @Param limit query int false "Limit for number of products to return (default is 4)" default(4)
// @Success 200 {object} Response{result=[]entity.ROM_V_ProductAll} "Search Results"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "No matching products found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /constants/search-product [get]
func (app *Application) SearchProduct(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	searchType := c.DefaultQuery("searchType", "")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))

	keyword = strings.TrimSpace(keyword) // ลบช่องว่างหน้าหลังข้อความกันการค้นหาผิดเพราะค่าว่าง
	if keyword == "" {
		app.Logger.Warn("[ keyword is required ]")
		handleError(c, Status.BadRequestError("[ keyword is required ]"))
		return
	}

	result, err := app.Service.Constant.SearchProduct(c.Request.Context(), keyword, searchType, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.Logger.Warn("[ No data found ]", zap.String("keyword", keyword))
			handleError(c, Status.NotFoundError("[ No data found for keyword: %s ]", keyword))
			return
		}
		app.Logger.Error("[ Failed to search ]", zap.String("keyword", keyword), zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]", zap.String("keyword", keyword))
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Search Product successfully ]", result, http.StatusOK)
}
