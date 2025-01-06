package repository

import (
	"boilerplate-backend-go/dto/request"
	"fmt"
	// "database/sql"
)

type ImportOrderRepository interface {
	InsertImageMetadataWithOrderNo(image request.Image) (int, error)
	// GetLastImageTypeByReturnID(returnID string) (int, error)
	CheckReturnIDExists(returnID string) (bool, error)
	// CheckDuplicateFileName(returnID, fileName string, imageType int) (bool, error)
	CheckDuplicateFileName(returnID, fileName string) (bool, error)
}

func (repo repositoryDB) InsertImageMetadataWithOrderNo(image request.Image) (int, error) {
    queryInsert := `
        INSERT INTO Images (ReturnID, SKU, OrderNo, FilePath, ImageTypeID, CreateBy, CreateDate)
        VALUES (
            :ReturnID, 
            :SKU, 
            (SELECT OrderNo FROM ReturnOrder WHERE ReturnID = :ReturnID), -- ดึง OrderNo ระหว่าง Insert
            :FilePath, 
            :ImageTypeID, 
            :CreateBy, 
            GETDATE()
        );
        SELECT SCOPE_IDENTITY();
    `
    params := map[string]interface{}{
        "ReturnID":    image.ReturnID,
        "SKU":         image.SKU,
        "FilePath":    image.FilePath,
        "ImageTypeID": image.ImageTypeID,
        "CreateBy":    image.CreateBy,
    }

    var imageID int
    rows, err := repo.db.NamedQuery(queryInsert, params)
    if err != nil {
        return 0, fmt.Errorf("error inserting image metadata with OrderNo: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&imageID)
        if err != nil {
            return 0, fmt.Errorf("error scanning imageID: %w", err)
        }
    }

    return imageID, nil
}


// func (repo repositoryDB) GetLastImageTypeByReturnID(returnID string) (int, error) {
//     var lastImageType int
//     query := `
//         SELECT TOP 1 ImageTypeID
//         FROM Images
//         WHERE ReturnID = :ReturnID
//         ORDER BY CreateDate DESC
//     `

//     // ใช้ NamedQuery และ Scan ค่า
//     rows, err := repo.db.NamedQuery(query, map[string]interface{}{"ReturnID": returnID})
//     if err != nil {
//         return 0, err
//     }
//     defer rows.Close()

//     if rows.Next() {
//         err = rows.Scan(&lastImageType)
//         if err != nil {
//             return 0, err
//         }
//     }
//     return lastImageType, nil
// }

func (repo repositoryDB) CheckReturnIDExists(returnID string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(1) 
        FROM ReturnOrder
        WHERE ReturnID = :ReturnID
    `
	// ใช้ NamedQuery เพื่อส่งพารามิเตอร์ในรูปแบบ map
	rows, err := repo.db.NamedQuery(query, map[string]interface{}{"ReturnID": returnID})
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// ดึงค่าจากผลลัพธ์
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}
	return count > 0, nil
}

func (repo repositoryDB) CheckDuplicateFileName(returnID, fileName string) (bool, error) {
    var count int
    query := `
        SELECT COUNT(1)
        FROM Images
        WHERE ReturnID = :ReturnID AND FilePath LIKE :FilePath
    `

    params := map[string]interface{}{
        "ReturnID": returnID,
        "FilePath": "%" + fileName, // ตรวจสอบเฉพาะชื่อไฟล์จริง
    }

    rows, err := repo.db.NamedQuery(query, params)
    if err != nil {
        return false, fmt.Errorf("error checking duplicate file name: %w", err)
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.Scan(&count)
        if err != nil {
            return false, fmt.Errorf("error scanning duplicate file name result: %w", err)
        }
    }

    return count > 0, nil
}





