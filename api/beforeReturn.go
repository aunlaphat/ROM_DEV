 package api

// import (
// 	"boilerplate-backend-go/dto/request"
// 	res "boilerplate-backend-go/dto/response"
// 	"boilerplate-backend-go/errors"
// 	"boilerplate-backend-go/utils"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/jwtauth"
// 	"go.uber.org/zap"
// )

// // ReturnOrderRoute defines the routes for return order operations
// func (app *Application) BeforeReturnRoute(apiRouter *chi.Mux) {
// 	apiRouter.Route("/before-return-order", func(r chi.Router) {
// 		r.Get("/list-orders", app.ListBeforeReturnOrders)
// 		r.Get("/list-lines", app.ListBeforeReturnOrderLines)
// 		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
// 		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)
// 		r.Post("/create", app.CreateBeforeReturnOrderWithLines)
// 		r.Patch("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines)

// 		// get real order
// 		r.Get("/get-orders", app.GetAllOrderDetails) // แสดงข้อมูลออเดอร์ head+line ทั้งหมดที่กรอกทำรายการเข้ามาในระบบ แบบ paginate                        
// 		r.Get("/get-orderbySO/{soNo}", app.GetOrderDetailBySO) // แสดงข้อมูล order ที่ทำการคืนมาโดยเลข SO               
// 		r.Delete("/delete-line/{orderNo}/{sku}", app.DeleteBeforeReturnOrderLine) // ลบรายการคืนแต่ละรายการ
// 	})

// 	apiRouter.Route("/sale-return", func(r chi.Router) {
// 		r.Use(jwtauth.Verifier(app.TokenAuth))
// 		r.Use(jwtauth.Authenticator)

// 		r.Get("/search", app.SearchOrder)
// 		r.Post("/create", app.CreateSaleReturn)
// 		r.Patch("/update/{orderNo}", app.UpdateSaleReturn)
// 		r.Post("/confirm/{orderNo}", app.ConfirmSaleReturn)
// 		r.Post("/cancel/{orderNo}", app.CancelSaleReturn)
// 	})

// 	apiRouter.Route("/draft-confirm", func(r chi.Router) {
// 		r.Use(jwtauth.Verifier(app.TokenAuth))
// 		r.Use(jwtauth.Authenticator)

// 		// Draft & Confirm ใช้เหมือนกันในส่วนของการเปิด Modal และดูรายละเอียดของ Order
// 		r.Get("/detail/{orderNo}", app.GetDraftConfirmOrderByOrderNo)

// 		// Draft
// 		r.Get("/list-drafts", app.ListDraftOrders)
// 		r.Get("/list-code-r", app.ListCodeR)
// 		r.Post("/code-r", app.AddCodeR)
// 		r.Delete("/code-r/{orderNo}/{sku}", app.DeleteCodeR)
// 		r.Patch("/update-draft/{orderNo}", app.UpdateDraftOrder)

// 		// Confirm
// 		r.Get("/list-confirms", app.ListConfirmOrders)
// 	})
// }

// // ListReturnOrders godoc
// // @Summary List all return orders
// // @Description Retrieve a list of all before return orders
// // @ID list-before-return-orders
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Success 200 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/list-orders [get]
// func (app *Application) ListBeforeReturnOrders(w http.ResponseWriter, r *http.Request) {
// 	result, err := app.Service.BeforeReturn.ListBeforeReturnOrders(r.Context())
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== All Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n📦 Order #%d 📦\n", i+1)
// 		utils.PrintOrderDetails(&order)
// 		for j, line := range order.BeforeReturnOrderLines {
// 			fmt.Printf("\n📦 Order Line #%d 📦\n", j+1)
// 			utils.PrintOrderLineDetails(&line)
// 		}
// 		fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(order.BeforeReturnOrderLines))
// 		fmt.Println("=====================================")
// 	}

// 	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
// }

// // CreateOrderWithLines godoc
// // @Summary Create a new return order with lines
// // @Description Create a new return order with the provided details
// // @ID create-before-return-order-with-lines
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Param body body request.BeforeReturnOrder true "Before return order details"
// // @Success 201 {object} api.Response
// // @Failure 400 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/create [post]
// func (app *Application) CreateBeforeReturnOrderWithLines(w http.ResponseWriter, r *http.Request) {
// 	var req request.BeforeReturnOrder
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	result, err := app.Service.BeforeReturn.CreateBeforeReturnOrderWithLines(r.Context(), req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Created Order ========== 📋\n")
// 	utils.PrintOrderDetails(result)
// 	for i, line := range result.BeforeReturnOrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Order created successfully ⭐", result, http.StatusCreated)
// }

// // UpdateBeforeReturnOrderWithLines godoc
// // @Summary Update an existing return order with lines
// // @Description Update an existing return order with the provided details
// // @ID update-return-order-with-lines
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Param body body request.BeforeReturnOrder true "Before return order details"
// // @Success 200 {object} api.Response
// // @Failure 400 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/update/{orderNo} [patch]
// func (app *Application) UpdateBeforeReturnOrderWithLines(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	var req request.BeforeReturnOrder
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	req.OrderNo = orderNo

// 	result, err := app.Service.BeforeReturn.UpdateBeforeReturnOrderWithLines(r.Context(), req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Updated Order ========== 📋\n")
// 	utils.PrintOrderDetails(result)
// 	for i, line := range result.BeforeReturnOrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Order updated successfully ⭐", result, http.StatusOK)
// }

// // GetBeforeReturnOrderByOrderNo godoc
// // @Summary Get return order by order number
// // @Description Retrieve the details of a specific return order by its order number
// // @ID get-before-return-order-by-order-no
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Success 200 {object} api.Response
// // @Failure 404 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/{orderNo} [get]
// func (app *Application) GetBeforeReturnOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	result, err := app.Service.BeforeReturn.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
// 	utils.PrintOrderDetails(result)
// 	for i, line := range result.BeforeReturnOrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines))
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Order retrieved successfully ⭐", result, http.StatusOK)
// }

// // ListBeforeReturnOrderLines godoc
// // @Summary List all return order lines
// // @Description Retrieve a list of all return order lines
// // @ID list-before-return-order-lines
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Success 200 {object} api.Response
// // @Failure 404 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/list-lines [get]
// func (app *Application) ListBeforeReturnOrderLines(w http.ResponseWriter, r *http.Request) {
// 	result, err := app.Service.BeforeReturn.ListBeforeReturnOrderLines(r.Context())
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== All Order Lines (%d) ========== 📋\n", len(result))
// 	for i, line := range result {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Order lines retrieved successfully ⭐", result, http.StatusOK)
// }

// // GetBeforeReturnOrderLineByOrderNo godoc
// // @Summary Get return order lines by order number
// // @Description Retrieve the details of all return order lines by order number
// // @ID get-before-return-order-line-by-order-no
// // @Tags Before Return Order
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Success 200 {object} api.Response
// // @Failure 404 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /before-return-order/line/{orderNo} [get]
// func (app *Application) GetBeforeReturnOrderLineByOrderNo(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	result, err := app.Service.BeforeReturn.GetBeforeReturnOrderLineByOrderNo(r.Context(), orderNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Order Lines for OrderNo: %s ========== 📋\n", orderNo)
// 	for i, line := range result {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Order lines retrieved successfully ⭐", result, http.StatusOK)
// }

// // SearchSaleOrder godoc
// // @Summary Search order by SO number or Order number
// // @Description Retrieve the details of a order by its SO number or Order number
// // @ID search-order
// // @Tags Sale Return
// // @Accept json
// // @Produce json
// // @Param soNo query string false "SO number"
// // @Param orderNo query string false "Order number"
// // @Success 200 {object} api.Response{data=response.SaleOrderResponse} "Order retrieved successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 404 {object} api.Response "Sale order not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /sale-return/search [get]
// func (app *Application) SearchOrder(w http.ResponseWriter, r *http.Request) {
// 	soNo := r.URL.Query().Get("soNo")
// 	orderNo := r.URL.Query().Get("orderNo")

// 	// Validate input parameters
// 	if soNo == "" && orderNo == "" {
// 		app.Logger.Warn("No search criteria provided")
// 		handleResponse(w, false, "Either SoNo or OrderNo is required", nil, http.StatusBadRequest)
// 		return
// 	}

// 	// Input sanitization (optional)
// 	soNo = strings.TrimSpace(soNo)
// 	orderNo = strings.TrimSpace(orderNo)

// 	// Authorization check
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		app.Logger.Error("Authorization failed", zap.Error(err))
// 		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
// 		return
// 	}

// 	// Call service layer with error handling
// 	result, err := app.Service.BeforeReturn.SearchOrder(r.Context(), soNo, orderNo)
// 	if err != nil {
// 		app.Logger.Error("Failed to search order",
// 			zap.Error(err),
// 			zap.String("soNo", soNo),
// 			zap.String("orderNo", orderNo))
// 		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		return
// 	}

// 	// Handle no results found
// 	if len(result) == 0 {
// 		handleResponse(w, false, "⚠️ No orders found ⚠️", nil, http.StatusNotFound)
// 		return
// 	}

// 	// Correctly populate soNo and orderNo in orderLines
// 	for i := range result {
// 		for j := range result[i].OrderLines {
// 			result[i].OrderLines[j].SoNo = result[i].SoNo
// 			result[i].OrderLines[j].OrderNo = result[i].OrderNo
// 		}
// 	}

// 	// Debug logging (always print for now, can be controlled by log level later)
// 	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
// 	for _, order := range result {
// 		utils.PrintSaleOrderDetails(&order)
// 		fmt.Printf("\n📋 ========== Order Line Details ========== 📋\n")
// 		for i, line := range order.OrderLines {
// 			fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 			utils.PrintSaleOrderLineDetails(&line)
// 		}
// 		fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(order.OrderLines))
// 		fmt.Println("=====================================")
// 	}

// 	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
// }

// // CreateSaleReturn godoc
// // @Summary Create a new sale return order
// // @Description Create a new sale return order based on the provided details
// // @ID create-sale-return
// // @Tags Sale Return
// // @Accept json
// // @Produce json
// // @Param saleReturn body request.BeforeReturnOrder true "Sale Return Order"
// // @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "Sale return order created successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /sale-return/create [post]
// func (app *Application) CreateSaleReturn(w http.ResponseWriter, r *http.Request) {
// 	// 1. Authentication check
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
// 		return
// 	}

// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		return
// 	}

// 	var req request.BeforeReturnOrder
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		app.Logger.Error("Failed to decode request", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		return
// 	}

// 	// Set user information from claims
// 	req.CreateBy = userID

// 	// 4. Call service
// 	result, err := app.Service.BeforeReturn.CreateSaleReturn(r.Context(), req)
// 	if err != nil {
// 		app.Logger.Error("Failed to create sale return",
// 			zap.Error(err),
// 			zap.String("orderNo", req.OrderNo))

// 		// Handle specific error cases
// 		switch {
// 		case strings.Contains(err.Error(), "validation failed"):
// 			handleResponse(w, false, err.Error(), nil, http.StatusBadRequest)
// 		case strings.Contains(err.Error(), "already exists"):
// 			handleResponse(w, false, err.Error(), nil, http.StatusConflict)
// 		default:
// 			handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		}
// 		return
// 	}

// 	// Debug logging (always print for now, can be controlled by log level later)
// 	fmt.Printf("\n📋 ========== Created Sale Return Order ========== 📋\n")
// 	utils.PrintOrderDetails(result)
// 	fmt.Printf("\n📋 ========== Sale Return Order Line Details ========== 📋\n")
// 	for i, line := range result.BeforeReturnOrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.BeforeReturnOrderLines)) // Add logging for the number of lines
// 	fmt.Println("=====================================")

// 	// Send successful response
// 	handleResponse(w, true, "⭐ Sale return order created successfully ⭐", result, http.StatusOK)
// }

// // UpdateSaleReturn godoc
// // @Summary Update the SR number for a sale return order
// // @Description Update the SR number for a sale return order based on the provided details
// // @ID update-sale-return
// // @Tags Sale Return
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Param request body request.UpdateSaleReturn true "SR number details"
// // @Success 200 {object} api.Response{data=response.BeforeReturnOrderResponse} "SR number updated successfully"
// // @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// // @Failure 404 {object} api.Response "Not Found - Order not found"
// // @Failure 401 {object} api.Response "Unauthorized - Missing or invalid token"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /sale-return/update/{orderNo} [patch]
// func (app *Application) UpdateSaleReturn(w http.ResponseWriter, r *http.Request) {
// 	// 1. รับและตรวจสอบ orderNo
// 	orderNo := chi.URLParam(r, "orderNo")
// 	if orderNo == "" {
// 		http.Error(w, "OrderNo is required", http.StatusBadRequest)
// 		return
// 	}

// 	// 2. รับและตรวจสอบ request body
// 	var req request.UpdateSaleReturn
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, fmt.Errorf("invalid request format: %v", err))
// 		return
// 	}

// 	// อัพเดทการตรวจสอบข้อมูล
// 	if req.SrNo == "" {
// 		http.Error(w, "SrNo is required", http.StatusBadRequest)
// 		return
// 	}

// 	// 3. ตรวจสอบว่า order มีอยู่จริง
// 	existingOrder, err := app.Service.BeforeReturn.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}
// 	if existingOrder == nil {
// 		handleResponse(w, false, "⚠️ Order not found ⚠️", nil, http.StatusNotFound)
// 		return
// 	}

// 	// ดึง userID จาก JWT token
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
// 		return
// 	}

// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// เรียกใช้ service พร้อมส่ง userID
// 	err = app.Service.BeforeReturn.UpdateSaleReturn(r.Context(), orderNo, req.SrNo, userID)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	response := res.UpdateSaleReturnResponse{
// 		OrderNo:    orderNo,
// 		SrNo:       req.SrNo,
// 		UpdateBy:   userID,
// 		UpdateDate: time.Now(),
// 	}

// 	handleResponse(w, true, "⭐ SR number updated successfully ⭐", response, http.StatusOK)
// }

// // ConfirmSaleReturn godoc
// // @Summary Confirm a sale return order
// // @Description Confirm a sale return order based on the provided details
// // @ID confirm-sale-return
// // @Tags Sale Return
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Success 200 {object} api.Response{data=response.ConfirmSaleReturnResponse} "Sale return order confirmed successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /sale-return/confirm/{orderNo} [post]
// func (app *Application) ConfirmSaleReturn(w http.ResponseWriter, r *http.Request) {
// 	// 1. รับค่า orderNo จาก URL parameter
// 	orderNo := chi.URLParam(r, "orderNo")
// 	if orderNo == "" {
// 		handleError(w, fmt.Errorf("order number is required"))
// 		return
// 	}

// 	// 2. ดึงค่า claims จาก JWT token
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
// 		return
// 	}

// 	// 3. ดึงค่า userID และ roleID จาก claims
// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// 4. เรียกใช้ service layer เพื่อดำเนินการ confirm
// 	err = app.Service.BeforeReturn.ConfirmSaleReturn(r.Context(), orderNo, userID)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// 5. สร้าง response และส่งกลับ
// 	response := res.ConfirmSaleReturnResponse{
// 		OrderNo:     orderNo,
// 		ConfirmBy:   userID,
// 		ConfirmDate: time.Now(),
// 	}

// 	handleResponse(w, true, "⭐ Sale return order confirmed successfully ⭐", response, http.StatusOK)
// }

// // CancelSaleReturn godoc
// // @Summary Cancel a sale return order
// // @Description Cancel a sale return order based on the provided details
// // @ID cancel-sale-return
// // @Tags Sale Return
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Param request body request.CancelSaleReturn true "Cancel Sale Return"
// // @Success 200 {object} api.Response{data=response.CancelSaleReturnResponse} "Sale return order canceled successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /sale-return/cancel/{orderNo} [post]
// func (app *Application) CancelSaleReturn(w http.ResponseWriter, r *http.Request) {
// 	// 1. Validation ข้อมูล
// 	orderNo := chi.URLParam(r, "orderNo") // รับค่า orderNo จาก URL
// 	if orderNo == "" {
// 		http.Error(w, "OrderNo is required", http.StatusBadRequest)
// 		return
// 	}

// 	// 2. ตรวจสอบว่า order มีอยู่จริง
// 	existingOrder, err := app.Service.BeforeReturn.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
// 	if err != nil || existingOrder == nil {
// 		handleResponse(w, false, "⚠️ Order not found ⚠️", nil, http.StatusNotFound)
// 		return
// 	}

// 	// 3. Authentication - ตรวจสอบ JWT token
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleError(w, fmt.Errorf("unauthorized"))
// 		return
// 	}

// 	// 4. ดึง userID จาก token
// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// 5. รับและตรวจสอบข้อมูล request
// 	var req request.CancelSaleReturn
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	// 6. ตรวจสอบ Remark
// 	if req.Remark == "" {
// 		http.Error(w, "Remark is required", http.StatusBadRequest)
// 		return
// 	}

// 	// 7. เรียกใช้ service
// 	err = app.Service.BeforeReturn.CancelSaleReturn(r.Context(), orderNo, userID, req.Remark)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// 8. สร้าง response
// 	response := res.CancelSaleReturnResponse{
// 		RefID:        orderNo,
// 		CancelStatus: true,
// 		CancelBy:     userID,
// 		Remark:       req.Remark,
// 		CancelDate:   time.Now(),
// 	}

// 	// 9. ส่ง response กลับ
// 	handleResponse(w, true, "⭐ Sale return order canceled successfully ⭐", response, http.StatusOK)
// }

// // ListDraftOrders godoc
// // @Summary List all draft orders
// // @Description Retrieve a list of all draft orders
// // @ID list-draft-orders
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Success 200 {object} api.Response{data=[]response.ListDraftConfirmOrdersResponse} "All Draft orders retrieved successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 404 {object} api.Response "Draft orders not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/list-drafts [get]
// func (app *Application) ListDraftOrders(w http.ResponseWriter, r *http.Request) {
// 	// Call service layer with error handling
// 	result, err := app.Service.BeforeReturn.ListDraftOrders(r.Context())
// 	if err != nil {
// 		app.Logger.Error("🚨 Failed to list draft orders 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	// Handle no results found
// 	if len(result) == 0 {
// 		handleResponse(w, false, "⚠️ No draft orders found ⚠️", nil, http.StatusOK)
// 		return
// 	}

// 	// Debug logging (always print for now, can be controlled by log level later)
// 	fmt.Printf("\n📋 ========== All Draft Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n📦 Draft Order #%d 📦\n", i+1)
// 		utils.PrintDraftConfirmOrderDetails(&order)
// 	}

// 	// Send successful response
// 	handleResponse(w, true, "⭐ Draft orders retrieved successfully ⭐", result, http.StatusOK)
// }

// // ListConfirmOrders godoc
// // @Summary List all confirm orders
// // @Description Retrieve a list of all confirm orders
// // @ID list-confirm-orders
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Success 200 {object} api.Response{data=[]response.ListDraftConfirmOrdersResponse} "All Confirm orders retrieved successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 404 {object} api.Response "Confirm orders not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/list-confirms [get]
// func (app *Application) ListConfirmOrders(w http.ResponseWriter, r *http.Request) {
// 	// Call service layer with error handling
// 	result, err := app.Service.BeforeReturn.ListConfirmOrders(r.Context())
// 	if err != nil {
// 		app.Logger.Error("🚨 Failed to list confirm orders 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	// Handle no results found
// 	if len(result) == 0 {
// 		handleResponse(w, false, "⚠️ No confirm orders found ⚠️", nil, http.StatusOK)
// 		return
// 	}

// 	// Debug logging (always print for now, can be controlled by log level later)
// 	fmt.Printf("\n📋 ========== All Confirm Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n📦 Confirm Order #%d 📦\n", i+1)
// 		utils.PrintDraftConfirmOrderDetails(&order)
// 	}

// 	// Send successful response
// 	handleResponse(w, true, "⭐ Confirm orders retrieved successfully ⭐", result, http.StatusOK)
// }

// // ListCodeR godoc
// // @Summary List all CodeR
// // @Description Retrieve a list of all codeR
// // @ID list-code-r
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Success 200 {object} api.Response{data=[]response.CodeRResponse} "CodeR retrieved successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/list-code-r [get]
// func (app *Application) ListCodeR(w http.ResponseWriter, r *http.Request) {
// 	// Call service layer with error handling
// 	result, err := app.Service.BeforeReturn.ListCodeR(r.Context())
// 	if err != nil {
// 		app.Logger.Error("🚨 Failed to get all CodeR 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ CodeR retrieved successfully ⭐", result, http.StatusOK)
// }

// // AddCodeR godoc
// // @Summary Add CodeR
// // @Description Add a new CodeR entry
// // @ID add-code-r
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Param body body request.CodeR true "CodeR details"
// // @Success 201 {object} api.Response{data=response.DraftLineResponse} "CodeR added successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/code-r [post]
// func (app *Application) AddCodeR(w http.ResponseWriter, r *http.Request) {
// 	var req request.CodeR
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		app.Logger.Error("🚨 Failed to decode request 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusBadRequest)
// 		return
// 	}

// 	// Extract userID from claims
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
// 		return
// 	}

// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		return
// 	}

// 	// Set CreateBy from claims
// 	req.CreateBy = userID

// 	result, err := app.Service.BeforeReturn.AddCodeR(r.Context(), req)
// 	if err != nil {
// 		app.Logger.Error("🚨 Failed to add CodeR 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ CodeR added successfully ⭐", result, http.StatusCreated)
// }

// // DeleteCodeR godoc
// // @Summary Delete CodeR
// // @Description Delete a CodeR entry by SKU and OrderNo
// // @ID delete-code-r
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Param sku path string true "SKU"
// // @Success 200 {object} api.Response "CodeR deleted successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/code-r/{orderNo}/{sku} [delete]
// func (app *Application) DeleteCodeR(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")
// 	sku := chi.URLParam(r, "sku")
// 	if orderNo == "" || sku == "" {
// 		handleResponse(w, false, "OrderNo and SKU are required", nil, http.StatusBadRequest)
// 		return
// 	}

// 	err := app.Service.BeforeReturn.DeleteCodeR(r.Context(), orderNo, sku)
// 	if err != nil {
// 		app.Logger.Error("🚨 Failed to delete CodeR 🚨", zap.Error(err))
// 		handleResponse(w, false, err.Error(), nil, http.StatusInternalServerError)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ CodeR deleted successfully ⭐", nil, http.StatusOK)
// }

// // GetDraftConfirmOrderByOrderNo godoc
// // @Summary Get draft order by order number
// // @Description Retrieve the details of a specific draft order by its order number
// // @ID get-draft-order-by-order-no
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Success 200 {object} api.Response{data=[]response.DraftHeadResponse} "Draft order retrieved successfully"
// // @Failure 404 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /draft-confirm/detail/{orderNo} [get]
// func (app *Application) GetDraftConfirmOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")
// 	result, err := app.Service.BeforeReturn.GetDraftConfirmOrderByOrderNo(r.Context(), orderNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Draft Order Details ========== 📋\n")
// 	utils.PrintDraftOrderDetails(result)
// 	fmt.Printf("\n📋 ========== Draft Order Line Details ========== 📋\n")
// 	for i, line := range result.OrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintDraftOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.OrderLines))
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Draft order retrieved successfully ⭐", result, http.StatusOK)
// }

// // UpdateDraftOrders godoc
// // @Summary Update draft orders
// // @Description Update draft orders and change status to Confirm and Booking
// // @ID update-draft-orders
// // @Tags Draft & Confirm
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Success 200 {object} api.Response{data=[]response.DraftHeadResponse} "Draft orders updated successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /draft-confirm/update-draft/{orderNo} [patch]
// func (app *Application) UpdateDraftOrder(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")
// 	if orderNo == "" {
// 		handleResponse(w, false, "Order number is required", nil, http.StatusBadRequest)
// 		return
// 	}

// 	// Extract userID from claims
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
// 		return
// 	}

// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
// 		return
// 	}

// 	err = app.Service.BeforeReturn.UpdateDraftOrder(r.Context(), orderNo, userID)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// Fetch updated order details
// 	result, err := app.Service.BeforeReturn.GetDraftConfirmOrderByOrderNo(r.Context(), orderNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Draft Orders Updated Successfully ========== 📋\n")
// 	utils.PrintDraftOrderDetails(result)
// 	fmt.Printf("\n📋 ========== Draft Order Line Details ========== 📋\n")
// 	for i, line := range result.OrderLines {
// 		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
// 		utils.PrintDraftOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n🚁 Total lines: %d 🚁\n", len(result.OrderLines))
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Draft orders updated successfully ⭐", result, http.StatusOK)
// }

// // review // แสดงข้อมูลออเดอร์ head+line ทั้งหมดที่กรอกทำรายการเข้ามาในระบบ แบบ paginate
// // @Summary 	Get Paginated Before Return Order
// // @Description Get all Before Return Order with pagination
// // @ID 			Get-BefReturnOrder-Paginated
// // @Tags 		Before Return Order
// // @Accept 		json
// // @Produce 	json
// // @Param       page  query int false "Page number" default(1)
// // @Param       limit query int false "Page size" default(10)
// // @Success 	200 {object} Response{result=[]response.OrderDetail} "Get Paginated Orders"
// // @Failure 	400 {object} Response "Bad Request"
// // @Failure 	404 {object} Response "Not Found"
// // @Failure 	500 {object} Response "Internal Server Error"
// // @Router 		/before-return-order/get-orders [get]
// func (app *Application) GetAllOrderDetails(w http.ResponseWriter, r *http.Request) {

// 	page, limit := utils.ParsePagination(r)

// 	result, err := app.Service.BeforeReturn.GetAllOrderDetails(r.Context(), page, limit)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ Orders retrieved successfully ⭐", result, http.StatusOK)
// }

// // review // แสดงข้อมูล order ที่ทำการคืนมาโดยเลข SO
// // @Summary      Get Before Return Order by SO
// // @Description  Get details of an order by its SO number
// // @ID           GetBySO-BefReturnOrder
// // @Tags         Before Return Order
// // @Accept       json
// // @Produce      json
// // @Param        soNo  path     string  true  "soNo"
// // @Success      200 	  {object} Response{result=[]response.OrderDetail} "Get by SO"
// // @Failure      400      {object} Response "Bad Request"
// // @Failure      404      {object} Response "not found endpoint"
// // @Failure      500      {object} Response "Internal Server Error"
// // @Router       /before-return-order/get-orderbySO/{soNo} [get]
// func (app *Application) GetOrderDetailBySO(w http.ResponseWriter, r *http.Request) {
// 	soNo := chi.URLParam(r, "soNo")
// 	if soNo == "" {
// 		handleError(w, errors.ValidationError("soNo is required"))
// 		return
// 	}

// 	result, err := app.Service.BeforeReturn.GetOrderDetailBySO(r.Context(), soNo)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ Orders retrieved by SO successfully ⭐", result, http.StatusOK)
// }

// // review // ลบรายการคืนแต่ละรายการ
// // @Summary 	Delete Order line
// // @Description Delete an order line
// // @ID 			delete-BeforeReturnOrderLine
// // @Tags 		Before Return Order
// // @Accept 		json
// // @Produce 	json
// // @Param 		orderNo path string true "Order No"
// // @Param 		sku path string true "SKU"
// // @Success 	200 {object} Response{result=string} "Before ReturnOrderLine Deleted"
// // @Failure 	404 {object} Response "Order Not Found"
// // @Failure 	422 {object} Response "Validation Error"
// // @Failure 	500 {object} Response "Internal Server Error"
// // @Router 		/before-return-order/delete-line/{orderNo}/{sku} [delete]
// func (app *Application) DeleteBeforeReturnOrderLine(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")
// 	sku := chi.URLParam(r, "sku")

// 	if err := app.Service.BeforeReturn.DeleteBeforeReturnOrderLine(r.Context(), orderNo, sku); err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	handleResponse(w, true, "⭐ Order lines deleted successfully ⭐", nil, http.StatusOK)
// }
