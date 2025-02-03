package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type BeforeReturnRepository interface {
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error
	CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)

	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error
	UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	// Create Return Order MKP 🚨//
	SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error)
	CreateSaleReturn(ctx context.Context, req request.CreateSaleReturnRequest) (*response.BeforeReturnOrderResponse, error)
	UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn, userID string) error
	ConfirmSaleReturn(ctx context.Context, orderNo string, statusReturnID, statusConfID int, userID string) error
	CancelSaleReturn(ctx context.Context, req request.CancelSaleReturn, userID string) error

	// Draft & Confirm MKP 🚨//
	ListDraftOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error)
	ListConfirmOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error)
	GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error)
	ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error)
	AddCodeR(ctx context.Context, req request.AddCodeR) ([]response.AddCodeRResponse, error)
	DeleteCodeR(ctx context.Context, orderNo string, sku string) (int64, error)
	UpdateOrderStatus(ctx context.Context, orderNo string, statusConfID int, statusReturnID int, userID string) (*response.UpdateOrderStatusResponse, error)

	// Get Real Order
	GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error)
	GetAllOrderDetails(ctx context.Context, offset, limit int) ([]response.OrderDetail, error)
	GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error)

	// Delete Line
	DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error

	// ************************ Trade Return ************************ //
	CheckBefOrderNoExists(ctx context.Context, orderNo string) (bool, error)
	CreateTradeReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	CreateTradeReturnLine(ctx context.Context, orderNo string, lines []request.OrderLines) error
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

// ตรวจสอบว่ามี OrderNo ใน BeforeReturnOrder หรือไม่
func (repo repositoryDB) CheckBefOrderNoExists(ctx context.Context, orderNo string) (bool, error) {
	var exists bool
	query := ` SELECT CASE 
			   WHEN EXISTS (SELECT 1 FROM BeforeReturnOrder WHERE OrderNo = @OrderNo) 
			   THEN 1 ELSE 0 
		       END `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}

	return exists, nil
}

func (repo repositoryDB) CheckBefOrderOrTrackingExists(ctx context.Context, identifier string) (bool, error) {
	var exists bool
	query := ` SELECT CASE 
               WHEN EXISTS (SELECT 1 FROM BeforeReturnOrder WHERE OrderNo = @Identifier OR TrackingNo = @Identifier) 
               THEN 1 ELSE 0 
               END `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("Identifier", identifier)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}

	return exists, nil
}

// search trackingNo by OrderNo
func (repo repositoryDB) GetTrackingNoByOrderNo(ctx context.Context, orderNo string) (string, error) {
	var trackingNo string
	query := ` SELECT TrackingNo
        	   FROM BeforeReturnOrder
               WHERE OrderNo = @OrderNo `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&trackingNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("order not found: %s", orderNo)
		}
		return "", fmt.Errorf("failed to fetch TrackingNo: %w", err)
	}
	return trackingNo, nil
}

func (repo repositoryDB) CreateTradeReturnLine(ctx context.Context, orderNo string, lines []request.OrderLines) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// ตรวจสอบว่า OrderNo มีอยู่ใน BeforeReturnOrder หรือไม่
		exists, err := repo.CheckBefOrderNoExists(ctx, orderNo)
		if err != nil {
			return fmt.Errorf("failed to check order existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("order not found: %s", orderNo)
		}

		// ดึง TrackingNo จาก BeforeReturnOrder
		trackingNo, err := repo.GetTrackingNoByOrderNo(ctx, orderNo)
		if err != nil {
			return fmt.Errorf("failed to fetch TrackingNo for OrderNo %s: %w", orderNo, err)
		}

		// สร้างข้อมูล BeforeReturnOrderLine สำหรับหลายรายการ
		query := `INSERT INTO BeforeReturnOrderLine 
					(OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo, CreateDate) 
				  VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo, GETDATE())`

		// เตรียมพารามิเตอร์สำหรับหลายรายการ
		var params []map[string]interface{}
		for _, line := range lines {
			params = append(params, map[string]interface{}{
				"OrderNo":    orderNo,
				"SKU":        line.SKU,
				"ItemName":   line.ItemName,
				"QTY":        line.QTY,
				"ReturnQTY":  line.ReturnQTY,
				"Price":      line.Price,
				"CreateBy":   line.CreateBy, // ใช้ CreateBy จากคำขอ
				"TrackingNo": trackingNo,
			})
		}

		// ใช้ NamedExecContext เพื่อแทรกรายการทั้งหมด
		for _, param := range params {
			_, err = tx.NamedExecContext(ctx, query, param)
			if err != nil {
				return fmt.Errorf("failed to create trade return line: %w", err)
			}
		}

		return nil
	})
}

/************** Confirm To ReturnOrder ****************/
// รวม func. UpdateStatusToSuccess + GetBeforeOrderDetails + UpdateReturnOrderAndLines + InsertReturnOrderLine in service

// step 1: update status BeforeReturnOrder, เก็บค่าผู้ updateBy Date เพื่อนำไปใช้เข้าใน CreateBy Date => ReturnOrder,Line
func (repo repositoryDB) UpdateStatusToSuccess(ctx context.Context, orderNo, updateBy string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		query := `
            UPDATE BeforeReturnOrder
            SET StatusReturnID = 6, -- success status
                UpdateBy = :UpdateBy, 
                UpdateDate = GETDATE()
            WHERE OrderNo = :OrderNo
        `
		stmt, err := tx.PrepareNamed(query)
		if err != nil {
			log.Printf("Error preparing statement for OrderNo %s: %v", orderNo, err)
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(map[string]interface{}{
			"OrderNo":  orderNo,
			"UpdateBy": updateBy,
		})
		if err != nil {
			log.Printf("Error updating status to success for OrderNo %s: %v", orderNo, err)
			return fmt.Errorf("failed to update status to success: %w", err)
		}

		return nil
	})
}

// step 2: Fetch ค่า Befod ออกมา เก็บค่าผู้ updateBy Date เพื่อนำไปใช้เข้าใน CreateBy Date => ReturnOrder,Line
func (repo repositoryDB) GetBeforeOrderDetails(ctx context.Context, orderNo string) (*response.ConfirmReturnOrderDetails, error) {
	query := ` SELECT UpdateBy, UpdateDate
        	   FROM BeforeReturnOrder
               WHERE OrderNo = :OrderNo `

	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
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

// step 3: update
func (repo repositoryDB) UpdateReturnOrderAndLines(ctx context.Context, req request.ConfirmToReturnRequest, returnOrderData *response.ConfirmReturnOrderDetails) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 2: อัปเดต ReturnOrder
		for _, head := range req.UpdateToReturn {
			queryUpdateReturnOrder := ` UPDATE ReturnOrder
                                        SET StatusCheckID = 2, --CONFIRM status
                                            SrNo = :SrNo, 
                                            UpdateBy = :UpdateBy, 
                                            UpdateDate = :UpdateDate,
											CheckBy = :CheckBy, 
                                            CheckDate = :CheckDate
                                        WHERE OrderNo = :OrderNo `
			stmt, err := tx.PrepareNamed(queryUpdateReturnOrder)
			if err != nil {
				return fmt.Errorf("failed to prepare statement: %w", err)
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
				return fmt.Errorf("failed to update ReturnOrder: %w", err)
			}
		}

		// Step 3: อัปเดต ReturnOrderLine
		for _, line := range req.ImportLinesActual { // COALESCE => ฟิลด์ที่ไม่ได้ใช้จะดึงค่าเดิมมาแทน

			queryUpdateReturnOrderLine := ` UPDATE ReturnOrderLine
											SET SKU = COALESCE(:SKU, SKU),
												ActualQTY = COALESCE(:ActualQTY, ActualQTY),
												Price = COALESCE(:Price, Price),
												StatusDelete = COALESCE(:StatusDelete, StatusDelete),
												UpdateBy = COALESCE(:UpdateBy, UpdateBy),
												UpdateDate = COALESCE(:UpdateDate, UpdateDate),
												DeleteBy = COALESCE(:DeleteBy, DeleteBy),
												DeleteDate = COALESCE(:DeleteDate, DeleteDate)
											WHERE OrderNo = :OrderNo AND SKU = :SKU `
			stmt, err := tx.PrepareNamed(queryUpdateReturnOrderLine)
			if err != nil {
				return fmt.Errorf("failed to prepare statement: %w", err)
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
				"ActualQTY":    sql.NullInt32{Int32: int32(line.ActualQTY), Valid: line.ActualQTY != 0}, // เมื่อส่งค่า ว่าง/0 มาให้ใช้ค่าเดิม
				"Price":        sql.NullFloat64{Float64: line.Price, Valid: line.Price != 0},            // เมื่อส่งค่า ว่าง/0 มาให้ใช้ค่าเดิม
				"StatusDelete": sql.NullBool{Bool: line.StatusDelete, Valid: line.StatusDelete},         // เมื่อส่งค่าว่างมาให้ใช้ค่าเดิม
				"UpdateBy":     returnOrderData.UpdateBy,
				"UpdateDate":   returnOrderData.UpdateDate,
				"DeleteBy":     deleteBy,
				"DeleteDate":   deleteDate,
			})
			if err != nil {
				return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
			}
		}

		// Step 4: Commit Transaction
		return nil
	})
}

// Create Return Order MKP
func (repo repositoryDB) SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error) {
	// ตรวจสอบว่า soNo และ orderNo ไม่ว่างทั้งคู่ ถ้าว่างทั้งคู่จะคืนค่าข้อผิดพลาด
	if soNo == "" && orderNo == "" {
		return nil, fmt.Errorf("🚩 Either SoNo or OrderNo must be provided 🚩")
	}

	// คำสั่ง SQL สำหรับดึงข้อมูล OrderHead
	queryHead := `
        SELECT SoNo, OrderNo, StatusMKP, SalesStatus, CreateDate
        FROM ROM_V_OrderHeadDetail
        WHERE (:SoNo = '' OR SoNo = :SoNo) 
        AND (:OrderNo = '' OR OrderNo = :OrderNo)
    `

	// คำสั่ง SQL สำหรับดึงข้อมูล OrderLine
	queryLines := `
        SELECT SKU, ItemName, QTY, Price
        FROM ROM_V_OrderLineDetail
        WHERE (:SoNo = '' OR SoNo = :SoNo) 
        AND (:OrderNo = '' OR OrderNo = :OrderNo)
        ORDER BY RecID
    `

	params := map[string]interface{}{
		"SoNo":    soNo,
		"OrderNo": orderNo,
	}

	// ดึงข้อมูล OrderHead จากฐานข้อมูล
	var orderHead response.SaleOrderResponse
	stmtHead, err := repo.db.PrepareNamed(queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare head query: %w", err)
	}
	err = stmtHead.GetContext(ctx, &orderHead, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order head: %w", err)
	}

	// ดึงข้อมูล OrderLine จากฐานข้อมูล
	var orderLines []response.SaleOrderLineResponse
	stmtLines, err := repo.db.PrepareNamed(queryLines)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare lines query: %w", err)
	}
	err = stmtLines.SelectContext(ctx, &orderLines, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order lines: %w", err)
	}

	// รวมข้อมูล OrderHead และ OrderLines
	orderHead.OrderLines = orderLines

	// คืนค่าข้อมูล OrderHead ที่รวมกับ OrderLines
	return &orderHead, nil
}

func (repo repositoryDB) CreateSaleReturn(ctx context.Context, order request.CreateSaleReturnRequest) (*response.BeforeReturnOrderResponse, error) {
	// ✅ Start Transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// ✅ Insert Header (BeforeReturnOrder)
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatus, MkpStatus, ReturnDate, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :SoNo, :ChannelID, :Reason, :CustomerID, ISNULL(:TrackingNo, ''), :Logistic, :WarehouseID, 
            ISNULL(:SoStatus, ''), ISNULL(:MkpStatus, ''), :ReturnDate, :CreateBy, GETDATE()
        )
    `

	_, err = tx.NamedExecContext(ctx, queryOrder, &order)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create BeforeReturnOrder: %w", err)
	}

	// ✅ Insert Lines (BeforeReturnOrderLine)
	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE(), ISNULL(:TrackingNo, '')
        )
    `

	// ✅ ใช้ `NamedExecContext` Insert ถ้ามีข้อมูล
	if len(order.OrderLines) > 0 {
		_, err = tx.NamedExecContext(ctx, queryLine, order.OrderLines)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to create BeforeReturnOrderLine: %w", err)
		}
	}

	// ✅ Commit Transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("❌ Failed to commit transaction: %w", err)
	}

	// ✅ Fetch Created Order
	createdOrder, err := repo.GetBeforeReturnOrderByOrderNo(ctx, order.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch created order: %w", err)
	}

	return createdOrder, nil
}

func (repo repositoryDB) UpdateSaleReturn(ctx context.Context, req request.UpdateSaleReturn, userID string) error {
	// ✅ Start Transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// ✅ Update SrNo in BeforeReturnOrder
	query := `
        UPDATE BeforeReturnOrder
        SET SrNo = :SrNo,
            UpdateBy = :UserID,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `

	// 🔎 Debug Params
	params := map[string]interface{}{
		"SrNo":    req.SrNo,
		"UserID":  userID, // ✅ ใช้ `UserID` แทน `UpdateBy`
		"OrderNo": req.OrderNo,
	}
	fmt.Println("🔍 Debug Params:", params) // ✅ Log Debugging Params

	result, err := tx.NamedExecContext(ctx, query, params)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("❌ Failed to update sale return: %w", err)
	}

	// ✅ Check Rows Affected
	rows, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("❌ Failed to get rows affected: %w", err)
	}
	if rows == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("⚠️ No rows updated for order: %s", req.OrderNo)
	}

	// ✅ Commit Transaction
	return tx.Commit()
}

func (repo repositoryDB) ConfirmSaleReturn(ctx context.Context, orderNo string, statusReturnID, statusConfID int, userID string) error {
	// ✅ Update Order Status
	updateQuery := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = :StatusReturnID,
            StatusConfID = :StatusConfID,
            ConfirmBy = :ConfirmBy,
            ConfirmDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `

	// 🔎 Debug
	fmt.Println("🔍 Debug Params:", orderNo, statusReturnID, statusConfID, userID) // ✅ Log Debugging Params

	res, err := repo.db.NamedExecContext(ctx, updateQuery, map[string]interface{}{
		"OrderNo":        orderNo,
		"StatusReturnID": statusReturnID,
		"StatusConfID":   statusConfID,
		"ConfirmBy":      userID,
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to update return order status: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("❌ Failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("⚠️ No rows updated for order: %s", orderNo)
	}

	return nil
}

func (repo repositoryDB) CancelSaleReturn(ctx context.Context, req request.CancelSaleReturn, userID string) error {
	// ✅ Start Transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("❌ Failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// ✅ Insert Cancel Status & Retrieve CancelID
	insertCancelStatus := `
        INSERT INTO CancelStatus (RefID, CancelStatus, Remark, CancelBy, CancelDate) 
        OUTPUT INSERTED.CancelID
        VALUES (:OrderNo, 1, :Remark, :CancelBy, GETDATE())
    `

	var cancelID int64
	params := map[string]interface{}{
		"OrderNo":  req.OrderNo,
		"Remark":   req.Remark,
		"CancelBy": userID,
	}

	// 🔎 Debug
	fmt.Println("🔍 Debug Params:", params) // ✅ Log Debugging Params

	stmt, err := tx.PrepareNamedContext(ctx, insertCancelStatus)
	if err != nil {
		return fmt.Errorf("❌ Failed to prepare cancel status query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, params).Scan(&cancelID)
	if err != nil {
		return fmt.Errorf("❌ Failed to create cancel status: %w", err)
	}

	// ✅ Update Order Status
	updateOrder := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 2,
            StatusConfID = 3,
            CancelID = :CancelID,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `

	_, err = tx.NamedExecContext(ctx, updateOrder, map[string]interface{}{
		"OrderNo":  req.OrderNo,
		"CancelID": cancelID,
		"UpdateBy": userID,
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to update order status: %w", err)
	}

	// ✅ Commit Transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("❌ Failed to commit transaction: %w", err)
	}

	return nil
}

// Draft & Confirm MKP 🚨//
// ListDraftOrders ดึงรายการ Draft Status Orders 🚗
func (repo repositoryDB) ListDraftOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error) {
	query := `
        SELECT TOP 100 OrderNo, SoNo, SrNo, CustomerID, TrackingNo, Logistic, ChannelID, CreateDate, WarehouseID
        FROM BeforeReturnOrder
        WHERE StatusConfID = 1 -- Draft status
        AND CreateDate BETWEEN CONVERT(DATETIME, :startDate, 120) AND CONVERT(DATETIME, :endDate, 120)
        ORDER BY CreateDate DESC
    `

	var orders []response.ListDraftConfirmOrdersResponse

	err := repo.db.SelectContext(ctx, &orders, query, map[string]interface{}{
		"startDate": startDate + " 00:00:00",
		"endDate":   endDate + " 23:59:59",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch draft orders: %w", err)
	}

	return orders, nil
}

// ListConfirmOrders ดึงรายการ Confirm Satus Orders 🚗
func (repo repositoryDB) ListConfirmOrders(ctx context.Context, startDate, endDate string) ([]response.ListDraftConfirmOrdersResponse, error) {
	query := `
        SELECT OrderNo, SoNo, SrNo, CustomerID, TrackingNo, Logistic, ChannelID, CreateDate, WarehouseID
        FROM BeforeReturnOrder
        WHERE StatusConfID = 2 -- Confirm status
        AND CreateDate BETWEEN CONVERT(DATETIME, :startDate, 120) AND CONVERT(DATETIME, :endDate, 120)
        ORDER BY CreateDate DESC
    `

	var orders []response.ListDraftConfirmOrdersResponse

	err := repo.db.SelectContext(ctx, &orders, query, map[string]interface{}{
		"startDate": startDate + " 00:00:00",
		"endDate":   endDate + " 23:59:59",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch confirm orders: %w", err)
	}

	return orders, nil
}

// GetDraftConfirmOrderByOrderNo ดึงข้อมูล Draft Status Order หรือ Confirm Status Order ตาม OrderNo 🚗
func (repo repositoryDB) GetDraftConfirmOrderByOrderNo(ctx context.Context, orderNo string) (*response.DraftHeadResponse, error) {
	var head response.DraftHeadResponse
	var lines []response.DraftLineResponse

	headQuery := `
        SELECT 
            OrderNo,
            SoNo,
            SrNo
        FROM BeforeReturnOrder
        WHERE OrderNo = :OrderNo
    `
	err := repo.db.GetContext(ctx, &head, repo.db.Rebind(headQuery), orderNo)
	if err != nil {
		return nil, err
	}

	lineQuery := `
        SELECT 
            SKU,
            ItemName,
            QTY,
            Price
        FROM BeforeReturnOrderLine
        WHERE OrderNo = :OrderNo
    `
	err = repo.db.SelectContext(ctx, &lines, repo.db.Rebind(lineQuery), orderNo)
	if err != nil {
		return nil, err
	}

	head.OrderLines = lines

	return &head, nil
}

// ListCodeR ดึงรายการ CodeR ที่ขึ้นต้นด้วย 'R' 🚗
func (repo repositoryDB) ListCodeR(ctx context.Context) ([]response.ListCodeRResponse, error) {
	query := `
        SELECT SKU, ItemName
        FROM ROM_V_ProductAll
        WHERE SKU LIKE 'R%' -- ดึงเฉพาะ SKU ที่ขึ้นต้นด้วย 'R'
        ORDER BY ItemName ASC
    `

	var codeRList []response.ListCodeRResponse

	err := repo.db.SelectContext(ctx, &codeRList, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeR list: %w", err)
	}

	return codeRList, nil
}

func (repo repositoryDB) AddCodeR(ctx context.Context, req request.AddCodeR) ([]response.AddCodeRResponse, error) {
	query := `
        INSERT INTO BeforeReturnOrderLine (OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, CreateDate)
        OUTPUT inserted.OrderNo, inserted.SKU, inserted.ItemName, inserted.QTY, inserted.ReturnQTY, inserted.Price, inserted.CreateBy, inserted.CreateDate
        VALUES (:OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, GETDATE())
    `

	var results []response.AddCodeRResponse

	err := repo.db.SelectContext(ctx, &results, query, req)
	if err != nil {
		return nil, fmt.Errorf("failed to insert CodeR: %w", err)
	}

	return results, nil
}

func (repo repositoryDB) DeleteCodeR(ctx context.Context, orderNo string, sku string) (int64, error) {
	query := `
        DELETE FROM BeforeReturnOrderLine
        WHERE OrderNo = :OrderNo AND SKU = :SKU
    `

	result, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo": orderNo,
		"SKU":     sku,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to delete CodeR: %w", err)
	}

	// ✅ ตรวจสอบจำนวนแถวที่ถูกลบ
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rowsAffected, nil
}

func (repo repositoryDB) UpdateOrderStatus(ctx context.Context, orderNo string, statusConfID int, statusReturnID int, userID string) (*response.UpdateOrderStatusResponse, error) {
	query := `
        UPDATE BeforeReturnOrder
        SET StatusConfID = :StatusConfID,
            StatusReturnID = :StatusReturnID,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        OUTPUT inserted.OrderNo, inserted.StatusConfID, inserted.StatusReturnID, inserted.UpdateBy, inserted.UpdateDate
        WHERE OrderNo = :OrderNo
    `
	params := map[string]interface{}{
		"OrderNo":        orderNo,
		"StatusConfID":   statusConfID,
		"StatusReturnID": statusReturnID,
		"UpdateBy":       userID,
	}

	var updatedOrder response.UpdateOrderStatusResponse
	err := repo.db.GetContext(ctx, &updatedOrder, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &updatedOrder, nil
}

func (repo repositoryDB) CheckBefLineSKUExists(ctx context.Context, identifier, sku string) (bool, error) {
	query := ` SELECT 1 FROM BeforeReturnOrderLine 
               WHERE SKU = :SKU AND (OrderNo = :Identifier OR TrackingNo = :Identifier) `
	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
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

func (repo repositoryDB) CheckReLineSKUExists(ctx context.Context, orderNo, sku string) (bool, error) {
	query := ` SELECT 1 FROM ReturnOrderLine 
               WHERE SKU = :SKU AND OrderNo = :OrderNo `
	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
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

/************** Confirm Receipt ****************/
// รวม func. UpdateBefToWaiting + GetBeforeReturnOrderData + InsertReturnOrder + InsertReturnOrderLine in service

// 1. Update สถานะใน BeforeReturnOrder to "WAITING" (Page: Confirm Trade)
func (repo repositoryDB) UpdateBefToWaiting(ctx context.Context, req request.ConfirmTradeReturnRequest, updateBy string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		queryUpdate := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 7, -- WAITING status
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :Identifier OR TrackingNo = :Identifier
    `
		params := map[string]interface{}{
			"Identifier": req.Identifier,
			"UpdateBy":   updateBy,
		}
		_, err := tx.NamedExecContext(ctx, queryUpdate, params)
		return err
	})
}

// 2. ดึงข้อมูลจาก BeforeReturnOrder fetch ออกมาเพื่อเอาเข้า ReturnOrder
func (repo repositoryDB) GetBeforeReturnOrderData(ctx context.Context, req request.ConfirmTradeReturnRequest) (*response.ConfirmReturnOrderDetails, error) {
	querySelectOrder := `
        SELECT OrderNo, SoNo, SrNo, TrackingNo, ChannelID, Reason,
			   UpdateBy AS CreateBy, UpdateDate AS CreateDate
        FROM BeforeReturnOrder
        WHERE OrderNo = :Identifier OR TrackingNo = :Identifier
    `
	var returnOrderData response.ConfirmReturnOrderDetails

	rows, err := repo.db.NamedQueryContext(ctx, querySelectOrder, map[string]interface{}{
		"Identifier": req.Identifier,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.StructScan(&returnOrderData); err != nil {
			return nil, fmt.Errorf("failed to scan BeforeReturnOrder: %w", err)
		}
	}

	return &returnOrderData, nil
}

// 3. Insert ข้อมูลลงใน ReturnOrder
func (repo repositoryDB) InsertReturnOrder(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		queryInsertOrder := `
        INSERT INTO ReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, Reason, TrackingNo, CreateBy, CreateDate, StatusCheckID
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :TrackingNo, :CreateBy, :CreateDate, :StatusCheckID
        )
    `
		_, err := tx.NamedExecContext(ctx, queryInsertOrder, returnOrderData)
		return err

	})
}

// 4. Insert ข้อมูลจาก importLines ลงใน ReturnOrderLine
func (repo repositoryDB) InsertReturnOrderLine(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		queryInsertLine := `
        INSERT INTO ReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, TrackingNo, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :TrackingNo, :CreateBy, :CreateDate
        )
    `
		for _, line := range req.ImportLines {
			lineParams := map[string]interface{}{
				"OrderNo":    returnOrderData.OrderNo,
				"SKU":        line.SKU,
				"QTY":        line.QTY,
				"ReturnQTY":  line.ReturnQTY,
				"Price":      line.Price,
				"TrackingNo": returnOrderData.TrackingNo,
				"CreateBy":   returnOrderData.CreateBy,
				"CreateDate": returnOrderData.CreateDate,
			}
			_, err := tx.NamedExecContext(ctx, queryInsertLine, lineParams)
			if err != nil {
				return fmt.Errorf("failed to insert into ReturnOrderLine: %w", err)
			}
		}
		return nil
	})
}

// InsertImages ฟังก์ชันที่ใช้เพิ่มข้อมูลภาพลงในฐานข้อมูล
func (repo repositoryDB) InsertImages(ctx context.Context, returnOrderData *response.ConfirmReturnOrderDetails, req request.ConfirmTradeReturnRequest) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		queryInsertImage := `
        INSERT INTO Images (
            OrderNo, ImageTypeID, SKU, FilePath, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :ImageTypeID, :SKU, :FilePath, :CreateBy, :CreateDate
        )
    `
		for _, line := range req.ImportLines {
			imageParams := map[string]interface{}{
				"OrderNo":     returnOrderData.OrderNo,
				"ImageTypeID": line.ImageTypeID,
				"SKU":         line.SKU,
				"FilePath":    line.FilePath,
				"CreateBy":    returnOrderData.CreateBy,
				"CreateDate":  returnOrderData.CreateDate,
			}
			_, err := tx.NamedExecContext(ctx, queryInsertImage, imageParams)
			if err != nil {
				return fmt.Errorf("failed to insert into Images: %w", err)
			}
		}
		return nil
	})
}

/************************** Delete Line *************************/

func (repo repositoryDB) DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// ลบ BeforeReturnOrderLine ตาม RecID
		deleteQuery := `
			DELETE FROM BeforeReturnOrderLine
			WHERE RecID = :RecID
		`

		_, err := tx.NamedExecContext(ctx, deleteQuery, map[string]interface{}{
			"RecID": recID,
		})
		if err != nil {
			log.Printf("Error deleting BeforeReturnOrderLine by RecID: %v", err)
			return fmt.Errorf("failed to delete BeforeReturnOrderLine: %w", err)
		}

		return nil
	})
}

/************************** Get Order Head+Line *************************/

func (repo repositoryDB) GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error) {
	var headDetails []response.OrderHeadDetail
	var lineDetails []response.OrderLineDetail

	// Query Order Head
	headQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
        ORDER BY OrderNo
    `
	err := repo.db.SelectContext(ctx, &headDetails, headQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying OrderHeadDetail: %w", err)
	}

	// Query Order Line
	lineQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, SKU, ItemName, QTY, Price, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderLineDetail
        ORDER BY OrderNo
    `
	err = repo.db.SelectContext(ctx, &lineDetails, lineQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying OrderLineDetail: %w", err)
	}

	// Map Order Lines to Order Heads
	orderLineMap := make(map[string][]response.OrderLineDetail)
	for _, line := range lineDetails {
		orderLineMap[line.OrderNo] = append(orderLineMap[line.OrderNo], line)
	}

	for i := range headDetails {
		headDetails[i].OrderLineDetail = orderLineMap[headDetails[i].OrderNo]
	}

	return []response.OrderDetail{
		{OrderHeadDetail: headDetails},
	}, nil
}

/************************** Get Order Head+Line : Paginate *************************/

func (repo repositoryDB) GetAllOrderDetails(ctx context.Context, offset, limit int) ([]response.OrderDetail, error) {
	var headDetails []response.OrderHeadDetail
	var lineDetails []response.OrderLineDetail

	// Query Order Head with Pagination
	headQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
        ORDER BY OrderNo
        OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
    `
	err := repo.db.SelectContext(ctx, &headDetails, headQuery, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		log.Printf("Error querying OrderHeadDetail: %v", err)
		return nil, fmt.Errorf("error querying OrderHeadDetail: %w", err)
	}

	// Query Order Line
	lineQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, SKU, ItemName, QTY, Price, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderLineDetail
        WHERE OrderNo IN (
            SELECT OrderNo
            FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
            ORDER BY OrderNo
            OFFSET @offset ROWS FETCH NEXT @limit ROWS ONLY
        )
        ORDER BY OrderNo
    `
	err = repo.db.SelectContext(ctx, &lineDetails, lineQuery, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		log.Printf("Error querying OrderLineDetail: %v", err)
		return nil, fmt.Errorf("error querying OrderLineDetail: %w", err)
	}

	// Map Order Lines to Order Heads
	orderLineMap := make(map[string][]response.OrderLineDetail)
	for _, line := range lineDetails {
		orderLineMap[line.OrderNo] = append(orderLineMap[line.OrderNo], line)
	}

	for i := range headDetails {
		headDetails[i].OrderLineDetail = orderLineMap[headDetails[i].OrderNo]
	}

	return []response.OrderDetail{
		{OrderHeadDetail: headDetails},
	}, nil
}

/************************** Search by SO *************************/

func (repo repositoryDB) GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error) {
	var headDetails []response.OrderHeadDetail
	var lineDetails []response.OrderLineDetail

	// Query Order Head
	headQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderHeadDetail
        WHERE SoNo = @SoNo
    `
	err := repo.db.SelectContext(ctx, &headDetails, headQuery, sql.Named("SoNo", soNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("Error querying OrderHeadDetail by SO: %v", err)
		return nil, fmt.Errorf("error querying OrderHeadDetail by SO: %w", err)
	}

	// Query Order Line
	lineQuery := `
        SELECT OrderNo, SoNo, StatusMKP, SalesStatus, SKU, ItemName, QTY, Price, CreateDate
        FROM Data_WebReturn.dbo.ROM_V_OrderLineDetail
        WHERE SoNo = @SoNo
    `
	err = repo.db.SelectContext(ctx, &lineDetails, lineQuery, sql.Named("SoNo", soNo))
	if err != nil {
		log.Printf("Error querying OrderLineDetail by SO: %v", err)
		return nil, fmt.Errorf("error querying OrderLineDetail by SO: %w", err)
	}

	// Map Order Lines to Order Heads
	orderLineMap := make(map[string][]response.OrderLineDetail)
	for _, line := range lineDetails {
		orderLineMap[line.OrderNo] = append(orderLineMap[line.OrderNo], line)
	}

	for i := range headDetails {
		headDetails[i].OrderLineDetail = orderLineMap[headDetails[i].OrderNo]
	}

	return &response.OrderDetail{
		OrderHeadDetail: headDetails,
	}, nil
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

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy
        )
    `
	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
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
		return nil, fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range order.BeforeReturnOrderLines {
		// Ensure TrackingNo is not NULL
		// trackingNo := line.TrackingNo
		// if trackingNo == "" {
		// 	trackingNo = "N/A" // Default value if TrackingNo is not provided
		// }

		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"ItemName":   line.ItemName,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   order.CreateBy,
			"TrackingNo": order.TrackingNo,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
		}
	}

	// 4. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Logging
	fmt.Println("Transaction committed")

	// 5. ดึงข้อมูลที่สร้างเสร็จแล้ว
	createdOrder, err := repo.GetBeforeReturnOrderByOrderNo(ctx, order.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created order: %w", err)
	}

	// Logging
	fmt.Println("Fetched created order")

	return createdOrder, nil
}

// CRUD Dafault //
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CancelID
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, :CancelID
        )
    `
	paramsOrder := map[string]interface{}{
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
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"CreateBy":       order.CreateBy,
		"CancelID":       order.CancelID,
	}

	_, err := repo.db.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	return nil
}

func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range lines {
		params := map[string]interface{}{
			"OrderNo":   orderNo,
			"SKU":       line.SKU,
			"ItemName":  line.ItemName,
			"QTY":       line.QTY,
			"ReturnQTY": line.ReturnQTY,
			"Price":     line.Price,
			"CreateBy":  line.CreateBy,
		}
		_, err := repo.db.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("failed to create order line: %w", err)
		}
	}
	return nil
}

// Implementation สำหรับ Transaction CreateBeforeReturnOrder & CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :Reason, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatus, :MkpStatus, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy
        )
    `
	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
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
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"CreateBy":       order.CreateBy,
	})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, ItemName, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :ItemName, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range order.BeforeReturnOrderLines {
		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo":   order.OrderNo,
			"SKU":       line.SKU,
			"ItemName":  line.ItemName,
			"QTY":       line.QTY,
			"ReturnQTY": line.ReturnQTY,
			"Price":     line.Price,
			"CreateBy":  order.CreateBy,
		})
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create BeforeReturnOrderLine: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	query := `
        SELECT OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder
        ORDER BY RecID ASC
    `

	var orders []response.BeforeReturnOrderResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &orders, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	for i := range orders {
		lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orders[i].OrderNo)
		if err != nil {
			return nil, err
		}
		orders[i].BeforeReturnOrderLines = lines
	}

	return orders, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
			ItemName,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine
        ORDER BY RecID
    `

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	return lines, nil
}

func (repo repositoryDB) ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
			ItemName,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine
        WHERE OrderNo = :OrderNo
        ORDER BY RecID
    `

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	return lines, nil
}

func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	query := `
        SELECT OrderNo, SoNo, SrNo, ChannelID, Reason, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatus, MkpStatus, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID, IsCNCreated, IsEdited 
        FROM BeforeReturnOrder
        WHERE OrderNo = :OrderNo
    `
	order := new(response.BeforeReturnOrderResponse)
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, order, map[string]interface{}{"OrderNo": orderNo})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch BeforeReturnOrder: %w", err)
	}

	lines, err := repo.ListBeforeReturnOrderLinesByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	order.BeforeReturnOrderLines = lines

	return order, nil
}

func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
			ItemName,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine
        WHERE OrderNo = :OrderNo
    `

	var lines []response.BeforeReturnOrderLineResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &lines, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}

	fmt.Printf("Fetched %d lines from the database for OrderNo: %s\n", len(lines), orderNo) // Add logging for the number of lines

	return lines, nil
}

func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	query := `
        UPDATE BeforeReturnOrder 
        SET SoNo = COALESCE(:SoNo, SoNo),
            SrNo = COALESCE(:SrNo, SrNo),
            ChannelID = COALESCE(:ChannelID, ChannelID),
            Reason = COALESCE(:Reason, Reason),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatus = COALESCE(:SoStatus, SoStatus),
            MkpStatus = COALESCE(:MkpStatus, MkpStatus),
            ReturnDate = COALESCE(:ReturnDate, ReturnDate),
            StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
            StatusConfID = COALESCE(:StatusConfID, StatusConfID),
            ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy)
        WHERE OrderNo = :OrderNo
    `
	params := map[string]interface{}{
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
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"UpdateBy":       order.UpdateBy,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	return nil
}

func (repo repositoryDB) UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error {
	query := `
        UPDATE BeforeReturnOrderLine 
        SET ItemName = COALESCE(:ItemName, ItemName),
			QTY = COALESCE(:QTY, QTY),
            ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
            Price = COALESCE(:Price, Price),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo)
        WHERE OrderNo = :OrderNo
          AND SKU = :SKU
    `
	params := map[string]interface{}{
		"OrderNo":    orderNo,
		"SKU":        line.SKU,
		"ItemName":   line.ItemName,
		"QTY":        line.QTY,
		"ReturnQTY":  line.ReturnQTY,
		"Price":      line.Price,
		"UpdateBy":   line.UpdateBy,
		"TrackingNo": line.TrackingNo,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
	}

	return nil
}

// Implementation สำหรับ UpdateBeforeReturnOrderWithTransaction
func (repo repositoryDB) UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
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

	// Update BeforeReturnOrderLine first
	queryLine := `
        UPDATE BeforeReturnOrderLine 
        SET ItemName = COALESCE(:ItemName, ItemName),
			QTY = COALESCE(:QTY, QTY),
            ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
            Price = COALESCE(:Price, Price),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo)
        WHERE OrderNo = :OrderNo
          AND SKU = :SKU
    `

	for _, line := range order.BeforeReturnOrderLines {
		paramsLine := map[string]interface{}{
			"OrderNo":    line.OrderNo,
			"SKU":        line.SKU,
			"ItemName":   line.ItemName,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"UpdateBy":   line.UpdateBy,
			"TrackingNo": line.TrackingNo,
		}

		result, err := tx.NamedExecContext(ctx, queryLine, paramsLine)
		if err != nil {
			return fmt.Errorf("failed to update BeforeReturnOrderLine: %w", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("no rows updated for OrderNo: %s, SKU: %s", line.OrderNo, line.SKU)
		}
	}

	// Update BeforeReturnOrder
	queryOrder := `
        UPDATE BeforeReturnOrder 
        SET SoNo = COALESCE(:SoNo, SoNo),
            SrNo = COALESCE(:SrNo, SrNo),
            ChannelID = COALESCE(:ChannelID, ChannelID),
            Reason = COALESCE(:Reason, Reason),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatus = COALESCE(:SoStatus, SoStatus),
            MkpStatus = COALESCE(:MkpStatus, MkpStatus),
            ReturnDate = COALESCE(:ReturnDate, ReturnDate),
            StatusReturnID = COALESCE(:StatusReturnID, StatusReturnID),
            StatusConfID = COALESCE(:StatusConfID, StatusConfID),
            ConfirmBy = COALESCE(:ConfirmBy, ConfirmBy),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            CancelID = COALESCE(:CancelID, CancelID)
        WHERE OrderNo = :OrderNo
    `

	paramsOrder := map[string]interface{}{
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
		"StatusReturnID": order.StatusReturnID,
		"StatusConfID":   order.StatusConfID,
		"ConfirmBy":      order.ConfirmBy,
		"UpdateBy":       order.UpdateBy,
		"CancelID":       order.CancelID,
	}

	_, err = tx.NamedExecContext(ctx, queryOrder, paramsOrder)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	return nil
}
