package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) BefRORoute(apiRouter *chi.Mux) {
	apiRouter.Route("/before-return-order", func(r chi.Router) {
		r.Get("/list-orders", app.ListBeforeReturnOrders)
		r.Get("/list-lines", app.ListBeforeReturnOrderLines) // Updated route for listing return order lines without orderNo
		r.Get("/get-order", app.GetAllOrderDetail)
		r.Get("/get-orders", app.GetAllOrderDetails) //with paginate

		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)

		r.Get("/get-orderbySO/{soNo}", app.GetOrderDetailBySO)
		r.Get("/search/{soNo}", app.SearchSaleOrder) // New route for searching sale order

		r.Post("/create", app.CreateBeforeReturnOrderWithLines)
		r.Post("/create-trade", app.CreateTradeReturn)

		r.Patch("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines) // New route for updating return order with lines

		r.Delete("/delete-befodline/{recID}", app.DeleteBeforeReturnOrderLine)

		/******** Trade Retrun ********/
		r.Post("/create-trade", app.CreateTradeReturn)
		r.Post("/add-line/{orderNo}", app.AddTradeReturnLine)
		r.Post("/confirm/{orderNo}", app.ConfirmTradeReturn)
		r.Post("/cancel/{orderNo}", app.CancelTradeReturn)

	})
}

// @Summary Create a new trade return order
// @Description Create a new trade return order with multiple order lines
// @ID create-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Trade Return Detail"
// @Success 201 {object} api.Response "Trade return created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /before-return-order/create-trade [post]
func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	result, err := app.Service.BefRO.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order created successfully", result, http.StatusCreated)
}

// @Summary Add a new trade return line to an existing order
// @Description Add a new trade return line based on the provided order number and line details
// @ID add-trade-return-line
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.TradeReturnLineRequest true "Trade Return Line Details"
// @Success 201 {object} api.Response "Trade return line created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /before-return-order/add-line/{orderNo} [post]
func (app *Application) AddTradeReturnLine(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, fmt.Errorf("OrderNo is required"))
		return
	}

	var req request.TradeReturnLineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, fmt.Errorf("invalid request format: %v", err))
		return
	}

	// เรียก service layer เพื่อสร้างข้อมูล
	err := app.Service.BefRO.CreateTradeReturnLine(r.Context(), orderNo, req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Trade return line created successfully", nil, http.StatusCreated)
}

// ConfirmSaleReturn godoc
// @Summary Confirm a sale return order
// @Description Confirm a sale return order based on the provided details
// @ID confirm-sale-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response{data=response.ConfirmReturnResponse} "Sale return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/confirm/{orderNo} [post]
func (app *Application) ConfirmTradeReturn(w http.ResponseWriter, r *http.Request) {
	// 1. รับค่า orderNo จาก URL parameter
	orderNo := chi.URLParam(r, "orderNo")

	// 2. ดึงข้อมูล claims (ข้อมูลผู้ใช้) จาก context (มาจาก JWT authentication)
	claims, ok := r.Context().Value("claims").(map[string]interface{})
	if !ok {
		handleError(w, fmt.Errorf("user claims are missing or invalid"))
		return
	}

	// 3. ดึง username จาก claims เพื่อใช้เป็น confirmBy
	confirmBy, ok := claims["username"].(string)
	if !ok || confirmBy == "" {
		handleError(w, fmt.Errorf("username is missing or invalid"))
		return
	}

	// 4. เรียกใช้ service layer เพื่อดำเนินการ confirm
	err := app.Service.BefRO.ConfirmReturn(r.Context(), orderNo, confirmBy)
	if err != nil {
		handleError(w, err)
		return
	}

	// 5. สร้าง response และส่งกลับ
	response := res.ConfirmReturnResponse{
		OrderNo:   orderNo,
		ConfirmBy: confirmBy,
	}
	handleResponse(w, true, "Trade Return Order confirm successfully", response, http.StatusOK)
}

// CancelSaleReturn godoc
// @Summary Cancel a sale return order
// @Description Cancel a sale return order based on the provided details
// @ID cancel-sale-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param cancelDetails body request.CancelReturnRequest true "Cancel details"
// @Success 200 {object} api.Response{data=response.CancelReturnResponse} "Sale return order canceled successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/cancel/{orderNo} [post]
func (app *Application) CancelTradeReturn(w http.ResponseWriter, r *http.Request) {
	// 1. รับค่า orderNo จาก URL parameter
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		http.Error(w, "OrderNo is required", http.StatusBadRequest)
		return
	}

	// 2. ตรวจสอบว่า order มีอยู่จริง
	existingOrder, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}
	if existingOrder == nil {
		handleResponse(w, false, "Order not found", nil, http.StatusNotFound)
		return
	}

	// 3. รับและตรวจสอบข้อมูล request body
	var req request.CancelReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 4. ตรวจสอบข้อมูลที่จำเป็น
	if req.CancelBy == "" {
		http.Error(w, "CancelBy is required", http.StatusBadRequest)
		return
	}
	if req.Remark == "" {
		http.Error(w, "Remark is required", http.StatusBadRequest)
		return
	}

	// 5. เรียกใช้ service
	err = app.Service.BefRO.CancelReturn(r.Context(), orderNo, req.CancelBy, req.Remark)
	if err != nil {
		handleError(w, err)
		return
	}

	// 6. ส่ง response
	response := res.CancelReturnResponse{
		RefID:    orderNo,
		CancelBy: req.CancelBy,
		Remark:   req.Remark,
	}
	handleResponse(w, true, "Trade Return Order cancel successfully", response, http.StatusOK)
}

// ListReturnOrders godoc
// @Summary List all return orders
// @Description Retrieve a list of all before return orders
// @ID list-before-return-orders
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/list-orders [get]
func (app *Application) ListBeforeReturnOrders(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.BefRO.ListBeforeReturnOrders(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// ListBeforeReturnOrderLines godoc
// @Summary List all return order lines
// @Description Retrieve a list of all return order lines
// @ID list-before-return-order-lines
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/list-lines [get]
func (app *Application) ListBeforeReturnOrderLines(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.BefRO.ListBeforeReturnOrderLines(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order lines retrieved successfully", result, http.StatusOK)
}

// GetBeforeReturnOrderByOrderNo godoc
// @Summary Get return order by order number
// @Description Retrieve the details of a specific return order by its order number
// @ID get-before-return-order-by-order-no
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/{orderNo} [get]
func (app *Application) GetBeforeReturnOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order retrieved successfully", result, http.StatusOK)
}

// GetBeforeReturnOrderLineByOrderNo godoc
// @Summary Get return order lines by order number
// @Description Retrieve the details of all return order lines by order number
// @ID get-before-return-order-line-by-order-no
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/line/{orderNo} [get]
func (app *Application) GetBeforeReturnOrderLineByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.BefRO.GetBeforeReturnOrderLineByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order lines retrieved successfully", result, http.StatusOK)
}

// @Summary 	Get Before Return Order
// @Description Get all Before Return Order
// @ID 			Allget-BefReturnOrder
// @Tags 		Before Return Order
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} Response{result=[]response.OrderDetail} "Get All"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "not found endpoint"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/before-return-order/get-order [get]
func (api *Application) GetAllOrderDetail(w http.ResponseWriter, r *http.Request) {
	result, err := api.Service.BefRO.GetAllOrderDetail(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// @Summary 	Get Paginated Before Return Order
// @Description Get all Before Return Order with pagination
// @ID 			Get-BefReturnOrder-Paginated
// @Tags 		Before Return Order
// @Accept 		json
// @Produce 	json
// @Param       page  query int false "Page number" default(1)
// @Param       limit query int false "Page size" default(10)
// @Success 	200 {object} Response{result=[]response.OrderDetail} "Get Paginated Orders"
// @Failure 	400 {object} Response "Bad Request"
// @Failure 	404 {object} Response "Not Found"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/before-return-order/get-orders [get]
func (api *Application) GetAllOrderDetails(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r) // ฟังก์ชันช่วยดึง page และ limit จาก Query Parameters

	result, err := api.Service.BefRO.GetAllOrderDetails(r.Context(), page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// Helper function: parsePagination
func parsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	page := parseInt(query.Get("page"), 1)    // Default page = 1
	limit := parseInt(query.Get("limit"), 10) // Default limit = 10
	return page, limit
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// @Summary      Get Before Return Order by SO
// @Description  Get details of an order by its SO number
// @ID           GetBySO-BefReturnOrder
// @Tags         Before Return Order
// @Accept       json
// @Produce      json
// @Param        soNo  path     string  true  "soNo"
// @Success      200 	  {object} Response{result=[]response.OrderDetail} "Get by SO"
// @Failure      400      {object} Response "Bad Request"
// @Failure      404      {object} Response "not found endpoint"
// @Failure      500      {object} Response "Internal Server Error"
// @Router       /before-return-order/get-orderbySO/{soNo} [get]
func (app *Application) GetOrderDetailBySO(w http.ResponseWriter, r *http.Request) {
	soNo := chi.URLParam(r, "soNo")
	if soNo == "" {
		handleError(w, errors.ValidationError("soNo is required"))
		return
	}

	result, err := app.Service.BefRO.GetOrderDetailBySO(r.Context(), soNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved by SO successfully", result, http.StatusOK)
}

// SearchSaleOrder godoc
// @Summary Search sale order by SO number
// @Description Retrieve the details of a sale order by its SO number
// @ID search-sale-order
// @Tags Search Sale Order
// @Accept json
// @Produce json
// @Param soNo path string true "SO number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/search/{soNo} [get]
func (app *Application) SearchSaleOrder(w http.ResponseWriter, r *http.Request) {
	soNo := chi.URLParam(r, "soNo")
	result, err := app.Service.BefRO.SearchSaleOrder(r.Context(), soNo)
	if err != nil {
		handleError(w, err)
		return
	}

	if result == nil {
		handleResponse(w, false, "Sale order not found", nil, http.StatusNotFound)
		return
	}

	handleResponse(w, true, "Sale order retrieved successfully", result, http.StatusOK)
}

// CreateOrderWithLines godoc
// @Summary Create a new return order with lines
// @Description Create a new return order with the provided details
// @ID create-before-return-order-with-lines
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Before return order details"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/create [post]
func (app *Application) CreateBeforeReturnOrderWithLines(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	result, err := app.Service.BefRO.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order created successfully", result, http.StatusCreated)
}

// UpdateBeforeReturnOrderWithLines godoc
// @Summary Update an existing return order with lines
// @Description Update an existing return order with the provided details
// @ID update-return-order-with-lines
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.BeforeReturnOrder true "Before return order details"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/update/{orderNo} [patch]
func (app *Application) UpdateBeforeReturnOrderWithLines(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	req.OrderNo = orderNo // Ensure the orderNo from the URL is used

	result, err := app.Service.BefRO.UpdateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order updated successfully", result, http.StatusOK)
}

// @Summary 	Delete Order line
// @Description Delete an order line
// @ID 			delete-BeforeReturnOrderLine
// @Tags 		Before Return Order
// @Accept 		json
// @Produce 	json
// @Param 		recID path string true "Rec ID"
// @Success 	200 {object} Response{result=string} "Before ReturnOrderLine Deleted"
// @Failure 	404 {object} Response "Order Not Found"
// @Failure 	422 {object} Response "Validation Error"
// @Failure 	500 {object} Response "Internal Server Error"
// @Router 		/before-return-order/delete-befodline/{recID} [delete]
func (api *Application) DeleteBeforeReturnOrderLine(w http.ResponseWriter, r *http.Request) {
	recID := chi.URLParam(r, "recID")
	if recID == "" {
		handleError(w, errors.ValidationError("ReturnID is required in the path"))
		return
	}

	if err := api.Service.BefRO.DeleteBeforeReturnOrderLine(r.Context(), recID); err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order lines deleted successfully", nil, http.StatusOK)
}
