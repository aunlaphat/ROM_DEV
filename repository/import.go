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
	log.Printf("Repository: Inserting image metadata for ReturnID: %s, OrderNo: %s", image.ReturnID, image.OrderNo)

	query := `
		INSERT INTO Images (ReturnID, SKU, FilePath, ImageTypeID, CreateBy, CreateDate, OrderNo)
		VALUES (@ReturnID, @SKU, @FilePath, @ImageTypeID, @CreateBy, GETDATE(), @OrderNo);
		SELECT SCOPE_IDENTITY();
	`

	var imageID int
	err := repo.db.QueryRowContext(ctx, query,
		sql.Named("ReturnID", image.ReturnID),
		sql.Named("SKU", image.SKU),
		sql.Named("FilePath", image.FilePath),
		sql.Named("ImageTypeID", image.ImageTypeID),
		sql.Named("CreateBy", image.CreateBy),
		sql.Named("OrderNo", image.OrderNo),
	).Scan(&imageID)
	if err != nil {
		log.Printf("Repository: Error inserting image metadata - %v", err)
		return 0, fmt.Errorf("failed to insert image metadata: %w", err)
	}

	log.Printf("Repository: Successfully inserted image metadata with ImageID: %d", imageID)
	return imageID, nil
}
