package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// เพิ่ม constant สำหรับ timeout
const (
	defaultTimeout = 10 * time.Second
)

// ReturnOrderRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูล
type BefRORepository interface {
	// Create
	CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error

	// Read
	ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error)
	GetBeforeReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.BeforeReturnOrderResponse, error)
	ListBeforeReturnOrderLines(ctx context.Context) ([]response.BeforeReturnOrderLineResponse, error)
	ListBeforeReturnOrderLinesByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)
	GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error)

	// Update
	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error
	UpdateSrNo(ctx context.Context, orderNo string, srNo string) error
	UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error

	// Transaction
	CreateReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error

	//Cancle

	//Search
	SearchSaleOrder(ctx context.Context, soNo string) (*response.SaleOrderResponse, error)
}

// Implementation สำหรับ CreateBeforeReturnOrder
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, CancelID
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, GETDATE(), :CancelID
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
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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

// Implementation สำหรับ GetBeforeReturnOrderLineByOrderNo
func (repo repositoryDB) GetBeforeReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.BeforeReturnOrderLineResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT 
            OrderNo,
            SKU,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine WITH (NOLOCK)
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
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
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
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT 
            OrderNo,
            SKU,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine WITH (NOLOCK)
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
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT 
            OrderNo,
            SKU,
            QTY,
            ReturnQTY,
            Price,
            TrackingNo,
            CreateDate
        FROM BeforeReturnOrderLine WITH (NOLOCK)
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

func (repo repositoryDB) ListBeforeReturnOrders(ctx context.Context) ([]response.BeforeReturnOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        SELECT OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate, UpdateBy, UpdateDate, CancelID
        FROM BeforeReturnOrder WITH (NOLOCK)
        ORDER BY CreateDate ASC
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

// Implementation สำหรับ BeginTransaction CreateBeforeReturnOrder & CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatusID, MkpStatusID, ReturnDate, StatusReturnID, StatusConfID, ConfirmBy, CreateBy, CreateDate
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatusID, :MkpStatusID, :ReturnDate, :StatusReturnID, :StatusConfID, :ConfirmBy, :CreateBy, GETDATE()
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

func (repo repositoryDB) UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

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
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
            CancelID = COALESCE(:CancelID, CancelID)
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
		"CancelID":       order.CancelID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update BeforeReturnOrder: %w", err)
	}

	return nil
}

func (repo repositoryDB) UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        UPDATE BeforeReturnOrderLine 
        SET QTY = COALESCE(:QTY, QTY),
            ReturnQTY = COALESCE(:ReturnQTY, ReturnQTY),
            Price = COALESCE(:Price, Price),
            UpdateBy = COALESCE(:UpdateBy, UpdateBy),
            UpdateDate = GETDATE(),
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
            UpdateDate = GETDATE(),
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
            UpdateDate = GETDATE(),
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

func (repo repositoryDB) SearchSaleOrder(ctx context.Context, soNo string) (*response.SaleOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	queryHead := `
        SELECT SoNo, OrderNo, StatusMKP, SalesStatus, CreateDate
        FROM ROM_V_OrderHeadDetail
        WHERE SoNo = :SoNo
    `

	queryLines := `
        SELECT SKU, ItemName, QTY, Price
        FROM ROM_V_OrderLineDetail
        WHERE SoNo = :SoNo
    `

	var orderHead response.SaleOrderResponse
	nstmtHead, err := repo.db.PrepareNamed(queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for order head: %w", err)
	}
	defer nstmtHead.Close()

	err = nstmtHead.GetContext(ctx, &orderHead, map[string]interface{}{"SoNo": soNo})
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

	err = nstmtLines.SelectContext(ctx, &orderLines, map[string]interface{}{"SoNo": soNo})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order lines: %w", err)
	}

	orderHead.OrderLines = orderLines

	return &orderHead, nil
}

// UpdateSaleReturnSR updates the SR(Sale Return) number for a given order
func (repo repositoryDB) UpdateSrNo(ctx context.Context, orderNo string, srNo string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `
        UPDATE BeforeReturnOrder
        SET SrNo = :SrNo
        WHERE OrderNo = :OrderNo
    `
	params := map[string]interface{}{
		"OrderNo": orderNo,
		"SrNo":    srNo,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update SaleReturn SR number: %w", err)
	}

	return nil
}

// ConfirmSaleReturn confirms a sale return order
func (repo repositoryDB) ConfirmSaleReturn(ctx context.Context, order request.BeforeReturnOrder) error {
	query := `
        UPDATE BeforeReturnOrder
        SET StatusConfID = :StatusConfID,
            ConfirmBy = :ConfirmBy,
            UpdateBy = :UpdateBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `
	_, err := repo.db.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo":      order.OrderNo,
		"StatusConfID": order.StatusConfID,
		"ConfirmBy":    order.ConfirmBy,
		"UpdateBy":     order.UpdateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to confirm BeforeReturnOrder: %w", err)
	}

	return nil
}

// UpdateDynamicFields updates specified fields for a given order
func (repo repositoryDB) UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	setClause := ""
	params := map[string]interface{}{"OrderNo": orderNo}
	for key, value := range fields {
		setClause += fmt.Sprintf("%s = :%s, ", key, key)
		params[key] = value
	}
	setClause = setClause[:len(setClause)-2] // Remove trailing comma and space

	query := fmt.Sprintf(`
        UPDATE BeforeReturnOrder
        SET %s
        WHERE OrderNo = :OrderNo
    `, setClause)

	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update dynamic fields: %w", err)
	}

	return nil
}
