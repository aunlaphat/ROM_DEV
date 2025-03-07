package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
)

type ImportOrderRepository interface {
	SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	SearchOrderORTrackingNo(ctx context.Context, search string) ([]response.ImportOrderResponse, error)
	GetOrderTracking(ctx context.Context) ([]response.ImportItem, error)
	ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error)
	FetchReturnDetailsBySaleOrder(ctx context.Context, soNo string) (string, error)
	CheckSearch(ctx context.Context, search string) (bool, error)

	// ยังไม่ใช้
	InsertImageMetadata(ctx context.Context, image request.Images) (int, error)
}

func (repo repositoryDB) SearchOrderORTracking(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	// *️⃣ จำนวนข้อมูลที่ต้องการดึงในแต่ละ chunk
	const chunkSize = 1000
	var orders []response.ImportOrderResponse
	offset := 0

	for {
		queryHead := `  SELECT OrderNo, SoNo, TrackingNo, CreateDate
						FROM ROM_V_OrderHeadDetail
						WHERE OrderNo = :Search OR TrackingNo = :Search
						ORDER BY OrderNo
						OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
					 `
		var orderHeadBatch []response.ImportOrderResponse
		nstmtHead, err := repo.db.PrepareNamed(queryHead)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement for order head: %w", err)
		}
		defer nstmtHead.Close()

		// *️⃣ ดึงข้อมูล Order Head ในแต่ละ batch
		err = nstmtHead.SelectContext(ctx, &orderHeadBatch, map[string]interface{}{
			"Search": search,
			"Limit":  chunkSize,
			"Offset": offset,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order head: %w", err)
		}

		// *️⃣ ถ้าไม่มีข้อมูลใน batch นี้ ให้หยุดการทำงาน
		if len(orderHeadBatch) == 0 {
			break
		}

		// *️⃣ ดึงข้อมูล Order Lines สำหรับแต่ละ Order Head ใน batch
		for _, orderHead := range orderHeadBatch {

			queryLines := ` SELECT SKU, ItemName, QTY, Price
							FROM ROM_V_OrderLineDetail
							WHERE OrderNo = :Search OR TrackingNo = :Search
						  `

			var orderLines []response.ImportOrderLineResponse
			nstmtLines, err := repo.db.PrepareNamed(queryLines)
			if err != nil {
				return nil, fmt.Errorf("failed to prepare statement for order lines: %w", err)
			}
			defer nstmtLines.Close()

			// *️⃣ ดึงข้อมูล Order Lines
			err = nstmtLines.SelectContext(ctx, &orderLines, map[string]interface{}{
				"Search": search,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch order lines: %w", err)
			}

			// *️⃣ ผูกข้อมูล Line กับ Order Head
			orderHead.OrderLines = orderLines

			// *️⃣ เติมค่าของ OrderLines (TrackingNo และ OrderNo)
			for i := range orderHead.OrderLines {
				orderHead.OrderLines[i].TrackingNo = orderHead.TrackingNo
				orderHead.OrderLines[i].OrderNo = orderHead.OrderNo
			}

			// *️⃣ เพิ่มข้อมูลที่ดึงมาในแต่ละ batch
			orders = append(orders, orderHead)
		}

		// *️⃣ เพิ่ม offset สำหรับ batch ถัดไป
		offset += chunkSize
	}

	return orders, nil
}

func (repo repositoryDB) SearchOrderORTrackingNo(ctx context.Context, search string) ([]response.ImportOrderResponse, error) {
	// *️⃣ จำนวนข้อมูลที่ต้องการดึงในแต่ละ chunk
	const chunkSize = 1000
	var orders []response.ImportOrderResponse
	offset := 0

	for {
		queryHead := `  SELECT OrderNo, SoNo, TrackingNo, CreateDate
						FROM BeforeReturnOrder
						WHERE OrderNo = :Search OR TrackingNo = :Search
						ORDER BY OrderNo
						OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
					 `
		var orderHeadBatch []response.ImportOrderResponse
		nstmtHead, err := repo.db.PrepareNamed(queryHead)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement for order head: %w", err)
		}
		defer nstmtHead.Close()

		// *️⃣ ดึงข้อมูล Order Head ในแต่ละ batch
		err = nstmtHead.SelectContext(ctx, &orderHeadBatch, map[string]interface{}{
			"Search": search,
			"Limit":  chunkSize,
			"Offset": offset,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order head: %w", err)
		}

		// *️⃣ ถ้าไม่มีข้อมูลใน batch นี้ ให้หยุดการทำงาน
		if len(orderHeadBatch) == 0 {
			break
		}

		// *️⃣ ดึงข้อมูล Order Lines สำหรับแต่ละ Order Head ใน batch
		for _, orderHead := range orderHeadBatch {

			queryLines := ` SELECT SKU, ItemName, QTY, Price
							FROM BeforeReturnOrderLine
							WHERE OrderNo = :Search OR TrackingNo = :Search
						  `

			var orderLines []response.ImportOrderLineResponse
			nstmtLines, err := repo.db.PrepareNamed(queryLines)
			if err != nil {
				return nil, fmt.Errorf("failed to prepare statement for order lines: %w", err)
			}
			defer nstmtLines.Close()

			// *️⃣ ดึงข้อมูล Order Lines
			err = nstmtLines.SelectContext(ctx, &orderLines, map[string]interface{}{
				"Search": search,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch order lines: %w", err)
			}

			// *️⃣ ผูกข้อมูล Line กับ Order Head
			orderHead.OrderLines = orderLines

			// *️⃣ เติมค่าของ OrderLines (TrackingNo และ OrderNo)
			for i := range orderHead.OrderLines {
				orderHead.OrderLines[i].TrackingNo = orderHead.TrackingNo
				orderHead.OrderLines[i].OrderNo = orderHead.OrderNo
			}

			// *️⃣ เพิ่มข้อมูลที่ดึงมาในแต่ละ batch
			orders = append(orders, orderHead)
		}

		// *️⃣ เพิ่ม offset สำหรับ batch ถัดไป
		offset += chunkSize
	}

	return orders, nil
}

func (repo repositoryDB) GetOrderTracking(ctx context.Context) ([]response.ImportItem, error) {
	var order []response.ImportItem

	query := `  SELECT OrderNo, TrackingNo
				FROM BeforeReturnOrder
				ORDER BY RecID
			 `

	err := repo.db.SelectContext(ctx, &order, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return order : %w", err)
	}

	return order, nil
}

func (repo repositoryDB) CheckSearch(ctx context.Context, search string) (bool, error) {
	query := `
		SELECT COUNT(1) 
		FROM ROM_V_OrderHeadDetail 
		WHERE OrderNo = @Search OR TrackingNo = @Search
	`

	var count int
	err := repo.db.GetContext(ctx, &count, query, sql.Named("Search", search))
	if err != nil {
		return false, fmt.Errorf("failed to check search existence: %w", err)
	}

	return count > 0, nil
}

func (repo repositoryDB) ValidateSKU(ctx context.Context, orderNo, sku string) (bool, error) {
	var exists bool

	query := `  SELECT CASE WHEN EXISTS (
					SELECT 1 FROM BeforeReturnOrderLine 
					WHERE OrderNo = @OrderNo AND SKU = @SKU
			    ) THEN 1 ELSE 0 END
		     `
	err := repo.db.GetContext(ctx, &exists, query, sql.Named("OrderNo", orderNo), sql.Named("SKU", sku))
	if err != nil {
		return false, fmt.Errorf("failed to validate SKU: %w", err)
	}

	return exists, nil
}

// FetchReturnDetailsBySaleOrder retrieves OrderNo from SoNo
func (repo repositoryDB) FetchReturnDetailsBySaleOrder(ctx context.Context, soNo string) (string, error) {
	query := `  SELECT OrderNo
				FROM ReturnOrder
				WHERE SoNo = @SoNo
			 `
	var orderNo string
	err := repo.db.QueryRowContext(ctx, query, sql.Named("SoNo", soNo)).Scan(&orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no records found for SoNo: %s", soNo)
		}
		return "", fmt.Errorf("database query error: %w", err)
	}

	return orderNo, nil
}

// ทำรอไว้ยังไม่ได้ใช้

// InsertImageMetadata inserts image metadata into the database
func (repo repositoryDB) InsertImageMetadata(ctx context.Context, image request.Images) (int, error) {
	query := `  INSERT INTO Images (SKU, OrderNo, FilePath, ImageTypeID, CreateBy, CreateDate)
				VALUES (:SKU, :OrderNo, :FilePath, :ImageTypeID, :CreateBy, GETDATE());
				SELECT SCOPE_IDENTITY();
			 `

	params := map[string]interface{}{
		"SKU":         image.SKU,
		"OrderNo":     image.OrderNo,
		"FilePath":    image.FilePath,
		"ImageTypeID": image.ImageTypeID,
		"CreateBy":    image.CreateBy,
	}

	var imageID int
	rows, err := repo.db.NamedQuery(query, params)
	if err != nil {
		return 0, fmt.Errorf("error inserting image metadata: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&imageID)
		if err != nil {
			return 0, fmt.Errorf("error scanning inserted image ID: %w", err)
		}
	}

	return imageID, nil
}
