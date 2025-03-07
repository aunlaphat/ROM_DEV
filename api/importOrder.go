package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	Status "boilerplate-backend-go/errors"
	"boilerplate-backend-go/middleware"
	"strings"

	// "boilerplate-backend-go/utils"
	// "encoding/json"
	"time"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ImportOrder => หน้ารับเข้าสินค้าหน้าคลัง โดยฝั่ง scan สินค้าหน้าคลังจะต้องเข้าทำหน้านี้เพื่อสแกนสินค้าแต่ละ order เข้าระบบ
func (app *Application) ImportOrderRoute(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/import-order")
	api.GET("/search", app.SearchOrderORTracking) // เสิร์ช orderNo, trackingNo เพื่อทำรายการรับเข้าต่อ (ต้องมีข้อมูลขึ้นจึงทำรายการต่อได้)
	api.GET("/search-order-tracking", app.SearchOrderORTrackingNo)
	api.GET("/get-order-tracking", app.GetOrderTracking)

	apiAuth := api.Group("/")
	apiAuth.Use(middleware.JWTMiddleware(app.TokenAuth))
	apiAuth.POST("/upload-photo", app.UploadPhotoHandler)            // รับการอัปโหลดภาพเข้าระบบ
	apiAuth.GET("/summary/:orderNo", app.GetSummaryImportOrder)      // หน้าสรุปรวมข้อมูลภาพที่ถ่ายเข้า ก่อนกดยืนยันรับเข้า (confirm receipt)
	apiAuth.POST("/validate-sku/:orderNo/:sku", app.ValidateSKU)     // ใช้ตรวจสอบ sku ที่รับเข้าหน้างานว่าตรงกับในระบบที่ user ทำรายการมา
	apiAuth.POST("/confirm-receipt/:identifier", app.ConfirmReceipt) // ยืนยันการรับเข้าหน้าคลัง (เมื่อถ่ายรับเข้าเสร็จทุกรายการ)

	// ยังไม่ใช้
	apiAuth.POST("/create-confirm-wh", app.ConfirmFromWH)

}

// SearchOrderORTracking godoc
// @Summary 	Search order by OrderNo or TrackingNo
// @Description Retrieve the details of an order by its OrderNo or TrackingNo
// @ID 			search-orderNo-or-trackingNo
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Param 		search query string true "OrderNo or TrackingNo"
// @Success 200 {object} api.Response{result=response.ImportOrderResponse} "Order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "OrderNo or TrackingNo not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router  /import-order/search [get]
func (app *Application) SearchOrderORTracking(c *gin.Context) {
	search := c.DefaultQuery("search", "")

	// *️⃣ ตรวจสอบ search ว่าเป็นค่าว่างหรือไม่
	search = strings.TrimSpace(search) // ลบช่องว่างหน้าหลังข้อความกันการค้นหาผิดเพราะค่าว่าง
	if search == "" {
		app.Logger.Warn("[ Search input is required ]")
		handleError(c, Status.BadRequestError("[ Search input is required ]"))
		return
	}

	result, err := app.Service.ImportOrder.SearchOrderORTracking(c.Request.Context(), search)
	if err != nil {
		app.Logger.Error("[ Failed to search ]", zap.String("Search", search), zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Found Orders retrieved successfully ]", result, http.StatusOK)
}

// SearchOrderORTracking godoc
// @Summary 	Search order by OrderNo or TrackingNo
// @Description Retrieve the details of an order by its OrderNo or TrackingNo
// @ID 			search-orderNo-or-tracking
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Param 		search query string true "OrderNo or TrackingNo"
// @Success 200 {object} api.Response{result=response.ImportOrderResponse} "Order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "OrderNo or TrackingNo not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router  /import-order/search-order-tracking [get]
func (app *Application) SearchOrderORTrackingNo(c *gin.Context) {
	search := c.DefaultQuery("search", "")

	// *️⃣ ตรวจสอบ search ว่าเป็นค่าว่างหรือไม่
	search = strings.TrimSpace(search) // ลบช่องว่างหน้าหลังข้อความกันการค้นหาผิดเพราะค่าว่าง
	if search == "" {
		app.Logger.Warn("[ Search input is required ]")
		handleError(c, Status.BadRequestError("[ Search input is required ]"))
		return
	}

	result, err := app.Service.ImportOrder.SearchOrderORTrackingNo(c.Request.Context(), search)
	if err != nil {
		app.Logger.Error("[ Failed to search ]", zap.String("Search", search), zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Found Orders retrieved successfully ]", result, http.StatusOK)
}

// @Summary 	Get Order and Tracking
// @Description Retrieve Order and Tracking
// @ID 			Get-Order-Tracking
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]response.ImportItem} "Get All"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Not Found Endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/import-order/get-order-tracking [get]
func (app *Application) GetOrderTracking(c *gin.Context) {
	result, err := app.Service.ImportOrder.GetOrderTracking(c.Request.Context())
	if err != nil {
		app.Logger.Error("[ Error fetching order ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Get Order successfully ]", result, http.StatusOK)
}

// UploadPhotoHandler godoc
// @Summary Upload Photo
// @Description Upload a photo for a return order
// @ID upload-photo
// @Tags Import Order
// @Accept multipart/form-data
// @Produce json
// @Param orderNo formData string true "Order No"
// @Param imageTypeID formData string true "ImageTypeID (1, 2, 3)"
// @Param sku formData string false "SKU (required if imageTypeID is '3')"
// @Param file formData file true "Photo file"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router  /import-order/upload-photo [post]
func (app *Application) UploadPhotoHandler(c *gin.Context) {
	orderNo := c.PostForm("orderNo")
	imageTypeID := c.PostForm("imageTypeID")
	sku := c.PostForm("sku")

	if orderNo == "" || imageTypeID == "" {
		app.Logger.Warn("[ OrderNo and ImageTypeID are required ]")
		handleError(c, Status.BadRequestError("[ OrderNo and ImageTypeID are required ]"))
		return
	}

	// *️⃣ หาก ImageTypeID เป็น 3 แต่ SKU ไม่ได้ถูกส่งมา
	if imageTypeID == "3" && sku == "" {
		app.Logger.Warn("[ SKU is required for 3 imageTypeID ]")
		handleError(c, Status.BadRequestError("[ SKU is required for 3 imageTypeID ]"))
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		app.Logger.Error("[ Failed to get file from request ]", zap.Error(err))
		handleError(c, Status.BadRequestError("[ Failed to get file from request ]"))
		return
	}

	// เปิดไฟล์เพื่ออ่านข้อมูล
	file, err := header.Open()
	if err != nil {
		app.Logger.Error("[ Failed to open file ]", zap.Error(err))
		handleError(c, Status.InternalError("[ Failed to open file ]"))
		return
	}
	defer file.Close() // ปิดไฟล์เมื่อใช้งานเสร็จ

	err = app.Service.ImportOrder.UploadPhotoHandler(c.Request.Context(), orderNo, imageTypeID, sku, file, header.Filename)
	if err != nil {
		app.Logger.Error("[ Failed to upload ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ File uploaded successfully ]", nil, http.StatusOK)
}

// @Summary 	Get Sum detail of Import Order
// @Description Retrieve the details of Receipt
// @ID 			GetSummary-ImportOrder
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Param orderNo path string true "Order No"
// @Success 200 {object} api.Response{result=[]response.ImportOrderSummary} "Get All"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Not Found Endpoint"
// @Failure 500 {object} Response "Internal Server Error"
// @Router 	/import-order/summary/{orderNo} [get]
func (app *Application) GetSummaryImportOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

	summary, err := app.Service.ImportOrder.GetSummaryImportOrder(c.Request.Context(), orderNo)
	if err != nil {
		app.Logger.Error("[ Failed to get order ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(summary) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Summary retrieved successfully ]", summary, http.StatusOK)
}

// ValidateSKU godoc
// @Summary 	Validate SKU
// @Description Validate SKU
// @ID 			validate-sku
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Param 		orderNo path string true "Order No"
// @Param 		sku path string true "SKU"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /import-order/validate-sku/{orderNo}/{sku} [post]
func (app *Application) ValidateSKU(c *gin.Context) {
	orderNo := c.Param("orderNo")
	sku := c.Param("sku")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

	if sku == "" {
		app.Logger.Warn("[ SKU is required ]")
		handleError(c, Status.BadRequestError("[ SKU is required ]"))
		return
	}

	valid, err := app.Service.ImportOrder.ValidateSKU(c.Request.Context(), orderNo, sku)
	if err != nil {
		handleError(c, err)
		return
	}
	if !valid {
		app.Logger.Warn("[ SKU mismatch ]", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
		handleError(c, Status.BadRequestError("[ SKU %s mismatch for this order %s ]", sku, orderNo))
		return
	}

	handleResponse(c, true, "[ Both match: Confirm Receipt ]", nil, http.StatusOK)
}

// ConfirmReceipt godoc
// @Summary 	Confirm Receipt from Ware House
// @Description Confirm a trade return order based on the provided identifier (OrderNo or TrackingNo) and input lines for ReturnOrderLine.
// @ID 			confirm-trade-return
// @Tags 		Import Order
// @Accept 		json
// @Produce 	json
// @Param 		identifier path string true "OrderNo or TrackingNo"
// @Param 		request body request.ConfirmTradeReturnRequest true "Trade return request details"
// @Success 200 {object} api.Response{result=response.ConfirmReceipt} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /import-order/confirm-receipt/{identifier} [post]
func (app *Application) ConfirmReceipt(c *gin.Context) {
	identifier := c.Param("identifier")

	// *️⃣ ตรวจสอบค่าว่าง
	if identifier == "" {
		app.Logger.Warn("[ OrderNo or TrackingNo are required ]")
		handleError(c, Status.BadRequestError("[ OrderNo or TrackingNo are required ]"))
		return
	}

	var req request.ConfirmTradeReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("[ Unauthorized - Missing UserID ]")
		handleError(c, Status.UnauthorizedError("[ Unauthorized - Missing UserID ]"))
		return
	}

	req.Identifier = identifier
	err := app.Service.BeforeReturn.ConfirmReceipt(c.Request.Context(), req, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	result := res.ConfirmReceipt{
		Identifier:     req.Identifier,
		StatusReturnID: "7 (WAITING)",
		StatusCheckID:  "1 (WAITING)",
		UpdateBy:       userID.(string),
		UpdateDate:     time.Now(),
	}

	handleResponse(c, true, "[ Confirmed from Ware House successfully => Status (waiting ✔️) ]", result, http.StatusOK)
}

// ยังไม่ใช้
// ConfirmFromWH godoc
// @Summary 	Import order
// @Description Upload multiple images and data for a specific SoNo
// @ID 			Import-Order
// @Tags 		Import Order
// @Accept 		multipart/form-data
// @Produce 	json
// @Param 		soNo formData string true "Sale Order Number"
// @Param 		imageTypeID formData int true "Type of the image (1, 2, or 3)"
// @Param 		skus formData string false "SKU (Optional)"
// @Param 		files formData file true "Files to upload"
// @Success 200 {object} api.Response{result=response.ImageResponse} "Successful"
// @Failure 400 {object} api.Response "Invalid input"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /import-order/create-confirm-wh [post]
func (app *Application) ConfirmFromWH(c *gin.Context) {
	// *️⃣ Parse Form Data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		handleError(c, Status.BadRequestError("Unable to parse form data"))
		return
	}

	// *️⃣ รับค่าจาก Form
	soNo := c.PostForm("soNo")
	imageTypeID, err := strconv.Atoi(c.PostForm("imageTypeID"))
	if err != nil {
		handleError(c, Status.BadRequestError("Invalid Image Type ID"))
		return
	}

	skus := c.PostForm("skus")
	files := c.Request.MultipartForm.File["files"]

	if soNo == "" {
		app.Logger.Warn("[ SoNo is required ]")
		handleError(c, Status.BadRequestError("[ SoNo is required ]"))
		return
	}

	if imageTypeID < 1 || imageTypeID > 3 {
		app.Logger.Warn("[ Image Type ID 1, 2, 3 only ]")
		handleError(c, Status.BadRequestError("[ Image Type ID 1, 2, 3 only ]"))
		return
	}

	if len(files) == 0 {
		app.Logger.Warn("[ no files uploaded ]")
		handleError(c, Status.BadRequestError("[ no files uploaded ]"))
		return
	}

	// *️⃣ เรียก Service เพื่อประมวลผล
	result, err := app.Service.ImportOrder.ConfirmFromWH(c.Request.Context(), soNo, imageTypeID, skus, files)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Data Insert successful ]", result, http.StatusOK)
}
