package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) OrderRoute(apiRouter *gin.RouterGroup) {
	order := apiRouter.Group("/order")
	order.GET("/search", app.SearchOrder)
}

// SearchSaleOrder godoc
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
	soNo := c.Query("soNo")
	orderNo := c.Query("orderNo")

	if soNo == "" && orderNo == "" {
		handleResponse(c, false, "⚠️ Either SoNo or OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	order, err := app.Service.Order.SearchOrder(c, soNo, orderNo)
	if err != nil {
		if err.Error() == "Sale order not found" {
			handleResponse(c, false, "⚠️ Sale order not found", nil, http.StatusNotFound)
			return
		}
		handleResponse(c, false, "🔥 Internal server error", nil, http.StatusInternalServerError)
		return
	}

	handleResponse(c, true, "⭐ Order retrieved successfully ⭐", order, http.StatusOK)
}
