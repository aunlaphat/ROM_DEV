package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type BeforeReturnRepository interface {
	// CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	// CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error
	// CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	// ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	// ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error)

	// UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	// UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error
	// UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	// // ************************ Create Sale Return ************************ //
	// SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error)
	// CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	// UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error
	// ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error
	// CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	// // ************************ Draft & Confirm ************************ //
	// ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error)
	// GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, []response.DraftLineResponse, error)
	// ListCodeR(ctx context.Context) ([]response.CodeRResponse, error)
	// AddCodeR(ctx context.Context, codeR request.CodeR) (*response.DraftLineResponse, error)
	// DeleteCodeR(ctx context.Context, orderNo string, sku string) error
	// UpdateOrderStatus(ctx context.Context, orderNo string, statusConfID int, statusReturnID int, userID string) error

	// Get Real Order
	GetAllOrderDetails(ctx context.Context, offset, limit int) ([]response.OrderDetail, error)
	SearchOrderDetail(ctx context.Context, soNo string) (*response.OrderDetail, error)

	// Delete Line
	DeleteBeforeReturnOrderLine(ctx context.Context, orderNo string, sku string) error

	// ************************ Trade Return ************************ //
	CreateTradeReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	CreateTradeReturnLine(ctx context.Context, orderNo string, lines []request.OrderLines) error
	CheckBefOrderNoExists(ctx context.Context, orderNo string) (bool, error)
	CheckBefLineSKUExists(ctx context.Context, identifier, sku string) (bool, error)
	// ConfirmToReturn(ctx context.Context, req request.ConfirmToReturnRequest, updateBy string) error

	// ************************ ImportOrder: Search Sale Return ************************ //
	GetTrackingNoByOrderNo(ctx context.Context, orderNo string) (string, error)

	// ************************ Confirm Return ************************ //
	UpdateStatusToSuccess(ctx context.Context, orderNo, updateBy string) error
	GetBeforeOrderDetails(ctx context.Context, orderNo string) (*response.ConfirmReturnOrderDetails, error)
	UpdateReturnOrderAndLines(ctx context.Context, req request.ConfirmToReturnRequest, returnOrderData *response.ConfirmReturnOrderDetails) error
	CheckReLineSKUExists(ctx context.Context, orderNo, sku string) (bool, error)

	// ************************ Confirm Receipt ************************ //
	InsertImages(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error
	InsertReturnOrderLine(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error
	InsertReturnOrder(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails) error
	GetBeforeReturnOrderData(ctx context.Context, req request.ConfirmTradeReturnRequest) (*response.ConfirmReturnOrderDetails, error)
	UpdateBefToWaiting(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error
	CheckBefOrderOrTrackingExists(ctx context.Context, identifier string) (bool, error)
}

func (repo repositoryDB) CreateTradeReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
	// 1. เริ่ม transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Logging
	fmt.Println("Transaction started")

	queryHead := `  INSERT INTO BeforeReturnOrder (
						OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
						SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy
					) VALUES (
						:OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
						:SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy
					)
				 `
	_, err = tx.NamedExecContext(ctx, queryHead, map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"SoNo":           order.SoNo,
		"SrNo":           order.SrNo,
		"ChannelID":      order.ChannelID,
		"Reason":         order.Reason,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatus":       order.SoStatus,
		"MkpStatus":      order.MkpStatus,
		"ReturnDate":     order.ReturnDate,
		"StatusReturnID": 3,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"CreateBy":       order.CreateBy,
	})
	if err != nil {
		return nil, fmt.Errorf("[ failed to create BeforeReturnOrder: %w ]", err)
	}

	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	// เตรียมข้อมูลทั้งหมดที่ต้องการ insert
	var params []map[string]interface{}
	for _, line := range order.BeforeReturnOrderLines {
		lineParams := map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"ItemName":   line.ItemName,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   order.CreateBy,
			"TrackingNo": order.TrackingNo,
		}
		params = append(params, lineParams)
	}

	_, err = tx.NamedExecContext(ctx, queryLine, params)
	if err != nil {
		return nil, fmt.Errorf("[ failed to create BeforeReturnOrderLine: %w ]", err)
	}

	// 4. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("[ failed to commit transaction: %w ]", err)
	}

	// Logging
	fmt.Println("Transaction committed")

	// 5. ดึงข้อมูลที่สร้างเสร็จแล้ว
	createdOrder, err := repo.GetBeforeReturnOrderByOrderNo(ctx, order.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("[ failed to fetch created order: %w ]", err)
	}

	return createdOrder, nil
}

func (repo repositoryDB) CreateTradeReturnLine(ctx context.Context, orderNo string, lines []request.OrderLines) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// *️⃣ ดึงค่า TrackingNo ด้วยเลขออเดอร์ จาก BeforeReturnOrder
		trackingNo, err := repo.GetTrackingNoByOrderNo(ctx, orderNo)
		if err != nil {
			return fmt.Errorf("[ failed to fetch TrackingNo for OrderNo: %w ]", err)
		}

		query := `INSERT INTO BeforeReturnOrderLine 
					(OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo, CreateDate) 
				  VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo, GETDATE())`

		// เตรียมพารามิเตอร์ทั้งหมดสำหรับหลายรายการในครั้งเดียว
		params := make([]map[string]interface{}, len(lines))
		for i, line := range lines {
			params[i] = map[string]interface{}{
				"OrderNo":    orderNo,
				"SKU":        line.SKU,
				"ItemName":   line.ItemName,
				"QTY":        line.QTY,
				"ReturnQTY":  line.ReturnQTY,
				"Price":      line.Price,
				"CreateBy":   line.CreateBy, // ใช้ CreateBy จาก userID
				"TrackingNo": trackingNo,
			}
		}

		_, err = tx.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("[ failed to create trade return lines: %w ]", err)
		}

		return nil
	})
}

// *️⃣ search trackingNo by OrderNo
func (repo repositoryDB) GetTrackingNoByOrderNo(ctx context.Context, orderNo string) (string, error) {
	var trackingNo string

	query := ` SELECT TrackingNo
        	   FROM BeforeReturnOrder
               WHERE OrderNo = @OrderNo 
			 `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&trackingNo)
	if err != nil {
		return "", fmt.Errorf("[ failed to fetch TrackingNo: %w ]", err)
	}

	return trackingNo, nil
}

/************** Confirm To ReturnOrder ****************/
// รวม func. UpdateStatusToSuccess + GetBeforeOrderDetails + UpdateReturnOrderAndLines + InsertReturnOrderLine in service

// *️⃣ step 1: update status BeforeReturnOrder, เก็บค่าผู้ updateBy Date เพื่อนำไปใช้เข้าใน CreateBy Date => ReturnOrder,Line
func (repo repositoryDB) UpdateStatusToSuccess(ctx context.Context, orderNo, updateBy string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		query := `  UPDATE BeforeReturnOrder
					SET StatusReturnID = 6, -- success status
						UpdateBy = :UpdateBy, 
						UpdateDate = GETDATE()
					WHERE OrderNo = :OrderNo
				 `
		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			return fmt.Errorf("[ error preparing statement for OrderNo: %w ]", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(map[string]interface{}{
			"OrderNo":  orderNo,
			"UpdateBy": updateBy,
		})
		if err != nil {
			return fmt.Errorf("[ error updating status to success for OrderNo: %w", err)
		}

		return nil
	})
}

// *️⃣ step 2: Fetch ค่า Befod ออกมา เก็บค่าผู้ updateBy Date เพื่อนำไปใช้เข้าใน CreateBy Date => ReturnOrder,Line
func (repo repositoryDB) GetBeforeOrderDetails(ctx context.Context, orderNo string) (*response.ConfirmReturnOrderDetails, error) {
	query := ` SELECT UpdateBy, UpdateDate
        	   FROM BeforeReturnOrder
               WHERE OrderNo = :OrderNo 
			 `

	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer stmt.Close()

	var returnOrderData response.ConfirmReturnOrderDetails
	err = stmt.QueryRowx(map[string]interface{}{"OrderNo": orderNo}).StructScan(&returnOrderData)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &returnOrderData, nil
}

// *️⃣ step 3: update return order (status,sr) + line (actualqty,price)
func (repo repositoryDB) UpdateReturnOrderAndLines(ctx context.Context, req request.ConfirmToReturnRequest, returnOrderData *response.ConfirmReturnOrderDetails) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// อัปเดต ReturnOrder
		for _, head := range req.UpdateToReturn {
			queryHead := `  UPDATE ReturnOrder
                            SET StatusCheckID = 2, --CONFIRM status
                                SrNo = :SrNo, 
                                UpdateBy = :UpdateBy, 
                                UpdateDate = :UpdateDate,
								CheckBy = :CheckBy, 
                                CheckDate = :CheckDate
                            WHERE OrderNo = :OrderNo 
						 `
			stmt, err := tx.PrepareNamed(queryHead)
			if err != nil {
				return fmt.Errorf("[ failed to prepare statement: %w ]", err)
			}
			defer stmt.Close()

			_, err = stmt.Exec(map[string]interface{}{
				"OrderNo":    req.OrderNo,
				"SrNo":       head.SrNo,
				"UpdateBy":   returnOrderData.UpdateBy,
				"UpdateDate": returnOrderData.UpdateDate,
				"CheckBy":    returnOrderData.UpdateBy,
				"CheckDate":  returnOrderData.UpdateDate,
			})
			if err != nil {
				return fmt.Errorf("[ failed to update ReturnOrder: %w ]", err)
			}
		}

		// อัปเดต ReturnOrderLine
		for _, line := range req.ImportLinesActual { // COALESCE => ฟิลด์ที่ไม่ได้ใช้จะดึงค่าเดิมมาแทน

			queryLine := `  UPDATE ReturnOrderLine
						    SET SKU = COALESCE(:SKU, SKU),
								ActualQTY = COALESCE(:ActualQTY, ActualQTY),
							    Price = COALESCE(:Price, Price),
								StatusDelete = COALESCE(:StatusDelete, StatusDelete),
								UpdateBy = COALESCE(:UpdateBy, UpdateBy),
								UpdateDate = COALESCE(:UpdateDate, UpdateDate),
								DeleteBy = COALESCE(:DeleteBy, DeleteBy),
								DeleteDate = COALESCE(:DeleteDate, DeleteDate)
							WHERE OrderNo = :OrderNo AND SKU = :SKU 
						 `
			stmt, err := tx.PrepareNamed(queryLine)
			if err != nil {
				return fmt.Errorf("[ failed to prepare statement: %w ]", err)
			}
			defer stmt.Close()

			deleteBy := sql.NullString{}
			deleteDate := sql.NullString{}
			if line.StatusDelete {
				deleteBy = sql.NullString{String: returnOrderData.UpdateBy, Valid: true}
				deleteDate = sql.NullString{String: returnOrderData.UpdateDate, Valid: true}
			}

			_, err = stmt.Exec(map[string]interface{}{
				"OrderNo":      req.OrderNo,
				"SKU":          line.SKU,
				"ActualQTY":    sql.NullInt32{Int32: int32(line.ActualQTY), Valid: line.ActualQTY != 0}, // เมื่อส่งค่า ว่าง หรือ 0 มาให้ใช้ค่าเดิม
				"Price":        sql.NullFloat64{Float64: line.Price, Valid: line.Price != 0},            // เมื่อส่งค่า ว่าง หรือ 0 มาให้ใช้ค่าเดิม
				"StatusDelete": sql.NullBool{Bool: line.StatusDelete, Valid: line.StatusDelete},         // เมื่อส่งค่าว่างมาให้ใช้ค่าเดิม
				"UpdateBy":     returnOrderData.UpdateBy,
				"UpdateDate":   returnOrderData.UpdateDate,
				"DeleteBy":     deleteBy,
				"DeleteDate":   deleteDate,
			})
			if err != nil {
				return fmt.Errorf("[ failed to update ReturnOrderLine: %w ]", err)
			}
		}

		// Commit Transaction
		return nil
	})
}

/************** Confirm Receipt ****************/
// *️⃣ รวม func. UpdateBefToWaiting + GetBeforeReturnOrderData + InsertReturnOrder + InsertReturnOrderLine in service

// 1. *️⃣Update สถานะใน BeforeReturnOrder to "WAITING" (Page: Confirm Trade)
func (repo repositoryDB) UpdateBefToWaiting(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		query := ` UPDATE BeforeReturnOrder
					SET StatusReturnID = 7, -- WAITING status
						UpdateBy = :UpdateBy,
						UpdateDate = GETDATE()
					WHERE OrderNo = :Identifier OR TrackingNo = :Identifier
				  `
		params := map[string]interface{}{
			"Identifier": req.Identifier,
			"UpdateBy":   updateBy,
		}
		_, err := tx.NamedExecContext(ctx, query, params)
		return err
	})
}

// 2. *️⃣ดึงข้อมูลจาก BeforeReturnOrder fetch ออกมาเพื่อเอาเข้า ReturnOrder
func (repo repositoryDB) GetBeforeReturnOrderData(ctx context.Context, req request.ConfirmTradeReturnRequest) (*response.ConfirmReturnOrderDetails, error) {
	query := `	SELECT OrderNo, SoNo, SrNo, TrackingNo, ChannelID, Reason,
						UpdateBy AS CreateBy, UpdateDate AS CreateDate
				FROM BeforeReturnOrder
				WHERE OrderNo = :Identifier OR TrackingNo = :Identifier
    		 `
	var returnOrderData response.ConfirmReturnOrderDetails

	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"Identifier": req.Identifier,
	})
	if err != nil {
		return nil, fmt.Errorf("[ failed to fetch BeforeReturnOrder: %w ]", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.StructScan(&returnOrderData); err != nil {
			return nil, fmt.Errorf("[ failed to scan BeforeReturnOrder: %w ]", err)
		}
	}

	return &returnOrderData, nil
}

// 3. *️⃣Insert ข้อมูลลงใน ReturnOrder
func (repo repositoryDB) InsertReturnOrder(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		query := `  INSERT INTO ReturnOrder (
						OrderNo, SoNo, SrNo, ChannelID, Reason, TrackingNo, CreateBy, CreateDate, StatusCheckID
					) VALUES (
						:OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :TrackingNo, :CreateBy, :CreateDate, :StatusCheckID
					) 
				 `
		_, err := tx.NamedExecContext(ctx, query, returnOrderData)

		return err
	})
}

// 4. *️⃣Insert ข้อมูลจาก importLines ลงใน ReturnOrderLine
func (repo repositoryDB) InsertReturnOrderLine(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		query := `  INSERT INTO ReturnOrderLine (
            			OrderNo, SKU, QTY, ReturnQTY, Price, TrackingNo, CreateBy, CreateDate
					) VALUES (
						:OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :TrackingNo, :CreateBy, :CreateDate
					) 
				 `
		// เตรียมข้อมูลทั้งหมดที่ต้องการ insert
		var params []map[string]interface{}
		for _, line := range req.ImportLines {
			params = append(params, map[string]interface{}{
				"OrderNo":    returnOrderData.OrderNo,
				"SKU":        line.SKU,
				"QTY":        line.QTY,
				"ReturnQTY":  line.ReturnQTY,
				"Price":      line.Price,
				"TrackingNo": returnOrderData.TrackingNo,
				"CreateBy":   returnOrderData.CreateBy,
				"CreateDate": returnOrderData.CreateDate,
			})
		}

		_, err := tx.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("[ failed to insert into ReturnOrderLine: %w ]", err)
		}

		return nil
	})
}

// *️⃣ InsertImages ฟังก์ชันที่ใช้เพิ่มข้อมูลภาพลงในฐานข้อมูล
func (repo repositoryDB) InsertImages(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		query := `	INSERT INTO Images (
						OrderNo, ImageTypeID, SKU, FilePath, CreateBy, CreateDate
					) VALUES (
						:OrderNo, :ImageTypeID, :SKU, :FilePath, :CreateBy, :CreateDate
					)
				 `
		// *️⃣ เตรียมข้อมูลทั้งหมดที่ต้องการ insert
		var params []map[string]interface{}
		for _, line := range req.ImportLines {
			params = append(params, map[string]interface{}{
				"OrderNo":     returnOrderData.OrderNo,
				"ImageTypeID": line.ImageTypeID,
				"SKU":         line.SKU,
				"FilePath":    line.FilePath,
				"CreateBy":    returnOrderData.CreateBy,
				"CreateDate":  returnOrderData.CreateDate,
			})
		}
		_, err := tx.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("[ failed to insert into Images: %w ]", err)
		}

		return nil
	})
}

/************************** Delete Line *************************/

func (repo repositoryDB) DeleteBeforeReturnOrderLine(ctx context.Context, orderNo string, sku string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// ลบ BeforeReturnOrderLine ตาม OrderNo และ SKU
		query := ` DELETE FROM BeforeReturnOrderLine
				   WHERE OrderNo = :OrderNo AND SKU = :SKU
				 `

		_, err := tx.NamedExecContext(ctx, query, map[string]interface{}{
			"OrderNo": orderNo,
			"SKU":     sku,
		})
		if err != nil {
			return fmt.Errorf("[ Error deleting BeforeReturnOrderLine: %w ]", err)
		}

		return nil
	})
}

/************************** Validate *************************/

// ตรวจสอบว่ามี OrderNo ใน BeforeReturnOrder หรือไม่
func (repo repositoryDB) CheckBefOrderNoExists(ctx context.Context, orderNo string) (bool, error) {
	var exists bool

	query := ` SELECT CASE 
			   WHEN EXISTS (SELECT 1 FROM BeforeReturnOrder 
			   				WHERE OrderNo = @OrderNo) 
			   THEN 1 ELSE 0 
		       END 
			 `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("[ failed to check order existence: %w ]", err)
	}

	return exists, nil
}

// ตรวจสอบว่ามี OrderNo, TrackingNo ใน BeforeReturnOrder หรือไม่
func (repo repositoryDB) CheckBefOrderOrTrackingExists(ctx context.Context, identifier string) (bool, error) {
	var exists bool

	query := ` SELECT CASE 
               WHEN EXISTS (SELECT 1 FROM BeforeReturnOrder 
			   				WHERE OrderNo = @Identifier OR TrackingNo = @Identifier) 
               THEN 1 ELSE 0 
               END 
			 `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("Identifier", identifier)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("[ failed to check order existence: %w ]", err)
	}

	return exists, nil
}

// ตรวจสอบว่ามี sku นี้ของ OrderNo,TrackingNo ใน BeforeReturnOrderLine หรือไม่
func (repo repositoryDB) CheckBefLineSKUExists(ctx context.Context, identifier, sku string) (bool, error) {
	query := ` SELECT 1 FROM BeforeReturnOrderLine 
               WHERE SKU = :SKU AND (OrderNo = :Identifier OR TrackingNo = :Identifier) 
			 `
	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return false, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer stmt.Close()

	var exists int

	err = stmt.QueryRowx(map[string]interface{}{"SKU": sku, "Identifier": identifier}).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// ตรวจสอบว่ามี sku นี้ของ OrderNo ใน ReturnOrderLine หรือไม่
func (repo repositoryDB) CheckReLineSKUExists(ctx context.Context, orderNo, sku string) (bool, error) {
	query := ` SELECT 1 FROM ReturnOrderLine 
               WHERE SKU = :SKU AND OrderNo = :OrderNo 
			 `

	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return false, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer stmt.Close()

	var exists int
	err = stmt.QueryRowx(map[string]interface{}{"SKU": sku, "OrderNo": orderNo}).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

/************************** Get Order Head+Line : Paginate *************************/

func (repo repositoryDB) GetAllOrderDetails(ctx context.Context, offset, limit int) ([]response.OrderDetail, error) {
	var headDetails []response.OrderHeadDetail

	queryHead := `	SELECT OrderNo, SoNo, StatusMKP, SalesStatus, CreateDate
					FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
					ORDER BY OrderNo 
					OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
				 `
	err := repo.db.SelectContext(ctx, &headDetails, queryHead, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("[ error querying OrderHeadDetail: %w ]", err)
	}

	// ถ้าไม่มี order ให้ return กลับเลย
	if len(headDetails) == 0 {
		return []response.OrderDetail{}, nil
	}

	// สร้าง slice ของ OrderNo เพื่อนำไปใช้ใน WHERE IN
	var orderNos []string
	for _, head := range headDetails {
		orderNos = append(orderNos, head.OrderNo)
	}

	// Batch OrderNo เป็น chunks (สูงสุด 1000 ต่อชุด)
	const chunkSize = 1000
	var allLineDetails []response.OrderLineDetail

	for i := 0; i < len(orderNos); i += chunkSize {
		end := i + chunkSize
		if end > len(orderNos) {
			end = len(orderNos)
		}

		// ดึง subset ของ OrderNo
		orderNoChunk := orderNos[i:end]

		// ใช้ sqlx.In เพื่อ binding ORDER IN
		queryLine, args, err := sqlx.In(`
				SELECT OrderNo, SoNo, StatusMKP, SalesStatus, SKU, ItemName, QTY, Price, CreateDate
				FROM Data_WebReturn.dbo.ROM_V_OrderLineDetail
				WHERE OrderNo IN (?)
				ORDER BY OrderNo `,
			orderNoChunk)

		if err != nil {
			return nil, fmt.Errorf("[ error building OrderLineDetail query: %w ]", err)
		}

		// ใช้ Rebind เพื่อให้รองรับ SQL Server
		queryLine = repo.db.Rebind(queryLine)

		var lineDetails []response.OrderLineDetail
		err = repo.db.SelectContext(ctx, &lineDetails, queryLine, args...)
		if err != nil {
			return nil, fmt.Errorf("[ error querying OrderLineDetail: %w ]", err)
		}

		// รวมผลลัพธ์ทั้งหมด
		allLineDetails = append(allLineDetails, lineDetails...)
	}

	// Map Order Lines to Order Heads
	orderLineMap := make(map[string][]response.OrderLineDetail)
	for _, line := range allLineDetails {
		orderLineMap[line.OrderNo] = append(orderLineMap[line.OrderNo], line)
	}

	// เชื่อมข้อมูล OrderLineDetail เข้า OrderHeadDetail
	for i := range headDetails {
		headDetails[i].OrderLineDetail = orderLineMap[headDetails[i].OrderNo]
	}

	// สร้างตัวแปรเพื่อเก็บผลลัพธ์
	allorder := []response.OrderDetail{
		{OrderHeadDetail: headDetails},
	}

	return allorder, nil
}

func (repo repositoryDB) SearchOrderDetail(ctx context.Context, soNo string) (*response.OrderDetail, error) {
	var headDetails []response.OrderHeadDetail

	// Query to get OrderHeadDetails filtered by SoNo
	queryHead := `
		SELECT OrderNo, SoNo, StatusMKP, SalesStatus, CreateDate
		FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
		WHERE SoNo = @SoNo
	`

	// Execute the query to get OrderHeadDetails
	err := repo.db.SelectContext(ctx, &headDetails, queryHead, sql.Named("SoNo", soNo))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("[ error querying OrderHeadDetail by SoNo: %w ]", err)
	}

	// ถ้าไม่มี OrderHeadDetails ให้ return nil
	if len(headDetails) == 0 {
		return nil, fmt.Errorf("[ no order details found ]")
	}

	// Retrieve OrderLineDetails if there are any OrderHeadDetails
	var lineDetails []response.OrderLineDetail
	if len(headDetails) > 0 {
		// สร้าง slice ของ OrderNo เพื่อนำไปใช้ใน WHERE IN
		var orderNos []string
		for _, head := range headDetails {
			orderNos = append(orderNos, head.OrderNo)
		}

		// Query to get OrderLineDetails
		queryLine, args, err := sqlx.In(`
			SELECT OrderNo, SoNo, StatusMKP, SalesStatus, SKU, ItemName, QTY, Price, CreateDate
			FROM Data_WebReturn.dbo.ROM_V_OrderLineDetail
			WHERE OrderNo IN (?)
			ORDER BY OrderNo
		`, orderNos)
		if err != nil {
			return nil, fmt.Errorf("[ error building OrderLineDetail query: %w ]", err)
		}

		// Rebind the query for SQL Server compatibility
		queryLine = repo.db.Rebind(queryLine)

		err = repo.db.SelectContext(ctx, &lineDetails, queryLine, args...)
		if err != nil {
			return nil, fmt.Errorf("[ error querying OrderLineDetail: %w ]", err)
		}
	}

	// Map Order Lines to Order Heads
	orderLineMap := make(map[string][]response.OrderLineDetail)
	for _, line := range lineDetails {
		orderLineMap[line.OrderNo] = append(orderLineMap[line.OrderNo], line)
	}

	// Add the OrderLineDetail to each OrderHeadDetail
	for i := range headDetails {
		headDetails[i].OrderLineDetail = orderLineMap[headDetails[i].OrderNo]
	}

	// Store the result in a variable before returning
	orderDetail := &response.OrderDetail{OrderHeadDetail: headDetails}

	// Return the result
	return orderDetail, nil
}

// // Implementation สำหรับ CreateBeforeReturnOrder
// func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
// 	queryOrder := `
//         INSERT INTO BeforeReturnOrder (
//             OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CancelID
//         ) VALUES (
//             :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, :CancelID
//         )
//     `
// 	paramsOrder := map[string]interface{}{
// 		"OrderNo":        order.OrderNo,
// 		"SoNo":           order.SoNo,
// 		"SrNo":           order.SrNo,
// 		"ChannelID":      order.ChannelID,
// 		"Reason":         order.Reason,
// 		"CustomerID":     order.CustomerID,
// 		"TrackingNo":     order.TrackingNo,
// 		"Logistic":       order.Logistic,
// 		"WarehouseID":    order.WarehouseID,
// 		"SoStatus":       order.SoStatus,
// 		"MkpStatus":      order.MkpStatus,
// 		"ReturnDate":     order.ReturnDate,
// 		"StatusReturnID": order.StatusReturnID,
// 		"StatusConfID":   order.StatusConfID,
// 		"ConfirmBy":      order.ConfirmBy,
// 		"CreateBy":       order.CreateBy,
// 		"CancelID":       order.CancelID,
// 	}

// 	_, err := repo.db.NamedExecContext(ctx, queryOrder, paramsOrder)
// 	if err != nil {
// 		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
// 	}

// 	return nil
// }

// // Implementation สำหรับ CreateBeforeReturnOrderLine
// func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
// 	query := `
//         INSERT INTO BeforeReturnOrderLine (
//             OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
//         ) VALUES (
//             :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
//         )
//     `
// 	for _, line := range lines {
// 		params := map[string]interface{}{
// 			"OrderNo":   orderNo,
// 			"SKU":       line.SKU,
// 			"ItemName":  line.ItemName,
// 			"QTY":       line.QTY,
// 			"ReturnQTY": line.ReturnQTY,
// 			"Price":     line.Price,
// 			"CreateBy":  line.CreateBy,
// 		}
// 		_, err := repo.db.NamedExecContext(ctx, query, params)
// 		if err != nil {
// 			return fmt.Errorf("failed to create order line: %w", err)
// 		}
// 	}
// 	return nil
// }

// // Implementation สำหรับ BeginTransaction CreateBeforeReturnOrder & CreateBeforeReturnOrderLine
// func (repo repositoryDB) CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}

// 	queryOrder := `
//         INSERT INTO BeforeReturnOrder (
//             OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID,
//             SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy
//         ) VALUES (
//             :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID,
//             :SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy
//         )
//     `
// 	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
// 		"OrderNo":        order.OrderNo,
// 		"SoNo":           order.SoNo,
// 		"SrNo":           order.SrNo,
// 		"ChannelID":      order.ChannelID,
// 		"Reason":         order.Reason,
// 		"CustomerID":     order.CustomerID,
// 		"TrackingNo":     order.TrackingNo,
// 		"Logistic":       order.Logistic,
// 		"WarehouseID":    order.WarehouseID,
// 		"SoStatus":       order.SoStatus,
// 		"MkpStatus":      order.MkpStatus,
// 		"ReturnDate":     order.ReturnDate,
// 		"StatusReturnID": order.StatusReturnID,
// 		"StatusConfID":   order.StatusConfID,
// 		"ConfirmBy":      order.ConfirmBy,
// 		"CreateBy":       order.CreateBy,
// 	})
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
// 	}

// 	queryLine := `
//         INSERT INTO BeforeReturnOrderLine (
//             OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
//         ) VALUES (
//             :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
//         )
//     `
// 	for _, line := range order.BeforeReturnOrderLines {
// 		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
// 			"OrderNo":   order.OrderNo,
// 			"SKU":       line.SKU,
// 			"ItemName":  line.ItemName,
// 			"QTY":       line.QTY,
// 			"ReturnQTY": line.ReturnQTY,
// 			"Price":     line.Price,
// 			"CreateBy":  order.CreateBy,
// 		})
// 		if err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return nil
// }

// Implementation สำหรับ GetBeforeReturnOrderLineByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error) {
	query := `  SELECT OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, TrackingNo, CreateDate
				FROM BeforeReturnOrderLine
				WHERE OrderNo = :OrderNo
				ORDER BY RecID
			 `
	var lines []response.BeforeReturnOrderItem
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("[ failed to get order lines: %w ]", err)
	}

	fmt.Printf("[ Fetched %d lines from the database for OrderNo: %s ]", len(lines), orderNo)

	return lines, nil
}

// Implementation สำหรับ GetBeforeReturnOrderByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	query := `  SELECT OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
					   SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, 
					   CreateDate, UpdateBy, UpdateDate, CancelID
				FROM BeforeReturnOrder
				WHERE OrderNo = :OrderNo
			`
	order := new(response.BeforeReturnOrderResponse)
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, order, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[ failed to fetch BeforeReturnOrder: %w ]", err)
	}

	lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	order.Items = lines

	return order, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderItem, error) {
	query := `  SELECT OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, TrackingNo, CreateDate
				FROM BeforeReturnOrderLine
				WHERE OrderNo = :OrderNo
				ORDER BY RecID
			 `
	var lines []response.BeforeReturnOrderItem
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("[ failed to prepare statement: %w ]", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("[ failed to get order lines: %w ]", err)
	}

	return lines, nil
}

// func (repo repositoryDB) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
// 	query := `
//         SELECT
//             OrderNo,
//             SKU,
// 			ItemName,
//             QTY,
//             ReturnQTY,
//             Price,
//             TrackingNo,
//             CreateDate
//         FROM BeforeReturnOrderLine
//         ORDER BY RecID
//     `

// 	var lines []response.BeforeReturnOrderLineResponse
// 	nstmt, err := repo.db.PrepareNamed(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to prepare statement: %w", err)
// 	}
// 	defer nstmt.Close()

// 	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get order lines: %w", err)
// 	}

// 	return lines, nil
// }

// // ฟังก์ชันพื้นฐานสำหรับการดึงข้อมูล
// func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
// 	query := `
//         SELECT OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
//         FROM BeforeReturnOrder
//         ORDER BY RecID ASC
//     `

// 	var orders []response.BeforeReturnOrderResponse
// 	nstmt, err := repo.db.PrepareNamed(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to prepare statement: %w", err)
// 	}
// 	defer nstmt.Close()

// 	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list orders: %w", err)
// 	}

// 	for i := range orders {
// 		lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orders[i].OrderNo)
// 		if err != nil {
// 			return nil, err
// 		}
// 		orders[i].BeforeReturnOrderLines = lines
// 	}

// 	return orders, nil
// }

// func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
// 	query := `
//         UPDATE BeforeReturnOrder
//         SET SoNo = COALESCE(:SoNo, SoNo),
//             SrNo = COALESCE(:SrNo, SrNo),
//             ChannelID = COALESCE(:ChannelID, ChannelID),
//             Reason = COALESCE(:Reason, Reason),
//             CustomerID = COALESCE(:CustomerID, CustomerID),
//             TrackingNo = COALESCE(:TrackingNo, TrackingNo),
//             Logistic = COALESCE(:Logistic, Logistic),
//             WarehouseID = COALESCE(:WarehouseID, WarehouseID),
//             SoStatus = COALESCE(:SoStatus, SoStatus),
//             MkpStatus = COALESCE(:MkpStatus, MkpStatus),
//             ReturnDate = COALESCE(:ReturnDate, ReturnDate),
//             StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
//             StatusConfID = COALESCE(:StatusConfID, StatusConfID),
//             ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
//             UpdateBy = COALESCE(:UpdateBy, UpdateBy)
//         WHERE OrderNo = :OrderNo
//     `
// 	params := map[string]interface{}{
// 		"OrderNo":        order.OrderNo,
// 		"SoNo":           order.SoNo,
// 		"SrNo":           order.SrNo,
// 		"ChannelID":      order.ChannelID,
// 		"Reason":         order.Reason,
// 		"CustomerID":     order.CustomerID,
// 		"TrackingNo":     order.TrackingNo,
// 		"Logistic":       order.Logistic,
// 		"WarehouseID":    order.WarehouseID,
// 		"SoStatus":       order.SoStatus,
// 		"MkpStatus":      order.MkpStatus,
// 		"ReturnDate":     order.ReturnDate,
// 		"StatusReturnID": order.StatusReturnID,
// 		"StatusConfID":   order.StatusConfID,
// 		"ConfirmBy":      order.ConfirmBy,
// 		"UpdateBy":       order.UpdateBy,
// 	}

// 	_, err := repo.db.NamedExecContext(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
// 	}

// 	return nil
// }

// func (repo repositoryDB) UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error {
// 	query := `
//         UPDATE BeforeReturnOrderLine
//         SET ItemName = COALESCE(:ItemName, ItemName),
// 			QTY = COALESCE(:QTY, QTY),
//             ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
//             Price = COALESCE(:Price, Price),
//             UpdateBy = COALESCE(:UpdateBy, UpdateBy),
//             TrackingNo = COALESCE(:TrackingNo, TrackingNo)
//         WHERE OrderNo = :OrderNo
//           AND SKU = :SKU
//     `
// 	params := map[string]interface{}{
// 		"OrderNo":    orderNo,
// 		"SKU":        line.SKU,
// 		"ItemName":   line.ItemName,
// 		"QTY":        line.QTY,
// 		"ReturnQTY":  line.ReturnQTY,
// 		"Price":      line.Price,
// 		"UpdateBy":   line.UpdateBy,
// 		"TrackingNo": line.TrackingNo,
// 	}

// 	_, err := repo.db.NamedExecContext(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
// 	}

// 	return nil
// }

// // Implementation สำหรับ UpdateBeforeReturnOrderWithTransaction
// func (repo repositoryDB) UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// Update BeforeReturnOrderLine first
// 	queryLine := `
//         UPDATE BeforeReturnOrderLine
//         SET ItemName = COALESCE(:ItemName, ItemName),
// 			QTY = COALESCE(:QTY, QTY),
//             ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
//             Price = COALESCE(:Price, Price),
//             UpdateBy = COALESCE(:UpdateBy, UpdateBy),
//             TrackingNo = COALESCE(:TrackingNo, TrackingNo)
//         WHERE OrderNo = :OrderNo
//           AND SKU = :SKU
//     `

// 	for _, line := range order.BeforeReturnOrderLines {
// 		paramsLine := map[string]interface{}{
// 			"OrderNo":    line.OrderNo,
// 			"SKU":        line.SKU,
// 			"ItemName":   line.ItemName,
// 			"QTY":        line.QTY,
// 			"ReturnQTY":  line.ReturnQTY,
// 			"Price":      line.Price,
// 			"UpdateBy":   line.UpdateBy,
// 			"TrackingNo": line.TrackingNo,
// 		}

// 		result, err := tx.NamedExecContext(ctx, queryLine, paramsLine)
// 		if err != nil {
// 			return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
// 		}

// 		rows, _ := result.RowsAffected()
// 		if rows == 0 {
// 			return fmt.Errorf("no rows updated for OrderNo: %s, SKU: %s", line.OrderNo, line.SKU)
// 		}
// 	}

// 	// Update BeforeReturnOrder
// 	queryOrder := `
//         UPDATE BeforeReturnOrder
//         SET SoNo = COALESCE(:SoNo, SoNo),
//             SrNo = COALESCE(:SrNo, SrNo),
//             ChannelID = COALESCE(:ChannelID, ChannelID),
//             Reason = COALESCE(:Reason, Reason),
//             CustomerID = COALESCE(:CustomerID, CustomerID),
//             TrackingNo = COALESCE(:TrackingNo, TrackingNo),
//             Logistic = COALESCE(:Logistic, Logistic),
//             WarehouseID = COALESCE(:WarehouseID, WarehouseID),
//             SoStatus = COALESCE(:SoStatus, SoStatus),
//             MkpStatus = COALESCE(:MkpStatus, MkpStatus),
//             ReturnDate = COALESCE(:ReturnDate, ReturnDate),
//             StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
//             StatusConfID = COALESCE(:StatusConfID, StatusConfID),
//             ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
//             UpdateBy = COALESCE(:UpdateBy, UpdateBy),
//             CancelID = COALESCE(:CancelID, CancelID)
//         WHERE OrderNo = :OrderNo
//     `

// 	paramsOrder := map[string]interface{}{
// 		"OrderNo":        order.OrderNo,
// 		"SoNo":           order.SoNo,
// 		"SrNo":           order.SrNo,
// 		"ChannelID":      order.ChannelID,
// 		"Reason":         order.Reason,
// 		"CustomerID":     order.CustomerID,
// 		"TrackingNo":     order.TrackingNo,
// 		"Logistic":       order.Logistic,
// 		"WarehouseID":    order.WarehouseID,
// 		"SoStatus":       order.SoStatus,
// 		"MkpStatus":      order.MkpStatus,
// 		"ReturnDate":     order.ReturnDate,
// 		"StatusReturnID": order.StatusReturnID,
// 		"StatusConfID":   order.StatusConfID,
// 		"ConfirmBy":      order.ConfirmBy,
// 		"UpdateBy":       order.UpdateBy,
// 		"CancelID":       order.CancelID,
// 	}

// 	_, err = tx.NamedExecContext(ctx, queryOrder, paramsOrder)
// 	if err != nil {
// 		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
// 	}

// 	return nil
// }

// // ************************ Create Sale Return ************************ //

// func (repo repositoryDB) SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error) {
// 	// 1. Optimize SQL query
// 	query := `
//         SELECT
//             h.SoNo,
//             h.OrderNo,
//             h.StatusMKP,
//             h.SalesStatus,
//             h.CreateDate,
//             l.SKU,
//             l.ItemName,
//             l.QTY,
//             l.Price
//         FROM ROM_V_OrderHeadDetail h
//         INNER JOIN ROM_V_OrderLineDetail l ON h.SoNo = l.SoNo AND h.OrderNo = l.OrderNo
//         WHERE ((:SoNo != '' AND h.SoNo = :SoNo)
//            OR (:OrderNo != '' AND h.OrderNo = :OrderNo))
//         ORDER BY l.SKU` // Add index-based ordering

// 	// 2. Input validation
// 	if soNo == "" && orderNo == "" {
// 		return nil, fmt.Errorf("🚩 Either SoNo or OrderNo must be provided 🚩")
// 	}

// 	// 3. Prepare query parameters
// 	params := map[string]interface{}{
// 		"SoNo":    soNo,
// 		"OrderNo": orderNo,
// 	}

// 	// 4. Execute query with timeout context
// 	rows, err := repo.db.NamedQueryContext(ctx, query, params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer rows.Close()

// 	// 5. Process results efficiently
// 	var (
// 		orderHead  response.SaleOrderResponse
// 		orderLines = make([]response.SaleOrderLineResponse, 0, 10)
// 		isFirst    = true
// 	)

// 	// 6. Scan results with error handling
// 	for rows.Next() {
// 		var line response.SaleOrderLineResponse
// 		scanData := struct {
// 			*response.SaleOrderResponse
// 			*response.SaleOrderLineResponse
// 		}{&orderHead, &line}

// 		if err := rows.StructScan(&scanData); err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}

// 		// Only copy header data once
// 		if isFirst {
// 			isFirst = false
// 		}
// 		orderLines = append(orderLines, line)
// 	}

// 	// 7. Check for errors after iteration
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error during row iteration: %w", err)
// 	}

// 	// 8. Handle no results case
// 	if len(orderLines) == 0 {
// 		return nil, nil
// 	}

// 	// 9. Construct final response
// 	orderHead.OrderLines = orderLines
// 	return &orderHead, nil
// }

// func (repo repositoryDB) CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
// 	// 1. เริ่ม transaction
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// Logging
// 	fmt.Println("Transaction started")

// 	// 2. Insert BeforeReturnOrder (Header)
// 	queryOrder := `
//         INSERT INTO BeforeReturnOrder (
//             OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID,
//             SoStatus, MkpStatus, ReturnDate, CreateBy, CreateDate
//         ) VALUES (
//             :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID,
//             :SoStatus, :MkpStatus, :ReturnDate, :CreateBy, GETDATE()
//         )
//     `
// 	_, err = tx.NamedExecContext(ctx, queryOrder, order)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
// 	}

// 	// Logging
// 	fmt.Println("Inserted BeforeReturnOrder")

// 	// 3. Insert BeforeReturnOrderLine (Lines)
// 	queryLine := `
//         INSERT INTO BeforeReturnOrderLine (
//             OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo
//         ) VALUES (
//             :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE(), :TrackingNo
//         )
//     `
// 	for _, line := range order.BeforeReturnOrderLines {
// 		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
// 			"OrderNo":    order.OrderNo,
// 			"SKU":        line.SKU,
// 			"ItemName":   line.ItemName,
// 			"QTY":        line.QTY,
// 			"ReturnQTY":  line.ReturnQTY,
// 			"Price":      line.Price,
// 			"CreateBy":   order.CreateBy,
// 			"TrackingNo": line.TrackingNo,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
// 		}
// 	}

// 	// 4. Commit transaction
// 	if err = tx.Commit(); err != nil {
// 		return nil, fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	// Logging
// 	fmt.Println("Transaction committed")

// 	// 5. ดึงข้อมูลที่สร้างเสร็จแล้ว
// 	createdOrder, err := repo.GetBeforeReturnOrderByOrderNo(ctx, order.OrderNo)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch created order: %w", err)
// 	}

// 	// Logging
// 	fmt.Println("Fetched created order")

// 	return createdOrder, nil
// }

// func (repo repositoryDB) UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error {
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// 1. SQL query สำหรับ update
// 	query := `
//         UPDATE BeforeReturnOrder
//         SET SrNo = :SrNo,
//             UpdateBy = :UpdateBy,
//             UpdateDate = GETDATE()
//         WHERE OrderNo = :OrderNo
//     `

// 	// 2. กำหนด parameters
// 	params := map[string]interface{}{
// 		"OrderNo":  orderNo,
// 		"SrNo":     srNo,
// 		"UpdateBy": updateBy,
// 	}

// 	// 3. Execute query
// 	result, err := repo.db.NamedExecContext(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to update SR number: %w", err)
// 	}

// 	// 4. ตรวจสอบว่ามีการอัพเดทจริง
// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %w", err)
// 	}
// 	if rows == 0 {
// 		return fmt.Errorf("order not found: %s", orderNo)
// 	}

// 	return nil
// }

// func (repo repositoryDB) ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error {
// 	// เริ่ม transaction
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// 1. กำหนด SQL query สำหรับ update สถานะ
// 	query := `
//         UPDATE BeforeReturnOrder
//         SET StatusReturnID = 1, -- Pending status
//             StatusConfID = 1,   -- Draft status
//             ConfirmBy = :ConfirmBy,
//             ComfirmDate = GETDATE()
//         WHERE OrderNo = :OrderNo
//     `
// 	// 2. กำหนด parameters สำหรับ query
// 	params := map[string]interface{}{
// 		"OrderNo":   orderNo,
// 		"ConfirmBy": confirmBy,
// 	}

// 	// 3. เตรียม statement
// 	nstmt, err := repo.db.PrepareNamed(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare statement for confirming sale return: %w", err)
// 	}
// 	defer nstmt.Close()

// 	// 4. execute query
// 	_, err = nstmt.ExecContext(ctx, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to confirm sale return: %w", err)
// 	}

// 	return nil
// }

// func (repo repositoryDB) CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error {
// 	tx, err := repo.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// 1. ตรวจสอบสถานะปัจจุบันของ order
// 	checkQuery := `
//         SELECT StatusConfID, StatusReturnID
//         FROM BeforeReturnOrder
//         WHERE OrderNo = @OrderNo
//     `
// 	var statusConfID, statusReturnID *int
// 	err = tx.QueryRowContext(ctx, checkQuery, sql.Named("OrderNo", orderNo)).Scan(&statusConfID, &statusReturnID)
// 	if err == sql.ErrNoRows {
// 		return fmt.Errorf("order not found: %s", orderNo)
// 	}
// 	if err != nil {
// 		return fmt.Errorf("failed to check order status: %w", err)
// 	}

// 	// ตรวจสอบว่าสามารถยกเลิกได้หรือไม่
// 	if statusConfID != nil && *statusConfID == 3 {
// 		return fmt.Errorf("order is already canceled")
// 	}
// 	if statusReturnID != nil && *statusReturnID == 2 {
// 		return fmt.Errorf("order is already canceled")
// 	}

// 	// 2. สร้าง CancelStatus และรับค่า CancelID
// 	insertCancelStatus := `
//         INSERT INTO CancelStatus (
//             RefID,
//             CancelStatus,
//             Remark,
//             CancelBy,
//             CancelDate
//         )
//         OUTPUT INSERTED.CancelID
//         VALUES (
//             @OrderNo,
//             1, -- สถานะยกเลิก
//             @Remark,
//             @CancelBy,
//             GETDATE()
//         )
//     `
// 	var cancelID int
// 	err = tx.QueryRowContext(ctx, insertCancelStatus,
// 		sql.Named("OrderNo", orderNo),
// 		sql.Named("Remark", remark),
// 		sql.Named("CancelBy", updateBy),
// 	).Scan(&cancelID)
// 	if err != nil {
// 		return fmt.Errorf("failed to create cancel status: %w", err)
// 	}

// 	// 3. อัพเดทสถานะการยกเลิกใน BeforeReturnOrder
// 	updateOrder := `
//         UPDATE BeforeReturnOrder
//         SET StatusReturnID = 2,
//             StatusConfID = 3,
//             CancelID = @CancelID,
//             UpdateBy = @UpdateBy,
//             UpdateDate = GETDATE()
//         WHERE OrderNo = @OrderNo
//     `
// 	result, err := tx.ExecContext(ctx, updateOrder,
// 		sql.Named("OrderNo", orderNo),
// 		sql.Named("CancelID", cancelID),
// 		sql.Named("UpdateBy", updateBy),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to update order status: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no rows updated for order: %s", orderNo)
// 	}

// 	// 4. Commit transaction
// 	if err = tx.Commit(); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return nil
// }

// // ************************ Draft & Confirm ************************ //

// func (repo repositoryDB) ListDraftOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
// 	query := `
//         SELECT OrderNo, SoNo, SrNo, CustomerID, TrackingNo, Logistic, ChannelID, CreateDate, WarehouseID
//         FROM BeforeReturnOrder
//         WHERE StatusConfID = 1 -- Draft status
// 		ORDER BY CreateDate DESC
//     `

// 	var orders []response.ListDraftConfirmOrdersResponse
// 	nstmt, err := repo.db.PrepareNamed(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to prepare statement: %w", err)
// 	}
// 	defer nstmt.Close()

// 	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list draft orders: %w", err)
// 	}

// 	return orders, nil
// }

// func (repo repositoryDB) ListConfirmOrders(ctx context.Context) ([]response.ListDraftConfirmOrdersResponse, error) {
// 	query := `
//         SELECT OrderNo, SoNo, SrNo, CustomerID, TrackingNo, Logistic, ChannelID, CreateDate, WarehouseID
//         FROM BeforeReturnOrder
//         WHERE StatusConfID = 2 -- Confirm status
// 		ORDER BY CreateDate DESC
//     `

// 	var orders []response.ListDraftConfirmOrdersResponse
// 	nstmt, err := repo.db.PrepareNamed(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to prepare statement: %w", err)
// 	}
// 	defer nstmt.Close()

// 	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list confirm orders: %w", err)
// 	}

// 	return orders, nil
// }

// func (repo repositoryDB) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, []response.DraftLineResponse, error) {
// 	var head response.DraftHeadResponse
// 	var lines []response.DraftLineResponse

// 	headQuery := `
//         SELECT
//             OrderNo,
//             SoNo,
//             SrNo
//         FROM BeforeReturnOrder
//         WHERE OrderNo = :OrderNo
//     `

// 	headQuery, args, err := sqlx.Named(headQuery, map[string]interface{}{"OrderNo": orderNo})
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to prepare head query: %w", err)
// 	}
// 	headQuery = repo.db.Rebind(headQuery)
// 	err = repo.db.GetContext(ctx, &head, headQuery, args...)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to get head data: %w", err)
// 	}

// 	lineQuery := `
//         SELECT
//             SKU,
//             ItemName,
//             QTY,
//             Price
//         FROM BeforeReturnOrderLine
//         WHERE OrderNo = :OrderNo
//     `
// 	lineQuery, args, err = sqlx.Named(lineQuery, map[string]interface{}{"OrderNo": orderNo})
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to prepare line query: %w", err)
// 	}
// 	lineQuery = repo.db.Rebind(lineQuery)
// 	err = repo.db.SelectContext(ctx, &lines, lineQuery, args...)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to get line data: %w", err)
// 	}

// 	return &head, lines, nil
// }

// // Implementation สำหรับ ListCodeR
// func (repo repositoryDB) ListCodeR(ctx context.Context) ([]response.CodeRResponse, error) {
// 	query := `
// 		SELECT SKU, NameAlias
// 		FROM ROM_V_ProductAll
// 		WHERE SKU LIKE 'R%'
// 	`

// 	var CodeR []response.CodeRResponse
// 	err := repo.db.SelectContext(ctx, &CodeR, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get CodeR: %w", err)
// 	}

// 	return CodeR, nil
// }

// func (repo repositoryDB) AddCodeR(ctx context.Context, CodeR request.CodeR) (*response.DraftLineResponse, error) {
// 	CodeR.ReturnQTY = CodeR.QTY

// 	query := `
//         INSERT INTO BeforeReturnOrderLine (OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate)
//         VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE())
//     `

// 	_, err := repo.db.NamedExecContext(ctx, query, CodeR)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to insert CodeR: %w", err)
// 	}

// 	result := &response.DraftLineResponse{
// 		SKU:      CodeR.SKU,
// 		ItemName: CodeR.ItemName,
// 		QTY:      CodeR.QTY,
// 		Price:    CodeR.Price,
// 	}

// 	return result, nil
// }

// func (repo repositoryDB) DeleteCodeR(ctx context.Context, orderNo string, sku string) error {
// 	query := `
//         DELETE FROM BeforeReturnOrderLine
//         WHERE OrderNo = :OrderNo AND SKU = :SKU
//     `

// 	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
// 		"OrderNo": orderNo,
// 		"SKU":     sku,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to delete CodeR: %w", err)
// 	}

// 	return nil
// }

// func (repo repositoryDB) UpdateOrderStatus(ctx context.Context, orderNo string, statusConfID int, statusReturnID int, userID string) error {
// 	query := `
//         UPDATE BeforeReturnOrder
//         SET StatusConfID = :StatusConfID,
//             StatusReturnID = :StatusReturnID,
//             UpdateBy = :UpdateBy,
//             UpdateDate = GETDATE()
//         WHERE OrderNo = :OrderNo
//     `
// 	params := map[string]interface{}{
// 		"OrderNo":        orderNo,
// 		"StatusConfID":   statusConfID,
// 		"StatusReturnID": statusReturnID,
// 		"UpdateBy":       userID,
// 	}

// 	_, err := repo.db.NamedExecContext(ctx, query, params)
// 	if err != nil {
// 		return fmt.Errorf("failed to update order status: %w", err)
// 	}

// 	return nil
// }
