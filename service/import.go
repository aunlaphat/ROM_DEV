package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/errors"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, string, error)
	SaveImageMetadata(ctx context.Context, image request.Images) (int, error)
}

func (srv service) SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	srv.logger.Info("🏁 Starting to search OrderNo or TrackingNo", zap.String("Search", search))

	// เรียก repository เพื่อค้นหาข้อมูล
	order, err := srv.importOrderRepo.SearchOrderORTracking(ctx, search)
	if err != nil {
		srv.logger.Error("❌ Failed to search OrderNo or TrackingNo", zap.Error(err))
		return nil, err
	}

	if order == nil {
		srv.logger.Info("❗ No OrderNo or TrackingNo order found", zap.String("Search", search))
		return nil, nil
	}

	srv.logger.Info("✅ Successfully searched OrderNo or TrackingNo", zap.String("Search", search))
	return []response.ImportOrderResponse{*order}, nil
}



// GetReturnDetailsFromSaleOrder retrieves ReturnID and OrderNo based on SoNo
func (srv service) GetReturnDetailsFromSaleOrder(ctx context.Context, soNo string) (string, string, error) {
	srv.logger.Info("Service: Fetching ReturnID and OrderNo from SoNo", zap.String("SoNo", soNo))

	returnID, orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, soNo)
	if err != nil {
		srv.logger.Error("Service: Error fetching ReturnID and OrderNo", zap.Error(err))
		return "", "", errors.ValidationError(fmt.Sprintf("Invalid SoNo: %s", soNo))
	}

	srv.logger.Info("Service: Successfully fetched ReturnID and OrderNo", zap.String("ReturnID", returnID), zap.String("OrderNo", orderNo))
	return returnID, orderNo, nil
}

// SaveImageMetadata saves image metadata to the database
func (srv service) SaveImageMetadata(ctx context.Context, image request.Images) (int, error) {
	srv.logger.Info("Service: Saving image metadata", zap.Any("Image", image))

	imageID, err := srv.importOrderRepo.InsertImageMetadata(ctx, image)
	if err != nil {
		srv.logger.Error("Service: Error saving image metadata", zap.Error(err))
		return 0, errors.InternalError("Failed to save image metadata")
	}

	srv.logger.Info("Service: Successfully saved image metadata", zap.Int("ImageID", imageID))
	return imageID, nil
}
