package repository

import (
	"boilerplate-backend-go/dto/request"
	"database/sql"

)

type ImportOrderRepository interface {
	InsertImageMetadata(image request.Image) (int, error)
	CheckReturnIDExists(returnID string) (bool, error)
	GetOrderNoByReturnID(returnID string) (string, error)
}

// func handleNull(value string) interface{} {
// 	if value == "" {
// 		return nil // คืนค่า NULL ให้ฐานข้อมูล
// 	}
// 	return value // คืนค่าปกติ
// }

func (repo repositoryDB) InsertImageMetadata(image request.Image) (int, error) {

	//log.Printf("📂 Saving image metadata: %+v", imageMetadata)
	query := `
		INSERT INTO Images (ReturnID, SKU, OrderNo, FilePath, ImageTypeID, CreateBy, CreateDate)
		VALUES (@ReturnID, @SKU, @OrderNo, @FilePath, @ImageTypeID, @CreateBy, GETDATE());
		SELECT SCOPE_IDENTITY();
	`
	var imageID int
	err := repo.db.QueryRow(query, 
	sql.Named("ReturnID", image.ReturnID),
    sql.Named("SKU", image.SKU),
	sql.Named("OrderNo", image.OrderNo), 
    sql.Named("FilePath", image.FilePath),
    sql.Named("ImageTypeID", image.ImageTypeID),
    sql.Named("CreateBy", image.CreateBy),
	).Scan(&imageID)

	return imageID, err
}

func (repo repositoryDB) GetOrderNoByReturnID(returnID string) (string, error) {
    var orderNo string
    query := `
        SELECT OrderNo
        FROM ReturnOrder
        WHERE ReturnID = @ReturnID
    `
    err := repo.db.QueryRow(query, sql.Named("ReturnID", returnID)).Scan(&orderNo)
    if err != nil {
        return "", err
    }
    return orderNo, nil
}

func (repo repositoryDB) CheckReturnIDExists(returnID string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(1) 
        FROM ReturnOrder
        WHERE ReturnID = @ReturnID
    `
	err := repo.db.QueryRow(query, sql.Named("ReturnID", returnID)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
