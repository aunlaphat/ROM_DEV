package repository

import (
	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// เพิ่ม constant สำหรับ timeout
const (
	defaultTimeout = 10 * time.Second
	txTimeout      = 30 * time.Second
)

// เพิ่ม constants สำหรับ status
const (
	StatusPending    = 1
	StatusInProgress = 2
	StatusCompleted  = 3
	StatusCancelled  = 4
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

	//SO
	GetAllOrderDetail(ctx context.Context) ([]response.OrderDetail, error)
	GetAllOrderDetails(ctx context.Context, offset, limit int) ([]response.OrderDetail, error)
	GetOrderDetailBySO(ctx context.Context, soNo string) (*response.OrderDetail, error)

	// Update
	UpdateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error
	UpdateBeforeReturnOrderLine(ctx context.Context, orderNo string, line request.BeforeReturnOrderLine) error
	UpdateBeforeReturnOrderWithTransaction(ctx context.Context, order request.BeforeReturnOrder) error
	// Update
	//UpdateDynamicFields(ctx context.Context, orderNo string, fields map[string]interface{}) error

	//Cancle
	DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error

	// ************************ Create Sale Return ************************ //
	//Search SoNo (Sale Order from Order - MKP)
	SearchSaleOrder(ctx context.Context, soNo string) (*response.SaleOrderResponse, error)
	//Create
	CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error)
	//Insert SrNo (SR Create from AX)
	UpdateSaleReturn(ctx context.Context, orderNo string, srNo string, updateBy string) error

	ConfirmSaleReturn(ctx context.Context, orderNo string, confirmBy string) error

	CancelSaleReturn(ctx context.Context, orderNo string, updateBy string, remark string) error

	CheckOrderNoExists(ctx context.Context, orderNo string) (bool, error)
	CreateTradeReturnLine(ctx context.Context, orderNo string, line request.TradeReturnLineRequest) error
	ConfirmOrderNo(ctx context.Context, orderNo string, confirmBy string) error
	CancelOrderNo(ctx context.Context, orderNo string, cancelBy string, remark string) error
	GetTrackingNoByOrderNo(ctx context.Context, orderNo string) (string, error)
}

// CheckOrderExists ตรวจสอบว่ามี OrderNo ใน BeforeReturnOrder หรือไม่
func (repo repositoryDB) CheckOrderNoExists(ctx context.Context, orderNo string) (bool, error) {
	var exists bool

	query := `
		SELECT CASE 
			WHEN EXISTS (SELECT 1 FROM BeforeReturnOrder WHERE OrderNo = @OrderNo) 
			THEN 1 ELSE 0 
		END
	`

	// ใช้ QueryRowContext เพื่อระบุ context
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}

	return exists, nil
}

func (repo repositoryDB) GetTrackingNoByOrderNo(ctx context.Context, orderNo string) (string, error) {
	var trackingNo string
	query := `
        SELECT TrackingNo
        FROM BeforeReturnOrder
        WHERE OrderNo = @OrderNo
    `
	err := repo.db.QueryRowContext(ctx, query, sql.Named("OrderNo", orderNo)).Scan(&trackingNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("order not found: %s", orderNo)
		}
		return "", fmt.Errorf("failed to fetch TrackingNo: %w", err)
	}
	return trackingNo, nil
}


// CreateTradeReturnLine สร้างข้อมูลใน BeforeReturnOrderLine
func (repo repositoryDB) CreateTradeReturnLine(ctx context.Context, orderNo string, line request.TradeReturnLineRequest) error {
	// ตรวจสอบว่า OrderNo มีอยู่ใน BeforeReturnOrder หรือไม่
	exists, err := repo.CheckOrderNoExists(ctx, orderNo)
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

	// สร้างข้อมูล BeforeReturnOrderLine ด้วย NamedExecContext
	query := `
        INSERT INTO BeforeReturnOrderLine (
            OrderNo, SKU, QTY, ReturnQTY, Price, CreateBy, TrackingNo, CreateDate
        ) VALUES (
            :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :CreateBy, :TrackingNo, GETDATE()
        )
    `
	params := map[string]interface{}{
		"OrderNo":    orderNo,
		"SKU":        line.SKU,
		"QTY":        line.QTY,
		"ReturnQTY":  line.ReturnQTY,
		"Price":      line.Price,
		"CreateBy":   "user", // กำหนด "user" สำหรับฟิลด์ CreateBy
		"TrackingNo": trackingNo,
	}

	_, err = repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to create trade return line: %w", err)
	}

	return nil
}

func (repo repositoryDB) ConfirmOrderNo(ctx context.Context, orderNo string, confirmBy string) error {
	// 1. กำหนด SQL query สำหรับ update สถานะ
	query := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 1, -- Pending status
            StatusConfID = 1,   -- Draft status
            ConfirmBy = :ConfirmBy
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

func (repo repositoryDB) CancelOrderNo(ctx context.Context, orderNo string, cancelBy string, remark string) error {
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

	// 2. อัพเดทสถานะการยกเลิก
	query := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 2, -- Cancel status
            StatusConfID = 3,   -- Cancel status
            CancelID = (SELECT ISNULL(MAX(CancelID), 0) + 1 FROM CancelStatus), -- สร้าง CancelID ใหม่
            UpdateBy = :CancelBy
        WHERE OrderNo = :OrderNo
    `
	result, err := tx.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo":  orderNo,
		"CancelBy": cancelBy,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// 3. บันทึกประวัติการยกเลิก
	queryCancelStatus := `
        INSERT INTO CancelStatus (
            RefID, CancelStatus, Remark, CancelBy, CancelDate
        ) VALUES (
            :OrderNo, 1, :Remark, :CancelBy, GETDATE()
        )
    `
	_, err = tx.NamedExecContext(ctx, queryCancelStatus, map[string]interface{}{
		"OrderNo":  orderNo,
		"Remark":   remark,
		"CancelBy": cancelBy,
	})
	if err != nil {
		return fmt.Errorf("failed to insert cancel status: %w", err)
	}

	return nil
}

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

// Implementation สำหรับ CreateBeforeReturnOrder
func (repo repositoryDB) CreateBeforeReturnOrder(ctx context.Context, order request.BeforeReturnOrder) error {
	log.Printf("🚀 Starting CreateBeforeReturnOrder for OrderNo: %s", order.OrderNo)
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
		"SoNo":      order.SoNo,
		"SrNo":     order.SrNo,
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
		log.Printf("❌ Failed to create BeforeReturnOrder: %v", err)
		return fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	log.Printf("✅ Successfully created BeforeReturnOrder for OrderNo: %s", order.OrderNo)
	return nil
}

// Implementation สำหรับ CreateBeforeReturnOrderLine
func (repo repositoryDB) CreateBeforeReturnOrderLine(ctx context.Context, orderNo string, lines []request.BeforeReturnOrderLine) error {
	log.Printf("🚀 Starting CreateBeforeReturnOrderLine for OrderNo: %s", orderNo)
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
			log.Printf("❌ Failed to create order line: %v", err)
			return fmt.Errorf("failed to create order line: %w", err)
		}
	}
	log.Printf("✅ Successfully created BeforeReturnOrderLine for OrderNo: %s", orderNo)
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

func (repo repositoryDB) DeleteBeforeReturnOrderLine(ctx context.Context, recID string) error {
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
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

func (repo repositoryDB) SearchSaleOrder(ctx context.Context, soNo string) (*response.SaleOrderResponse, error) {
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

func (repo repositoryDB) CreateSaleReturn(ctx context.Context, order request.BeforeReturnOrder) (*response.BeforeReturnOrderResponse, error) {
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

	// Insert into BeforeReturnOrder
	queryOrder := `
        INSERT INTO BeforeReturnOrder (
            OrderNo, SoNo, SrNo, ChannelID, ReturnType, CustomerID, TrackingNo, Logistic, WarehouseID, 
            SoStatusID, MkpStatusID, ReturnDate, CreateBy
        ) VALUES (
            :OrderNo, :SoNo, :SrNo, :ChannelID, :ReturnType, :CustomerID, :TrackingNo, :Logistic, :WarehouseID, 
            :SoStatusID, :MkpStatusID, :ReturnDate, :CreateBy
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
		"CreateBy":       order.CreateBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create BeforeReturnOrder: %w", err)
	}

	// Insert into BeforeReturnOrderLine
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

	// Fetch the created order to return as response
	createdOrder, err := repo.GetBeforeReturnOrderByOrderNo(ctx, order.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created order: %w", err)
	}

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

	// 2. อัพเดทสถานะการยกเลิก
	query := `
        UPDATE BeforeReturnOrder
        SET StatusReturnID = 2, -- Cancel status
            StatusConfID = 3,   -- Cancel status
            CancelID = (SELECT ISNULL(MAX(CancelID), 0) + 1 FROM CancelStatus), -- สร้าง CancelID ใหม่
            UpdateBy = :CancelBy,
            UpdateDate = GETDATE()
        WHERE OrderNo = :OrderNo
    `
	result, err := tx.NamedExecContext(ctx, query, map[string]interface{}{
		"OrderNo":  orderNo,
		"UpdateBy": updateBy,
		"CancelBy": updateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found: %s", orderNo)
	}

	// 3. บันทึกประวัติการยกเลิก
	queryCancelStatus := `
        INSERT INTO CancelStatus (
            RefID, CancelStatus, Remark, CancelBy, CancelDate
        ) VALUES (
            :OrderNo, 1, :Remark, :CancelBy, GETDATE()
        )
    `
	_, err = tx.NamedExecContext(ctx, queryCancelStatus, map[string]interface{}{
		"OrderNo":  orderNo,
		"Remark":   remark,
		"CancelBy": updateBy,
	})
	if err != nil {
		return fmt.Errorf("failed to insert cancel status: %w", err)
	}

	return nil
}
