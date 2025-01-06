package service

import (
	"boilerplate-backend-go/dto/request"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type ImportOrderService interface {
	SaveFile(file io.Reader, path string) (*os.File, error)
	SaveImageMetadata(image request.Image) (int, error)
	ValidateReturnID(returnID string) bool
	ValidateImageType(imageType int, returnID, sku string) error
	// ValidateImageSequence(returnID string, imageType int) error
	// ValidateDuplicateFileName(returnID, fileName string, imageType int) error
	ValidateDuplicateFileName(returnID, fileName string) error
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
    // ตรวจสอบไฟล์ซ้ำ
    if err := srv.ValidateDuplicateFileName(image.ReturnID, filepath.Base(image.FilePath)); err != nil {
        return 0, err
    }

    // Insert Metadata
    imageID, err := srv.importOrderRepo.InsertImageMetadataWithOrderNo(image)
    if err != nil {
        srv.logger.Error("Error saving image metadata with OrderNo", zap.Error(err))
        return 0, fmt.Errorf("failed to save image metadata for ReturnID %s: %w", image.ReturnID, err)
    }

    return imageID, nil
}

func (srv service) ValidateReturnID(returnID string) bool {
	exists, err := srv.importOrderRepo.CheckReturnIDExists(returnID)
	if err != nil {
		srv.logger.Error("Error validating ReturnID", zap.Error(err))
		return false
	}

	return exists
}

// ValidateImageType ตรวจสอบลำดับ imageType และข้อมูลที่ต้องการ
func (srv service) ValidateImageType(imageType int, returnID, sku string) error {
    // ตรวจสอบ imageType มีค่า 1, 2, หรือ 3
    if imageType < 1 || imageType > 3 {
        return fmt.Errorf("invalid imageType: %d", imageType)
    }

    // ตรวจสอบว่า imageType = 3 ต้องมี SKU
    if imageType == 3 && (sku == "" || len(sku) == 0) {
        return fmt.Errorf("imageType 3 requires a valid SKU")
    }

    // ตรวจสอบลำดับ imageType ในระบบ (ฐานข้อมูล)
    // if err := srv.ValidateImageSequence(returnID, imageType); err != nil {
    //     return fmt.Errorf("image type sequence validation failed: %w", err)
    // }

    return nil
}

// ตรวจสอบลำดับของ imageType เช่น อัพโหลด imageType=3 โดยไม่มี imageType=1 หรือ imageType=2
// func (srv service) ValidateImageSequence(returnID string, imageType int) error {
//     lastImageType, err := srv.importOrderRepo.GetLastImageTypeByReturnID(returnID)
//     if err != nil {
//         return fmt.Errorf("failed to fetch last image type: %w", err)
//     }

//     // ตรวจสอบเงื่อนไขลำดับ imageType
//     switch imageType {
//     case 1:
//         if lastImageType != 0 {
//             return fmt.Errorf("imageType 1 must be the first image")
//         }
//     case 2:
//         if lastImageType != 1 {
//             return fmt.Errorf("imageType 2 requires a previous imageType 1")
//         }
//     case 3:
//         if lastImageType != 2 && lastImageType != 3 {
//             return fmt.Errorf("imageType 3 requires a previous imageType 2 or 3")
//         }
//     }
//     return nil
// }

func (srv service) ValidateDuplicateFileName(returnID, fileName string) error {
    isDuplicate, err := srv.importOrderRepo.CheckDuplicateFileName(returnID, fileName)
    if err != nil {
        srv.logger.Error("Error checking duplicate file name", zap.Error(err))
        return fmt.Errorf("failed to check duplicate file name: %w", err)
    }

    if isDuplicate {
        return fmt.Errorf("duplicate file name '%s' is not allowed under ReturnID %s", fileName, returnID)
    }

    return nil
}




