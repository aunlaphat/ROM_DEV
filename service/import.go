package service

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/errors"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type ImportOrderService interface {
	GetReturnDetailsFromSaleOrder(ctx context.Context, saleOrder string) (string, string, error)
	SaveImageMetadata(ctx context.Context, image request.Images) (int, error)
}

// GetReturnDetailsFromSaleOrder retrieves ReturnID and OrderNo based on SaleOrder
func (srv service) GetReturnDetailsFromSaleOrder(ctx context.Context, saleOrder string) (string, string, error) {
	srv.logger.Info("Service: Fetching ReturnID and OrderNo from SaleOrder", zap.String("SaleOrder", saleOrder))

	returnID, orderNo, err := srv.importOrderRepo.FetchReturnDetailsBySaleOrder(ctx, saleOrder)
	if err != nil {
		srv.logger.Error("Service: Error fetching ReturnID and OrderNo", zap.Error(err))
		return "", "", errors.ValidationError(fmt.Sprintf("Invalid SaleOrder: %s", saleOrder))
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
