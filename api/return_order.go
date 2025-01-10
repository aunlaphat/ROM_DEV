package api

// ตัวกลางของ http request ที่คอยรับส่งข้อมูลไปมา
// รับมาในรูป request ส่งออกในรูป response
// ส่งคำขอไปยัง service เพื่อขอให้ดึงข้อมูลที่ต้องการออกมา แต่ service จะทำการ validation ก่อนเพื่อตรวจข้อผิดพลาดก่อนจะดึง query ออกมาจาก repo ให้
// api handle error about res req send error to client ex. 400 500 401
// service handle error about business logic/ relation data send error to api
// ตรวจสอบส่วนหน้า ส่วนที่รับมาจาก client เช่น input ที่ถูกป้อนเข้ามา
import (
	request "boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/errors"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) ReturnOrders(apiRouter *chi.Mux) {
	apiRouter.Route("/reorder", func(r chi.Router) {
		r.Get("/allget", app.AllGetReturnOrder)                             // GET /reorder/allget
		r.Get("/getbyID/{returnID}", app.GetReturnOrderID)                  // GET /reorder/getbyID/{returnID}
		r.Get("/allgetline", app.GetAllReturnOrderLines)                    // GET /reorder/allgetline
		r.Get("/getlinebyID/{returnID}", app.GetReturnOrderLinesByReturnID) // GET /reorder/getlinebyID/{returnID}
		r.Post("/create", app.CreateReturnOrder)                            // POST /reorder/create
		r.Patch("/update/{returnID}", app.UpdateReturnOrder)                // PUT /reorder/update/{returnID}
		r.Delete("/delete/{returnID}", app.DeleteReturnOrder)               // DELETE /reorder/delete/{returnID}
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

	result, err := api.Service.ReturnOrder.AllGetReturnOrder(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// @Summary      Get Return Order by ID
// @Description  Get details of an order by its return id
// @ID           GetByID-ReturnOrder
// @Tags         ReturnOrder
// @Accept       json
// @Produce      json
// @Param        returnID  path     string  true  "Return ID"
// @Success      200 	  {object} Response{result=[]entity.ReturnOrder} "Get by ID"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "not found endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /reorder/getbyID/{returnID} [get]
func (app *Application) GetReturnOrderID(w http.ResponseWriter, r *http.Request) {

	returnID := chi.URLParam(r, "returnID")
	if returnID == "" { 
		handleError(w, errors.ValidationError("ReturnID is required"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderByID(r.Context(), returnID)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order retrieved successfully", result, http.StatusOK)
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

	handleResponse(w, true, "Order lines retrieved successfully", result, http.StatusOK)
}

// @Summary      Get Return Order Line by ID
// @Description  Get details of an order line by its return id
// @ID           GetLineByID-ReturnOrder
// @Tags         ReturnOrder
// @Accept       json
// @Produce      json
// @Param        returnID  path     string  true  "Return ID"
// @Success      200 	  {object} Response{result=[]entity.ReturnOrderLine} "Get by ID"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "not found endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /reorder/getlinebyID/{returnID} [get]
func (app *Application) GetReturnOrderLinesByReturnID(w http.ResponseWriter, r *http.Request) {
	returnID := chi.URLParam(r, "returnID")
	if returnID == "" {
		handleError(w, errors.ValidationError("ReturnID is required"))
		return
	}

	result, err := app.Service.ReturnOrder.GetReturnOrderLinesByReturnID(r.Context(), returnID)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order lines retrieved successfully", result, http.StatusOK)

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
	var req request.CreateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// Validate required fields
	if req.ReturnID == "" || req.OrderNo == "" {
		handleError(w, errors.BadRequestError("ReturnID or OrderNo are required"))
		return
	}

	// จำนวนสินค้าที่คืนไม่ถูกต้อง เช่น ค่าติดลบหรือเป็น 0
	for _, line := range req.ReturnOrderLine {
		if line.ReturnQTY <= 0 {
			handleError(w, errors.BadRequestError("Invalid ReturnQTY in ReturnOrderLine"))
			return
		}
	}

	err := api.Service.ReturnOrder.CreateReturnOrder(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Return order created successfully", nil, http.StatusCreated)
}

// @Summary Update Order
// @Description Update an existing return order using returnID in the path
// @ID Update-ReturnOrder
// @Tags ReturnOrder
// @Accept json
// @Produce json
// @Param returnID path string true "Return ID"
// @Param order body request.UpdateReturnOrder true "Updated Order Data"
// @Success 200 {object} Response "ReturnOrder Updated Successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Order Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /reorder/update/{returnID} [patch]
func (api *Application) UpdateReturnOrder(w http.ResponseWriter, r *http.Request) {
	// ดึง ReturnID จากพาธ
	returnID := chi.URLParam(r, "returnID")
	if returnID == "" {
		handleError(w, errors.ValidationError("ReturnID is required in the path"))
		return
	}

	// Decode JSON Payload เพื่ออัปเดตฟิลด์ที่กำหนด
	var req request.UpdateReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, errors.BadRequestError("Invalid JSON format"))
		return
	}

	// กำหนด ReturnID จากพาธ
	req.ReturnID = returnID
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
// @Param 		returnID path string true "Return ID"
// @Success 	200 {object} Response{result=string} "ReturnOrder Deleted"
// @Success 	204 {object} Response "No Content, Order Delete Successfully"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/reorder/delete/{returnID} [delete]
func (api *Application) DeleteReturnOrder(w http.ResponseWriter, r *http.Request) {
	returnID := chi.URLParam(r, "returnID")
	if returnID == "" {
		handleError(w, errors.ValidationError("ReturnID is required in the path"))
		return
	}

	err := api.Service.ReturnOrder.DeleteReturnOrder(r.Context(), returnID)
    if err != nil {
        handleError(w, err)
        return
    }

	handleResponse(w, true, "Return order deleted successfully", nil, http.StatusOK)
}
