package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) TradeReturnRoute(apiRouter *chi.Mux) {
	apiRouter.Post("/login", app.Login)

	apiRouter.Route("/trade-return", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		/******** Trade Retrun ********/
		r.Post("/create-trade", app.CreateTradeReturn)
		r.Post("/add-line/{orderNo}", app.AddTradeReturnLine)
		r.Post("/confirm-receipt/{identifier}", app.ConfirmReceipt)
		r.Patch("/confirm-return/{orderNo}", app.ConfirmReturn)
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Error("Failed to decode request", zap.Error(err))
		handleResponse(w, false, "Invalid request format", nil, http.StatusBadRequest)
		return
	}

	// ตรวจสอบว่า OrderNo ถูกตั้งค่า
	if req.OrderNo == "" {
		handleResponse(w, false, "OrderNo is required", nil, http.StatusBadRequest)
		return
	}

	existingOrder, err := app.Service.BefRO.GetBeforeReturnOrderByOrderNo(r.Context(), req.OrderNo)
	if err != nil {
		handleError(w, err) // Handle error if GetBeforeReturnOrderByOrderNo fails
		return
	}
	if existingOrder != nil {
		handleResponse(w, false, "Order already exists", nil, http.StatusConflict)
		return
	}

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

	// Set user information from claims
	req.CreateBy = userID

	// 4. Call service
	result, err := app.Service.BefRO.CreateTradeReturn(r.Context(), req)
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
// @Param body body request.TradeReturnLine true "Trade Return Line Details"
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
	err = app.Service.BefRO.CreateTradeReturnLine(r.Context(), orderNo, req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Trade return line created successfully", nil, http.StatusCreated)
}

var validate = validator.New()

// ConfirmTradeReturn godoc
// @Summary Confirm a trade return order
// @Description Confirm a trade return order based on the provided identifier (OrderNo or TrackingNo) and input lines for ReturnOrderLine.
// @ID confirm-trade-return
// @Tags Trade Return
// @Accept json
// @Produce json
// @Param identifier path string true "OrderNo or TrackingNo"
// @Param request body request.ConfirmTradeReturnRequest true "Trade return request details"
// @Success 200 {object} api.Response{data=response.ConfirmReceipt} "Trade return order confirmed successfully"
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
	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	// เรียก service layer เพื่อดำเนินการ confirm
	err = app.Service.BefRO.ConfirmReceipt(r.Context(), req, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := res.ConfirmReceipt{
		Identifier: req.Identifier,
		UpdateBy:   userID,
		UpdateDate: time.Now(),
	}

	handleResponse(w, true, "Trade return order confirmed successfully", response, http.StatusOK)
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

	userID, err := getUserIDFromClaims(claims)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := app.Service.BefRO.ConfirmReturn(r.Context(), req, userID); err != nil {
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
