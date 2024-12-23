package api

// ตัวกลางของ http request ที่คอยรับส่งข้อมูลไปมา
// รับมาในรูป request ส่งออกในรูป response
// ส่งคำขอไปยัง service เพื่อทำการ validation ตรวจข้อผิดพลาดก่อนจะดึง query ออกมาจาก repo
import (
	request "boilerplate-backend-go/dto/request"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) ReturnOrders(apiRouter *chi.Mux) {
	apiRouter.Route("/reorder", func(r chi.Router) {
		r.Get("/allget", app.AllGetReturnOrder)               // GET /reorder/allget
		r.Get("/getbyID/{returnID}", app.GetReturnOrderID)    // GET /reorder/getbyID/{returnID}
		r.Post("/create", app.CreateReturnOrder)           	  // POST /reorder/create
		r.Put("/update/{returnID}", app.UpdateReturnOrder)    // PUT /reorder/update/{returnID}
		r.Delete("/delete/{returnID}", app.DeleteReturnOrder) // DELETE /reorder/delete/{returnID}

		
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

	res, err := api.Service.ReturnOrder.AllGetReturnOrder()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
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

	returnID := chi.URLParam(r, "returnID") //รับค่าจากพาทเพื่อดึงข้อมูล returnid ตาม ReturnID in db
	res, err := app.Service.ReturnOrder.GetReturnOrderByID(returnID)
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
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
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ReturnID == "" || req.OrderNo == "" {
		http.Error(w, "ReturnID and OrderNo are required", http.StatusBadRequest)
		return
	}

	// Validate nested fields
	for _, line := range req.ReturnOrderLine {
		if line.ReturnQTY <= 0 {
			http.Error(w, "Invalid ReturnQTY in ReturnOrderLine", http.StatusBadRequest)
			return
		}
	}

	err := api.Service.ReturnOrder.CreateReturnOrder(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Return order created successfully",
	})
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
// @Router /reorder/update/{returnID} [put]
func (api *Application) UpdateReturnOrder(w http.ResponseWriter, r *http.Request) {
    // ดึง ReturnID จากพาธ
    returnID := chi.URLParam(r, "returnID")
    if returnID == "" {
        http.Error(w, "ReturnID is required in the path", http.StatusBadRequest)
        return
    }

    // Decode JSON Payload เพื่ออัปเดตฟิลด์ที่กำหนด
    var req request.UpdateReturnOrder
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // กำหนด ReturnID จากพาธ
    req.ReturnID = returnID

    // เรียก Service เพื่ออัปเดตข้อมูล
    if err := api.Service.ReturnOrder.UpdateReturnOrder(req); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ตอบกลับเมื่อสำเร็จ
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Return order updated successfully",
    })
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
		http.Error(w, "ReturnID is required in the path", http.StatusBadRequest)
		return
	}

	if err := api.Service.ReturnOrder.DeleteReturnOrder(returnID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Return order deleted successfully",
	})
}
