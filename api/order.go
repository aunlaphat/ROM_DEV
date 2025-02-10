package api

import (
	"boilerplate-backend-go/dto/request"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) OrderRoute(apiRouter *gin.RouterGroup) {
	order := apiRouter.Group("/order")
	order.GET("/search", app.SearchOrder)
	order.POST("/create", app.CreateBeforeReturnOrder)
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
// @Failure 500 {object} Response
// @Router /order/create [post]
func (app *Application) CreateBeforeReturnOrder(c *gin.Context) {
	var req request.CreateBeforeReturnOrder

	// 🔹 ตรวจสอบ JSON Request
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	// 🔹 ดึง UserID จาก JWT (สมมติว่าเก็บใน Cookie)
	userID, exists := c.Get("UserID")
	if !exists {
		handleResponse(c, false, "⚠️ Unauthorized - Missing UserID", nil, http.StatusUnauthorized)
		return
	}

	// 🔹 เรียก Service Layer
	resp, err := app.Service.Order.CreateBeforeReturnOrder(c.Request.Context(), req, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Return order created successfully ⭐", resp, http.StatusCreated)
}
