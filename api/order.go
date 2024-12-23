package api

import (
	"boilerplate-backend-go/dto/request"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) ReturnOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/return-order", func(r chi.Router) {
		r.Get("/list-orders", app.ListOrders)
		r.Post("/create-return-order", app.CreateOrder)
		r.Get("/{orderNo}", app.GetOrder)
		//r.Put("/{orderNo}/status", app.UpdateStatus)
		//r.Post("/{orderNo}/cancel", app.CancelOrder)
	})
}

// ListOrders godoc
// @Summary List all return orders
// @Description Retrieve a list of all return orders with pagination
// @ID list-orders
// @Tags Return Order
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/list-orders [get]
func (app *Application) ListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	result, err := app.Service.ReturnOrder.ListOrders(r.Context(), page, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// GetOrder godoc
// @Summary Get return order by order number
// @Description Retrieve the details of a specific return order by its order number
// @ID get-order
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/{orderNo} [get]
func (app *Application) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.ReturnOrder.GetOrder(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order retrieved successfully", result, http.StatusOK)
}

// CreateOrder godoc
// @Summary Create a new return order
// @Description Create a new return order with the provided details
// @ID create-return-order
// @Tags Return Order
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Return order details"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/create-return-order [post]
func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.Logger.Error("Failed to decode request",
			zap.String("requestID", requestID),
			zap.Error(err))
		handleError(w, err)
		return
	}

	result, err := app.Service.ReturnOrder.CreateOrder(r.Context(), req)
	if err != nil {
		app.Logger.Error("Failed to create order",
			zap.String("requestID", requestID),
			zap.Error(err))
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order created successfully", result, http.StatusCreated)
}

// UpdateStatus godoc
// @Summary Update return order status
// @Description Update the status of a specific return order
// @ID update-order-status
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.UpdateStatusRequest true "Status update details"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/{orderNo}/status [put]
/* func (app *Application) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	var req request.UpdateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	err := app.Service.ReturnOrder.UpdateStatus(r.Context(), orderNo, req.StatusID, req.UpdateBy)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order status updated successfully", nil, http.StatusOK)
} */

// CancelOrder godoc
// @Summary Cancel return order
// @Description Cancel a specific return order
// @ID cancel-order
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Param body body request.CancelOrderRequest true "Cancel details"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/{orderNo}/cancel [post]
/* func (app *Application) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	var req request.CancelOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	err := app.Service.ReturnOrder.CancelOrder(r.Context(), orderNo, req.CancelBy)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order cancelled successfully", nil, http.StatusOK)
} */
