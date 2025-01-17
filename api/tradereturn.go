package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"

	"encoding/json"
	"fmt"
	"net/http"
	
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) TradeReturnRoute(apiRouter *chi.Mux) {
	apiRouter.Post("/login", app.Login)

	apiRouter.Route("/trade-return", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		/******** Trade Retrun ********/
		r.Post("/confirm/{orderNo}", app.ConfirmToReturn)
		r.Post("/create-trade", app.CreateTradeReturn)
		r.Post("/add-line/{orderNo}", app.AddTradeReturnLine)
		r.Post("/confirm/{identifier}", app.ConfirmTradeReturn)
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
// @Router /trade-return/create-trade [post]
func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder

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

// ConfirmToReturn godoc
// @Summary Confirm a trade return order
// @Description Confirm a trade return order based on the provided order number (OrderNo) and input lines for ReturnOrderLine.
// @ID confirm-to-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param orderNo path string true "OrderNo"
// @Param request body request.ConfirmToReturnRequest true "Updated trade return request details"
// @Success 200 {object} api.Response{data=response.ConfirmToReturnOrder} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/confirm/{orderNo} [post]
func (app *Application) ConfirmToReturn(w http.ResponseWriter, r *http.Request) {
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

    userID, err := getUserIDFromClaims(claims)
    if err != nil {
        handleError(w, err)
        return
    }

    if err := app.Service.BefRO.ConfirmToReturn(r.Context(), req, userID); err != nil {
        handleError(w, err)
        return
    }

    response := res.ConfirmToReturnOrder{
        OrderNo:    req.OrderNo,
        UpdateBy:   userID,
        UpdateDate: time.Now(),
    }
    handleResponse(w, true, "return order confirmed successfully", response, http.StatusOK)
}


// ConfirmTradeReturn godoc
// @Summary Confirm a trade return order
// @Description Confirm a trade return order based on the provided identifier (OrderNo or TrackingNo) and input lines for ReturnOrderLine.
// @ID confirm-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param identifier path string true "OrderNo or TrackingNo"
// @Param request body request.ConfirmTradeReturnRequest true "Trade return request details"
// @Success 200 {object} api.Response{data=response.ConfirmToReturnOrder} "Trade return order confirmed successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /trade-return/confirm/{identifier} [post]
func (app *Application) ConfirmTradeReturn(w http.ResponseWriter, r *http.Request) {
	// 1. รับค่า identifier จาก URL parameter
	identifier := chi.URLParam(r, "identifier")
	if identifier == "" {
		handleError(w, fmt.Errorf("identifier (OrderNo or TrackingNo) is required"))
		return
	}

	// 1. รับค่า identifier จาก request body
	var req request.ConfirmTradeReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: %w", http.StatusBadRequest)
		return
	}

	// 3. กำหนดค่า identifier จาก path parameter
	req.Identifier = identifier

	// 2. ดึงค่า claims จาก JWT token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleError(w, fmt.Errorf("unauthorized: missing or invalid token"))
		return
	}

	// 3. ดึงค่า userID จาก claims
	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	// 4. เรียกใช้ service layer เพื่อดำเนินการ confirm
	err = app.Service.BefRO.ConfirmTradeReturn(r.Context(), req, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	// 5. สร้าง response และส่งกลับ
	response := res.ConfirmToReturnOrder{
		OrderNo:    req.Identifier,
		UpdateBy:   userID,
		UpdateDate: time.Now(),
	}

	handleResponse(w, true, "Trade return order confirmed successfully", response, http.StatusOK)
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
	err = app.Service.BefRO.CancelBeforeReturn(r.Context(), orderNo, req.CancelBy, req.Remark)
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