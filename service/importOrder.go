package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error)
	SaveImageMetadata(ctx context.Context, image request.Images) (int, error)

	ConfirmFromWH(ctx context.Context, soNo string, imageTypeID int, skus string, files []*multipart.FileHeader) ([]response.ImageResponse, error)
	SaveImage(file *multipart.FileHeader) (string, error)
}

func (srv service) SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	srv.logger.Info("🏁 Starting to search OrderNo or TrackingNo", zap.String("Search", search))

	order, err := srv.importOrderRepo.SearchOrderORTracking(ctx, search)
	if err != nil {
		srv.logger.Error("❌ Failed to search OrderNo or TrackingNo",
			zap.String("Search", search),
			zap.Error(err),
		)
		return nil, err
	}

	if order == nil {
		srv.logger.Info("❗ No OrderNo or TrackingNo order found",
			zap.String("Search", search),
		)
		return nil, nil
	}

	srv.logger.Info("✅ Successfully searched OrderNo or TrackingNo",
		zap.String("Search", search),
		zap.Any("Order", order),
	)
	return []response.ImportOrderResponse{*order}, nil
}

// retrieves ReturnID and OrderNo based on SoNo
func (srv service) GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, error) {
	srv.logger.Info("🏁 Service: Fetching OrderNo from SoNo", zap.String("SoNo", soNo))

	// Validate SoNo
	if soNo == "" {
		srv.logger.Error("❌ SoNo is empty")
		return "", errors.ValidationError("SoNo is required")
	}

	// Fetch data from repository
	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("❌ Error fetching OrderNo",
			zap.String("SoNo", soNo),
			zap.Error(err),
		)
		return "", errors.InternalError("Failed to fetch OrderNo")
	}

	srv.logger.Info("✅ Successfully fetched ReturnID and OrderNo",
		zap.String("OrderNo", orderNo),
	)
	return orderNo, nil
}

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
	srv.logger.Info("🏁 Processing Image Upload", zap.String("SoNo", soNo))

	orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
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
			srv.logger.Error("❌ Error saving image metadata", zap.Any("Image", image), zap.Error(err))
			return nil, errors.InternalError("Failed to save image metadata")
		}

		result = append(result, response.ImageResponse{ImageID: imageID, FilePath: filePath})
	}

	srv.logger.Info("✅ Successfully processed image upload", zap.Int("Total Images", len(result)))
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
