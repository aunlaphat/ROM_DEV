package api

import (
	"boilerplate-backend-go/dto/request"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) ReturnOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/return-order", func(r chi.Router) {
		r.Get("/list-orders", app.ListReturnOrders)
		r.Post("/create", app.CreateOrderWithLines)
		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/list-lines/{orderNo}", app.ListBeforeReturnOrderLines)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)
	})
}

// ListReturnOrders godoc
// @Summary List all return orders
// @Description Retrieve a list of all return orders
// @ID list-return-orders
// @Tags Return Order
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/list-orders [get]
func (app *Application) ListReturnOrders(w http.ResponseWriter, r *http.Request) {
	result, err := app.Service.ReturnOrder.ListBeforeReturnOrders(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}

// CreateOrderWithLines godoc
// @Summary Create a new return order with lines
// @Description Create a new return order with the provided details
// @ID create-return-order-with-lines
// @Tags Return Order
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Return order details"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/create [post]
func (app *Application) CreateOrderWithLines(w http.ResponseWriter, r *http.Request) {
	var req request.BeforeReturnOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, err)
		return
	}

	result, err := app.Service.ReturnOrder.CreateOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order created successfully", result, http.StatusCreated)
}

// GetBeforeReturnOrderByOrderNo godoc
// @Summary Get return order by order number
// @Description Retrieve the details of a specific return order by its order number
// @ID get-before-return-order-by-order-no
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/{orderNo} [get]
func (app *Application) GetBeforeReturnOrderByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.ReturnOrder.GetBeforeReturnOrderByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order retrieved successfully", result, http.StatusOK)
}

// ListBeforeReturnOrderLines godoc
// @Summary List all return order lines by order number
// @Description Retrieve a list of all return order lines by order number
// @ID list-before-return-order-lines
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/list-lines/{orderNo} [get]
func (app *Application) ListBeforeReturnOrderLines(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.ReturnOrder.ListBeforeReturnOrderLines(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order lines retrieved successfully", result, http.StatusOK)
}

// GetBeforeReturnOrderLineByOrderNo godoc
// @Summary Get return order line by order number and line number
// @Description Retrieve the details of a specific return order line by its order number and line number
// @ID get-before-return-order-line-by-order-no
// @Tags Return Order
// @Accept json
// @Produce json
// @Param orderNo path string true "Order number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /return-order/line/{orderNo} [get]
func (app *Application) GetBeforeReturnOrderLineByOrderNo(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	result, err := app.Service.ReturnOrder.GetBeforeReturnOrderLineByOrderNo(r.Context(), orderNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, "Order line retrieved successfully", result, http.StatusOK)
}
