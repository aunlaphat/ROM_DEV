package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const UploadDir = "uploads"

func (app *Application) ImportOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/images", func(r chi.Router) {
		r.Post("/upload", app.UploadImages)
	})
}

// UploadImages godoc
// @Summary Upload images for return order
// @Description Handle image upload for return process
// @ID upload-images
// @Tags Import Order
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Image files to upload"
// @Param returnID formData string true "number of order"
// @Param skus formData string false "SKU for one same model in img"
// @Param imageTypeID formData int true "type image"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /images/upload [post]
func (app *Application) UploadImages(w http.ResponseWriter, r *http.Request) {
	if err := app.createUploadDir(); err != nil {
		log.Printf("Error setting up upload directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing form data: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	returnID, imageTypeID, skus, err := app.validateUploadRequest(r)
	if err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		log.Println("No files uploaded")
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var uploadedImages []res.ImageResponse
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}
		defer src.Close()

		originalFileName := filepath.Base(file.Filename)
		if err := app.Service.ImportOrder.ValidateDuplicateFileName(returnID, originalFileName); err != nil {
			log.Printf("Duplicate file name error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filename := time.Now().Format("20060102_150405") + "_" + originalFileName
		filePath := filepath.Join(UploadDir, filename)
		filePath = filepath.ToSlash(filePath)

		if err := app.saveUploadedFile(src, filePath); err != nil {
			log.Printf("Error saving file: %v", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		log.Printf("File uploaded successfully: %s", filePath)
		image := request.Image{
			ReturnID:    returnID,
			SKU:         skus,
			FilePath:    filePath,
			ImageTypeID: imageTypeID,
			CreateBy:    "user",
		}
		imageID, err := app.Service.ImportOrder.SaveImageMetadata(image)
		if err != nil {
			log.Printf("Error saving image metadata: %v", err)
			if removeErr := os.Remove(filePath); removeErr != nil {
				log.Printf("Error removing file: %v", removeErr)
			}
			http.Error(w, "Failed to save image metadata", http.StatusInternalServerError)
			return
		}

		uploadedImages = append(uploadedImages, res.ImageResponse{
			ImageID:  imageID,
			FilePath: filePath,
		})
	}

	response := res.BaseResponse{
		Success: true,
		Data:    uploadedImages,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (app *Application) createUploadDir() error {
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		if mkdirErr := os.Mkdir(UploadDir, os.ModePerm); mkdirErr != nil {
			return fmt.Errorf("failed to create uploads directory: %w", mkdirErr)
		}
	}
	return nil
}

func (app *Application) saveUploadedFile(file io.ReadCloser, filePath string) error {
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return fmt.Errorf("error saving file: %w", err)
	}
	return nil
}

func (app *Application) validateUploadRequest(r *http.Request) (string, int, string, error) {
	returnID := r.FormValue("returnID")
	imageTypeID, err := strconv.Atoi(r.FormValue("imageTypeID"))
	if err != nil || imageTypeID < 1 || imageTypeID > 3 {
		return "", 0, "", fmt.Errorf("invalid imageTypeID: %v", err)
	}
	skus := r.FormValue("skus")
	return returnID, imageTypeID, skus, nil
}
