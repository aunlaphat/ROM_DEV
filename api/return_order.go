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
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) ReturnOrders(apiRouter *chi.Mux) {
	apiRouter.Route("/reorder", func(r chi.Router) {
		r.Get("/allget", app.AllGetReturnOrder)                            // GET /reorder/allget
		r.Get("/getbyID/{orderNo}", app.GetReturnOrderID)                  // GET /reorder/getbyID/{orderNo}
		r.Get("/allgetline", app.GetAllReturnOrderLines)                   // GET /reorder/allgetline
		r.Get("/getlinebyID/{orderNo}", app.GetReturnOrderLinesByReturnID) // GET /reorder/getlinebyID/{orderNo}
		r.Post("/create", app.CreateReturnOrder)                           // POST /reorder/create
		r.Patch("/update/{orderNo}", app.UpdateReturnOrder)                // PUT /reorder/update/{orderNo}
		r.Delete("/delete/{orderNo}", app.DeleteReturnOrder)               // DELETE /reorder/delete/{orderNo}
	})
}

// @Summary 	Get Return Order
// @Description Get all Return Order
// @ID 			Allget-ReturnOrder
// @Tags 		ReturnOrder
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]entity.ReturnOrder} "Get All"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "not found endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/reorder/allget [get]
func (api *Application) AllGetReturnOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: เรียก Service เพื่อดึงข้อมูล Return Order ทั้งหมด
	result, err := api.Service.ReturnOrder.AllGetReturnOrder(r.Context())
	if err != nil {
		// Step 2: หากพบข้อผิดพลาด ให้ส่งข้อผิดพลาดไปยัง Client
		handleError(w, err)
		return
	}

	// Step 2.1: หากไม่มีข้อมูล Return Orders ให้ส่ง [] กลับไปพร้อมข้อความ
	if len(result) == 0 {
		handleResponse(w, true, "No return orders found", []response.ReturnOrder{}, http.StatusOK)
		return
	}

	// Step 3: หากสำเร็จ ให้ส่งข้อมูลกลับไปยัง Client
	handleResponse(w, true, "Get Order successfully", result, http.StatusOK)
}

// @Summary      Get Return Order by ID
// @Description  Get details of an order by its return id
// @ID           GetByID-ReturnOrder
// @Tags         ReturnOrder
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Return ID"
// @Success      200 	  {object} Response{result=[]entity.ReturnOrder} "Get by ID"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "not found endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /reorder/getbyID/{orderNo} [get]
func (app *Application) GetReturnOrderID(w http.ResponseWriter, r *http.Request) {
	// Step 1: ดึง OrderNo จาก URL Parameter
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		// Step 2: ตรวจสอบว่า OrderNo ว่างหรือไม่
		handleError(w, errors.ValidationError("OrderNo is required"))
		return
	}

	// Step 3: เรียก Service เพื่อดึงข้อมูล Return Order ที่ตรงกับ OrderNo
	result, err := app.Service.ReturnOrder.GetReturnOrderByID(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	// Step 3.1: หากไม่มี ReturnOrderLines ในคำสั่งซื้อ ให้ส่งข้อความแจ้ง
	if len(result.ReturnOrderLine) == 0 {
		handleResponse(w, true, "No lines found for this return order", result, http.StatusOK)
		return
	}

	handleResponse(w, true, "Get Order by ID successfully", result, http.StatusOK)
}

// @Summary 	Get Return Order Line
// @Description Get all Return Order Line
// @ID 			Allget-ReturnOrderLine
// @Tags 		ReturnOrder
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]entity.ReturnOrderLine} "Get Order Line All"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "not found endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/reorder/allgetline [get]
func (app *Application) GetAllReturnOrderLines(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.ReturnOrder.GetAllReturnOrderLines(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Get Order Line successfully", result, http.StatusOK)
}

// @Summary      Get Return Order Line by ID
// @Description  Get details of an order line by its return id
// @ID           GetLineByID-ReturnOrder
// @Tags         ReturnOrder
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Return ID"
// @Success      200 	  {object} Response{result=[]entity.ReturnOrderLine} "Get by ID"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "not found endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /reorder/getlinebyID/{orderNo} [get]
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

	handleResponse(w, true, "Get Order Line by ID successfully", result, http.StatusOK)

}

// @Summary 	Create Order
// @Description Create a new order
// @ID 			Create-ReturnOrder
// @Tags 		ReturnOrder
// @Accept 		json
// @Produce 	json
// @Param 		CreateReturnOrder body request.CreateReturnOrder true "ReturnOrder Data"
// @Success 	201 {object} Response{result=[]request.CreateReturnOrder} "ReturnOrder Created"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/reorder/create [post]
func (api *Application) CreateReturnOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode JSON Payload เป็นโครงสร้าง CreateReturnOrder
	var req request.CreateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Step 2: หากพบข้อผิดพลาดในการ Decode ให้ส่งข้อผิดพลาดไปยัง Client
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// Step 3: ตรวจสอบฟิลด์ที่จำเป็น OrderNo และ OrderNo ว่าว่างหรือไม่
	if req.OrderNo == "" {
		handleError(w, errors.BadRequestError("OrderNo are required"))
		return
	}

	// Step 4: ตรวจสอบค่าจำนวนสินค้าใน ReturnOrderLine
	for _, line := range req.ReturnOrderLine {
		if line.ReturnQTY <= 0 {
			handleError(w, errors.BadRequestError("Invalid ReturnQTY in ReturnOrderLine"))
			return
		}
	}

	// ตรวจสอบว่า ReturnOrderLine มีข้อมูลอย่างน้อย 1 รายการ
	if len(req.ReturnOrderLine) == 0 {
		handleError(w, errors.ValidationError("ReturnOrderLine cannot be empty"))
		return
	}

	// Step 5: เรียก Service เพื่อสร้างข้อมูล Return Order
	err := api.Service.ReturnOrder.CreateReturnOrder(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Return order created successfully", nil, http.StatusCreated)
}

// @Summary Update Order
// @Description Update an existing return order using orderNo in the path
// @ID Update-ReturnOrder
// @Tags ReturnOrder
// @Accept json
// @Produce json
// @Param orderNo path string true "Return ID"
// @Param order body request.UpdateReturnOrder true "Updated Order Data"
// @Success 200 {object} Response "ReturnOrder Updated Successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Order Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /reorder/update/{orderNo} [patch]
func (api *Application) UpdateReturnOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: ดึง OrderNo จาก URL Parameter
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		// Step 2: ตรวจสอบว่า OrderNo ว่างหรือไม่
		handleError(w, errors.ValidationError("OrderNo is required in the path"))
		return
	}

	// Step 3: Decode JSON Payload เป็นโครงสร้าง UpdateReturnOrder
	var req request.UpdateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// Step 5: ระบุ OrderNo ที่ได้จาก URL ลงในโครงสร้างข้อมูล
	req.OrderNo = orderNo

	// Step 6: เรียก Service เพื่ออัปเดตข้อมูล Return Order
	err := api.Service.ReturnOrder.UpdateReturnOrder(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Return order updated successfully", nil, http.StatusOK)
}

// @Summary 	Delete Order
// @Description Delete an order
// @ID 			delete-ReturnOrder
// @Tags 		ReturnOrder
// @Accept 		json
// @Produce 	json
// @Param 		orderNo path string true "Return ID"
// @Success 	200 {object} Response{result=string} "ReturnOrder Deleted"
// @Success 	204 {object} Response "No Content, Order Delete Successfully"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/reorder/delete/{orderNo} [delete]
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

	handleResponse(w, true, "Return order deleted successfully", nil, http.StatusOK)
}
