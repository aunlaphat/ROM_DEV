package api

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/errors"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ReturnOrderRoute defines the routes for return order operations
func (app *Application) BefRORoute(apiRouter *chi.Mux) {
	apiRouter.Route("/before-return-order", func(r chi.Router) {
		r.Get("/list-orders", app.ListBeforeReturnOrders)
		r.Get("/list-lines", app.ListBeforeReturnOrderLines) // Updated route for listing return order lines without orderNo
		r.Get("/get-order", app.GetAllOrderDetail)
		r.Get("/get-orders", app.GetAllOrderDetails) //with paginate


		r.Get("/{orderNo}", app.GetBeforeReturnOrderByOrderNo)
		r.Get("/line/{orderNo}", app.GetBeforeReturnOrderLineByOrderNo)

		r.Get("/get-orderbySO/{soNo}", app.GetOrderDetailBySO)
		r.Get("/search/{soNo}", app.SearchSaleOrder) // New route for searching sale order

		r.Post("/create", app.CreateBeforeReturnOrderWithLines)
		r.Post("/create-trade", app.CreateTradeReturn)

		r.Patch("/update/{orderNo}", app.UpdateBeforeReturnOrderWithLines) // New route for updating return order with lines
		
		r.Delete("/delete-befodline/{recID}", app.DeleteBeforeReturnOrderLine)
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

	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
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
	page := parseInt(query.Get("page"), 1)   // Default page = 1
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
	res, err := app.Service.BefRO.GetOrderDetailBySO(r.Context(), soNo)
	if err != nil {
		handleError(w, err)
		return
	}

	handleResponse(w, true, response, res, http.StatusOK)
}

// SearchSaleOrder godoc
// @Summary Search sale order by SO number
// @Description Retrieve the details of a sale order by its SO number
// @ID search-sale-order
// @Tags Search Sale Order
// @Accept json
// @Produce json
// @Param soNo path string true "SO number"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/search/{soNo} [get]
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

	handleResponse(w, true, "Order created successfully", result, http.StatusCreated)
}

// @Summary Create a new return order
// @Description Create a new return order
// @ID create-trade-return
// @Tags Before Return Order
// @Accept json
// @Produce json
// @Param body body request.BeforeReturnOrder true "Trade Return Detail"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /before-return-order/create-trade [post]
func (app *Application) CreateTradeReturn(w http.ResponseWriter, r *http.Request) {
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

	handleResponse(w, true, "created trade return successfully", result, http.StatusCreated)
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

	handleResponse(w, true, "Return order deleted successfully", nil, http.StatusOK)
}
