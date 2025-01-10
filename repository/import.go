package repository

import (
	"boilerplate-backend-go/dto/request"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type ImportOrderRepository interface {
	FetchReturnDetailsBySaleOrder(ctx context.Context, saleOrder string) (string, string, error)
	InsertImageMetadata(ctx context.Context, image request.Image) (int, error)
}

// FetchReturnDetailsBySaleOrder retrieves ReturnID and OrderNo from SaleOrder
func (repo repositoryDB) FetchReturnDetailsBySaleOrder(ctx context.Context, saleOrder string) (string, string, error) {
	log.Printf("Repository: Fetching ReturnID and OrderNo for SaleOrder: %s", saleOrder)

	query := `
		SELECT ReturnID, OrderNo
		FROM ReturnOrder
		WHERE SaleOrder = @SaleOrder
	`

	var returnID, orderNo string
	err := repo.db.QueryRowContext(ctx, query, sql.Named("SaleOrder", saleOrder)).Scan(&returnID, &orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Repository: No records found for SaleOrder: %s", saleOrder)
			return "", "", fmt.Errorf("no records found for SaleOrder: %s", saleOrder)
		}
		log.Printf("Repository: Error querying database - %v", err)
		return "", "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Repository: Successfully fetched ReturnID: %s, OrderNo: %s for SaleOrder: %s", returnID, orderNo, saleOrder)
	return returnID, orderNo, nil
}

// InsertImageMetadata inserts image metadata into the database
func (repo repositoryDB) InsertImageMetadata(ctx context.Context, image request.Image) (int, error) {
	query := `
		INSERT INTO Images (ReturnID, SKU, OrderNo, FilePath, ImageTypeID, CreateBy, CreateDate)
		VALUES (:ReturnID, :SKU, :OrderNo, :FilePath, :ImageTypeID, :CreateBy, GETDATE());
		SELECT SCOPE_IDENTITY();
	`

	params := map[string]interface{}{
		"ReturnID":    image.ReturnID,
		"SKU":         image.SKU,
		"OrderNo":     image.OrderNo,
		"FilePath":    image.FilePath,
		"ImageTypeID": image.ImageTypeID,
		"CreateBy":    image.CreateBy,
	}

	var imageID int
	rows, err := repo.db.NamedQuery(query, params)
	if err != nil {
		log.Printf("Repository: Error inserting image metadata - %v", err)
		return 0, fmt.Errorf("error inserting image metadata: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&imageID)
		if err != nil {
			log.Printf("Error scanning inserted image ID - %v", err)
			return 0, fmt.Errorf("error scanning inserted image ID: %w", err)
		}
	}

	log.Printf("Repository: Successfully inserted image metadata with ImageID: %d", imageID)
	return imageID, nil
}

