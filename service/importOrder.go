package service

import (
	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/dto/response"
	"boilerplate-back-go-2411/errors"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	SearchOrderORTrackingNo(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	GetOrderTracking(ctx context.Context) ([]response.ImportItem, error)
	UploadPhotoHandler(ctx context.Context, orderNo, imageTypeID, sku string, file io.Reader, filename string) error
	GetSummaryImportOrder(ctx context.Context, orderNo string) ([]response.ImportOrderSummary, error)
	ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error)

	// ยังไม่ใช้
	GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error)
	SaveImageMetadata(ctx context.Context, image request.Images) (int, error)
	ConfirmFromWH(ctx context.Context, soNo string, imageTypeID int, skus string, files []*multipart.FileHeader) ([]response.ImageResponse, error)
	SaveImage(file *multipart.FileHeader) (string, error)
}

func (srv service) SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	srv.logger.Info("[ Starting search order process ]", zap.String("Search", search))

	// *️⃣ ตรวจสอบว่า search มีอยู่จริงในฐานข้อมูลหรือไม่
	exists, err := srv.importOrderRepo.CheckSearch(ctx, search)
	if err != nil {
		srv.logger.Error("[ Failed to check OrderNo or TrackingNo existence ]", zap.String("Search", search), zap.Error(err))
		return nil, errors.InternalError("[ Failed to check OrderNo or TrackingNo existence: %v ]", err)
	}
	if !exists {
		srv.logger.Warn("[ No orders found ]", zap.String("Search", search))
		return nil, errors.NotFoundError("[ No orders found : %s ]", search)
	}

	// *️⃣ ค้นหา order จาก repository (เรียกใช้แบบ chunking)
	orders, err := srv.importOrderRepo.SearchOrderORTracking(ctx, search)
	if err != nil {
		srv.logger.Error("[ Failed to search OrderNo or TrackingNo ]", zap.String("Search", search), zap.Error(err))
		return nil, errors.InternalError("[ Failed to search OrderNo or TrackingNo: %v ]", err)
	}

	srv.logger.Info("[ Successfully search order detail ]")
	return orders, nil
}

func (srv service) SearchOrderORTrackingNo(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	srv.logger.Info("[ Starting search order process ]", zap.String("Search", search))

	// *️⃣ ค้นหา order จาก repository (เรียกใช้แบบ chunking)
	orders, err := srv.importOrderRepo.SearchOrderORTrackingNo(ctx, search)
	if err != nil {
		srv.logger.Error("[ Failed to search OrderNo or TrackingNo ]", zap.String("Search", search), zap.Error(err))
		return nil, errors.InternalError("[ Failed to search OrderNo or TrackingNo: %v ]", err)
	}

	srv.logger.Info("[ Successfully search order detail ]")
	return orders, nil
}

func (srv service) GetOrderTracking(ctx context.Context) ([]response.ImportItem, error) {
	order, err := srv.importOrderRepo.GetOrderTracking(ctx)
	if err != nil {
		srv.logger.Error("[ Error fetching rders ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch orders: %v ]", err)
	}

	srv.logger.Info("[ Fetched all orders ]", zap.Int("Total amount of data", len(order)))
	return order, nil
}

var (
	photoData = make(map[string][]response.ImportOrderSummary) // ข้อมูลภาพ+sku จะบันทึกลงตัวแปรนี้เพื่อนำไปแสดงที่ GetSummaryImportOrder
	mu        sync.Mutex
)

func (srv service) UploadPhotoHandler(ctx context.Context, orderNo, imageTypeID, sku string, file io.Reader, filename string) error {
	srv.logger.Info("[ Starting upload photo process ]", zap.String("OrderNo", orderNo), zap.String("ImageTypeID", imageTypeID), zap.String("SKU", sku), zap.String("Filename", filename))

	// *️⃣ สร้าง path สำหรับบันทึกไฟล์
	dirPath := filepath.Join("uploads/images", orderNo, imageTypeID)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		srv.logger.Error("[ Failed to create directory ]", zap.Error(err))
		return errors.InternalError("[ Failed to create directory: %v ]", err)
	}

	filePath := filepath.Join(dirPath, filename)
	out, err := os.Create(filePath)
	if err != nil {
		srv.logger.Error("[ Failed to create file ]", zap.Error(err))
		return errors.InternalError("[ Failed to create file: %v ]", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		srv.logger.Error("[ Failed to save file ]", zap.Error(err))
		return errors.InternalError("[ Failed to save file: %v ]", err)
	}

	// *️⃣ บันทึกข้อมูลรูปภาพในหน่วยความจำ
	if imageTypeID == "3" {
		mu.Lock()
		defer mu.Unlock()
		photoData[orderNo] = append(photoData[orderNo], response.ImportOrderSummary{
			OrderNo: orderNo,
			SKU:     sku,
			Photo:   filename,
		})
	}

	srv.logger.Info("[ Successfully upload photo process ]", zap.String("OrderNo", orderNo), zap.String("ImageTypeID", imageTypeID), zap.String("SKU", sku), zap.String("Filename", filename))
	return nil
}

func (srv service) GetSummaryImportOrder(ctx context.Context, orderNo string) ([]response.ImportOrderSummary, error) {
	srv.logger.Info("[ Starting get summary import order process ]", zap.String("OrderNo", orderNo))

	mu.Lock()
	defer mu.Unlock()

	summary, exists := photoData[orderNo]
	if !exists {
		srv.logger.Warn("[ no data found ]", zap.String("OrderNo", orderNo))
		return nil, errors.ValidationError("[ no data found for orderNo: %s] ", orderNo)
	}

	srv.logger.Info("[ Successfully get summary import order ]", zap.String("OrderNo", orderNo))
	return summary, nil
}

// *️⃣ เช็คว่า sku ที่รับเข้าค่าตรงกับที่มีในระบบของออเดอร์นั้น หากตรงกันจึงจะยืนยันการรับเข้าได้สำเร็จ
func (srv service) ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error) {
	srv.logger.Info("[ Starting validate SKU process ]", zap.String("OrderNo", orderNo), zap.String("SKU", sku))

	valid, err := srv.importOrderRepo.ValidateSKU(ctx, orderNo, sku)
	if err != nil {
		srv.logger.Error("[ Failed to validate SKU ]", zap.Error(err))
		return false, errors.InternalError("[ Failed to validate SKU: %v ]", err)
	}

	srv.logger.Info("[ Both match: Confirm Receipt]", zap.String("OrderNo", orderNo), zap.String("SKU", sku))
	return valid, nil
}

// ทำรอไว้ยังไม่ได้ข้อสรุปว่าหน้าที่จัดการจะเป็นหน้าอย่างไร ทุกฟังก์ชันหลังบรรทัดนี้

// *️⃣ retrieves ReturnID and OrderNo based on SoNo
func (srv service) GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error) {
	srv.logger.Info("[ Starting get return order process ]", zap.String("SoNo", soNo))

	if soNo == "" {
		srv.logger.Warn("[ SoNo is required ]")
		return "", errors.ValidationError("[ SoNo is required ]")
	}

	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("[ Error fetching OrderNo ]", zap.String("SoNo", soNo), zap.Error(err))
		return "", errors.InternalError("[ Failed to fetch OrderNo: %v ]", err)
	}

	srv.logger.Info("[ Successfully get return order ]", zap.String("SoNo", soNo))
	return orderNo, nil
}

// saves image metadata to the database
func (srv service) SaveImageMetadata(ctx context.Context, image request.Images) (int, error) {
	srv.logger.Info("[ Staring Saving image metadata ]", zap.Any("Image", image))

	// Validate image metadata
	if image.FilePath == "" {
		srv.logger.Warn("[ Invalid image metadata ]")
		return 0, errors.ValidationError("[ FileName and FilePath are required ]")
	}

	// Save to database
	imageID, err := srv.importOrderRepo.InsertImageMetadata(ctx, image)
	if err != nil {
		srv.logger.Error("[ Error saving image metadata ]", zap.Any("Image", image), zap.Error(err))
		return 0, errors.InternalError("[ Failed to save image metadata ]")
	}

	srv.logger.Info("[ Successfully saved image metadata ]", zap.Int("ImageID", imageID))
	return imageID, nil
}

func (srv service) ConfirmFromWH(ctx context.Context, soNo string, imageTypeID int, skus string, files []*multipart.FileHeader) ([]response.ImageResponse, error) {
	srv.logger.Info("[ Starting confirm from warehouse process ]", zap.String("soNo", soNo), zap.Int("imageTypeID", imageTypeID))

	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("[ Error fetching OrderNo ]", zap.String("SoNo", soNo), zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch OrderNo: %v ]", err)
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
			srv.logger.Error("[ Error saving image metadata ]", zap.Any("Image", image), zap.Error(err))
			return nil, errors.InternalError("[ Failed to save image metadata: %v ]", err)
		}

		result = append(result, response.ImageResponse{ImageID: imageID, FilePath: filePath})
	}

	srv.logger.Info("[ Successfully processed image upload ]", zap.Int("Total Images", len(result)))
	return result, nil
}

// Function to save the uploaded file
func (srv service) SaveImage(file *multipart.FileHeader) (string, error) {
	srv.logger.Info("[ Starting save image process ]")

	src, err := file.Open()
	if err != nil {
		return "", errors.InternalError("[ Unable to read file: %v ]", err)
	}
	defer src.Close()

	filename := time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)
	filePath := filepath.Join("uploads", filename)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("uploads", os.ModePerm); err != nil {
			return "", errors.InternalError("[ Failed to create uploads directory: %v ]", err)
		}
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.InternalError("[ Failed to create file: %v ]", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.InternalError("[ Failed to save file data: %v ]", err)
	}

	srv.logger.Info("[ Successfully save image upload ]")
	return filePath, nil
}
