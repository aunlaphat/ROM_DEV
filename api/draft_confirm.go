package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) DraftConfirmRoute(apiRouter *gin.RouterGroup) {
	draft := apiRouter.Group("/draft-confirm")

	// ✅ ไม่ต้องใช้ JWT (Public API)
	draft.GET("/list-codeR", app.ListCodeR) // ดึงข้อมูล CodeR

	// ✅ ใช้ JWT Middleware (ต้อง Auth)
	draft.Use(middleware.JWTMiddleware(app.TokenAuth))

	// ✅ ดึงรายการ Draft & Confirm Orders
	draft.GET("/orders", app.GetOrders)

	// ✅ ดึงรายละเอียดออเดอร์ + รายการสินค้า
	draft.GET("/order/details", app.GetOrderWithItems)

	// ✅ จัดการสินค้าใน Draft Order
	draft.POST("/add-item/:orderNo", app.AddItemToDraftOrder)                // เพิ่มสินค้าเข้า Draft Order
	draft.DELETE("/remove-item/:orderNo/:sku", app.RemoveItemFromDraftOrder) // ลบสินค้าออกจาก Draft Order

	// ✅ อัปเดตสถานะจาก Draft → Confirm
	draft.PATCH("/update-status/:orderNo", app.ConfirmDraftOrder)
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
// @Router /draft-confirm/orders [get]
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
// @Summary Get Order Details with Items
// @Description Retrieves details of an order including items
// @ID get-order-with-items
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Param statusConfID query int true "StatusConfID (1 = Draft, 2 = Confirm)"
// @Param orderNo query string true "Order Number"
// @Success 200 {object} Response{data=response.DraftConfirmResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /draft-confirm/order/details [get]
func (app *Application) GetOrderWithItems(c *gin.Context) {
	// ✅ รับค่า `statusConfID` และ `orderNo` ผ่าน Query Parameter
	statusConfID, err := strconv.Atoi(c.Query("statusConfID"))
	if err != nil || (statusConfID != 1 && statusConfID != 2) {
		handleResponse(c, false, "⚠️ Invalid statusConfID (must be 1 or 2)", nil, http.StatusBadRequest)
		return
	}

	orderNo := c.Query("orderNo")
	if orderNo == "" {
		handleResponse(c, false, "⚠️ orderNo is required", nil, http.StatusBadRequest)
		return
	}

	order, err := app.Service.DraftConfirm.GetOrderWithItems(c.Request.Context(), orderNo, statusConfID)
	if err != nil {
		handleResponse(c, false, "⚠️ Order not found", nil, http.StatusNotFound)
		return
	}

	handleResponse(c, true, "⭐ Order retrieved successfully ⭐", order, http.StatusOK)
}

// GetListofCodeR godoc
// @Summary Get List of CodeR
// @Description Retrieves a list of CodeR (SKU starting with 'R')
// @ID get-list-codeR
// @Tags Draft & Confirm MKP
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]response.ListCodeRResponse}
// @Failure 500 {object} Response
// @Router /draft-confirm/list-codeR [get]
func (app *Application) ListCodeR(c *gin.Context) {
	codeRList, err := app.Service.DraftConfirm.ListCodeR(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	handleResponse(c, true, "⭐ CodeR list retrieved successfully ⭐", codeRList, http.StatusOK)
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
// @Success 200 {object} Response{data=[]response.AddItemResponse}
// @Failure 400 {object} Response
// @Router /draft-confirm/add-item/{orderNo} [post]
func (app *Application) AddItemToDraftOrder(c *gin.Context) {
	var req request.AddItem
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err)
		return
	}

	// ✅ ดึง `UserID` จาก JWT Middleware
	userID := c.MustGet("UserID").(string)
	req.OrderNo = c.Param("orderNo") // ✅ ใช้ `OrderNo` จาก Path Parameter

	// ✅ ส่งไปที่ Service Layer
	results, err := app.Service.DraftConfirm.AddItemToDraftOrder(c.Request.Context(), req, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Item added successfully ⭐", results, http.StatusOK)
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
// @Success 200 {object} Response{data=response.UpdateOrderStatusResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /draft-confirm/update-status/{orderNo} [patch]
func (app *Application) ConfirmDraftOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	userID := c.MustGet("UserID").(string) // ✅ ดึง UserID จาก JWT Middleware

	app.Logger.Info("📦 Request to Confirm Draft Order",
		zap.String("OrderNo", orderNo),
		zap.String("UserID", userID),
	)

	// ✅ เรียก Service เพื่อเปลี่ยนสถานะจาก Draft → Confirm
	updateResponse, err := app.Service.DraftConfirm.ConfirmDraftOrder(c.Request.Context(), orderNo, userID)
	if err != nil {
		app.Logger.Error("❌ Failed to confirm draft order", zap.Error(err))
		handleError(c, err)
		return
	}

	handleResponse(c, true, "⭐ Order confirmed successfully ⭐", updateResponse, http.StatusOK)
}
