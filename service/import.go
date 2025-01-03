package service

import (
	"boilerplate-backend-go/dto/request"
	"fmt"
	// "boilerplate-backend-go/repository"
	"io"
	"os"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SaveFile(file io.Reader, path string) (*os.File, error)
	SaveImageMetadata(image request.Image) (int, error)
	ValidateReturnID(returnID string) bool
}

func (srv service) SaveFile(file io.Reader, path string) (*os.File, error) {
	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (srv service) SaveImageMetadata(image request.Image) (int, error) {
	// ดึง OrderNo โดยใช้ ReturnID
	orderNo, err := srv.importOrderRepo.GetOrderNoByReturnID(image.ReturnID)
	if err != nil {
		srv.logger.Error("Error fetching OrderNo from ReturnID", zap.Error(err))
		return 0, fmt.Errorf("failed to fetch OrderNo for ReturnID %s: %w", image.ReturnID, err)
	}

	// ใส่ OrderNo ลงใน struct image
	image.OrderNo = orderNo

	return srv.importOrderRepo.InsertImageMetadata(image)
}

func (srv service) ValidateReturnID(returnID string) bool {
	// ตรวจสอบกับ Repository Layer ว่า ReturnID มีอยู่จริงหรือไม่
	exists, err := srv.importOrderRepo.CheckReturnIDExists(returnID)
	if err != nil {
		srv.logger.Error("Error validating ReturnID", zap.Error(err))
		return false
	}

	return exists
}
