package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"

	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	})
}

// SearchSaleOrder godoc
// @Summary Search order by OrderNo or TrackingNo
// @Description Retrieve the details of an order by its OrderNo or TrackingNo using a single input
// @ID search-orderNo-or-trackingNo-single
// @Tags Import Order
// @Accept json
// @Produce json
// @Param search query string true "OrderNo or TrackingNo"
// @Success 200 {object} api.Response{data=response.ImportOrderResponse} "Order retrieved successfully"
// @Failure 400 {object} api.Response "Bad Request"
// @Failure 404 {object} api.Response "OrderNo or TrackingNo not found"
// @Failure 500 {object} api.Response "Internal Server Error"
// @Router /import-order/search [get]
func (app *Application) SearchOrderORTracking(w http.ResponseWriter, r *http.Request) {
	// รับค่าจาก query parameter `search`
	search := r.URL.Query().Get("search")
	if search == "" {
		handleResponse(w, false, "Search input is required (OrderNo or TrackingNo)", nil, http.StatusBadRequest)
		return
	}

	// Trim input
	search = strings.TrimSpace(search)

	// ตรวจสอบ JWT Token (Authorization)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		handleResponse(w, false, "Unauthorized access", nil, http.StatusUnauthorized)
		return
	}

	// เรียกใช้ service layer
	result, err := app.Service.ImportOrder.SearchOrderORTracking(r.Context(), search)
	if err != nil {
		handleResponse(w, false, "Internal server error", nil, http.StatusInternalServerError)
		return
	}

	// หากไม่พบข้อมูล
	if result == nil || len(result) == 0 {
		handleResponse(w, false, "No orders found for the given input", nil, http.StatusNotFound)
		return
	}

	// ส่งข้อมูลกลับ
	handleResponse(w, true, "Orders retrieved successfully", result, http.StatusOK)
}



// UploadImages handles image upload requests
// UploadImagesHandler godoc
// @Summary Upload images for import order
// @Description Upload multiple images for a specific SoNo
// @Tags Import Order
// @Accept multipart/form-data
// @Produce json
// @Param soNo formData string true "Sale Order Number"
// @Param imageTypeID formData int true "Type of the image (1, 2, or 3)"
// @Param sku formData string false "SKU (Optional)"
// @Param files formData file true "Files to upload"
// @Success 200 {object} Response "Successful upload"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /importorder/salereturn [post]
func (app *Application) UploadImages(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		handleError(w, errors.ValidationError("Unable to parse form data"))
		return
	}

	soNo := r.FormValue("soNo")
	if soNo == "" {
		handleError(w, errors.ValidationError("SoNo is required"))
		return
	}

	imageTypeID, err := strconv.Atoi(r.FormValue("imageTypeID"))
	if err != nil || imageTypeID < 1 || imageTypeID > 3 {
		handleError(w, errors.ValidationError("Invalid Image Type ID"))
		return
	}

	skus := r.FormValue("skus")
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		handleError(w, errors.ValidationError("No files uploaded"))
		return
	}

	// Get ReturnID and OrderNo from SoNo
	returnID, orderNo, err := app.Service.ImportOrder.GetReturnDetailsFromSaleOrder(r.Context(), soNo)
	if err != nil {
		handleError(w, err)
		return
	}

	var uploadedImages []res.ImageResponse
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			handleError(w, errors.InternalError("Unable to read file"))
			return
		}
		defer src.Close()

		filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
		filePath := filepath.Join("uploads", filename)
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", os.ModePerm); err != nil {
				handleError(w, errors.InternalError("Failed to create uploads directory"))
				return
			}
		}

		dst, err := os.Create(filePath)
		if err != nil {
			handleError(w, errors.InternalError("Failed to create file"))
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			handleError(w, errors.InternalError("Failed to save file data"))
			return
		}

		image := request.Images{
			ReturnID:    returnID,
			OrderNo:     orderNo,
			FilePath:    filePath,
			ImageTypeID: imageTypeID,
			SKU:         skus,
			CreateBy:    "user",
		}
		imageID, err := app.Service.ImportOrder.SaveImageMetadata(r.Context(), image)
		if err != nil {
			handleError(w, errors.InternalError("Failed to save image metadata"))
			return
		}

		uploadedImages = append(uploadedImages, res.ImageResponse{
			ImageID:  imageID,
			FilePath: filePath,
		})
	}

	handleResponse(w, true, "Upload successful", uploadedImages, http.StatusOK)
}
