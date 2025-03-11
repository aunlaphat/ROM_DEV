package api

import (
	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/dto/response"
	Status "boilerplate-back-go-2411/errors"
	"boilerplate-back-go-2411/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// *️⃣ ReturnOrder => ข้อมูลรับเข้าจากส่วนหน้าคลังทั้งหมด ข้อมูลสินค้าที่ถูกส่งคืนมาทั้งหมด
func (app *Application) ReturnOrder(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/return-order")
	api.GET("/get-all", app.GetAllReturnOrder)                      // แสดงข้อมูลรับเข้ารวม
	api.GET("/get-all/:orderNo", app.GetReturnOrderByOrderNo)       // แสดงข้อมูลรับเข้าด้วย orderNo
	api.GET("/get-lines", app.GetAllReturnOrderLines)               // แสดงรายการคืนของรวม
	api.GET("/get-lines/:orderNo", app.GetReturnOrderLineByOrderNo) // แสดงรายการคืนของโดย orderNo

	apiAuth := api.Group("/")
	apiAuth.Use(middleware.JWTMiddleware(app.TokenAuth))
	apiAuth.POST("/create", app.CreateReturnOrder)            // สร้างข้อมูลของที่ถูกส่งคืนมา
	apiAuth.PATCH("/update/:orderNo", app.UpdateReturnOrder)  // อัพเดทข้อมูลของที่ถูกส่งคืน
	apiAuth.DELETE("/delete/:orderNo", app.DeleteReturnOrder) // ลบ order ที่ทำการคืนมาออกหมด head+line
}

// @Summary 	Get Return Order
// @Description Retrieve the details of Return Order
// @ID 			GetAll-ReturnOrder
// @Tags 		Return Order
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]response.ReturnOrder} "Get All"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Not Found Endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/return-order/get-all [get]
func (app *Application) GetAllReturnOrder(c *gin.Context) {
	result, err := app.Service.ReturnOrder.GetAllReturnOrder(c.Request.Context())
	if err != nil {
		app.Logger.Error("[ Error fetching return order ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Get Return Order successfully ]", result, http.StatusOK)
}

// @Summary      Get Return Order by OrderNo
// @Description  Get details return order by order no
// @ID           GetAllByOrderNo-ReturnOrder
// @Tags         Return Order
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Order No"
// @Success      200 	  {object} Response{result=[]response.ReturnOrder} "Get All by OrderNo"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "Not Found Endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /return-order/get-all/{orderNo} [get]
func (app *Application) GetReturnOrderByOrderNo(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Return Order by OrderNo successfully ]", result, http.StatusOK)
}

// @Summary 	Get Return Order Line
// @Description Get all Return Order Line
// @ID 			GetAllLines-ReturnOrderLine
// @Tags 		Return Order
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]response.ReturnOrderLine} "Get Return Order Lines"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Not Found Endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/return-order/get-lines [get]
func (app *Application) GetAllReturnOrderLines(c *gin.Context) {
	result, err := app.Service.ReturnOrder.GetAllReturnOrderLines(c.Request.Context())
	if err != nil {
		app.Logger.Error("[ Error fetching return order lines ]", zap.Error(err))
		handleError(c, err)
		return
	}

	if len(result) == 0 {
		app.Logger.Info("[ No data found ]")
		handleResponse(c, true, "[ No data found ]", nil, http.StatusOK)
		return
	}

	handleResponse(c, true, "[ Get Return Order Lines successfully ]", result, http.StatusOK)
}

// @Summary      Get Return Order Line by OrderNo
// @Description  Get details of an order line by its order no
// @ID           GetLineByOrderNo-ReturnOrder
// @Tags         Return Order
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Order No"
// @Success      200 	  {object} Response{result=[]response.ReturnOrderLine} "Get Lines by OrderNo"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "Not Found Endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /return-order/get-lines/{orderNo} [get]
func (app *Application) GetReturnOrderLineByOrderNo(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderLineByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Get Return Order Line by OrderNo successfully ]", result, http.StatusOK)

}

// @Summary 	Create Return Order
// @Description Create a new return order
// @ID 			Create-ReturnOrder
// @Tags 		Return Order
// @Accept 		json
// @Produce 	json
// @Param 		CreateReturnOrder body request.CreateReturnOrder true "Return Order"
// @Success 	200 {object} Response{result=[]response.CreateReturnOrder} "Return Order Created"
// @Success 	201 {object} Response{result=[]response.CreateReturnOrder} "Return Order Created"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/return-order/create [post]
func (app *Application) CreateReturnOrder(c *gin.Context) {
	var req request.CreateReturnOrder

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
		handleError(c, Status.UnauthorizedError("[ Unauthorized - Missing UserID ]"))
		return
	}

	req.CreateBy = userID.(string)
	result, err := app.Service.ReturnOrder.CreateReturnOrder(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Created successfully ]", result, http.StatusOK)
}

// @Summary Update Return Order
// @Description Update an existing return order using orderNo in the path
// @ID Update-ReturnOrder
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order No"
// @Param order body request.UpdateReturnOrder true "Updated Order Data"
// @Success 200 {object} Response{result=[]response.UpdateReturnOrder} "Return Order Update"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Order Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /return-order/update/{orderNo} [patch]
func (app *Application) UpdateReturnOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	var req request.UpdateReturnOrder

	if req.OrderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

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
		handleError(c, Status.UnauthorizedError("[ Unauthorized - Missing UserID ]"))
		return
	}

	req.OrderNo = orderNo

	// *️⃣ ตรวจสอบ nil ก่อนใช้ UpdateBy
	if req.UpdateBy == nil {
		req.UpdateBy = new(string)
	}
	*req.UpdateBy = userID.(string)

	result, err := app.Service.ReturnOrder.UpdateReturnOrder(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	handleResponse(c, true, "[ Updated successfully ]", result, http.StatusOK)
}

// @Summary 	Delete Order
// @Description Delete an order
// @ID 			delete-ReturnOrder
// @Tags 		Return Order
// @Accept 		json
// @Produce 	json
// @Param 		orderNo path string true "Order No"
// @Success 	200 {object} Response{result=[]response.DeleteReturnOrder} "ReturnOrder Deleted"
// @Success 	204 {object} Response "No Content, Order Delete Successfully"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/return-order/delete/{orderNo} [delete]
func (app *Application) DeleteReturnOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")

	if orderNo == "" {
		app.Logger.Warn("[ OrderNo is required ]")
		handleError(c, Status.BadRequestError("[ OrderNo is required ]"))
		return
	}

	err := app.Service.ReturnOrder.DeleteReturnOrder(c.Request.Context(), orderNo)
	if err != nil {
		handleError(c, err)
		return
	}

	result := response.DeleteReturnOrder{
		OrderNo: orderNo,
	}

	handleResponse(c, true, "[ Deleted successfully ]", result, http.StatusOK)
}
