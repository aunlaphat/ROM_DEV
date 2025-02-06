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

	// ✅ Bind Query Parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		handleResponse(c, false, "⚠️ Invalid request parameters", nil, http.StatusBadRequest)
		return
	}

	// ✅ Validate required parameters
	if req.SoNo == "" && req.OrderNo == "" {
		handleResponse(c, false, "⚠️ Either SoNo or OrderNo must be provided", nil, http.StatusBadRequest)
		return
	}

	// 🛠 Call Service Layer (Logging will be handled there)
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
