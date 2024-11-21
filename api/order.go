package api

import (
	request "boilerplate-backend-go/dto/request"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Orders(apiRouter *chi.Mux) {
	apiRouter.Route("/orders", func(r chi.Router) {  
		r.Get("/allgetorder", app.AllGetOrder)			// GET /orders/allgetorder
		r.Get("/getbyID/{orderNo}", app.GetOrderID) 	// GET /orders/getbyID/{orderNo}
		r.Post("/create-order", app.CreateOrder) 		// POST /orders
		r.Put("/update/{orderNo}", app.UpdateOrder)     // PUT /orders/{orderNo}
		r.Delete("/delete/{orderNo}", app.DeleteOrder)  // DELETE /orders/{orderNo}
	})
}

// @Summary 	Get Order
// @Description Get all Order
// @ID 			Allget-order
// @Tags 		Orders
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]entity.Order} "Order Get"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Province not found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/orders/allgetorder [get]
func (api *Application) AllGetOrder(w http.ResponseWriter, r *http.Request) {
	res, err := api.Service.Order.AllGetOrder()
	if err != nil {
		HandleError(w, err)
		return
	}
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary      Get Order by ID
// @Description  Get details of an order by its order number
// @ID           get-order-by-id
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        orderNo  path     string  true  "Order Number"
// @Success      200 	  {object} Response{result=[]entity.Order} "Order Get by ID"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "Order not found"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /orders/getbyID/{orderNo} [get]
func (app *Application) GetOrderID(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")

	// เรียกใช้ Service เพื่อประมวลผล ตรวจสอบ ก่อนที่จะเข้า method query ตาม api ที่ส่งไป
	res, err := app.Service.Order.GetOrderID(orderNo)
	if err != nil {
		HandleError(w, err)
		return
	}
	
	// ส่งคืนข้อมูลคำสั่งซื้อ
	handleResponse(w, true, response, res, http.StatusOK)
}

// @Summary 	Create Order
// @Description Create a new order
// @ID 			create-order
// @Tags 		Orders
// @Accept 		json
// @Produce 	json
// @Param 		order body request.CreateOrderRequest true "Order Data"
// @Success 	200 {object} Response{result=entity.order} "Order Created"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/orders/create-order [post]
func (api *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req request.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := api.Service.Order.CreateOrder(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

// @Summary 	Update Order
// @Description Update an existing order
// @ID 			update-order
// @Tags 		Orders
// @Accept 		json
// @Produce 	json
// @Param 		orderNo path string true "Order Number"
// @Param 		order body request.UpdateOrderRequest true "Updated Order Data"
// @Success 	200 {object} Response{result=string} "Order Updated"
// @Success 	204 {object} Response "No Content, Order Updated Successfully"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/orders/update/{orderNo} [put]
func (api *Application) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var req request.UpdateOrderRequest
	orderNo := chi.URLParam(r, "orderNo")
	req.OrderNo = orderNo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.Service.Order.UpdateOrder(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary 	Delete Order
// @Description Delete an order
// @ID 			delete-order
// @Tags 		Orders
// @Accept 		json
// @Produce 	json
// @Param 		orderNo path string true "Order Number"
// @Success 	200 {object} Response{result=string} "Order Deleted"
// @Success 	204 {object} Response "No Content, Order Delete Successfully"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/orders/delete/{orderNo} [delete]
func (api *Application) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	if err := api.Service.Order.DeleteOrder(orderNo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}



// // @Summary 	Get Order
// // @Description Get all Order
// // @ID 			get-order
// // @Tags 		Orders
// // @Accept 		json
// // @Produce 	json
// // @Success 	200 {object} Response{result=[]entity.Order} "AllOrder"
// // @Failure 	400 {object} Response "Bad Request"
// // @Failure 	404 {object} Response "Province not found"
// // @Failure 	500 {object} Response "Internal Server Error"
// // @Router 		/orders/get-order [get]
// func (app *Application) GetOrder(w http.ResponseWriter, r *http.Request) {
// 	res, err := app.Service.Order.GetOrder()
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 	handleResponse(w, true, response, res, http.StatusOK)
// }