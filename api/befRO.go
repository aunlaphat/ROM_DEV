package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) BefRORoute(apiRouter *chi.Mux) {
	apiRouter.Route("/before-return-order", func(r chi.Router) {
		//r.Use(middleware.AuthMiddleware(app.Logger.Logger, "TRADE_CONSIGN", "WAREHOUSE", "VIEWER", "ACCOUNTING", "SYSTEM_ADMIN"))
		r.Get("/list-orders", app.ListBeforeReturnOrders)
		r.Post("/create", app.CreateBeforeReturnOrderWithLines)
		r.Put("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines)
		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/list-lines", app.ListBeforeReturnOrderLines)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)
	})

	apiRouter.Route("/sale-return", func(r chi.Router) {
		//r.Use(middleware.AuthMiddleware(app.Logger.Logger, "TRADE_CONSIGN", "WAREHOUSE", "VIEWER", "ACCOUNTING", "SYSTEM_ADMIN"))
		r.Get("/search/{soNo}", app.SearchSaleOrder)
		r.Post("/create", app.CreateSaleReturn)
		r.Post("/confirm", app.ConfirmSaleReturn)
	})

	apiRouter.Route("/draft-confirm", func(r chi.Router) {
		//r.Use(middleware.AuthMiddleware(app.Logger.Logger, "TRADE_CONSIGN", "WAREHOUSE", "VIEWER", "ACCOUNTING", "SYSTEM_ADMIN"))
		r.Get("/list-drafts", app.ListDrafts)
		r.Put("/edit-order/{orderNo}", app.EditDraftCF)
		r.Post("/confirm-order", app.ConfirmOrder)
	})
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
// @Router /before-return-order/update/{orderNo} [put]
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

	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully retrieved order",
		zap.String("OrderNo", result.OrderNo))
	handleResponse(w, true, "Order retrieved successfully", result, http.StatusOK)
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
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validation
	if req.OrderNo == "" {
		http.Error(w, "OrderNo is required", http.StatusBadRequest)
		return
	}
	if req.SoNo == "" {
		http.Error(w, "SoNo is required", http.StatusBadRequest)
		return
	}
	if req.ChannelID == 0 {
		http.Error(w, "ChannelID is required", http.StatusBadRequest)
		return
	}
	if req.CustomerID == "" {
		http.Error(w, "CustomerID is required", http.StatusBadRequest)
		return
	}
	if req.WarehouseID == 0 {
		http.Error(w, "WarehouseID is required", http.StatusBadRequest)
		return
	}
	if req.ReturnType == "" {
		http.Error(w, "ReturnType is required", http.StatusBadRequest)
		return
	}
	if req.TrackingNo == "" {
		http.Error(w, "TrackingNo is required", http.StatusBadRequest)
		return
	}
	if len(req.BeforeReturnOrderLines) == 0 {
		http.Error(w, "At least one order line is required", http.StatusBadRequest)
		return
	}
	for _, line := range req.BeforeReturnOrderLines {
		if line.SKU == "" {
			http.Error(w, "SKU is required for all order lines", http.StatusBadRequest)
			return
		}
		if line.QTY <= 0 {
			http.Error(w, "QTY must be greater than 0 for all order lines", http.StatusBadRequest)
			return
		}
		if line.ReturnQTY < 0 {
			http.Error(w, "ReturnQTY cannot be negative for all order lines", http.StatusBadRequest)
			return
		}
		if line.Price < 0 {
			http.Error(w, "Price cannot be negative for all order lines", http.StatusBadRequest)
			return
		}
	}

	// Check if the order already exists
	existingOrder, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), req.OrderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	if existingOrder != nil {
		// If the order already exists, update the SR number
		srNo := "SR123456" // Generate SR number (this is a placeholder, replace with actual SR number generation logic)
		err = app.Service.BefRO.UpdateDynamicFields(r.Context(), req.OrderNo, map[string]interface{}{"SrNo": srNo})
		if err != nil {
			handleError(w, err)
			return
		}

		// Update the result with the new SR number
		existingOrder.SrNo = srNo

		fmt.Printf("\n📋 ========== Updated Sale Return Order ========== 📋\n")
		printOrderDetails(existingOrder)
		handleResponse(w, true, "Sale return order updated successfully", existingOrder, http.StatusOK)
		return
	}

	// If the order does not exist, create a new one
	result, err := app.Service.BefRO.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	// Generate SR number (this is a placeholder, replace with actual SR number generation logic)
	srNo := "SR123456"

	// Update the SR number in the database
	err = app.Service.BefRO.UpdateDynamicFields(r.Context(), result.OrderNo, map[string]interface{}{"SrNo": srNo})
	if err != nil {
		handleError(w, err)
		return
	}

	// Update the result with the new SR number
	result.SrNo = srNo

	// Check user role
	/* userRole := r.Context().Value(middleware.ContextUserRole).(string)
	if userRole == "ACCOUNTING" {
		// Show "Create CN" button for accounting role
		handleResponse(w, true, "Sale return order created successfully. You can create CN.", result, http.StatusOK)
	} else {
		// Do not show "Create CN" button for other roles
		handleResponse(w, true, "Sale return order created successfully", result, http.StatusOK)
	} */
}

// ConfirmSaleReturn godoc
// @Summary Confirm a sale return order
// @Description Confirm a sale return order based on the provided details
// @ID confirm-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param saleReturn body request.BeforeReturnOrder true "Sale Return Order"
// @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "Sale return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/confirm [post]
func (app *Application) ConfirmSaleReturn(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fields := map[string]interface{}{
		"StatusConfID": req.StatusConfID,
		"ConfirmBy":    req.ConfirmBy,
		"UpdateBy":     req.UpdateBy,
	}

	err := app.Service.BefRO.UpdateDynamicFields(r.Context(), req.OrderNo, fields)
	if err != nil {
		handleError(w, err)
		return
	}

	result, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), req.OrderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Confirmed Sale Return Order ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	handleResponse(w, true, "Sale return order confirmed successfully", result, http.StatusOK)
}

// ListDrafts godoc
// @Summary List all draft orders
// @Description Retrieve a list of all draft orders
// @ID list-drafts
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /draft-confirm/list-drafts [get]
func (app *Application) ListDrafts(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.BefRO.ListDrafts(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== All Drafts (%d) ========== 📋\n", len(result))
	for i, draft := range result {
		fmt.Printf("\n📦 Draft #%d:\n", i+1)
		printDraftDetails(&draft)
		for i := range draft.BeforeReturnOrderLines {
			fmt.Printf("\n📦 Order Line #%d:\n", i+1)
			printDraftLineDetails(&draft.BeforeReturnOrderLines[i])
		}
	}
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully retrieved all drafts",
		zap.Int("totalDrafts", len(result)))
	handleResponse(w, true, "Drafts retrieved successfully", result, http.StatusOK)
}

// EditOrder godoc
// @Summary Edit an existing order
// @Description Edit an existing order with the provided details
// @ID edit-order
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.EditOrderRequest true "Edit order details"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /draft-confirm/edit-order/{orderNo} [put]
func (app *Application) EditDraftCF(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	var req request.EditOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	req.OrderNo = orderNo // Ensure the orderNo from the URL is used

	result, err := app.Service.BefRO.EditOrder(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Edited Order ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	app.Logger.Info("✅ Successfully edited order",
		zap.String("OrderNo", result.OrderNo))
	handleResponse(w, true, "Order edited successfully", result, http.StatusOK)
}

// ConfirmOrder godoc
// @Summary Confirm an order
// @Description Confirm an order based on the provided details
// @ID confirm-order
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param confirmOrder body request.ConfirmOrderRequest true "Confirm order details"
// @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "Order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/confirm-order [post]
func (app *Application) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	var req request.ConfirmOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := app.Service.BefRO.ConfirmOrder(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Confirmed Order ========== 📋\n")
	printOrderDetails(result)
	// fmt.Println("=====================================")

	handleResponse(w, true, "Order confirmed successfully", result, http.StatusOK)
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
