package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

func (app *Application) TradeReturnRoute(apiRouter *chi.Mux) {
	apiRouter.Post("/login", app.Login)

	apiRouter.Route("/trade-return", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		/******** Trade Retrun ********/
		r.Get("/get-waiting", app.GetStatusWaitingDetail) // แสดงข้อมูล ReturnOrder เฉพาะสถานะของ StatusCheckID =1
		r.Get("/get-confirm", app.GetStatusConfirmDetail) // แสดงข้อมูล ReturnOrder เฉพาะสถานะของ StatusCheckID =2
		r.Get("/search-waiting", app.SearchStatusWaitingDetail) // แสดงข้อมูล ReturnOrder เฉพาะสถานะของ StatusCheckID =1 ตามช่วงวันที่สร้าง(CreateDate)ที่เลือก วันที่เริ่มต้น-สิ้นสุด แสดงข้อมูลจำนวนตามวันที่นั้น
		r.Get("/search-confirm", app.SearchStatusConfirmDetail)  // แสดงข้อมูล ReturnOrder เฉพาะสถานะของ StatusCheckID =2 ตามช่วงวันที่สร้าง(CreateDate)ที่เลือก วันที่เริ่มต้น-สิ้นสุด แสดงข้อมูลจำนวนตามวันที่นั้น
		r.Post("/create-trade", app.CreateTradeReturn)
		r.Post("/add-line/{orderNo}", app.CreateTradeReturnLine)
		r.Post("/confirm-receipt/{identifier}", app.ConfirmReceipt)
		r.Patch("/confirm-return/{orderNo}", app.ConfirmReturn)
	})

}

// @Summary Get Return Orders with StatusCheckID = 1
// @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting)
// @ID Get-Waiting-ReturnOrder
// @Tags Trade Return
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]response.ReturnOrder} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /trade-return/get-waiting [get]
func (api *Application) GetStatusWaitingDetail(w http.ResponseWriter, r *http.Request) {
	result, err := api.Service.ReturnOrder.GetReturnOrdersByStatus(r.Context(), 1) // StatusCheckID = 1
	if err != nil {
		handleError(w, err)
		return
	}

	if len(result) == 0 {
		handleResponse(w, true, "No return orders found with StatusCheckID = 1", []res.ReturnOrder{}, http.StatusOK)
		return
	}

	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n======== Order #%d ========\n", i+1)
		utils.PrintReturnOrderDetails(&order)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 1 retrieved successfully ⭐", result, http.StatusOK)
}

// @Summary Get Return Orders with StatusCheckID = 2
// @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed)
// @ID Get-Confirm-ReturnOrder
// @Tags Trade Return
// @Accept json
// @Produce json
// @Success 200 {object} Response{result=[]response.ReturnOrder} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /trade-return/get-confirm [get]
func (api *Application) GetStatusConfirmDetail(w http.ResponseWriter, r *http.Request) {
	result, err := api.Service.ReturnOrder.GetReturnOrdersByStatus(r.Context(), 2) // StatusCheckID = 2
	if err != nil {
		handleError(w, err)
		return
	}

	if len(result) == 0 {
		handleResponse(w, true, "No return orders found with StatusCheckID = 2", []res.ReturnOrder{}, http.StatusOK)
		return
	}

	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n======== Order #%d ========\n", i+1)
		utils.PrintReturnOrderDetails(&order)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 2 retrieved successfully ⭐", result, http.StatusOK)
}

// @Summary Search Return Orders with StatusCheckID = 1 by Date Range
// @Description Retrieve Return Orders with StatusCheckID = 1 (Waiting) within a specific date range
// @ID Search-Waiting-ReturnOrder
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} Response{result=[]response.ReturnOrder} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /trade-return/search-waiting [get]
func (api *Application) SearchStatusWaitingDetail(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	result, err := api.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(r.Context(), 1, startDate, endDate) // StatusCheckID = 1
	if err != nil {
		handleError(w, err)
		return
	}

	if len(result) == 0 {
		handleResponse(w, true, "No return orders found with StatusCheckID = 1 within the specified date range", []res.ReturnOrder{}, http.StatusOK)
		return
	}

	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n======== Order #%d ========\n", i+1)
		utils.PrintReturnOrderDetails(&order)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 1 retrieved successfully ⭐", result, http.StatusOK)
}

// @Summary Search Return Orders with StatusCheckID = 2 by Date Range
// @Description Retrieve Return Orders with StatusCheckID = 2 (Confirmed) within a specific date range
// @ID Search-Confirm-ReturnOrder
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param startDate query string true "Start Date (YYYY-MM-DD)"
// @Param endDate query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} Response{result=[]response.ReturnOrder} "Success"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /trade-return/search-confirm [get]
func (api *Application) SearchStatusConfirmDetail(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	result, err := api.Service.ReturnOrder.GetReturnOrdersByStatusAndDateRange(r.Context(), 2, startDate, endDate) // StatusCheckID = 2
	if err != nil {
		handleError(w, err)
		return
	}

	if len(result) == 0 {
		handleResponse(w, true, "No return orders found with StatusCheckID = 2 within the specified date range", []res.ReturnOrder{}, http.StatusOK)
		return
	}

	fmt.Printf("\n📋 ========== All Return Orders (%d) ========== 📋\n", len(result))
	for i, order := range result {
		fmt.Printf("\n======== Order #%d ========\n", i+1)
		utils.PrintReturnOrderDetails(&order)
	}
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Return Orders with StatusCheckID = 2 retrieved successfully ⭐", result, http.StatusOK)
}

// @Summary Create a new trade return order
// @Description Create a new trade return order with multiple order lines
// @ID create-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Trade Return Detail"
// @Success 201 {object} api.Response{result=response.BeforeReturnOrderResponse} "Trade return created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/create-trade [post]
func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {

	var req request.BeforeReturnOrder

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Error("Failed to decode request", zap.Error(err))
		handleResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// ตรวจสอบว่ามี OrderNo
	if req.OrderNo == "" {
		handleResponse(w, false, "OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	existingOrder, err := app.Service.BeforeReturn.GetBeforeReturnOrderByOrderNo(r.Context(), req.OrderNo)
	if err != nil {
		handleError(w, err)
		return
	}
	if existingOrder != nil { // แจ้งเตือนถ้ามี OrderNo อยู่แล้ว
		handleResponse(w, false, "Order already exists", nil, http.StatusConflict)
		return
	}

	// Authentication check
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleResponse(w, false, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	// Set user information from claims
	req.CreateBy = userID

	// Call service
	result, err := app.Service.BeforeReturn.CreateTradeReturn(r.Context(), req)
	if err != nil {
		app.Logger.Error("Failed to create trade return",
			zap.Error(err),
			zap.String("orderNo", req.OrderNo))

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

	fmt.Printf("\n📋 ========== Created Trade Return Order ========== 📋\n")
	fmt.Printf("\n📋 ========== StatusReturn => 3 (booking) ========== 📋\n\n")
	utils.PrintOrderDetails(result)
	for i, line := range result.BeforeReturnOrderLines {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result.BeforeReturnOrderLines))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Created trade return order successfully ⭐", result, http.StatusOK)
}

// @Summary Add a new trade return line to an existing order
// @Description Add a new trade return line based on the provided order number and line details
// @ID add-trade-return-line
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.TradeReturnLine true "Trade Return Line Details"
// @Success 201 {object} api.Response{result=response.BeforeReturnOrderLineResponse} "Trade return line created successfully"
// @Failure 400 {object} api.Response "Bad Request - Invalid input or missing required fields"
// @Failure 404 {object} api.Response "Not Found - Order not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/add-line/{orderNo} [post]
func (app *Application) CreateTradeReturnLine(w http.ResponseWriter, r *http.Request) {

	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, fmt.Errorf("OrderNo is required"))
		return
	}

	var req request.TradeReturnLine

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, fmt.Errorf("invalid request format: %v", err))
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
	for i := range req.TradeReturnLine {
		req.TradeReturnLine[i].CreateBy = userID
	}

	// เรียก service layer เพื่อสร้างข้อมูล
	result, err := app.Service.BeforeReturn.CreateTradeReturnLine(r.Context(), orderNo, req)
	if err != nil {
		app.Logger.Error("Failed to create trade return",
			zap.Error(err),
			zap.String("orderNo", orderNo))

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

	fmt.Printf("\n📋 ========== Created Trade Return Line Order ========== 📋\n")
	for i, line := range result {
		fmt.Printf("\n📦 Order Line #%d 📦\n", i+1)
		utils.PrintOrderLineDetails(&line)
	}
	fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(result))
	fmt.Println("=====================================")

	handleResponse(w, true, "⭐ Created trade return line successfully ⭐", result, http.StatusCreated)
}

// ConfirmTradeReturn godoc
// @Summary Confirm Receipt from Ware House
// @Description Confirm a trade return order based on the provided identifier (OrderNo or TrackingNo) and input lines for ReturnOrderLine.
// @ID confirm-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param identifier path string true "OrderNo or TrackingNo"
// @Param request body request.ConfirmTradeReturnRequest true "Trade return request details"
// @Success 200 {object} api.Response{result=response.ConfirmReceipt} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/confirm-receipt/{identifier} [post]
func (app *Application) ConfirmReceipt(w http.ResponseWriter, r *http.Request) {
	// รับค่า identifier จาก URL parameter
	identifier := chi.URLParam(r, "identifier")
	if identifier == "" {
		handleError(w, fmt.Errorf("identifier (OrderNo or TrackingNo) is required"))
		return
	}

	// แปลงข้อมูล JSON
	var req request.ConfirmTradeReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// กำหนดค่า identifier
	req.Identifier = identifier

	// รับข้อมูล claims จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	// ดึง userID จาก claims
	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	// เรียก service layer เพื่อดำเนินการ confirm
	err = app.Service.BeforeReturn.ConfirmReceipt(r.Context(), req, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := res.ConfirmReceipt{
		Identifier:     req.Identifier,
		StatusReturnID: "7 (WAITING)",
		StatusCheckID:  "1 (WAITING)",
		UpdateBy:       userID,
		UpdateDate:     time.Now(),
	}

	handleResponse(w, true, "⭐ Confirmed from Ware House successfully ⭐", response, http.StatusOK)
}

// ConfirmToReturn godoc
// @Summary Confirm Return Order to Success
// @Description Confirm a trade return order based on the provided order number (OrderNo) and input lines for ReturnOrderLine.
// @ID confirm-to-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "OrderNo"
// @Param request body request.ConfirmToReturnRequest true "Updated trade return request details"
// @Success 200 {object} api.Response{result=response.ConfirmToReturnOrder} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/confirm-return/{orderNo} [patch]
func (app *Application) ConfirmReturn(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		handleError(w, fmt.Errorf("OrderNo is required"))
		return
	}

	var req request.ConfirmToReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	req.OrderNo = orderNo

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	userID, err := utils.GetUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := app.Service.BeforeReturn.ConfirmReturn(r.Context(), req, userID); err != nil {
		handleError(w, err)
		return
	}

	response := res.ConfirmToReturnOrder{
		OrderNo:        req.OrderNo,
		StatusReturnID: "6 (success)",
		StatusCheckID:  "2 (CONFIRM)",
		UpdateBy:       userID,
		UpdateDate:     time.Now(),
	}
	handleResponse(w, true, "⭐ Confirmed to Return Order successfully ⭐", response, http.StatusOK)
}

// // @Summary Confirm the return order and upload image
// // @Description This API confirms the return order and allows uploading an image.
// // @Tags Trade Return
// // @Accept multipart/form-data
// // @Produce json
// // @Param identifier path string true "OrderNo or TrackingNo"
// // @Param request body request.ConfirmTradeReturnRequest true "Confirm Return Order Data"
// // @Param image formData file true "Image File to Upload"
// // @Success 200 {object} api.Response{data=response.ConfirmReceipt} "Trade return order confirmed successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /trade-return/confirm-receipt/{identifier} [post]
// func (app *Application) ConfirmReceipt(w http.ResponseWriter, r *http.Request) {
// 	// รับค่า identifier จาก URL parameter
// 	identifier := chi.URLParam(r, "identifier")
// 	if identifier == "" {
// 		handleError(w, fmt.Errorf("identifier (OrderNo or TrackingNo) is required"))
// 		return
// 	}

// 	// parse multipart form (รับข้อมูล json + ไฟล์)
// 	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
// 	if err != nil {
// 		handleError(w, fmt.Errorf("unable to parse multipart form: %w", err))
// 		return
// 	}

// 	// รับข้อมูล JSON
// 	var req request.ConfirmTradeReturnRequest
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		handleError(w, fmt.Errorf("failed to read request body: %w", err))
// 		return
// 	}
// 	fmt.Println("Request Body:", string(body)) // ดูข้อมูลที่รับมาจาก body

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, fmt.Errorf("invalid request body: %w", err))
// 		return
// 	}

// 	// รับไฟล์ภาพจากฟอร์ม
// 	files := r.MultipartForm.File["images"] // key 'images' ใช้ในการส่งไฟล์
// 	if len(files) == 0 {
// 		handleError(w, fmt.Errorf("no images uploaded"))
// 		return
// 	}

// 	// กำหนดค่า identifier
// 	req.Identifier = identifier

// 	// รับข้อมูล claims จาก JWT token
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
// 		return
// 	}

// 	// ดึง userID จาก claims
// 	userID, err := utils.GetUserIDFromClaims(claims)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// อัปโหลดไฟล์และรับเส้นทาง
// 	filePaths := []string{}
// 	for _, file := range files {
// 		filePath, err := uploadImageFile(file)
// 		if err != nil {
// 			handleError(w, err)
// 			return
// 		}
// 		filePaths = append(filePaths, filePath)
// 	}

// 	// เรียก service layer เพื่อดำเนินการ confirm
// 	err = app.Service.BefRO.ConfirmReceipt(r.Context(), req, userID, filePaths)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	response := res.ConfirmReceipt{
// 		Identifier: req.Identifier,
// 		UpdateBy:   userID,
// 		UpdateDate: time.Now(),
// 	}

// 	handleResponse(w, true, "Trade return order confirmed successfully", response, http.StatusOK)
// }

// func uploadImageFile(file *multipart.FileHeader) (string, error) {
// 	// กำหนดที่อยู่โฟลเดอร์ที่จะเก็บไฟล์
// 	uploadDir := "uploads/images/"

// 	// สร้างโฟลเดอร์หากยังไม่มี
// 	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
// 		return "", fmt.Errorf("failed to create upload directory: %w", err)
// 	}

// 	// สร้างชื่อไฟล์ใหม่ (เพิ่ม timestamp เพื่อหลีกเลี่ยงการซ้ำชื่อไฟล์)
// 	timestamp := time.Now().UnixNano()
// 	fileName := fmt.Sprintf("%d-%s", timestamp, file.Filename)

// 	// สร้าง path ของไฟล์ที่จะเก็บ
// 	filePath := filepath.Join(uploadDir, fileName)

// 	// เปิดไฟล์ที่อัปโหลด
// 	srcFile, err := file.Open()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open uploaded file: %w", err)
// 	}
// 	defer srcFile.Close()

// 	// สร้างไฟล์เป้าหมายที่จะแนบไฟล์
// 	destFile, err := os.Create(filePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create file: %w", err)
// 	}
// 	defer destFile.Close()

// 	// คัดลอกข้อมูลจากไฟล์ต้นทางไปยังไฟล์เป้าหมาย
// 	_, err = destFile.ReadFrom(srcFile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to copy file data: %w", err)
// 	}

// 	return filePath, nil
// }
