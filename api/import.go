package api

import (
	"boilerplate-backend-go/dto/request"
	res "boilerplate-backend-go/dto/response"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// ImageRoute สำหรับกำหนดเส้นทาง API
func (app *Application) ImportOrderRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/images", func(r chi.Router) {
		r.Post("/upload", app.UploadImages) // API สำหรับอัปโหลดภาพ
	})
}

// UploadImages godoc
// @Summary Upload images for product return
// @Description Handle image upload for return process
// @ID upload-images
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Image files to upload"
// @Param returnID formData string true "Return ID"
// @Param skus formData string false "Comma-separated SKUs for product images"
// @Param imageTypeID formData int true "Image Type ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /images/upload [post]
func (app *Application) UploadImages(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบโฟลเดอร์สำหรับบันทึกไฟล์
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if mkdirErr := os.Mkdir("uploads", os.ModePerm); mkdirErr != nil {
			log.Printf("Error creating uploads directory: %v", mkdirErr)
			http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
			return
		}
	}

	// Parse Form Data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing form data: %v", err)
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	returnID := r.FormValue("returnID")
	imageTypeID, err := strconv.Atoi(r.FormValue("imageTypeID"))
	if err != nil {
		log.Printf("Invalid imageTypeID: %v", err)
		http.Error(w, "Invalid Image Type ID", http.StatusBadRequest)
		return
	}
	skus := r.FormValue("skus")
	log.Printf("Received Upload Request: ReturnID=%s, ImageTypeID=%d, SKUs=%s", returnID, imageTypeID, skus)

	// ตรวจสอบว่า ReturnID มีอยู่ในระบบหรือไม่
	if !app.Service.ImportOrder.ValidateReturnID(returnID) {
		http.Error(w, "Invalid Return ID", http.StatusBadRequest)
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
	
		filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
		filePath := filepath.Join("uploads", filename)
		filePath = filepath.ToSlash(filePath)
	
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
	
		// เพิ่ม defer ลบไฟล์เมื่อเกิดข้อผิดพลาด
		defer func() {
			if err != nil {
				os.Remove(filePath)
			}
		}()
	
		_, err = io.Copy(dst, src)
		if err != nil {
			log.Printf("Error saving file: %v", err)
			http.Error(w, "Failed to save file data", http.StatusInternalServerError)
			return
		}
		dst.Close()
	
		// Save Metadata to Database
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
	
			// ลบไฟล์เมื่อบันทึกฐานข้อมูลล้มเหลว
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
}