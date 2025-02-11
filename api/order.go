package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/middleware"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ✅ ตั้งค่าเส้นทาง API
func (app *Application) OrderRoute(apiRouter *gin.RouterGroup) {
	order := apiRouter.Group("/order")

	order.GET("/search", app.SearchOrder)

	orderAuth := order.Group("/")
	orderAuth.Use(middleware.JWTMiddleware(app.TokenAuth))
	orderAuth.POST("/create", app.CreateBeforeReturnOrder)
	orderAuth.POST("/update-sr/:orderNo", app.UpdateSrNo)
	orderAuth.POST("/update-status/:orderNo", app.UpdateOrderStatus)

	// Frontend อัปเดต IsEdited = true เมื่อมีการแก้ไขข้อมูล (เช่น เปลี่ยน QTY, ลบรายการ, หรือเพิ่มสินค้าใหม่)
	orderAuth.PATCH("/mark-edited/:orderNo", app.MarkOrderAsEdited)
}

// SearchOrder godoc
// @Summary Search order by SO number or Order number
// @Description Retrieve the details of an order by its SO number or Order number
// @ID search-order
// @Tags Return Order MKP
// @Accept json
// @Produce json
// @Param soNo query string false "SO number"
// @Param orderNo query string false "Order number"
// @Success 200 {object} Response{data=response.SearchOrderResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /order/search [get]
func (app *Application) SearchOrder(c *gin.Context) {
	var req request.SearchOrder

	if err := c.ShouldBindQuery(&req); err != nil {
		handleResponse(c, false, "⚠️ Invalid request parameters", nil, http.StatusBadRequest)
		return
	}

	if req.SoNo == "" && req.OrderNo == "" {
		handleResponse(c, false, "⚠️ Either SoNo or OrderNo must be provided", nil, http.StatusBadRequest)
		return
	}

	order, err := app.Service.Order.SearchOrder(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleResponse(c, false, "⚠️ Order not found", nil, http.StatusNotFound)
			return
		}
		handleResponse(c, false, "🔥 Internal server error", nil, http.StatusInternalServerError)
		return
	}

	handleResponse(c, true, "⭐ Order retrieved successfully ⭐", order, http.StatusOK)
}

// CreateBeforeReturnOrder godoc
// @Summary Create a new return order
// @Description Creates a new return order including order head and order lines
// @ID create-return-order
// @Tags Return Order MKP
// @Accept json
// @Produce json
// @Param request body request.CreateBeforeReturnOrder true "Return Order Data"
// @Success 201 {object} Response{data=response.BeforeReturnOrderResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response
// @Router /order/create [post]
func (app *Application) CreateBeforeReturnOrder(c *gin.Context) {
	var req request.CreateBeforeReturnOrder

	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("⚠️ Unauthorized - Missing UserID")
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		app.Logger.Warn("❌ Invalid UserID format in token", zap.Any("UserID", userID))
		handleResponse(c, false, "❌ Unauthorized - Invalid UserID format", nil, http.StatusUnauthorized)
		return
	}

	app.Logger.Info("📝 Creating BeforeReturnOrder",
		zap.String("UserID", userIDStr),
		zap.String("OrderNo", req.OrderNo),
		zap.String("SoNo", req.SoNo),
		zap.Int("TotalItems", len(req.Items)),
	)

	resp, err := app.Service.Order.CreateBeforeReturnOrder(c.Request.Context(), req, userIDStr)
	if err != nil {
		app.Logger.Error("❌ Failed to create BeforeReturnOrder", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Return order created successfully ⭐", resp, http.StatusCreated)
}

// UpdateSrNo godoc
// @Summary Update SrNo (Sale Return Number)
// @Description Generates SrNo and updates it in the database
// @ID update-sr-no
// @Tags Return Order MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} Response{data=response.UpdateSrNoResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /order/update-sr/{orderNo} [post]
func (app *Application) UpdateSrNo(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		handleResponse(c, false, "⚠️ OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("⚠️ Unauthorized - Missing UserID")
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	resp, err := app.Service.Order.UpdateSrNo(c.Request.Context(), orderNo, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ SrNo updated successfully ⭐", resp, http.StatusOK)
}

// UpdateOrderStatus godoc
// @Summary Update order status for return confirmation
// @Description Updates order status based on RoleID (Accounting/Warehouse)
// @ID update-order-status
// @Tags Return Order MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} Response{data=response.UpdateOrderStatusResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /order/update-status/{orderNo} [post]
func (app *Application) UpdateOrderStatus(c *gin.Context) {
	orderNo := c.Param("orderNo")

	// ✅ ดึง UserID และ RoleID จาก Claims
	userID, exists := c.Get("UserID")
	if !exists {
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	roleID, exists := c.Get("RoleID")
	if !exists {
		handleResponse(c, false, "⚠️ Unauthorized - Missing RoleID", nil, http.StatusUnauthorized)
		return
	}

	// ✅ แปลง RoleID เป็น int
	roleIDInt, ok := roleID.(int)
	if !ok {
		handleResponse(c, false, "⚠️ Unauthorized - Invalid RoleID format", nil, http.StatusUnauthorized)
		return
	}

	app.Logger.Info("🔄 Updating Order Status...",
		zap.String("OrderNo", orderNo),
		zap.String("RequestedBy", userID.(string)),
		zap.Int("RoleID", roleID.(int)),
	)

	resp, err := app.Service.Order.UpdateOrderStatus(c.Request.Context(), orderNo, userID.(string), roleIDInt)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Order status updated successfully ⭐", resp, http.StatusOK)
}

// MarkOrderAsEdited godoc
// @Summary Mark order as edited
// @Description Marks the order as edited when there are modifications
// @ID mark-order-as-edited
// @Tags Return Order MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /order/mark-edited/{orderNo} [patch]
func (app *Application) MarkOrderAsEdited(c *gin.Context) {
	orderNo := c.Param("orderNo")

	userID, exists := c.Get("UserID")
	if !exists {
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	err := app.Service.Order.MarkOrderAsEdited(c.Request.Context(), orderNo, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Order marked as edited ⭐", nil, http.StatusOK)
}
