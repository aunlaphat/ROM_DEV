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
}

// 📌 **Search Order**
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

	// ✅ ตรวจสอบค่า query string
	if err := c.ShouldBindQuery(&req); err != nil {
		handleResponse(c, false, "⚠️ Invalid request parameters", nil, http.StatusBadRequest)
		return
	}

	if req.SoNo == "" && req.OrderNo == "" {
		handleResponse(c, false, "⚠️ Either SoNo or OrderNo must be provided", nil, http.StatusBadRequest)
		return
	}

	// ✅ ค้นหา Order
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

// 📌 **Create Before Return Order**
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

	// ✅ ตรวจสอบ JSON Request
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	// ✅ ดึง UserID จาก JWT Middleware
	userID, exists := c.Get("UserID")
	if !exists {
		app.Logger.Warn("⚠️ Unauthorized - Missing UserID")
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	// ✅ ตรวจสอบว่า UserID เป็น `string`
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

	// ✅ เรียก Service Layer เพื่อสร้าง BeforeReturnOrder
	resp, err := app.Service.Order.CreateBeforeReturnOrder(c.Request.Context(), req, userIDStr)
	if err != nil {
		app.Logger.Error("❌ Failed to create BeforeReturnOrder", zap.Error(err))
		handleError(c, err)
		return
	}

	app.Logger.Info("✅ BeforeReturnOrder created successfully",
		zap.String("OrderNo", resp.OrderNo),
		zap.String("SoNo", resp.SoNo),
		zap.Int("TotalItems", len(resp.Items)),
	)

	handleResponse(c, true, "⭐ Return order created successfully ⭐", resp, http.StatusCreated)
}
