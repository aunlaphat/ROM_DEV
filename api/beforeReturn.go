package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) BeforeReturnRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/before-return-order", func(r chi.Router) {
		r.Get("/list-orders", app.ListBeforeReturnOrders)
		r.Get("/list-lines", app.ListBeforeReturnOrderLines)
		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)
		r.Post("/create", app.CreateBeforeReturnOrderWithLines)
		r.Patch("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines)

		// get real order
		//r.Get("/get-order", app.GetAllOrderDetail)                             // get Order of ROM_V_OrderDetail
		//r.Get("/get-orders", app.GetAllOrderDetails)                           // get Order of ROM_V_OrderDetail with paginate
		r.Get("/get-orderbySO/{soNo}", app.GetOrderDetailBySO)                 // search by SO of ROM_V_OrderDetail
		r.Delete("/delete-befodline/{recID}", app.DeleteBeforeReturnOrderLine) // delete line by recID of BeforeReturnOrder
	})

	apiRouter.Route("/sale-return", func(r chi.Router) {
		// ✅ ไม่ต้องใช้ JWT สำหรับการค้นหา
		r.Get("/search", app.SearchOrder)

		// ✅ ใช้ Middleware JWT สำหรับทุกเส้นทางที่แก้ไขข้อมูล
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(app.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/create", app.CreateSaleReturn)
			r.Patch("/update", app.UpdateSaleReturn)
			r.Patch("/confirm/{orderNo}", app.ConfirmSaleReturn)
			r.Post("/cancel/{orderNo}", app.CancelSaleReturn)
		})
	})

	apiRouter.Route("/draft-confirm", func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		// 📌 Draft & Confirm (ใช้ดูรายละเอียดของ Order)
		//r.Get("/detail/{orderNo}", app.GetDraftConfirmOrderByOrderNo)

		// 📌 Draft Status Orders
		r.Route("/drafts", func(draft chi.Router) {
			draft.Get("/", app.ListDraftOrders)
			draft.Get("/code-r", app.ListCodeR)
			draft.Post("/code-r", app.AddCodeR)
			draft.Delete("/code-r/{orderNo}/{sku}", app.DeleteCodeR)
			draft.Patch("/{orderNo}", app.UpdateDraftOrder)
		})

		// 📌 Confirm Status Orders
		r.Route("/confirms", func(confirm chi.Router) {
			confirm.Get("/", app.ListConfirmOrders)
		})
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
	result, err := app.Service.BeforeReturn.ListBeforeReturnOrders(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== All Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n📦 Order #%d 📦\n", i+1)
		utils.PrintOrderDetails(&order)
		for j, line := range order.BeforeReturnOrderLines {
			fmt.Printf("\n📦 Order Line #%d 📦\n", j+1)
			utils.PrintOrderLineDetails(&line)
		}
		fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(order.BeforeReturnOrderLines))
		fmt.Println("=====================================")
	}

	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
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

	result, err := app.Service.BeforeReturn.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Created Order ========== 📋\n")
	utils.PrintOrderDetails(result)
	for i, line := range result.BeforeReturnOrderLines {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Order created successfully ⭐", result, http.StatusCreated)
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

	req.OrderNo = orderNo

	result, err := app.Service.BeforeReturn.UpdateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Updated Order ========== 📋\n")
	utils.PrintOrderDetails(result)
	for i, line := range result.BeforeReturnOrderLines {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Order updated successfully ⭐", result, http.StatusOK)
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
	result, err := app.Service.BeforeReturn.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
	utils.PrintOrderDetails(result)
	for i, line := range result.BeforeReturnOrderLines {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Order retrieved successfully ⭐", result, http.StatusOK)
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
	result, err := app.Service.BeforeReturn.ListBeforeReturnOrderLines(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== All Order Lines (%d) ========== 📋\n", len(result))
	for i, line := range result {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Order lines retrieved successfully ⭐", result, http.StatusOK)
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
	result, err := app.Service.BeforeReturn.GetBeforeReturnOrderLineByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Order Lines for OrderNo: %s ========== 📋\n", orderNo)
	for i, line := range result {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Order lines retrieved successfully ⭐", result, http.StatusOK)
}

// SearchSaleOrder godoc
// @Summary 🔎 Search order by SO number or Order number
// @Description Retrieve the details of an order by its SO number or Order number
// @ID search-order
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param soNo query string false "SO number"
// @Param orderNo query string false "Order number"
// @Success 200 {object} api.Response{data=response.SaleOrderResponse} "⭐ Order retrieved successfully ⭐"
// @Failure 400 {object} api.Response "⚠️ Bad Request"
// @Failure 404 {object} api.Response "❌ Sale order not found"
// @Failure 500 {object} api.Response "🔥 Internal Server Error"
// @Router /sale-return/search [get]
func (app *Application) SearchOrder(w http.ResponseWriter, r *http.Request) {
	// ✅ รับค่า Query Parameters
	soNo := r.URL.Query().Get("soNo")
	orderNo := r.URL.Query().Get("orderNo")

	// 🚨 ตรวจสอบว่าอย่างน้อยต้องมีค่าใดค่าหนึ่ง
	if soNo == "" && orderNo == "" {
		app.Logger.Warn("⚠️ Missing search criteria")
		handleResponse(w, false, "⚠️ Either SoNo or OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	// 🔎 Log ค้นหาข้อมูลคำสั่งขาย
	app.Logger.Info("🔎 Searching for Sale Order...",
		zap.String("SoNo", soNo),
		zap.String("OrderNo", orderNo),
	)

	// 🛠 เรียกใช้ Service Layer
	order, err := app.Service.BeforeReturn.SearchOrder(r.Context(), soNo, orderNo)
	if err != nil {
		errMsg := err.Error()

		// ❌ กรณีไม่พบข้อมูล
		if errMsg == "ไม่พบข้อมูลคำสั่งซื้อสินค้า" {
			app.Logger.Warn("⚠️ No Sale Order found",
				zap.String("SoNo", soNo),
				zap.String("OrderNo", orderNo),
				zap.String("Error", errMsg),
			)
			handleResponse(w, false, "⚠️ Sale order not found", nil, http.StatusNotFound)
			return
		}

		// 🔥 กรณีเกิดข้อผิดพลาดอื่น ๆ
		app.Logger.Error("🔥 Failed to search order",
			zap.String("SoNo", soNo),
			zap.String("OrderNo", orderNo),
			zap.String("Error", errMsg),
			zap.Error(err),
		)
		handleResponse(w, false, "🔥 Internal server error", nil, http.StatusInternalServerError)
		return
	}

	// ✅ คืนค่าผลลัพธ์สำเร็จ
	app.Logger.Info("✅ Order retrieved successfully",
		zap.String("SoNo", order.SoNo),
		zap.String("OrderNo", order.OrderNo),
		zap.Int("TotalItems", len(order.OrderLines)),
	)
	handleResponse(w, true, "⭐ Order retrieved successfully ⭐", order, http.StatusOK)
}

// CreateSaleReturn godoc
// @Summary Create a new sale return order
// @Description Create a new sale return order with order details and items
// @ID create-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param request body request.CreateSaleReturnOrder true "Create Sale Return Request"
// @Success 201 {object} api.Response{data=response.BeforeReturnOrderResponse} "Sale Return Order created successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/create [post]
func (app *Application) CreateSaleReturn(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract claims from JWT
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		errMsg := "Unauthorized access - Missing or invalid JWT claims"
		app.Logger.Error("🚷 "+errMsg, zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract userID from claims
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		errMsg := "Unauthorized access - userID extraction failed"
		app.Logger.Error("🚷 "+errMsg, zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Decode request body
	var req request.CreateSaleReturnOrder
	body, _ := io.ReadAll(r.Body) // อ่าน JSON ก่อน Decode
	app.Logger.Info("📥 Received Request Body", zap.String("body", string(body)))

	if err := json.Unmarshal(body, &req); err != nil {
		errMsg := "Invalid request format"
		app.Logger.Warn("⚠️ "+errMsg, zap.Error(err))
		handleResponse(w, false, errMsg, nil, http.StatusBadRequest)
		return
	}

	// ✅ Call Service Layer
	createdOrder, err := app.Service.BeforeReturn.CreateSaleReturn(r.Context(), req, userID)
	if err != nil {
		errMsg := "Failed to create Sale Return Order"
		app.Logger.Error("❌ "+errMsg, zap.Error(err))
		handleResponse(w, false, errMsg, nil, http.StatusInternalServerError)
		return
	}

	// ✅ Return Success Response
	handleResponse(w, true, "⭐ Sale Return Order created successfully ⭐", createdOrder, http.StatusCreated)
}

// UpdateSaleReturn godoc
// @Summary Update Sale Return order
// @Description Update Sale Return order with SrNo
// @ID update-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param request body request.UpdateSaleReturn true "Update Sale Return Request"
// @Success 200 {object} api.Response{data=response.UpdateSaleReturnResponse} "Sale Return Order updated successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/update [patch]
func (app *Application) UpdateSaleReturn(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract claims from JWT
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		app.Logger.Error("🚷 Unauthorized access - Missing or invalid JWT claims", zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract userID from claims
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		app.Logger.Error("🚷 Unauthorized access - userID extraction failed", zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Decode request body
	var req request.UpdateSaleReturn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Warn("⚠️ Invalid request format", zap.Error(err))
		handleResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// ✅ Call Service Layer
	updatedOrder, err := app.Service.BeforeReturn.UpdateSaleReturn(r.Context(), req, userID)
	if err != nil {
		app.Logger.Error("❌ Failed to update Sale Return Order", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ✅ Return Success Response
	handleResponse(w, true, "⭐ Sale Return Order updated successfully ⭐", updatedOrder, http.StatusOK)
}

// ConfirmSaleReturn godoc
// @Summary Confirm Sale Return order
// @Description Confirm Sale Return order by updating StatusReturnID and StatusConfID
// @ID confirm-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response{data=response.ConfirmSaleReturnResponse} "Sale Return Order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 403 {object} api.Response "Forbidden - Insufficient permissions"
// @Failure 404 {object} api.Response "Sale order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/confirm/{orderNo} [patch]
func (app *Application) ConfirmSaleReturn(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract claims from JWT
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		app.Logger.Error("🚷 Unauthorized access - Missing or invalid JWT claims", zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract userID and roleID
	userID, roleID, err := utils.GetUserInfoFromClaims(claims)
	if err != nil {
		app.Logger.Error("🚷 Unauthorized access - userID or roleID extraction failed", zap.Error(err))
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract `orderNo` from path
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		app.Logger.Warn("⚠️ Missing OrderNo in request path")
		handleResponse(w, false, "OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	// ✅ Call Service Layer
	confirmedOrder, err := app.Service.BeforeReturn.ConfirmSaleReturn(r.Context(), orderNo, roleID, userID)
	if err != nil {
		app.Logger.Error("❌ Failed to confirm Sale Return Order", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	handleResponse(w, true, "⭐ Sale Return Order confirmed successfully ⭐", confirmedOrder, http.StatusOK)
}

// CancelSaleReturn godoc
// @Summary Cancel a sale return order
// @Description Cancel a sale return order based on the provided details
// @ID cancel-sale-return
// @Tags Sale Return
// @Accept json
// @Produce json
// @Param request body request.CancelSaleReturn true "Cancel Sale Return"
// @Success 200 {object} api.Response{data=response.CancelSaleReturnResponse} "Sale return order canceled successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /sale-return/cancel/{orderNo} [post]
func (app *Application) CancelSaleReturn(w http.ResponseWriter, r *http.Request) {
	// ✅ 1. Extract Order Number from URL
	orderNo := strings.TrimSpace(chi.URLParam(r, "orderNo"))
	if orderNo == "" {
		app.Logger.Warn("❌ Missing orderNo in request")
		handleResponse(w, false, "❌ OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	// ✅ 2. Authenticate User (JWT)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "🚷 Unauthorized Access 🚷", nil, http.StatusUnauthorized)
		return
	}

	// ✅ 3. Extract User ID from Token
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleResponse(w, false, "🔑 Invalid UserID in Token Claims 🔑", nil, http.StatusUnauthorized)
		return
	}

	// ✅ 4. Decode & Validate Request Body
	var req request.CancelSaleReturn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Warn("❌ Invalid request payload", zap.Error(err))
		handleResponse(w, false, "❌ Invalid request payload", nil, http.StatusBadRequest)
		return
	}

	req.Remark = strings.TrimSpace(req.Remark)
	if req.Remark == "" {
		app.Logger.Warn("❌ Missing remark in request")
		handleResponse(w, false, "❌ Remark is required", nil, http.StatusBadRequest)
		return
	}

	// ✅ 5. Log Request Data
	app.Logger.Info("🛑 CancelSaleReturn requested",
		zap.String("OrderNo", orderNo),
		zap.String("CanceledBy", userID),
		zap.String("Remark", req.Remark),
	)

	// ✅ 6. Call Service Layer (Ensuring Correct Response Handling)
	result, err := app.Service.BeforeReturn.CancelSaleReturn(r.Context(), req, userID)
	if err != nil {
		app.Logger.Error("❌ Failed to cancel sale return", zap.Error(err))
		handleError(w, err)
		return
	}

	// ✅ 7. Return JSON Response
	app.Logger.Info("✅ Sale return order canceled successfully",
		zap.String("OrderNo", orderNo),
		zap.String("CanceledBy", userID),
	)

	handleResponse(w, true, "⭐ Sale return order canceled successfully ⭐", result, http.StatusOK)
}

// ListDraftOrders godoc
// @Summary List all draft orders
// @Description Retrieve a list of all draft orders within a date range
// @ID list-draft-orders
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} api.Response{data=[]response.ListDraftConfirmOrdersResponse} "All Draft orders retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "Draft orders not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/list-drafts [get]
func (app *Application) ListDraftOrders(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	// 📌 ตรวจสอบค่าที่รับเข้ามา
	if startDate == "" || endDate == "" {
		app.Logger.Warn("⚠️ Missing required query parameters ⚠️")
		handleResponse(w, false, "⚠️ Missing startDate or endDate parameters ⚠️", nil, http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบรูปแบบวันที่
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		app.Logger.Warn("⚠️ Invalid startDate format ⚠️", zap.String("startDate", startDate))
		handleResponse(w, false, "⚠️ Invalid startDate format (YYYY-MM-DD) ⚠️", nil, http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		app.Logger.Warn("⚠️ Invalid endDate format ⚠️", zap.String("endDate", endDate))
		handleResponse(w, false, "⚠️ Invalid endDate format (YYYY-MM-DD) ⚠️", nil, http.StatusBadRequest)
		return
	}

	if start.After(end) {
		app.Logger.Warn("⚠️ startDate cannot be after endDate ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		handleResponse(w, false, "⚠️ startDate cannot be after endDate ⚠️", nil, http.StatusBadRequest)
		return
	}

	// 📌 เรียกใช้งาน Service Layer
	result, err := app.Service.BeforeReturn.ListDraftOrders(r.Context(), startDate, endDate)
	if err != nil {
		app.Logger.Error("🚨 Failed to list draft orders 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ⚠️ ถ้าไม่มีรายการ ส่ง response 404
	if len(result) == 0 {
		app.Logger.Warn("⚠️ No draft orders found ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		handleResponse(w, false, "⚠️ No draft orders found ⚠️", nil, http.StatusNotFound)
		return
	}

	// ✅ Debug logging
	app.Logger.Debug("📋 Retrieved draft orders",
		zap.Int("count", len(result)),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)

	// ✅ ส่ง response กลับไป
	handleResponse(w, true, "⭐ Draft orders retrieved successfully ⭐", result, http.StatusOK)
}

// ListConfirmOrders godoc
// @Summary List all confirm orders
// @Description Retrieve a list of all confirm orders within a date range
// @ID list-confirm-orders
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} api.Response{data=[]response.ListDraftConfirmOrdersResponse} "All Confirm orders retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "Confirm orders not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/list-confirms [get]
func (app *Application) ListConfirmOrders(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	// 📌 ตรวจสอบค่าที่รับเข้ามา
	if startDate == "" || endDate == "" {
		app.Logger.Warn("⚠️ Missing required query parameters ⚠️")
		handleResponse(w, false, "⚠️ Missing startDate or endDate parameters ⚠️", nil, http.StatusBadRequest)
		return
	}

	// ✅ ตรวจสอบรูปแบบวันที่
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		app.Logger.Warn("⚠️ Invalid startDate format ⚠️", zap.String("startDate", startDate))
		handleResponse(w, false, "⚠️ Invalid startDate format (YYYY-MM-DD) ⚠️", nil, http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		app.Logger.Warn("⚠️ Invalid endDate format ⚠️", zap.String("endDate", endDate))
		handleResponse(w, false, "⚠️ Invalid endDate format (YYYY-MM-DD) ⚠️", nil, http.StatusBadRequest)
		return
	}

	if start.After(end) {
		app.Logger.Warn("⚠️ startDate cannot be after endDate ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		handleResponse(w, false, "⚠️ startDate cannot be after endDate ⚠️", nil, http.StatusBadRequest)
		return
	}

	// 📌 เรียกใช้งาน Service Layer
	result, err := app.Service.BeforeReturn.ListConfirmOrders(r.Context(), startDate, endDate)
	if err != nil {
		app.Logger.Error("🚨 Failed to list confirm orders 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ⚠️ ถ้าไม่มีรายการ ส่ง response 404
	if len(result) == 0 {
		app.Logger.Warn("⚠️ No confirm orders found ⚠️",
			zap.String("startDate", startDate),
			zap.String("endDate", endDate),
		)
		handleResponse(w, false, "⚠️ No confirm orders found ⚠️", nil, http.StatusNotFound)
		return
	}

	// ✅ Debug logging
	app.Logger.Debug("📋 Retrieved confirm orders",
		zap.Int("count", len(result)),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
	)

	// ✅ ส่ง response กลับไป
	handleResponse(w, true, "⭐ Confirm orders retrieved successfully ⭐", result, http.StatusOK)
}

/* // GetDraftConfirmOrderByOrderNo godoc
// @Summary Get Draft Confirm Order by OrderNo
// @Description Retrieve Draft Confirm Order Head and Lines
// @ID get-draft-confirm-order
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param orderNo path string true "Order Number"
// @Success 200 {object} api.Response{data=response.DraftHeadResponse} "Draft Confirm Order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/detail/{orderNo} [get]
func (app *Application) GetDraftConfirmOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")

	// 📌 ใช้ Logger ที่มี `orderNo` ติดอยู่
	logger := app.Logger.With(zap.String("orderNo", orderNo))

	// ✅ Log API Call Start
	logFinish := logger.LogAPICall(r.Context(), "GetDraftConfirmOrderByOrderNo")
	defer logFinish("Completed", nil)

	// 📌 เรียก Service Layer และรับ Response + Error
	order, err := app.Service.BeforeReturn.GetDraftConfirmOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleResponse(w, err, logger)
		return
	}

	// ✅ ส่ง Response กลับไป
	handleResponse(w, true, "⭐ Draft Confirm Order retrieved successfully ⭐", nil, http.StatusOK)
} */

// ListCodeR godoc
// @Summary List all CodeR (SKU, ItemName) where SKU starts with 'R'
// @Description Retrieve a list of CodeR from ROM_V_ProductAll where SKU starts with 'R'
// @ID list-code-r
// @Tags CodeR
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=[]response.ListCodeRResponse} "All CodeR retrieved successfully"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /code-r/list [get]
func (app *Application) ListCodeR(w http.ResponseWriter, r *http.Request) {
	// 📌 เรียกใช้งาน Service Layer
	result, err := app.Service.BeforeReturn.ListCodeR(r.Context())
	if err != nil {
		app.Logger.Error("🚨 Failed to list CodeR 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ⚠️ ถ้าไม่มีรายการ CodeR ส่ง response 404
	if len(result) == 0 {
		app.Logger.Warn("⚠️ No CodeR found (WHERE SKU LIKE 'R%') ⚠️")
		handleResponse(w, false, "⚠️ No CodeR found ⚠️", nil, http.StatusNotFound)
		return
	}

	// ✅ Debug logging
	app.Logger.Debug("📋 Retrieved CodeR list", zap.Int("count", len(result)))

	// ✅ ส่ง response กลับไป
	handleResponse(w, true, "⭐ CodeR list retrieved successfully ⭐", result, http.StatusOK)
}

// AddCodeR godoc
// @Summary Add CodeR
// @Description Add a new CodeR entry
// @ID add-code-r
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param body body request.AddCodeR true "CodeR details"
// @Success 201 {object} api.Response{data=[]response.AddCodeRResponse} "CodeR added successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/code-r [post]
func (app *Application) AddCodeR(w http.ResponseWriter, r *http.Request) {
	var req request.AddCodeR

	// ✅ Decode JSON Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Error("🚨 Failed to decode request 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// ✅ Extract JWT Claims from Context
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		app.Logger.Warn("🚷 Unauthorized access attempt")
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract userID from Claims
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		app.Logger.Warn("🚷 Failed to extract userID from claims")
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	app.Logger.Info("👤 User authenticated", zap.String("userID", userID))

	// ✅ Call Service Layer
	results, err := app.Service.BeforeReturn.AddCodeR(r.Context(), req, userID)
	if err != nil {
		app.Logger.Error("🚨 Failed to add CodeR 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Failed to add CodeR", nil, http.StatusInternalServerError)
		return
	}

	// ✅ Return Success Response
	handleResponse(w, true, "⭐ CodeR added successfully ⭐", results, http.StatusCreated)
}

// DeleteCodeR godoc
// @Summary Delete CodeR
// @Description Delete a CodeR entry by SKU and OrderNo
// @ID delete-code-r
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param sku path string true "SKU"
// @Success 200 {object} api.Response "CodeR deleted successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 404 {object} api.Response "Not Found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/code-r/{orderNo}/{sku} [delete]
func (app *Application) DeleteCodeR(w http.ResponseWriter, r *http.Request) {
	// ✅ รับค่า `orderNo` และ `sku` จาก URL Path
	orderNo := chi.URLParam(r, "orderNo")
	sku := chi.URLParam(r, "sku")

	// ✅ ตรวจสอบว่าค่าถูกต้อง
	if orderNo == "" || sku == "" {
		app.Logger.Warn("⚠️ Missing required parameters: OrderNo and SKU")
		handleResponse(w, false, "⚠️ OrderNo and SKU are required", nil, http.StatusBadRequest)
		return
	}

	// ✅ Extract JWT Claims
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		app.Logger.Warn("🚷 Unauthorized access attempt")
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Extract userID from Claims
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		app.Logger.Warn("🚷 Failed to extract userID from claims")
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	app.Logger.Info("🗑️ User deleting CodeR",
		zap.String("userID", userID),
		zap.String("orderNo", orderNo),
		zap.String("sku", sku),
	)

	// ✅ Call Service Layer
	err = app.Service.BeforeReturn.DeleteCodeR(r.Context(), orderNo, sku, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			app.Logger.Warn("⚠️ CodeR not found", zap.String("orderNo", orderNo), zap.String("sku", sku))
			handleResponse(w, false, "⚠️ CodeR not found", nil, http.StatusNotFound)
			return
		}

		app.Logger.Error("🚨 Failed to delete CodeR 🚨", zap.Error(err))
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ✅ Return Success Response
	handleResponse(w, true, "⭐ Draft order retrieved successfully ⭐", nil, http.StatusOK)
}

// UpdateDraftOrders godoc
// @Summary Update draft orders
// @Description Update draft orders and change status to Confirm and Booking
// @ID update-draft-orders
// @Tags Draft & Confirm
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response{data=response.UpdateOrderStatusResponse} "Draft order updated successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 401 {object} api.Response "Unauthorized"
// @Failure 404 {object} api.Response "Not Found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /draft-confirm/update-draft/{orderNo} [patch]
func (app *Application) UpdateDraftOrder(w http.ResponseWriter, r *http.Request) {
	// ✅ รับค่า `orderNo` จาก URL Path
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleResponse(w, false, "⚠️ Order number is required", nil, http.StatusBadRequest)
		return
	}

	// ✅ Extract userID from claims
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	// ✅ Call Service Layer
	updatedOrder, err := app.Service.BeforeReturn.UpdateDraftOrder(r.Context(), orderNo, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			handleResponse(w, false, "⚠️ Draft order not found", nil, http.StatusNotFound)
			return
		}
		handleResponse(w, false, "❌ Internal Server Error", nil, http.StatusInternalServerError)
		return
	}

	// ✅ Return Success Response
	handleResponse(w, true, "⭐ Draft order updated successfully ⭐", updatedOrder, http.StatusOK)
}

/*
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

	result, err := api.Service.BeforeReturn.GetAllOrderDetail(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
} */
/*
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

	page, limit := utils.ParsePagination(r)

	result, err := api.Service.BeforeReturn.GetAllOrderDetails(r.Context(), page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
} */

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

	result, err := app.Service.BeforeReturn.GetOrderDetailBySO(r.Context(), soNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "⭐ Orders retrieved by SO successfully ⭐", result, http.StatusOK)
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
		handleError(w, errors.ValidationError("RecID is required in the path"))
		return
	}

	if err := api.Service.BeforeReturn.DeleteBeforeReturnOrderLine(r.Context(), recID); err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "⭐ Order lines deleted successfully ⭐", nil, http.StatusOK)
}
