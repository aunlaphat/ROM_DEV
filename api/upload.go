package api

import (
	"github.com/go-chi/chi/v5"
)

// FileServerRoute sets up the file server routes
func (app *Application) FileServerRoute(apiRouter *chi.Mux) {
	apiRouter.Route("/upload", func(r chi.Router) {
		// r.Post("/files-po", app.UploadFilesPo)
		// r.Post("/files-delivery", app.UploadFilesDelivery)
		// r.Post("/store-files-po", app.StoreFilesPo)
		// r.Post("/store-files-delivery", app.StoreFilesDelivery)
	})
}

// // @Summary Upload files
// // @Description Uploads files based on the route and ID.
// // @ID upload-files-po
// // @Tags Upload
// // @Accept multipart/form-data
// // @Produce json
// // @Param files formData []file true "List of files to upload" type=array items=file
// // @Success 200 {object} Response{result=[]string} "Returns the path of the uploaded file"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 404 {object} Response "Article not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /upload/files-po [post]
// func (app *Application) UploadFilesPo(w http.ResponseWriter, r *http.Request) {
// 	app.handleFileUpload(w, r, "uploads/files/po")
// }

// // @Summary Upload files
// // @Description Uploads files based on the route and ID.
// // @ID upload-files-delivery
// // @Tags Upload
// // @Accept multipart/form-data
// // @Produce json
// // @Param files formData []file true "List of files to upload" type=array items=file
// // @Success 200 {object} Response{result=[]string} "Returns the path of the uploaded file"
// // @Failure 400 {object} Response "Bad Request"
// // @Failure 404 {object} Response "Article not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /upload/files-delivery [post]
// func (app *Application) UploadFilesDelivery(w http.ResponseWriter, r *http.Request) {
// 	app.handleFileUpload(w, r, "uploads/files/delivery")
// }

// // @Summary Upload files
// // @Description Uploads files for a specific purchase order (PO) based on the provided file paths and order number.
// // @ID store-files-po
// // @Tags Upload
// // @Accept json
// // @Produce json
// // @Param Input body service.InputFilePath true "Order number and list of file paths to upload"
// // @Success 200 {object} Response{result=[]string} "Returns the paths of the uploaded files"
// // @Failure 400 {object} Response "Bad Request - Invalid input"
// // @Failure 404 {object} Response "Order not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /upload/store-files-po [post]
// func (app *Application) StoreFilesPo(w http.ResponseWriter, r *http.Request) {
// 	// if r.Header.Get("content-type") != "application/json" {
// 	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 	// 	return
// 	// }
// 	req := service.InputFilePath{}
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	res, err := app.Service.Upload.StoreFilesPo(req)
// 	if err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("failed to store file paths: %s", err)})
// 		return
// 	}

// 	// Respond with the retrieved SKUs
// 	handleResponse(w, true, "Upload successed", res, http.StatusOK)
// }

// // @Summary Upload files
// // @Description Uploads files for a specific order based on the provided file paths and order number.
// // @ID store-files-delivery
// // @Tags Upload
// // @Accept json
// // @Produce json
// // @Param Input body service.InputFilePath true "Order number and list of file paths to upload"
// // @Success 200 {object} Response{result=[]string} "Returns the paths of the uploaded files"
// // @Failure 400 {object} Response "Bad Request - Invalid input"
// // @Failure 404 {object} Response "Order not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /upload/store-files-delivery [post]
// func (app *Application) StoreFilesDelivery(w http.ResponseWriter, r *http.Request) {
// 	// if r.Header.Get("content-type") != "application/json" {
// 	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 	// 	return
// 	// }
// 	req := service.InputFilePath{}
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	res, err := app.Service.Upload.StoreFilesDelivery(req)
// 	if err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("failed to store file paths: %s", err)})
// 		return
// 	}

// 	// Respond with the retrieved SKUs
// 	handleResponse(w, true, "Upload successed", res, http.StatusOK)
// }

// func (app *Application) handleFileUpload(w http.ResponseWriter, r *http.Request, directory string) {
// 	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the entire form
// 	if err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("failed to parse form: %s", err)})
// 		return
// 	}

// 	files := r.MultipartForm.File["files"]
// 	if len(files) == 0 {
// 		render.Status(r, http.StatusBadRequest)
// 		render.JSON(w, r, map[string]interface{}{"error": "no files uploaded"})
// 		return
// 	}

// 	if err := utils.CreateDirectoryIfNotExists(directory); err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("create directory error: %s", err)})
// 		return
// 	}

// 	allowedExtensions := []string{".pdf", ".png", ".jpg"} // Adjust as per your requirements

// 	fileReaders := make(map[string]io.Reader)
// 	for _, fileHeader := range files {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			render.Status(r, http.StatusInternalServerError)
// 			render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("failed to open file: %s", err)})
// 			return
// 		}

// 		fileReaders[fileHeader.Filename] = file
// 	}

// 	savedPaths, err := utils.SaveFiles(fileReaders, directory, allowedExtensions)
// 	if err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]interface{}{"error": fmt.Sprintf("failed to save files: %s", err)})
// 		return
// 	}

// 	handleResponse(w, true, "Uploaded success", savedPaths, http.StatusOK)
// }
