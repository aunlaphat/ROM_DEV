package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
)

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type BefRORepository interface {
	// Create
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error
	CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error
	// Read
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)

	//Update
	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error
	UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error
	// Update
	//UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error

	//Cancle

	// ************************ Create Sale Return ************************ //
	//Search Order (SoNo and OrderNo from Order - MKP)
	SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error)
	//Create
	CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	//Insert SrNo (SR Create from AX)
	UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error

	ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error

	CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	// Draft & Confirm
	//ListDrafts(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	//EditOrder(ctx context.Context, req request.EditOrderRequest) error
}

// Implementation สำหรับ CreateBeforeReturnOrder
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CancelID
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, :CancelID
        )
    `
	paramsOrder := map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"SoNo":           order.SoNo,
		"SrNo":           order.SrNo,
		"ChannelID":      order.ChannelID,
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
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

// Implementation สำหรับ CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range lines {
		trackingNo := line.TrackingNo
		if trackingNo == "" {
			trackingNo = "N/A"
		}

		params := map[string]interface{}{
			"OrderNo":    orderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   line.CreateBy,
			"TrackingNo": trackingNo,
		}
		_, err := repo.db.NamedExecContext(ctx, query, params)
		if err != nil {
			return fmt.Errorf("failed to create order line: %w", err)
		}
	}
	return nil
}

// Implementation สำหรับ BeginTransaction CreateBeforeReturnOrder & CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy
        )
    `
	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
		"OrderNo":        order.OrderNo,
		"SoNo":           order.SoNo,
		"SrNo":           order.SrNo,
		"ChannelID":      order.ChannelID,
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
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
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range order.BeforeReturnOrderLines {
		// Ensure TrackingNo is not NULL
		trackingNo := line.TrackingNo
		if trackingNo == "" {
			trackingNo = "N/A" // Default value if TrackingNo is not provided
		}

		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   "SYSTEM",
			"TrackingNo": trackingNo,
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

// Implementation สำหรับ GetBeforeReturnOrderLineByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
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

	return lines, nil
}

// Implementation สำหรับ GetBeforeReturnOrderByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error) {
	query := `
        SELECT OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
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

func (repo repositoryDB) ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
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

func (repo repositoryDB) ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error) {
	query := `
        SELECT 
            OrderNo,
            SKU,
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

// ฟังก์ชันพื้นฐานสำหรับการดึงข้อมูล
func (repo repositoryDB) listBeforeReturnOrders(ctx context.Context, condition string, params map[string]interface{}) ([]response.BeforeReturnOrderResponse, error) {
	query := fmt.Sprintf(`
        SELECT OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder
        WHERE %s
        ORDER BY CreateDate ASC
    `, condition)

	var orders []response.BeforeReturnOrderResponse
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &orders, params)
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

// ฟังก์ชันสำหรับดึงข้อมูล BeforeReturnOrders ทั้งหมด
func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	return repo.listBeforeReturnOrders(ctx, "1=1", map[string]interface{}{})
}

// ฟังก์ชันสำหรับดึงข้อมูล Drafts
func (repo repositoryDB) ListDrafts(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	return repo.listBeforeReturnOrders(ctx, "StatusConfID = 1", map[string]interface{}{})
}

func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	query := `
        UPDATE BeforeReturnOrder 
        SET SoNo = COALESCE(:SoNo, SoNo),
            SrNo = COALESCE(:SrNo, SrNo),
            ChannelID = COALESCE(:ChannelID, ChannelID),
            ReturnType = COALESCE(:ReturnType, ReturnType),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatusID = COALESCE(:SoStatusID, SoStatusID),
            MkpStatusID = COALESCE(:MkpStatusID, MkpStatusID),
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
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
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
        SET QTY = COALESCE(:QTY, QTY),
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
        SET QTY = COALESCE(:QTY, QTY),
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
            ReturnType = COALESCE(:ReturnType, ReturnType),
            CustomerID = COALESCE(:CustomerID, CustomerID),
            TrackingNo = COALESCE(:TrackingNo, TrackingNo),
            Logistic = COALESCE(:Logistic, Logistic),
            WarehouseID = COALESCE(:WarehouseID, WarehouseID),
            SoStatusID = COALESCE(:SoStatusID, SoStatusID),
            MkpStatusID = COALESCE(:MkpStatusID, MkpStatusID),
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
		"ReturnType":     order.ReturnType,
		"CustomerID":     order.CustomerID,
		"TrackingNo":     order.TrackingNo,
		"Logistic":       order.Logistic,
		"WarehouseID":    order.WarehouseID,
		"SoStatusID":     order.SoStatusID,
		"MkpStatusID":    order.MkpStatusID,
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

func (repo repositoryDB) SearchOrder(ctx context.Context, soNo, orderNo string) (*response.SaleOrderResponse, error) {
	queryHead := `
        SELECT SoNo, OrderNo, StatusMKP, SalesStatus, CreateDate
        FROM ROM_V_OrderHeadDetail
        WHERE SoNo = :SoNo OR OrderNo = :OrderNo
    `

	queryLines := `
        SELECT SKU, ItemName, QTY, Price
        FROM ROM_V_OrderLineDetail
        WHERE SoNo = :SoNo OR OrderNo = :OrderNo
    `

	var orderHead response.SaleOrderResponse
	nstmtHead, err := repo.db.PrepareNamed(queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for order head: %w", err)
	}
	defer nstmtHead.Close()

	err = nstmtHead.GetContext(ctx, &orderHead, map[string]interface{}{"SoNo": soNo, "OrderNo": orderNo})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order head: %w", err)
	}

	var orderLines []response.SaleOrderLineResponse
	nstmtLines, err := repo.db.PrepareNamed(queryLines)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for order lines: %w", err)
	}
	defer nstmtLines.Close()

	err = nstmtLines.SelectContext(ctx, &orderLines, map[string]interface{}{"SoNo": soNo, "OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order lines: %w", err)
	}

	orderHead.OrderLines = orderLines

	return &orderHead, nil
}

func (repo repositoryDB) CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
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

	// 2. Insert BeforeReturnOrder (Header)
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatusID, MkpStatusID, ReturnDate, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatusID, :MkpStatusID, :ReturnDate, :CreateBy, GETDATE()
        )
    `
	_, err = tx.NamedExecContext(ctx, queryOrder, map[string]interface{}{
		"OrderNo":     order.OrderNo,
		"SoNo":        order.SoNo,
		"SrNo":        order.SrNo,
		"ChannelID":   order.ChannelID,
		"ReturnType":  order.ReturnType,
		"CustomerID":  order.CustomerID,
		"TrackingNo":  order.TrackingNo,
		"Logistic":    order.Logistic,
		"WarehouseID": order.WarehouseID,
		"SoStatusID":  order.SoStatusID,
		"MkpStatusID": order.MkpStatusID,
		"ReturnDate":  order.ReturnDate,
		"CreateBy":    order.CreateBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	// Logging
	fmt.Println("Inserted BeforeReturnOrder")

	// 3. Insert BeforeReturnOrderLine (Lines)
	queryLine := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo
        )
    `
	for _, line := range order.BeforeReturnOrderLines {
		_, err = tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo":    order.OrderNo,
			"SKU":        line.SKU,
			"QTY":        line.QTY,
			"ReturnQTY":  line.ReturnQTY,
			"Price":      line.Price,
			"CreateBy":   line.CreateBy,
			"TrackingNo": line.TrackingNo,
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

func (repo repositoryDB) UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error {
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

	// 1. SQL query สำหรับ update
	query := `
        UPDATE BeforeReturnOrder
        SET SrNo = :SrNo,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `

	// 2. กำหนด parameters
	params := map[string]interface{}{
		"OrderNo":  orderNo,
		"SrNo":     srNo,
		"UpdateBy": updateBy,
	}

	// 3. Execute query
	result, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update SR number: %w", err)
	}

	// 4. ตรวจสอบว่ามีการอัพเดทจริง
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	return nil
}

func (repo repositoryDB) ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error {
	// เริ่ม transaction
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

	// 1. กำหนด SQL query สำหรับ update สถานะ
	query := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 1, -- Pending status
            StatusConfID = 1,   -- Draft status
            ConfirmBy = :ConfirmBy,
            ComfirmDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `
	// 2. กำหนด parameters สำหรับ query
	params := map[string]interface{}{
		"OrderNo":   orderNo,
		"ConfirmBy": confirmBy,
	}

	// 3. เตรียม statement
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for confirming sale return: %w", err)
	}
	defer nstmt.Close()

	// 4. execute query
	_, err = nstmt.ExecContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to confirm sale return: %w", err)
	}

	return nil
}

func (repo repositoryDB) CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error {
	// 1. เริ่ม transaction
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

	// 2. สร้าง CancelID ใหม่
	var newCancelID int
	queryCancelID := `
        INSERT INTO CancelStatus (
            RefID, CancelStatus, Remark, CancelBy, CancelDate
        ) VALUES (
            :OrderNo, 1, :Remark, :CancelBy, GETDATE()
        );
        SELECT SCOPE_IDENTITY() AS CancelID;
    `
	err = tx.GetContext(ctx, &newCancelID, queryCancelID, map[string]interface{}{
		"OrderNo":  orderNo,
		"Remark":   remark,
		"CancelBy": updateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to create cancel status: %w", err)
	}

	// 3. อัพเดทสถานะการยกเลิกใน BeforeReturnOrder ด้วย CancelID ที่สร้างใหม่
	queryUpdate := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 2,    -- Cancel status
            StatusConfID = 3,      -- Cancel status
            CancelID = :CancelID,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `
	result, err := tx.NamedExecContext(ctx, queryUpdate, map[string]interface{}{
		"OrderNo":  orderNo,
		"CancelID": newCancelID,
		"UpdateBy": updateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// 4. ตรวจสอบว่ามีการอัพเดทจริง
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	return nil
}
