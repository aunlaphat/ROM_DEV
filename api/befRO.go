package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) BefRORoute(apiRouter *chi.Mux) {
	apiRouter.Route("/before-return-order", func(r chi.Router) {
		//r.Use(middleware.AuthMiddleware(app.Logger.Logger, "TRADE_CONSIGN", "WAREHOUSE", "VIEWER", "ACCOUNTING", "SYSTEM_ADMIN"))
		r.Get("/list-orders", app.ListBeforeReturnOrders)
		r.Get("/list-lines", app.ListBeforeReturnOrderLines) // Updated route for listing return order lines without orderNo
		r.Get("/get-order", app.GetAllOrderDetail)
		r.Get("/get-orders", app.GetAllOrderDetails) //with paginate

		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)

		r.Get("/get-orderbySO/{soNo}", app.GetOrderDetailBySO)

		r.Post("/create", app.CreateBeforeReturnOrderWithLines)

		r.Patch("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines) // New route for updating return order with lines

		r.Delete("/delete-befodline/{recID}", app.DeleteBeforeReturnOrderLine)
	})

	apiRouter.Route("/trade-return", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		/******** Trade Retrun ********/
		r.Post("/create-trade", app.CreateTradeReturn)
		r.Post("/add-line/{orderNo}", app.AddTradeReturnLine)
		r.Post("/confirm/{orderNo}", app.ConfirmTradeReturn)
		r.Post("/cancel/{orderNo}", app.CancelTradeReturn)
	})

	apiRouter.Route("/sale-return", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/search/{soNo}", app.SearchSaleOrder)
		r.Post("/create", app.CreateSaleReturn)
		r.Put("/update/{orderNo}", app.UpdateSaleReturn)
		r.Post("/confirm/{orderNo}", app.ConfirmSaleReturn)
		r.Post("/cancel/{orderNo}", app.CancelSaleReturn)
	})

	apiRouter.Post("/login", app.Login)

	/* 	apiRouter.Route("/draft-confirm", func(r chi.Router) {
	//r.Use(middleware.AuthMiddleware(app.Logger.Logger, "TRADE_CONSIGN", "WAREHOUSE", "VIEWER", "ACCOUNTING", "SYSTEM_ADMIN"))
	r.Get("/list-drafts", app.ListDrafts)
	r.Put("/edit-order/{orderNo}", app.EditDraftCF)
	r.Post("/confirm-order", app.ConfirmOrder)
	*/
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
// @Router /trade-return/create-trade [post]
func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 	return
	// }

	// Validation
	// if req.OrderNo == "" {
	// 	http.Error(w, "OrderNo is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.SoNo == "" {
	// 	http.Error(w, "SoNo is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.ChannelID == 0 {
	// 	http.Error(w, "ChannelID is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.CustomerID == "" {
	// 	http.Error(w, "CustomerID is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.WarehouseID == 0 {
	// 	http.Error(w, "WarehouseID is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.ReturnType == "" {
	// 	http.Error(w, "ReturnType is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.TrackingNo == "" {
	// 	http.Error(w, "TrackingNo is required", http.StatusBadRequest)
	// 	return
	// }
	// if len(req.BeforeReturnOrderLines) == 0 {
	// 	http.Error(w, "At least one order line is required", http.StatusBadRequest)
	// 	return
	// }
	// for _, line := range req.BeforeReturnOrderLines {
	// 	if line.SKU == "" {
	// 		http.Error(w, "SKU is required for all order lines", http.StatusBadRequest)
	// 		return
	// 	}
	// 	if line.QTY <= 0 {
	// 		http.Error(w, "QTY must be greater than 0 for all order lines", http.StatusBadRequest)
	// 		return
	// 	}
	// 	if line.ReturnQTY < 0 {
	// 		http.Error(w, "ReturnQTY cannot be negative for all order lines", http.StatusBadRequest)
	// 		return
	// 	}
	// 	if line.Price < 0 {
	// 		http.Error(w, "Price cannot be negative for all order lines", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// เช็คว่า orderNo ที่สร้างไม่ซ้ำกับตัวที่มีอยู่แล้ว
	existingOrder, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), req.OrderNo)
	if err != nil {
		handleError(w, err)
		return
	}
	if existingOrder != nil {
		handleResponse(w, false, "Order already exists", nil, http.StatusConflict)
		return
	}

	// ดึงค่า claims จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		handleError(w, fmt.Errorf("unauthorized: invalid user information"))
		return
	}

	// Set CreateBy จาก claims
	req.CreateBy = userID

	// Create a new order
	result, err := app.Service.BefRO.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Created Trade Return Order ========== 📋\n")
	printOrderDetails(result)
	fmt.Printf("\n📋 ========== Trade Return Order Line Details ========== 📋\n")
	for _, line := range result.BeforeReturnOrderLines {
		printOrderLineDetails(&line)
	}

	handleResponse(w, true, "Trade return order created successfully", result, http.StatusOK)
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
// @Router /trade-return/add-line/{orderNo} [post]
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
// @ID confirm-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response{data=response.ConfirmReturnResponse} "Sale return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/confirm/{orderNo} [post]
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
// @ID cancel-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param cancelDetails body request.CancelReturnRequest true "Cancel details"
// @Success 200 {object} api.Response{data=response.CancelReturnResponse} "Sale return order canceled successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/cancel/{orderNo} [post]
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

	fmt.Printf("\n📋 ========== All Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n📦 Order #%d:\n", i+1)
		printOrderDetails(&order)
		for j, line := range order.BeforeReturnOrderLines {
			fmt.Printf("\n📦 Order Line #%d:\n", j+1)
			printOrderLineDetails(&line)
		}
	}
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully retrieved all orders",
		zap.Int("totalOrders", len(result)))
	handleResponse(w, true, "📚 Orders retrieved successfully", result, http.StatusOK)
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

	fmt.Printf("\n📋 ========== All Order Lines (%d) ========== 📋\n", len(result))
	for i, line := range result {
		fmt.Printf("\n📦 Order Line #%d:\n", i+1)
		printOrderLineDetails(&line)
	}
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully retrieved all order lines",
		zap.Int("totalOrderLines", len(result)))
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

	fmt.Printf("\n📋 ========== Order Lines for OrderNo: %s ========== 📋\n", orderNo)
	for i, line := range result {
		fmt.Printf("\n📦 Order Line #%d:\n", i+1)
		printOrderLineDetails(&line)
	}
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully retrieved order lines",
		zap.String("OrderNo", orderNo),
		zap.Int("totalOrderLines", len(result)))
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
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param soNo path string true "SO number"
// @Success 200 {object} api.Response{data=response.SaleOrderResponse} "Sale order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "Sale order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/search/{soNo} [get]
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

	fmt.Printf("\n📋 ========== Sale Order Details ========== 📋\n")
	for _, order := range result {
		printSaleOrderDetails(&order)
		fmt.Printf("\n📋 ========== Sale Order Line Details ========== 📋\n")
		for _, line := range order.OrderLines {
			printSaleOrderLineDetails(&line)
		}
	}
	// fmt.Println("=====================================")

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

	fmt.Printf("\n📋 ========== Created Order ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully created order",
		zap.String("OrderNo", result.OrderNo))
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

	fmt.Printf("\n📋 ========== Updated Order ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully updated order",
		zap.String("OrderNo", result.OrderNo))
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

// CreateSaleReturn godoc
// @Summary Create a new sale return order
// @Description Create a new sale return order based on the provided details
// @ID create-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param saleReturn body request.BeforeReturnOrder true "Sale Return Order"
// @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "Sale return order created successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/create [post]
func (app *Application) CreateSaleReturn(w http.ResponseWriter, r *http.Request) {
	// 1. Authentication check
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Error("Failed to decode request", zap.Error(err))
		handleResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// Set user information from claims
	req.CreateBy = userID

	// 4. Call service
	result, err := app.Service.BefRO.CreateBeforeReturn(r.Context(), req)
	if err != nil {
		app.Logger.Error("Failed to create sale return",
			zap.Error(err),
			zap.String("orderNo", req.OrderNo))

		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			handleResponse(w, false, err.Error(), nil, http.StatusBadRequest)
		case strings.Contains(err.Error(), "already exists"):
			handleResponse(w, false, err.Error(), nil, http.StatusConflict)
		default:
			handleResponse(w, false, "Internal server error", nil, http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("\n📋 ========== Created Sale Return Order ========== 📋\n")
	printOrderDetails(result)
	fmt.Printf("\n📋 ========== Sale Return Order Line Details ========== 📋\n")
	for _, line := range result.BeforeReturnOrderLines {
		printOrderLineDetails(&line)
	}

	handleResponse(w, true, "Sale return order created successfully", result, http.StatusOK)
}

// UpdateSaleReturn godoc
// @Summary Update the SR number for a sale return order
// @Description Update the SR number for a sale return order based on the provided details
// @ID update-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderNo path string true "Order number"
// @Param request body request.UpdateSaleReturnRequest true "SR number details"
// @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "SR number updated successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 401 {object} api.Response "Unauthorized - Missing or invalid token"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/update/{orderNo} [put]
func (app *Application) UpdateSaleReturn(w http.ResponseWriter, r *http.Request) {
	// 1. รับและตรวจสอบ orderNo
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		http.Error(w, "OrderNo is required", http.StatusBadRequest)
		return
	}

	// 2. รับและตรวจสอบ request body
	var req request.UpdateSaleReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, fmt.Errorf("invalid request format: %v", err))
		return
	}

	// อัพเดทการตรวจสอบข้อมูล
	if req.SrNo == "" {
		http.Error(w, "SrNo is required", http.StatusBadRequest)
		return
	}

	// 3. ตรวจสอบว่า order มีอยู่จริง
	existingOrder, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}
	if existingOrder == nil {
		handleResponse(w, false, "Order not found", nil, http.StatusNotFound)
		return
	}

	// ดึง userID จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		handleError(w, fmt.Errorf("unauthorized: invalid user information"))
		return
	}

	// เรียกใช้ service พร้อมส่ง userID
	err = app.Service.BefRO.UpdateSaleReturn(r.Context(), orderNo, req.SrNo, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := res.UpdateSaleReturnResponse{
		OrderNo:    orderNo,
		SrNo:       req.SrNo,
		UpdateBy:   userID,
		UpdateDate: time.Now(),
	}

	handleResponse(w, true, "SR number updated successfully", response, http.StatusOK)
}

// ConfirmSaleReturn godoc
// @Summary Confirm a sale return order
// @Description Confirm a sale return order based on the provided details
// @ID confirm-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response{data=response.ConfirmSaleReturnResponse} "Sale return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/confirm/{orderNo} [post]
func (app *Application) ConfirmSaleReturn(w http.ResponseWriter, r *http.Request) {
	// 1. รับค่า orderNo จาก URL parameter
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, fmt.Errorf("order number is required"))
		return
	}

	// 2. ดึงค่า claims จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	// 3. ดึงค่า userID จาก claims
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		handleError(w, fmt.Errorf("unauthorized: invalid user information"))
		return
	}

	// 4. เรียกใช้ service layer เพื่อดำเนินการ confirm
	err = app.Service.BefRO.ConfirmSaleReturn(r.Context(), orderNo, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	// 5. สร้าง response และส่งกลับ
	response := res.ConfirmSaleReturnResponse{
		OrderNo:     orderNo,
		ConfirmBy:   userID,
		ConfirmDate: time.Now(),
	}

	handleResponse(w, true, "Sale return order confirmed successfully", response, http.StatusOK)
}

// CancelSaleReturn godoc
// @Summary Cancel a sale return order
// @Description Cancel a sale return order based on the provided details
// @ID cancel-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param request body request.CancelSaleReturnRequest true "Cancel Sale Return"
// @Success 200 {object} api.Response{data=response.CancelSaleReturnResponse} "Sale return order canceled successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/cancel/{orderNo} [post]
func (app *Application) CancelSaleReturn(w http.ResponseWriter, r *http.Request) {
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

	// 2. ดึงค่า claims จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	// 3. ดึงค่า userID จาก claims
	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		handleError(w, fmt.Errorf("unauthorized: invalid user information"))
		return
	}

	// 3. รับและตรวจสอบข้อมูล request body
	var req request.CancelSaleReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 4. ตรวจสอบข้อมูลที่จำเป็น
	if req.Remark == "" {
		http.Error(w, "Remark is required", http.StatusBadRequest)
		return
	}

	// 5. เรียกใช้ service
	err = app.Service.BefRO.CancelSaleReturn(r.Context(), orderNo, userID, req.Remark)
	if err != nil {
		handleError(w, err)
		return
	}

	// 6. ส่ง response
	response := res.CancelSaleReturnResponse{
		RefID:        orderNo,
		CancelStatus: true,
		CancelBy:     userID,
		Remark:       req.Remark,
		CancelDate:   time.Now(),
	}
	handleResponse(w, true, "Sale return order canceled successfully", response, http.StatusOK)
}

func printOrderDetails(order *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", order.SrNo)
	fmt.Printf("📡 ChannelID: %d\n", order.ChannelID)
	fmt.Printf("🔙 ReturnType: %s\n", order.ReturnType)
	fmt.Printf("👤 CustomerID: %s\n", order.CustomerID)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", order.Logistic)
	fmt.Printf("🏢 WarehouseID: %d\n", order.WarehouseID)
	fmt.Printf("📄 SoStatusID: %v\n", order.SoStatusID)
	fmt.Printf("📊 MkpStatusID: %v\n", order.MkpStatusID)
	fmt.Printf("📅 ReturnDate: %v\n", order.ReturnDate)
	fmt.Printf("🔖 StatusReturnID: %d\n", order.StatusReturnID)
	fmt.Printf("✅ StatusConfID: %d\n", order.StatusConfID)
	fmt.Printf("👤 ConfirmBy: %v\n", order.ConfirmBy)
	fmt.Printf("👤 CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
	fmt.Printf("👤 UpdateBy: %v\n", order.UpdateBy)
	fmt.Printf("📅 UpdateDate: %v\n", order.UpdateDate)
	fmt.Printf("❌ CancelID: %v\n", order.CancelID)
}

func printOrderLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("🔢 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}

func printSaleOrderDetails(order *res.SaleOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🔢 SoNo: %s\n", order.SoNo)
	fmt.Printf("📊 StatusMKP: %s\n", order.StatusMKP)
	fmt.Printf("📊 SalesStatus: %s\n", order.SalesStatus)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
}

func printSaleOrderLineDetails(line *res.SaleOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🚩 ItemName: %s\n", line.ItemName)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
}

func printDraftDetails(draft *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("👤 Customer: %s\n", draft.CustomerID)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
	fmt.Printf("📦 TrackingNo: %s\n", draft.TrackingNo)
	fmt.Printf("📡 Channel: %d\n", draft.ChannelID)
	fmt.Printf("📅 CreateDate: %v\n", draft.CreateDate)
	fmt.Printf("🏢 Warehouse: %d\n", draft.WarehouseID)
}

func printDraftLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("🔢 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}
