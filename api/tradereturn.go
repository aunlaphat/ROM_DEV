package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/middleware"
	"boilerplate-backend-go/utils"

	// "boilerplate-backend-go/utils"
	// "encoding/json"
	// "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TradeReturn => ส่วนการคืนของทางฝั่ง offline หาก order ค้างคลัง ต้องการเปลี่ยนรุ่นใหม่ ลูกค้าหน้าสาขาต้องการคืนของ จะทำรายการที่เทรด
func (app *Application) TradeReturnRoute(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/trade-return")
	api.GET("/get-waiting", app.GetStatusWaitingDetail)       // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ waiting
	api.GET("/get-confirm", app.GetStatusConfirmDetail)       // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ confirm
	api.GET("/search-waiting", app.SearchStatusWaitingDetail) // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ waiting ตามช่วงวันที่สร้าง(CreateDate) เริ่มต้น-สิ้นสุด จะแสดงตามวันที่นั้น
	api.GET("/search-confirm", app.SearchStatusConfirmDetail) // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ confirm ตามช่วงวันที่สร้าง(CreateDate) เริ่มต้น-สิ้นสุด จะแสดงตามวันที่นั้น

	apiAuth := api.Group("/")
	apiAuth.Use(middleware.JWTMiddleware(app.TokenAuth))
	apiAuth.POST("/create-trade", app.CreateTradeReturn)          // สร้างฟอร์มทำรายการคืนของเข้าระบบ
	apiAuth.POST("/add-line/:orderNo", app.CreateTradeReturnLine) // สร้างรายการคืนแต่ละรายการของ order นั้นเพิ่ม
	apiAuth.PATCH("/confirm-return/:orderNo", app.ConfirmReturn)  // การยันยืนรับเข้าโดยฝั่งบัญชี จะทำการตรวจเช็คจากข้อมูล confirmReceipt + ข้อมูลในระบบ เมื่อเช็คว่าจำนวนคืนตรงกันจะถือว่าทำรายการคืนสำเร็จ

}

// @Summary     Get Return Orders with StatusCheckID = 1
// @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting)
// @ID      Get-Waiting-ReturnOrder
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router  /trade-return/get-waiting [get]
func (app *Application) GetStatusWaitingDetail(c *gin.Context) {
	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatus(c.Request.Context(), 1) // StatusCheckID = 1
	if err != nil {
		app.Logger.Error("[ Failed to fetch Return Orders ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Return Orders with StatusCheckID = 1 (WAITING) retrieved successfully ]", result, http.StatusOK)
}

// @Summary     Get Return Orders with StatusCheckID = 2
// @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed)
// @ID      Get-Confirm-ReturnOrder
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router  /trade-return/get-confirm [get]
func (app *Application) GetStatusConfirmDetail(c *gin.Context) {
	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatus(c.Request.Context(), 2) // StatusCheckID = 2
	if err != nil {
		app.Logger.Error("[ Failed to fetch Return Orders ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Return Orders with StatusCheckID = 2 (CONFIRM) retrieved successfully ]", result, http.StatusOK)
}

// @Summary     Search Return Orders with StatusCheckID = 1 by Date Range
// @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting) within a specific date range
// @ID      Search-Waiting-ReturnOrder
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Param   startDate query string true "Start Date (YYYY-MM-DD)"
// @Param   endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router  /trade-return/search-waiting [get]
func (app *Application) SearchStatusWaitingDetail(c *gin.Context) {
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(c.Request.Context(), 1, startDate, endDate) // StatusCheckID = 1
	if err != nil {
		app.Logger.Error("[ Failed to fetch Return Orders ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Return Orders of StatusCheckID = 1 retrieved successfully ]", result, http.StatusOK)
}

// @Summary     Search Return Orders with StatusCheckID = 2 by Date Range
// @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed) within a specific date range
// @ID      Search-Confirm-ReturnOrder
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Param   startDate query string true "Start Date (YYYY-MM-DD)"
// @Param   endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router  /trade-return/search-confirm [get]
func (app *Application) SearchStatusConfirmDetail(c *gin.Context) {
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(c.Request.Context(), 2, startDate, endDate) // StatusCheckID = 2
	if err != nil {
		app.Logger.Error("[ Failed to fetch Return Orders ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Return Orders of StatusCheckID = 2 retrieved successfully ]", result, http.StatusOK)
}

// @Summary     Create a new trade return order
// @Description Create a new trade return order with multiple order lines
// @ID      create-trade-return
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Param   body body request.BeforeReturnOrder true "Trade Return Detail"
// @Success 201 {object} api.Response{result=response.BeforeReturnOrderResponse} "Trade return created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router  /trade-return/create-trade [post]
func (app *Application) CreateTradeReturn(c *gin.Context) {
	var req request.BeforeReturnOrder

	// *️⃣ ดึง Request JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Logger.Warn("[ Invalid request payload ]", zap.Error(err))
		handleValidationError(c, err)
		return
	}

	// *️⃣ ตรวจสอบว่า ReturnOrderLine ต้องไม่เป็นค่าว่าง (มีอย่างน้อย 1 รายการ)
	if len(req.BeforeReturnOrderLines) == 0 {
		app.Logger.Warn("[ sku information can't empty must be > 0 line ]")
		handleError(c, errors.BadRequestError("[ sku information can't empty must be > 0 line ]"))
		return
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateTradeReturn(req); err != nil {
		app.Logger.Warn("[ Validation failed ]", zap.Error(err))
		handleError(c, errors.BadRequestError("[ Validation failed: %v ]", err))
		return
	}

	// *️⃣ ดึง userID จาก JWT token
	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("[ Unauthorized - Missing UserID ]")
		handleResponse(c, false, "[ Unauthorized - Missing UserID ]", nil, http.StatusUnauthorized)
		return
	}

	req.CreateBy = userID.(string)
	result, err := app.Service.BeforeReturn.CreateTradeReturn(c.Request.Context(), req)
	if err != nil {
		app.Logger.Error("[ Failed to create trade return order ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Created trade return order successfully => Status (booking ✔️) ]", result, http.StatusOK)
}

// @Summary     Add a new trade return line to an existing order
// @Description Add a new trade return line based on the provided order number and line details
// @ID      add-trade-return-line
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Param   orderNo path string true "Order number"
// @Param   body body request.TradeReturnLine true "Trade Return Line Details"
// @Success 201 {object} api.Response{result=response.BeforeReturnOrderItem} "Trade return line created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router  /trade-return/add-line/{orderNo} [post]
func (app *Application) CreateTradeReturnLine(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, errors.BadRequestError("[ OrderNo is required ]"))
		return
	}

	var req request.TradeReturnLine

	// *️⃣ ดึง Request JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Logger.Warn("[ Invalid request payload ]", zap.Error(err))
		handleValidationError(c, err)
		return
	}

	// *️⃣ Validate request ที่ส่งมา
	if err := utils.ValidateCreateTradeReturnLine(req.TradeReturnLine); err != nil {
		app.Logger.Warn("[ Validation failed ]", zap.Error(err))
		handleError(c, errors.BadRequestError("[ Validation failed: %v ]", err))
		return
	} 

	// *️⃣ ดึง userID จาก JWT token
	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("[ Unauthorized - Missing UserID ]")
		handleResponse(c, false, "[ Unauthorized - Missing UserID ]", nil, http.StatusUnauthorized)
		return
	}

	//req.CreateBy = userID.(string)
	// Set CreateBy จาก claims
	for i := range req.TradeReturnLine {
		req.TradeReturnLine[i].CreateBy = userID.(string)
	}

	result, err := app.Service.BeforeReturn.CreateTradeReturnLine(c.Request.Context(), orderNo, req)
	if err != nil {
		app.Logger.Error("[ Failed to create trade return line ]", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Created trade return line successfully ]", result, http.StatusCreated)
}

// ConfirmToReturn godoc
// @Summary     Confirm Return Order to Success
// @Description Confirm a trade return order based on the provided order number (OrderNo) and input lines for ReturnOrderLine.
// @ID      confirm-to-return
// @Tags    Trade Return
// @Accept  json
// @Produce json
// @Param   orderNo path string true "OrderNo"
// @Param   request body request.ConfirmToReturnRequest true "Updated trade return request details"
// @Success 200 {object} api.Response{result=response.ConfirmToReturnOrder} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router  /trade-return/confirm-return/{orderNo} [Patch]
func (app *Application) ConfirmReturn(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, errors.BadRequestError("[ OrderNo is required ]"))
		return
	}

	var req request.ConfirmToReturnRequest

	// *️⃣ ดึง Request JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Logger.Warn("[ Invalid request payload ]", zap.Error(err))
		handleValidationError(c, err)
		return
	}

	// *️⃣ ดึง userID จาก JWT token
	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("[ Unauthorized - Missing UserID ]")
		handleResponse(c, false, "[ Unauthorized - Missing UserID ]", nil, http.StatusUnauthorized)
		return
	}

	req.OrderNo = orderNo
	err := app.Service.BeforeReturn.ConfirmReturn(c.Request.Context(), req, userID.(string))
	if err != nil {
		app.Logger.Error("[ Failed to comfirm return ]", zap.Error(err))
		handleError(c, err)
		return
	}

	result := response.ConfirmToReturnOrder{
		OrderNo:        req.OrderNo,
		StatusReturnID: "6 (success)",
		StatusCheckID:  "2 (CONFIRM)",
		UpdateBy:       userID.(string),
		UpdateDate:     time.Now(),
	}

	handleResponse(c, true, "[ Confirmed to Return Order successfully => Status (success ✔️, confirm ✔️) ]", result, http.StatusOK)
}
