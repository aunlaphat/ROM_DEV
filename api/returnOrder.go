package api

// ตัวกลางของ http request ที่คอยรับส่งข้อมูลไปมา
// รับมาในรูป request ส่งออกในรูป response
// ส่งคำขอไปยัง service เพื่อขอให้ดึงข้อมูลที่ต้องการออกมา แต่ service จะทำการ validation ก่อนเพื่อตรวจข้อผิดพลาดก่อนจะดึง query ออกมาจาก repo ให้
// api handle error about res req send error to client ex. 400 500 401
// service handle error about business logic/ relation data send error to api
// ตรวจสอบส่วนหน้า ส่วนที่รับมาจาก client เช่น input ที่ถูกป้อนเข้ามา
import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) ReturnOrders(apiRouter *chi.Mux) {
	/* apiRouter.Route("/return-order", func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/get-all", app.GetAllReturnOrder)                         // GET /return-order/get-all
		r.Get("/get-all/{orderNo}", app.GetReturnOrderByOrderNo)         // GET /return-order/getbyID/{orderNo}
		r.Get("/get-lines", app.GetAllReturnOrderLines)                  // GET /return-order/allgetline
		r.Get("/get-lines/{orderNo}", app.GetReturnOrderLinesByReturnID) // GET /return-order/getlinebyID/{orderNo}
		r.Post("/create", app.CreateReturnOrder)                         // POST /return-order/create
		r.Patch("/update/{orderNo}", app.UpdateReturnOrder)              // PATCH /return-order/update/{orderNo}
		r.Delete("/delete/{orderNo}", app.DeleteReturnOrder)             // DELETE /return-order/delete/{orderNo}
	}) */
}

/*
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
func (api *Application) GetAllReturnOrder(w http.ResponseWriter, r *http.Request) {

	result, err := api.Service.ReturnOrder.GetAllReturnOrder(r.Context())
	// เช็คเมื่อเกิดข้อผิดพลาดขึ้น
	if err != nil {
		handleError(w, err)
		return
	}
	// เช็คเมื่อไม่มีข้อมูลในคำสั่งซื้อ
	if len(result) == 0 {
		handleResponse(w, true, "No return orders found", []response.ReturnOrder{}, http.StatusOK)
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
	handleResponse(w, true, "⭐ Get Return Order successfully ⭐", result, http.StatusOK)
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
func (app *Application) GetReturnOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	// เช็คค่าว่าง
	if orderNo == "" {
		handleError(w, errors.ValidationError("OrderNo is required"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}
	// เช็คเมื่อไม่มีข้อมูลในคำสั่งซื้อ
	if len(result.ReturnOrderLine) == 0 {
		handleResponse(w, true, "No lines found for this return order", result, http.StatusOK)
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
func (app *Application) GetAllReturnOrderLines(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.ReturnOrder.GetAllReturnOrderLines(r.Context())
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

// @Summary      Get Return Order Line by OrderNo
// @Description  Get details of an order line by its order no
// @ID           GetLineByID-ReturnOrder
// @Tags         Return Order
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Order No"
// @Success      200 	  {object} Response{result=[]response.ReturnOrderLine} "Get Lines by OrderNo"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "Not Found Endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /return-order/get-lines/{orderNo} [get]
func (app *Application) GetReturnOrderLinesByReturnID(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, errors.ValidationError("OrderNo is required"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderLinesByReturnID(r.Context(), orderNo)
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
func (app *Application) CreateReturnOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Authentication check
	_, claims, err := jwtauth.FromContext(r.Context())
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
	// ตรวจสอบว่า ReturnOrderLine มีข้อมูลอย่างน้อย 1 รายการ
	if len(req.ReturnOrderLine) == 0 {
		handleError(w, errors.ValidationError("ReturnOrderLine cannot be empty"))
		return
	}

	// Set user information from claims
	req.CreateBy = userID

	result, err := app.Service.ReturnOrder.CreateReturnOrder(r.Context(), req)
	if err != nil {
		app.Logger.Error("Failed to create return",
			zap.Error(err),
			zap.String("orderNo", req.OrderNo))

		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			handleResponse(w, false, err.Error(), nil, http.StatusBadRequest)
		case strings.Contains(err.Error(), "already exists"):
			handleResponse(w, false, err.Error(), nil, http.StatusConflict)
		default:
			handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
		}
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
func (app *Application) UpdateReturnOrder(w http.ResponseWriter, r *http.Request) {

	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, errors.ValidationError("OrderNo is required in the path"))
		return
	}

	// Decode JSON Payload เป็นโครงสร้าง UpdateReturnOrder
	var req request.UpdateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// ระบุ OrderNo ที่ได้จาก URL ลงในโครงสร้างข้อมูล
	req.OrderNo = orderNo

	// ดึง userID จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	// Step 6: เรียก Service เพื่ออัปเดตข้อมูล Return Order
	result, err := app.Service.ReturnOrder.UpdateReturnOrder(r.Context(), req, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Updated Order ========== 📋\n")
	utils.PrintUpdateReturnOrder(result)
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Updated successfully ⭐", result, http.StatusOK)
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
func (api *Application) DeleteReturnOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: ดึง OrderNo จาก URL Parameter
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		// Step 2: ตรวจสอบว่า OrderNo ว่างหรือไม่
		handleError(w, errors.ValidationError("OrderNo is required in the path"))
		return
	}

	// Step 3: เรียก Service เพื่อทำการลบข้อมูล Return Order
	err := api.Service.ReturnOrder.DeleteReturnOrder(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	response := response.DeleteReturnOrder{
		OrderNo: orderNo,
	}

	handleResponse(w, true, "⭐ Deleted successfully ⭐", response, http.StatusOK)
}
*/
