package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"

	"github.com/go-chi/chi/v5"
)

// SaleReturnRoute defines the routes for sale return operations
func (app *Application) SaleReturnRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/sale-return", func(r chi.Router) {
		r.Get("/search/{soNo}", app.SearchSaleOrder)
		r.Post("/create", app.CreateSaleReturn)
		//r.Post("/create-sale-return/{soNo}", app.CreateSrAX)
		r.Post("/confirm", app.ConfirmSaleReturn)
		//r.With(RoleBasedAuthorization("Accounting")).Post("/create-cn", app.CreateCreditNote)
	})
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

	fmt.Printf("\n📋 ========== Sale Order Details ==========\n")
	for _, order := range result {
		printSaleReturnOrderDetails(&order)
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

	result, err := app.Service.BefRO.CreateBeforeReturnOrderWithLines(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Created Sale Return Order ==========\n")
	printSaleReturnOrderDetails(result)
	// fmt.Println("=====================================")

	handleResponse(w, true, "Sale return order created successfully", result, http.StatusOK)
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

	result, err := app.Service.BefRO.ConfirmSaleReturn(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("\n📋 ========== Confirmed Sale Return Order ========== 📋\n")
	printSaleReturnOrderDetails(result)
	// fmt.Println("=====================================")

	handleResponse(w, true, "Sale return order confirmed successfully", result, http.StatusOK)
}

func printSaleReturnOrderDetails(order *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SaleOrder: %s\n", order.SaleOrder)
	fmt.Printf("🔄 SaleReturn: %s\n", order.SaleReturn)
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
	fmt.Printf("📦 BeforeReturnOrderLines: %v\n", order.BeforeReturnOrderLines)
}

func printSaleReturnOrderLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("📦 OrderNo: %s\n", line.OrderNo)
	fmt.Printf("🔢 SKU: %s\n", line.SKU)
	fmt.Printf("🔢 QTY: %d\n", line.QTY)
	fmt.Printf("🔢 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}
