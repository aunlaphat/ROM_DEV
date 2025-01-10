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
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *Application) ImportOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/importorder", func(r chi.Router) {
		r.Post("/salereturn", app.UploadImages)
	})
}

// UploadImages handles image upload requests
// UploadImagesHandler godoc
// @Summary Upload images for import order
// @Description Upload multiple images for a specific SaleOrder
// @Tags Import Order
// @Accept multipart/form-data
// @Produce json
// @Param saleOrder formData string true "Sale Order Number"
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

	saleOrder := r.FormValue("saleOrder")
	if saleOrder == "" {
		handleError(w, errors.ValidationError("SaleOrder is required"))
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

	// Get ReturnID and OrderNo from SaleOrder
	returnID, orderNo, err := app.Service.ImportOrder.GetReturnDetailsFromSaleOrder(r.Context(), saleOrder)
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

		image := request.Image{
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
