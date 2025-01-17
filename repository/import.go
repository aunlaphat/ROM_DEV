package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type ImportOrderRepository interface {
	SearchOrderORTracking(ctx context.Context, search string) (*response.ImportOrderResponse, error)

	FetchReturnDetailsBySaleOrder(ctx context.Context, soNo string) (string, string, error)
	InsertImageMetadata(ctx context.Context, image request.Images) (int, error)
}

func (repo repositoryDB) SearchOrderORTracking(ctx context.Context, search string) (*response.ImportOrderResponse, error) {
	// Query สำหรับดึงข้อมูล Order Head
	queryHead := `
        SELECT OrderNo, SoNo, TrackingNo, CreateDate
        FROM ROM_V_OrderHeadDetail
        WHERE OrderNo = :Search OR TrackingNo = :Search
    `

	// Query สำหรับดึงข้อมูล Order Lines
	queryLines := `
        SELECT SKU, ItemName, QTY, Price
        FROM ROM_V_OrderLineDetail
        WHERE OrderNo = :Search OR TrackingNo = :Search
    `

	// ดึงข้อมูล Order Head
	var orderHead response.ImportOrderResponse
	nstmtHead, err := repo.db.PrepareNamed(queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for order head: %w", err)
	}
	defer nstmtHead.Close()

	err = nstmtHead.GetContext(ctx, &orderHead, map[string]interface{}{
		"Search": search,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order head: %w", err)
	}

	// ดึงข้อมูล Order Lines
	var orderLines []response.ImportOrderLineResponse
	nstmtLines, err := repo.db.PrepareNamed(queryLines)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for order lines: %w", err)
	}
	defer nstmtLines.Close()

	err = nstmtLines.SelectContext(ctx, &orderLines, map[string]interface{}{
		"Search": search,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order lines: %w", err)
	}

	// ผูกข้อมูล Line กับ Order Head
	orderHead.OrderLines = orderLines

	return &orderHead, nil
}



// FetchReturnDetailsBySaleOrder retrieves ReturnID and OrderNo from SoNo
func (repo repositoryDB) FetchReturnDetailsBySaleOrder(ctx context.Context, soNo string) (string, string, error) {
	log.Printf("Repository: Fetching ReturnID and OrderNo for SoNo: %s", soNo)

	query := `
		SELECT ReturnID, OrderNo
		FROM ReturnOrder
		WHERE SoNo = @SoNo
	`

	var returnID, orderNo string
	err := repo.db.QueryRowContext(ctx, query, sql.Named("SoNo", soNo)).Scan(&returnID, &orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Repository: No records found for SoNo: %s", soNo)
			return "", "", fmt.Errorf("no records found for SoNo: %s", soNo)
		}
		log.Printf("Repository: Error querying database - %v", err)
		return "", "", fmt.Errorf("database query error: %w", err)
	}

	log.Printf("Repository: Successfully fetched ReturnID: %s, OrderNo: %s for SoNo: %s", returnID, orderNo, soNo)
	return returnID, orderNo, nil
}

// InsertImageMetadata inserts image metadata into the database
func (repo repositoryDB) InsertImageMetadata(ctx context.Context, image request.Images) (int, error) {
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
