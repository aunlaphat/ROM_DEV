package api

// import (
// 	"boilerplate-backend-go/dto/request"
// 	res "boilerplate-backend-go/dto/response"
// 	"boilerplate-backend-go/errors"
// 	"boilerplate-backend-go/utils"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"net/http"
// 	"strconv"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/jwtauth"
// 	"go.uber.org/zap"
// )

// // ImportOrder => หน้ารับเข้าสินค้าหน้าคลัง โดยฝั่ง scan สินค้าหน้าคลังจะต้องเข้าทำหน้านี้เพื่อสแกนสินค้าแต่ละ order เข้าระบบ
// func (app *Application) ImportOrderRoute(apiRouter *chi.Mux) {
// 	apiRouter.Post("/login", app.Login)

// 	apiRouter.Route("/import-order", func(r chi.Router) {
// 		// Add auth middleware for protected routes
// 		r.Use(jwtauth.Verifier(app.TokenAuth))
// 		r.Use(jwtauth.Authenticator)

// 		r.Get("/search", app.SearchOrderORTracking) // เสิร์ช orderNo, trackingNo เพื่อทำรายการรับเข้าต่อ (ต้องมีข้อมูลขึ้นจึงทำรายการต่อได้)
// 		r.Post("/upload-photo", app.UploadPhotoHandler) // รับการอัปโหลดภาพเข้าระบบ
// 		r.Get("/summary/{orderNo}", app.GetSummaryImportOrder) // หน้าสรุปรวมข้อมูลภาพที่ถ่ายเข้า ก่อนกดยืนยันรับเข้า (confirm receipt)
// 		r.Post("/validate-sku/{orderNo}/{sku}", app.ValidateSKU) // ใช้ตรวจสอบ sku ที่รับเข้าหน้างานว่าตรงกับในระบบที่ user ทำรายการมา
// 		r.Post("/confirm-receipt/{identifier}", app.ConfirmReceipt) // ยืนยันการรับเข้าหน้าคลัง (เมื่อถ่ายรับเข้าเสร็จทุกรายการ)

// 		r.Post("/create-confirm-wh", app.ConfirmFromWH) // ยังไม่ใช้

// 	})
// }

// // SearchOrderORTracking godoc
// // @Summary 	Search order by OrderNo or TrackingNo
// // @Description Retrieve the details of an order by its OrderNo or TrackingNo
// // @ID 			search-orderNo-or-trackingNo
// // @Tags 		Import Order
// // @Accept 		json
// // @Produce 	json
// // @Param 		search query string true "OrderNo or TrackingNo"
// // @Success 200 {object} api.Response{result=response.ImportOrderResponse} "Order retrieved successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 404 {object} api.Response "OrderNo or TrackingNo not found"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router  /import-order/search [get]
// func (app *Application) SearchOrderORTracking(w http.ResponseWriter, r *http.Request) {
// 	search := r.URL.Query().Get("search")

// 	// ตรวจสอบ JWT Token (Authorization)
// 	_, claims, err := jwtauth.FromContext(r.Context())
// 	if err != nil || claims == nil {
// 		handleResponse(w, false, "🚷 Unauthorized access", nil, http.StatusUnauthorized)
// 		return
// 	}

// 	result, err := app.Service.ImportOrder.SearchOrderORTracking(r.Context(), search)
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		if err.Error() == "❌ Search input is required (OrderNo or TrackingNo)" {
// 			statusCode = http.StatusBadRequest
// 		} else if err.Error() == "❗ No OrderNo or TrackingNo order found" {
// 			statusCode = http.StatusNotFound
// 		}
// 		handleResponse(w, false, err.Error(), nil, statusCode)
// 		return
// 	}

// 	// Debug logging (always print for now, can be controlled by log level later)
// 	fmt.Printf("\n📋 ========== Order Details ========== 📋\n")
// 	for _, order := range result {
// 		utils.PrintImportOrderDetails(&order)
// 		fmt.Printf("\n📋 ========== Order Line Details ========== 📋\n")
// 		for i, line := range order.OrderLines {
// 			fmt.Printf("\n======== Order Line #%d ========\n", i+1)
// 			utils.PrintImportOrderLineDetails(&line)
// 		}
// 		fmt.Printf("\n✳️  Total lines: %d ✳️\n", len(order.OrderLines))
// 		fmt.Println("=====================================")
// 	}

// 	// ส่งข้อมูลกลับ
// 	handleResponse(w, true, "⭐ Found Orders retrieved successfully ⭐", result, http.StatusOK)
// }

//
// // UploadPhotoHandler godoc
// // @Summary Upload Photo
// // @Description Upload a photo for a return order
// // @ID upload-photo
// // @Tags Import Order
// // @Accept multipart/form-data
// // @Produce json
// // @Param orderNo formData string true "Order No"
// // @Param imageTypeID formData string true "ImageTypeID (1, 2, 3)"
// // @Param sku formData string false "SKU (required if imageTypeID is '3')"
// // @Param file formData file true "Photo file"
// // @Success 200 {object} api.Response
// // @Failure 400 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router  /import-order/upload-photo [post]
// func (app *Application) UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
// 	orderNo := r.FormValue("orderNo")
// 	imageTypeID := r.FormValue("imageTypeID")
// 	sku := r.FormValue("sku")

// 	if orderNo == "" || imageTypeID == "" {
// 		handleResponse(w, false, "OrderNo and ImageTypeID are required", nil, http.StatusBadRequest)
// 		return
// 	}

// 	if imageTypeID == "3" && sku == "" {
// 		handleResponse(w, false, "SKU is required for 3 imageTypeID", nil, http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		app.Logger.Error("Failed to get file from request", zap.Error(err))
// 		handleResponse(w, false, "Failed to get file from request", nil, http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	err = app.Service.ImportOrder.UploadPhotoHandler(r.Context(), orderNo, imageTypeID, sku, file, header.Filename)
// 	if err != nil {
// 		app.Logger.Error("Failed to upload photo", zap.Error(err))
// 		handleResponse(w, false, "Failed to upload photo", nil, http.StatusInternalServerError)
// 		return
// 	}

// 	handleResponse(w, true, "File uploaded successfully", nil, http.StatusOK)
// }

//
// // @Summary 	Get Sum detail of Import Order
// // @Description Retrieve the details of Receipt
// // @ID 			GetSummary-ImportOrder
// // @Tags 		Import Order
// // @Accept 		json
// // @Produce 	json
// // @Param orderNo path string true "Order No"
// // @Success 200 {object} api.Response{result=[]response.ImportOrderSummary} "Get All"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 404 {object} Response "Not Found Endpoint"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router 	/import-order/summary/{orderNo} [get]
// func (app *Application) GetSummaryImportOrder(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")

// 	summary, err := app.Service.ImportOrder.GetSummaryImportOrder(r.Context(), orderNo)
// 	if err != nil {
// 		app.Logger.Error("Failed to get summary", zap.Error(err))
// 		handleResponse(w, false, "Failed to get summary", nil, http.StatusInternalServerError)
// 		return
// 	}

// 	handleResponse(w, true, "Summary retrieved successfully", summary, http.StatusOK)
// }

//
// // ValidateSKU godoc
// // @Summary 	Validate SKU
// // @Description Validate SKU
// // @ID 			validate-sku
// // @Tags 		Import Order
// // @Accept 		json
// // @Produce 	json
// // @Param 		orderNo path string true "Order No"
// // @Param 		sku path string true "SKU"
// // @Success 201 {object} api.Response
// // @Failure 400 {object} api.Response
// // @Failure 500 {object} api.Response
// // @Router /import-order/validate-sku/{orderNo}/{sku} [post]
// func (app *Application) ValidateSKU(w http.ResponseWriter, r *http.Request) {
// 	orderNo := chi.URLParam(r, "orderNo")
// 	sku := chi.URLParam(r, "sku")

// 	valid, err := app.Service.ImportOrder.ValidateSKU(r.Context(), orderNo, sku)
// 	if err != nil {
// 		app.Logger.Error("Failed to validate SKU", zap.Error(err))
// 		handleResponse(w, false, "Failed to validate SKU", nil, http.StatusInternalServerError)
// 		return
// 	}

// 	if !valid {
// 		handleResponse(w, false, "SKU not found in Order", nil, http.StatusNotFound)
// 		return
// 	}

// 	handleResponse(w, true, "SKU is valid", nil, http.StatusOK)
// }

//
// // ConfirmReceipt godoc
// // @Summary 	Confirm Receipt from Ware House
// // @Description Confirm a trade return order based on the provided identifier (OrderNo or TrackingNo) and input lines for ReturnOrderLine.
// // @ID 			confirm-trade-return
// // @Tags 		Import Order
// // @Accept 		json
// // @Produce 	json
// // @Param 		identifier path string true "OrderNo or TrackingNo"
// // @Param 		request body request.ConfirmTradeReturnRequest true "Trade return request details"
// // @Success 200 {object} api.Response{result=response.ConfirmReceipt} "Trade return order confirmed successfully"
// // @Failure 400 {object} api.Response "Bad Request"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /import-order/confirm-receipt/{identifier} [post]
// func (app *Application) ConfirmReceipt(w http.ResponseWriter, r *http.Request) {
// 	// รับค่า identifier จาก URL parameter
// 	identifier := chi.URLParam(r, "identifier")

// 	// แปลงข้อมูล JSON
// 	var req request.ConfirmTradeReturnRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		handleError(w, fmt.Errorf("invalid request body: %w", err))
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

// 	// เรียก service layer เพื่อดำเนินการ confirm
// 	err = app.Service.BeforeReturn.ConfirmReceipt(r.Context(), req, userID)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	response := res.ConfirmReceipt{
// 		Identifier:     req.Identifier,
// 		StatusReturnID: "7 (WAITING)",
// 		StatusCheckID:  "1 (WAITING)",
// 		UpdateBy:       userID,
// 		UpdateDate:     time.Now(),
// 	}

// 	handleResponse(w, true, "⭐ Confirmed from Ware House successfully => Status [waiting ✔️] ⭐", response, http.StatusOK)
// }

// // ยังไม่ใช้
// // ConfirmFromWH godoc
// // @Summary 	Import order
// // @Description Upload multiple images and data for a specific SoNo
// // @ID 			Import-Order
// // @Tags 		Import Order
// // @Accept 		multipart/form-data
// // @Produce 	json
// // @Param 		soNo formData string true "Sale Order Number"
// // @Param 		imageTypeID formData int true "Type of the image (1, 2, or 3)"
// // @Param 		skus formData string false "SKU (Optional)"
// // @Param 		files formData file true "Files to upload"
// // @Success 200 {object} api.Response{result=response.ImageResponse} "Successful"
// // @Failure 400 {object} api.Response "Invalid input"
// // @Failure 500 {object} api.Response "Internal Server Error"
// // @Router /import-order/create-confirm-wh [post]
// func (app *Application) ConfirmFromWH(w http.ResponseWriter, r *http.Request) {
// 	// ✅ Parse Form Data
// 	if err := r.ParseMultipartForm(10 << 20); err != nil {
// 		handleError(w, errors.ValidationError("Unable to parse form data"))
// 		return
// 	}

// 	// ✅ รับค่าจาก Form
// 	soNo := r.FormValue("soNo")
// 	imageTypeID, err := strconv.Atoi(r.FormValue("imageTypeID"))
// 	if err != nil {
// 		handleError(w, errors.ValidationError("Invalid Image Type ID"))
// 		return
// 	}

// 	skus := r.FormValue("skus")
// 	files := r.MultipartForm.File["files"]

// 	// ✅ เรียก Service เพื่อประมวลผล
// 	result, err := app.Service.ImportOrder.ConfirmFromWH(r.Context(), soNo, imageTypeID, skus, files)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	// ✅ ส่ง Response กลับไป
// 	handleResponse(w, true, "⭐ Data Insert successful ⭐", result, http.StatusOK)
// }
