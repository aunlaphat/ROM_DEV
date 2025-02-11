package api

// ตัวกลางของ http request ที่คอยรับส่งข้อมูลไปมา
// รับมาในรูป request ส่งออกในรูป response
// ส่งคำขอไปยัง service เพื่อขอให้ดึงข้อมูลที่ต้องการออกมา แต่ service จะทำการ validation ก่อนเพื่อตรวจข้อผิดพลาดก่อนจะดึง query ออกมาจาก repo ให้
// api handle error about res req send error to client ex. 400 500 401
// service handle error about business logic/ relation data send error to api
// ตรวจสอบส่วนหน้า ส่วนที่รับมาจาก client เช่น input ที่ถูกป้อนเข้ามา
import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

// ReturnOrder => ข้อมูลรับเข้าจากส่วนหน้าคลังทั้งหมด ข้อมูลสินค้าที่ถูกส่งคืนมาทั้งหมด
func (app *Application) ReturnOrder(apiRouter *gin.RouterGroup) {
	api := apiRouter.Group("/return-order")
	
	// ใช้ JWT Middleware ของ Gin แทน Go-Chi
	api.Use(func(c *gin.Context) {
		_, claims, err := jwtauth.FromContext(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	})

	api.GET("/get-all", app.GetAllReturnOrder)                      // แสดงข้อมูลรับเข้ารวม                  
	api.GET("/get-all/:orderNo", app.GetReturnOrderByOrderNo)       // แสดงข้อมูลรับเข้าด้วย orderNo    
	api.GET("/get-lines", app.GetAllReturnOrderLines)               // แสดงรายการคืนของรวม
	api.GET("/get-lines/:orderNo", app.GetReturnOrderLineByOrderNo) // แสดงรายการคืนของโดย orderNo
	api.POST("/create", app.CreateReturnOrder)                      // สร้างข้อมูลของที่ถูกส่งคืนมา
	api.PATCH("/update/:orderNo", app.UpdateReturnOrder)            // อัพเดทข้อมูลของที่ถูกส่งคืน
	api.DELETE("/delete/:orderNo", app.DeleteReturnOrder)           // ลบ order ที่ทำการคืนมาออกหมด head+line
}

// review
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
		handleError(c, err)
		return
	}

	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n======== Order #%d ========\n", i+1)
		utils.PrintReturnOrderDetails(&order)
		for j, line := range order.ReturnOrderLine {
			fmt.Printf("\n======== Order Line #%d ========\n", j+1)
			utils.PrintReturnOrderLineDetails(&line)
		}
		fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(order.ReturnOrderLine))
		fmt.Println("=====================================")
	}

	// ส่งข้อมูลกลับไปยัง Client
	handleResponse(c, true, "⭐ Get Return Order successfully ⭐", result, http.StatusOK)
}

// review
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
	orderNo := chi.URLParam(r, "orderNo")

	result, err := app.Service.ReturnOrder.GetReturnOrderByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Return Order by OrderNo Details ========== 📋\n\n")
	utils.PrintReturnOrderDetails(result)
	for i, line := range result.ReturnOrderLine {
		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
		utils.PrintReturnOrderLineDetails(&line)
	}
	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result.ReturnOrderLine))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Get Return Order by OrderNo successfully ⭐", result, http.StatusOK)
}

// review
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
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Return Order Lines (%d) ========== 📋\n", len(result))
	for i, line := range result {
		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
		utils.PrintReturnOrderLineDetails(&line)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Get Return Order Lines successfully ⭐", result, http.StatusOK)
}

// review
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
	orderNo := chi.URLParam(r, "orderNo")

	result, err := app.Service.ReturnOrder.GetReturnOrderLineByOrderNo(c.Request.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Return Order Line of OrderNo: %s ========== 📋\n", orderNo)
	for i, line := range result {
		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
		utils.PrintReturnOrderLineDetails(&line)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Get Return Order Line by OrderNo successfully ⭐", result, http.StatusOK)

}

// review
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
	// 1. Authentication check
	_, claims, err := jwtauth.FromContext(c.Request.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	// Step 1: Decode JSON Payload เป็นโครงสร้าง CreateReturnOrder
	var req request.CreateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// พบข้อผิดพลาดในการ Decode ให้ส่งข้อผิดพลาดไปยัง Client
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// Set user information from claims
	req.CreateBy = userID

	result, err := app.Service.ReturnOrder.CreateReturnOrder(c.Request.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Created Return Order ========== 📋\n\n")
	utils.PrintCreateReturnOrder(result)
	fmt.Printf("\n📋 ========== Return Order Line Details ========== 📋\n")
	for i, line := range result.ReturnOrderLine {
		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
		utils.PrintReturnOrderLineDetails(&line)
	}
	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result.ReturnOrderLine))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Created successfully ⭐", result, http.StatusOK)
}

// review
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
	orderNo := chi.URLParam(r, "orderNo")

	// Decode JSON Payload เป็นโครงสร้าง UpdateReturnOrder
	var req request.UpdateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// Set OrderNo ที่ได้จาก URL ลงในโครงสร้างข้อมูล
	req.OrderNo = orderNo

	// ดึง userID จาก JWT token
	_, claims, err := jwtauth.FromContext(c.Request.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	result, err := app.Service.ReturnOrder.UpdateReturnOrder(c.Request.Context(), req, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Updated Order ========== 📋\n")
	utils.PrintUpdateReturnOrder(result)
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Updated successfully ⭐", result, http.StatusOK)
}

// review
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
	orderNo := chi.URLParam(r, "orderNo")

	err := app.Service.ReturnOrder.DeleteReturnOrder(c.Request.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	response := response.DeleteReturnOrder{
		OrderNo: orderNo,
	}

	handleResponse(w, true, "⭐ Deleted successfully ⭐", response, http.StatusOK)
}
