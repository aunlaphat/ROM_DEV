 package api

// import (
// 	"boilerplate-backend-go/dto/request"
// 	res "boilerplate-backend-go/dto/response"
// 	"boilerplate-backend-go/utils"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/jwtauth"
// 	"go.uber.org/zap"
// )

// // TradeReturn => ส่วนการคืนของทางฝั่ง offline หาก order ค้างคลัง ต้องการเปลี่ยนรุ่นใหม่ ลูกค้าหน้าสาขาต้องการคืนของ จะทำรายการที่เทรด
// func (app *Application) TradeReturnRoute(apiRouter *chi.Mux) {
// 	apiRouter.Post("/login", app.Login)

// 	apiRouter.Route("/trade-return", func(r chi.Router) {
// 		// Add auth middleware for protected routes
// 		r.Use(jwtauth.Verifier(app.TokenAuth))
// 		r.Use(jwtauth.Authenticator)

// 		/******** Trade Retrun ********/
// 		r.Get("/get-waiting", app.GetStatusWaitingDetail)       // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ waiting
// 		r.Get("/get-confirm", app.GetStatusConfirmDetail)       // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ confirm
// 		r.Get("/search-waiting", app.SearchStatusWaitingDetail) // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ waiting ตามช่วงวันที่สร้าง(CreateDate) เริ่มต้น-สิ้นสุด จะแสดงตามวันที่นั้น
// 		r.Get("/search-confirm", app.SearchStatusConfirmDetail) // แสดงข้อมูลของที่คืนมายังคลังแล้ว สถานะ confirm ตามช่วงวันที่สร้าง(CreateDate) เริ่มต้น-สิ้นสุด จะแสดงตามวันที่นั้น
// 		r.Post("/create-trade", app.CreateTradeReturn)			// สร้างฟอร์มทำรายการคืนของเข้าระบบ
// 		r.Post("/add-line/{orderNo}", app.CreateTradeReturnLine)// สร้างรายการคืนแต่ละรายการของ order นั้นเพิ่ม
// 		r.Patch("/confirm-return/{orderNo}", app.ConfirmReturn) // การยันยืนรับเข้าโดยฝั่งบัญชี จะทำการตรวจเช็คจากข้อมูล confirmReceipt + ข้อมูลในระบบ เมื่อเช็คว่าจำนวนคืนตรงกันจะถือว่าทำรายการคืนสำเร็จ
// 	})
// }


// // @Summary Get Return Orders with StatusCheckID = 1
// // @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting)
// // @ID Get-Waiting-ReturnOrder
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /trade-return/get-waiting [get]
// func (app *Application) GetStatusWaitingDetail(w http.ResponseWriter, r *http.Request) {
// 	result, err :=app.Service.ReturnOrder.GetReturnOrdersByStatus(r.Context(), 1) // StatusCheckID = 1
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n======== Order #%d ========\n", i+1)
// 		utils.PrintDraftTradeOrder(&order)
// 	}
// 	fmt.Println("===============================================")

// 	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 1 (WAITING) retrieved successfully ⭐", result, http.StatusOK)
// }


// // @Summary Get Return Orders with StatusCheckID = 2
// // @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed)
// // @ID Get-Confirm-ReturnOrder
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /trade-return/get-confirm [get]
// func (app *Application) GetStatusConfirmDetail(w http.ResponseWriter, r *http.Request) {
// 	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatus(r.Context(), 2) // StatusCheckID = 2
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n======== Order #%d ========\n", i+1)
// 		utils.PrintDraftTradeOrder(&order)
// 	}
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 2 (CONFIRM) retrieved successfully ⭐", result, http.StatusOK)
// }


// // @Summary Search Return Orders with StatusCheckID = 1 by Date Range
// // @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting) within a specific date range
// // @ID Search-Waiting-ReturnOrder
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Param startDate query string true "Start Date (YYYY-MM-DD)"
// // @Param endDate query string true "End Date (YYYY-MM-DD)"
// // @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /trade-return/search-waiting [get]
// func (app *Application) SearchStatusWaitingDetail(w http.ResponseWriter, r *http.Request) {
// 	// รับค่า startDate และ endDate จาก URL query
// 	startDate := r.URL.Query().Get("startDate")
// 	endDate := r.URL.Query().Get("endDate")

// 	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(r.Context(), 1, startDate, endDate) // StatusCheckID = 1
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// result debug
// 	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n======== Order #%d ========\n", i+1)
// 		utils.PrintDraftTradeOrder(&order)
// 	}
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Return Orders of StatusCheckID = 1 retrieved successfully ⭐", result, http.StatusOK)
// }


// // @Summary Search Return Orders with StatusCheckID = 2 by Date Range
// // @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed) within a specific date range
// // @ID Search-Confirm-ReturnOrder
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Param startDate query string true "Start Date (YYYY-MM-DD)"
// // @Param endDate query string true "End Date (YYYY-MM-DD)"
// // @Success 200 {object} Response{result=[]response.DraftTradeDetail} "Success"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /trade-return/search-confirm [get]
// func (app *Application) SearchStatusConfirmDetail(w http.ResponseWriter, r *http.Request) {
// 	// รับค่า startDate และ endDate จาก URL query
// 	startDate := r.URL.Query().Get("startDate")
// 	endDate := r.URL.Query().Get("endDate")

// 	result, err := app.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(r.Context(), 2, startDate, endDate) // StatusCheckID = 2
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
// 	for i, order := range result {
// 		fmt.Printf("\n======== Order #%d ========\n", i+1)
// 		utils.PrintDraftTradeOrder(&order)
// 	}
// 	fmt.Println("=====================================")

// 	handleResponse(w, true, "⭐ Return Orders of StatusCheckID = 2 retrieved successfully ⭐", result, http.StatusOK)
// }


// // @Summary Create a new trade return order
// // @Description Create a new trade return order with multiple order lines
// // @ID create-trade-return
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Param body body request.BeforeReturnOrder true "Trade Return Detail"
// // @Success 201 {object} api.Response{result=response.BeforeReturnOrderResponse} "Trade return created successfully"
// // @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// // @Failure 404 {object} api.Response "Not Found - Order not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /trade-return/create-trade [post]
// func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {
// 	var req request.BeforeReturnOrder

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		app.Logger.Error("Failed to decode request", zap.Error(err))
// 		handleResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
// 		return
// 	}

// 	// Authentication check
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

// 	// Set user information from claims
// 	req.CreateBy = userID

// 	// Call service
// 	result, err := app.Service.BeforeReturn.CreateTradeReturn(r.Context(), req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Created Trade Return Order ========== 📋\n")
// 	fmt.Printf("\n📋 ========== StatusReturn => 3 (booking) ========= 📋\n\n")
// 	utils.PrintOrderDetails(result)
// 	for i, line := range result.BeforeReturnOrderLines {
// 		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result.BeforeReturnOrderLines))
// 	fmt.Println("===============================================")

// 	handleResponse(w, true, "⭐ Created trade return order successfully => Status [booking ✔️]⭐", result, http.StatusOK)
// }


// // @Summary Add a new trade return line to an existing order
// // @Description Add a new trade return line based on the provided order number and line details
// // @ID add-trade-return-line
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "Order number"
// // @Param body body request.TradeReturnLine true "Trade Return Line Details"
// // @Success 201 {object} api.Response{result=response.BeforeReturnOrderLineResponse} "Trade return line created successfully"
// // @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// // @Failure 404 {object} api.Response "Not Found - Order not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /trade-return/add-line/{orderNo} [post]
// func (app *Application) CreateTradeReturnLine(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	var req request.TradeReturnLine

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, fmt.Errorf("invalid request format: %v", err))
// 		return
// 	}

// 	// ดึงค่า claims จาก JWT token
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
// 		return
// 	}

// 	userID, ok := claims["userID"].(string)
// 	if !ok || userID == "" {
// 		handleError(w, fmt.Errorf("unauthorized: invalid user information"))
// 		return
// 	}

// 	// Set CreateBy จาก claims
// 	for i := range req.TradeReturnLine {
// 		req.TradeReturnLine[i].CreateBy = userID
// 	}

// 	// เรียก service layer เพื่อสร้างข้อมูล
// 	result, err := app.Service.BeforeReturn.CreateTradeReturnLine(r.Context(), orderNo, req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	fmt.Printf("\n📋 ========== Created Trade Return Line Order ========== 📋\n")
// 	fmt.Printf("\n📋 ========== Trade Return Line Order: Latest ========== 📋\n\n")
// 	for i, line := range result {
// 		fmt.Printf("\n======== Order Line #%d ========\n", i+1)
// 		utils.PrintOrderLineDetails(&line)
// 	}
// 	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result))
// 	fmt.Println("===============================================")

// 	handleResponse(w, true, "⭐ Created trade return line successfully ⭐", result, http.StatusCreated)
// }


// // ConfirmToReturn godoc
// // @Summary Confirm Return Order to Success
// // @Description Confirm a trade return order based on the provided order number (OrderNo) and input lines for ReturnOrderLine.
// // @ID confirm-to-return
// // @Tags Trade Return
// // @Accept json
// // @Produce json
// // @Param orderNo path string true "OrderNo"
// // @Param request body request.ConfirmToReturnRequest true "Updated trade return request details"
// // @Success 200 {object} api.Response{result=response.ConfirmToReturnOrder} "Trade return order confirmed successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /trade-return/confirm-return/{orderNo} [Patch]
// func (app *Application) ConfirmReturn(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	var req request.ConfirmToReturnRequest

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, fmt.Errorf("invalid request body: %w", err))
// 		return
// 	}

// 	req.OrderNo = orderNo

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

// 	if err := app.Service.BeforeReturn.ConfirmReturn(r.Context(), req, userID); err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	response := res.ConfirmToReturnOrder{
// 		OrderNo:        req.OrderNo,
// 		StatusReturnID: "6 (success)",
// 		StatusCheckID:  "2 (CONFIRM)",
// 		UpdateBy:       userID,
// 		UpdateDate:     time.Now(),
// 	}

// 	handleResponse(w, true, "⭐ Confirmed to Return Order successfully => Status [success ✔️, confirm ✔️] ⭐", response, http.StatusOK)
// }
