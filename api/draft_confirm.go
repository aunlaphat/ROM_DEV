package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *Application) DraftConfirmRoute(apiRouter *gin.RouterGroup) {
	draft := apiRouter.Group("/draft-confirm")
	draft.Use(middleware.JWTMiddleware(app.TokenAuth)) // ใช้ JWT Middleware

	draft.GET("/", app.GetOrders)
	draft.GET("/:orderNo", app.GetOrderWithItems)
	draft.POST("/add-item/:orderNo", app.AddItemToDraftOrder)
	draft.DELETE("/remove-item/:orderNo/:sku", app.RemoveItemFromDraftOrder)
	draft.POST("/update-status/:orderNo", app.ConfirmDraftOrder)
}

// GetOrders godoc
// @Summary Get Draft or Confirm Orders
// @Description Retrieves all orders filtered by StatusConfID and Date Range
// @ID get-orders
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param statusConfID query int true "StatusConfID (1 = Draft, 2 = Confirm)"
// @Param startDate query string false "Start Date (YYYY-MM-DD)"
// @Param endDate query string false "End Date (YYYY-MM-DD)"
// @Success 200 {object} Response{data=[]response.OrderHeadResponse}
// @Failure 404 {object} Response
// @Router /draft-confirm [get]
func (app *Application) GetOrders(c *gin.Context) {
	statusConfID, err := strconv.Atoi(c.Query("statusConfID"))
	if err != nil {
		handleResponse(c, false, "⚠️ Invalid StatusConfID", nil, http.StatusBadRequest)
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	orders, err := app.Service.DraftConfirm.GetOrders(c.Request.Context(), statusConfID, startDate, endDate)
	if err != nil {
		handleResponse(c, false, "⚠️ No orders found", nil, http.StatusNotFound)
		return
	}

	handleResponse(c, true, "⭐ Orders retrieved successfully ⭐", orders, http.StatusOK)
}

// GetOrderWithItems godoc
// @Summary Get Order Details (Draft & Confirm)
// @Description Retrieves details of an order including items
// @ID get-order-with-items
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} Response{data=response.DraftConfirmResponse}
// @Failure 404 {object} Response
// @Router /draft-confirm/{orderNo} [get]
func (app *Application) GetOrderWithItems(c *gin.Context) {
	orderNo := c.Param("orderNo")
	order, err := app.Service.DraftConfirm.GetOrderWithItems(c.Request.Context(), orderNo)
	if err != nil {
		handleResponse(c, false, "⚠️ Order not found", nil, http.StatusNotFound)
		return
	}
	handleResponse(c, true, "⭐ Order retrieved successfully ⭐", order, http.StatusOK)
}

// AddItemToDraftOrder godoc
// @Summary Add Item to Draft Order
// @Description Adds an item to an existing draft order
// @ID add-item-to-draft
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Param request body request.AddItem true "Item Data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /draft-confirm/add-item/{orderNo} [post]
func (app *Application) AddItemToDraftOrder(c *gin.Context) {
	var req request.AddItem
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}
	orderNo := c.Param("orderNo")
	userID := c.MustGet("UserID").(string)

	err := app.Service.DraftConfirm.AddItemToDraftOrder(c.Request.Context(), orderNo, req, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	handleResponse(c, true, "⭐ Item added successfully ⭐", nil, http.StatusOK)
}

// RemoveItemFromDraftOrder godoc
// @Summary Remove Item from Draft Order
// @Description Removes an item from a draft order
// @ID remove-item-from-draft
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Param sku path string true "SKU"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /draft-confirm/remove-item/{orderNo}/{sku} [delete]
func (app *Application) RemoveItemFromDraftOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	sku := c.Param("sku")

	err := app.Service.DraftConfirm.RemoveItemFromDraftOrder(c.Request.Context(), orderNo, sku)
	if err != nil {
		handleError(c, err)
		return
	}
	handleResponse(c, true, "⭐ Item removed successfully ⭐", nil, http.StatusOK)
}

// ConfirmDraftOrder godoc
// @Summary Confirm Draft Order
// @Description Updates a draft order to confirm status
// @ID confirm-draft-order
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /update-status/{orderNo} [post]
func (app *Application) ConfirmDraftOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	userID := c.MustGet("UserID").(string) // ✅ ดึง UserID จาก JWT

	err := app.Service.DraftConfirm.ConfirmDraftOrder(c.Request.Context(), orderNo, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	handleResponse(c, true, "⭐ Order confirmed successfully ⭐", nil, http.StatusOK)
}
