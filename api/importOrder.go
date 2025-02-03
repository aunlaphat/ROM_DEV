package api

import (
	"boilerplate-backend-go/errors"
	"boilerplate-backend-go/utils"
	"fmt"

	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

func (app *Application) ImportOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Post("/login", app.Login)

	apiRouter.Route("/import-order", func(r chi.Router) {
		// Add auth middleware for protected routes
		r.Use(jwtauth.Verifier(app.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/search", app.SearchOrderORTracking)
		r.Post("/create-confirm-wh", app.ConfirmFromWH)

	})
}

// SearchOrderORTracking godoc
// @Summary Search order by OrderNo or TrackingNo
// @Description Retrieve the details of an order by its OrderNo or TrackingNo using a single input
// @ID search-orderNo-or-trackingNo-single
// @Tags Import Order
// @Accept json
// @Produce json
// @Param search query string true "OrderNo or TrackingNo"
// @Success 200 {object} api.Response{result=response.ImportOrderResponse} "Order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "OrderNo or TrackingNo not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /import-order/search [get]
func (app *Application) SearchOrderORTracking(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	// ตรวจสอบ JWT Token (Authorization)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	result, err := app.Service.ImportOrder.SearchOrderORTracking(r.Context(), search)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "❌ Search input is required (OrderNo or TrackingNo)" {
			statusCode = http.StatusBadRequest
		} else if err.Error() == "❗ No OrderNo or TrackingNo order found" {
			statusCode = http.StatusNotFound
		}
		handleResponse(w, false, err.Error(), nil, statusCode)
		return
	}

	// Debug logging (always print for now, can be controlled by log level later)
	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
	for _, order := range result {
		utils.PrintImportOrderDetails(&order)
		fmt.Printf("\n📋 ========== Order Line Details ========== 📋\n")
		for i, line := range order.OrderLines {
			fmt.Printf("\n======== Order Line #%d ========\n", i+1)
			utils.PrintImportOrderLineDetails(&line)
		}
		fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(order.OrderLines))
		fmt.Println("=====================================")
	}

	// ส่งข้อมูลกลับ
	handleResponse(w, true, "⭐ Found Orders retrieved successfully ⭐", result, http.StatusOK)
}

// ConfirmFromWH godoc
// @Summary Import order
// @Description Upload multiple images and data for a specific SoNo
// @ID Import-Order
// @Tags Import Order
// @Accept multipart/form-data
// @Produce json
// @Param soNo formData string true "Sale Order Number"
// @Param imageTypeID formData int true "Type of the image (1, 2, or 3)"
// @Param skus formData string false "SKU (Optional)"
// @Param files formData file true "Files to upload"
// @Success 200 {object} api.Response{result=response.ImageResponse} "Successful"
// @Failure 400 {object} api.Response "Invalid input"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /import-order/create-confirm-wh [post]
func (app *Application) ConfirmFromWH(w http.ResponseWriter, r *http.Request) {
	// ✅ Parse Form Data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		handleError(w, errors.ValidationError("Unable to parse form data"))
		return
	}

	// ✅ รับค่าจาก Form
	soNo := r.FormValue("soNo")
	imageTypeID, err := strconv.Atoi(r.FormValue("imageTypeID"))
	if err != nil {
		handleError(w, errors.ValidationError("Invalid Image Type ID"))
		return
	}

	skus := r.FormValue("skus")
	files := r.MultipartForm.File["files"]

	// ✅ เรียก Service เพื่อประมวลผล
	result, err := app.Service.ImportOrder.ConfirmFromWH(r.Context(), soNo, imageTypeID, skus, files)
	if err != nil {
		handleError(w, err)
		return
	}

	// ✅ ส่ง Response กลับไป
	handleResponse(w, true, "⭐ Data Insert successful ⭐", result, http.StatusOK)
}
