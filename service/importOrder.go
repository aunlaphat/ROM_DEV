package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	UploadPhotoHandler(ctx context.Context, orderNo, imageTypeID, sku string, file io.Reader, filename string) error 
	GetSummaryImportOrder(ctx context.Context, orderNo string) ([]response.ImportOrderSummary, error)
	ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error)

	// ยังไม่ใช้
	GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error)
	SaveImageMetadata(ctx context.Context, image request.Images) (int, error)
	ConfirmFromWH(ctx context.Context, soNo string, imageTypeID int, skus string, files []*multipart.FileHeader) ([]response.ImageResponse, error)
	SaveImage(file *multipart.FileHeader) (string, error)
}

// review
func (srv service) SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "SearchOrderORTracking", zap.String("Search", search))
	defer logFinish("Completed", nil)
	srv.logger.Info("🔎 Starting search order process", zap.String("Search", search))

	// ตรวจสอบ search ว่าเป็นค่าว่างหรือไม่
	search = strings.TrimSpace(search)
	if search == "" {
		err := fmt.Errorf("❌ Search input is required (OrderNo or TrackingNo)")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return nil, err
	}

	// ค้นหา order จาก repository
	order, err := srv.importOrderRepo.SearchOrderORTracking(ctx, search)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to search OrderNo or TrackingNo", zap.String("Search", search), zap.Error(err))
		return nil, fmt.Errorf("failed to search OrderNo or TrackingNo: %w", err)
	}

	// หากไม่พบข้อมูล
	if order == nil {
		err := fmt.Errorf("❗ No OrderNo or TrackingNo order found")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return nil, err
	}

	// เติมค่าของ OrderLines (TrackingNo และ OrderNo)
	for i := range order.OrderLines {
		order.OrderLines[i].TrackingNo = order.TrackingNo
		order.OrderLines[i].OrderNo = order.OrderNo
	}

	logFinish("Success", nil)
	return []response.ImportOrderResponse{*order}, nil
}

var (
    photoData = make(map[string][]response.ImportOrderSummary) // ข้อมูลภาพ+sku จะบันทึกลงตัวแปรนี้เพื่อนำไปแสดงที่ GetSummaryImportOrder
    mu        sync.Mutex
)

// review
func (srv service) UploadPhotoHandler(ctx context.Context, orderNo, imageTypeID, sku string, file io.Reader, filename string) error {
    logFinish := srv.logger.LogAPICall(ctx, "UploadPhoto", zap.String("OrderNo", orderNo), zap.String("ImageTypeID", imageTypeID), zap.String("SKU", sku), zap.String("Filename", filename))
    defer logFinish("Completed", nil)
    srv.logger.Info("🔎 Starting upload photo process 🔎", zap.String("OrderNo", orderNo), zap.String("ImageTypeID", imageTypeID), zap.String("SKU", sku), zap.String("Filename", filename))

    // สร้าง path สำหรับบันทึกไฟล์
    dirPath := filepath.Join("uploads/images", orderNo, imageTypeID)
    if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
        logFinish("Failed", err)
        srv.logger.Error("Failed to create directory", zap.Error(err))
        return err
    }

    filePath := filepath.Join(dirPath, filename)
    out, err := os.Create(filePath)
    if err != nil {
        logFinish("Failed", err)
        srv.logger.Error("Failed to create file", zap.Error(err))
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        logFinish("Failed", err)
        srv.logger.Error("Failed to save file", zap.Error(err))
        return err
    }

    // บันทึกข้อมูลรูปภาพในหน่วยความจำ
    if imageTypeID == "3" {
        mu.Lock()
        defer mu.Unlock()
        photoData[orderNo] = append(photoData[orderNo], response.ImportOrderSummary{
            OrderNo:  orderNo,
            SKU:      sku,
            Photo:    filename,
        })
    }

    logFinish("Success", nil)
    return nil
}

// review
func (srv service) GetSummaryImportOrder(ctx context.Context, orderNo string) ([]response.ImportOrderSummary, error) {
    logFinish := srv.logger.LogAPICall(ctx, "GetSummaryImportOrder", zap.String("OrderNo", orderNo))
    defer logFinish("Completed", nil)
    srv.logger.Info("🔎 Starting get summary import order process 🔎", zap.String("OrderNo", orderNo))

    mu.Lock()
    defer mu.Unlock()

    summary, exists := photoData[orderNo]
    if !exists {
        return nil, fmt.Errorf("no data found for orderNo: %s", orderNo)
    }

    logFinish("Success", nil)
    return summary, nil
}

// review
func (srv service) ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error) {
	logFinish := srv.logger.LogAPICall(ctx, "ValidateSKU", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	defer logFinish("Completed", nil)
	srv.logger.Info("🔎 Starting validate SKU process 🔎", zap.String("OrderNo", orderNo), zap.String("SKU", sku))

	if orderNo == "" || sku == "" {
		err := errors.ValidationError("OrderNo and SKU are required")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return false, err
	}

	valid, err := srv.importOrderRepo.ValidateSKU(ctx, orderNo, sku)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Failed to validate SKU", zap.Error(err))
		return false, errors.InternalError("failed to validate SKU")
	}

	logFinish("Success", nil)
	return valid, nil
}

// retrieves ReturnID and OrderNo based on SoNo
func (srv service) GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error) {
	logFinish := srv.logger.LogAPICall(ctx, "GetReturnDetailsFromSaleOrder", zap.String("SoNo", soNo))
	defer logFinish("Completed", nil)
	srv.logger.Info("🔎 Starting get return order process 🔎", zap.String("SoNo", soNo))

	// Validate SoNo
	if soNo == "" {
		err := fmt.Errorf("❗ SoNo is required")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return "", err
	}

	// Fetch data from repository
	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Error fetching OrderNo", zap.String("SoNo", soNo), zap.Error(err))
		return "", fmt.Errorf("failed to fetch OrderNo: %w", err)
	}

	logFinish("Success", nil)
	return orderNo, nil
}

// ยังไม่ได้ใช้ ทำรอไว้เผื่อใช้
// saves image metadata to the database
func (srv service) SaveImageMetadata(ctx context.Context, image request.Images) (int, error) {
	srv.logger.Info("🏁 Service: Saving image metadata", zap.Any("Image", image))

	// Validate image metadata
	if image.FilePath == "" {
		srv.logger.Error("❌ Invalid image metadata",
			zap.Any("Image", image),
		)
		return 0, errors.ValidationError("FileName and FilePath are required")
	}

	// Save to database
	imageID, err := srv.importOrderRepo.InsertImageMetadata(ctx, image)
	if err != nil {
		srv.logger.Error("❌ Error saving image metadata",
			zap.Any("Image", image),
			zap.Error(err),
		)
		return 0, errors.InternalError("Failed to save image metadata")
	}

	srv.logger.Info("✅ Service: Successfully saved image metadata",
		zap.Int("ImageID", imageID),
	)
	return imageID, nil
}

func (srv service) ConfirmFromWH(ctx context.Context, soNo string, imageTypeID int, skus string, files []*multipart.FileHeader) ([]response.ImageResponse, error) {
	logFinish := srv.logger.LogAPICall(ctx, "ConfirmFromWH")
	defer logFinish("Completed", nil)
	srv.logger.Info("🔎 Starting confirm from WH process 🔎", zap.String("soNo", soNo), zap.Int("imageTypeID", imageTypeID))

	if soNo == "" {
		err := fmt.Errorf("SoNo is required")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return nil, err
	}

	if imageTypeID < 1 || imageTypeID > 3 {
		err := fmt.Errorf("invalid Image Type ID")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return nil, err
	}

	if len(files) == 0 {
		err := fmt.Errorf("no files uploaded")
		logFinish("Failed", err)
		srv.logger.Error(err)
		return nil, err
	}

	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		logFinish("Failed", err)
		srv.logger.Error("❌ Error fetching OrderNo", zap.String("SoNo", soNo), zap.Error(err))
		return nil, errors.InternalError("Failed to fetch OrderNo")
	}

	var result []response.ImageResponse

	for _, file := range files {
		filePath, err := srv.SaveImage(file)
		if err != nil {
			return nil, err
		}

		image := request.Images{
			OrderNo:     orderNo,
			FilePath:    filePath,
			ImageTypeID: imageTypeID,
			SKU:         skus,
			CreateBy:    "user",
		}

		imageID, err := srv.importOrderRepo.InsertImageMetadata(ctx, image)
		if err != nil {
			logFinish("Failed", err)
			srv.logger.Error("❌ Error saving image metadata", zap.Any("Image", image), zap.Error(err))
			return nil, errors.InternalError("Failed to save image metadata")
		}

		result = append(result, response.ImageResponse{ImageID: imageID, FilePath: filePath})
	}

	srv.logger.Info("✅ Successfully processed image upload", zap.Int("Total Images", len(result)))
	logFinish("Success", nil)
	return result, nil
}

// Function to save the uploaded file
func (srv service) SaveImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", errors.InternalError("Unable to read file")
	}
	defer src.Close()

	filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
	filePath := filepath.Join("uploads", filename)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("uploads", os.ModePerm); err != nil {
			return "", errors.InternalError("Failed to create uploads directory")
		}
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.InternalError("Failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.InternalError("Failed to save file data")
	}

	return filePath, nil
}
