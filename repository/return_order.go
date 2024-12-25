package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	// "boilerplate-backend-go/repository/transaction"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// defines methods for CRUD operations on orders
type ReturnOrderRepository interface {
	AllGetReturnOrder() ([]response.ReturnOrder, error)
	GetReturnOrderByID(returnID string) (*response.ReturnOrder, error)
	CreateReturnOrder(req request.CreateReturnOrder) error
	UpdateReturnOrder(req request.UpdateReturnOrder) error
	DeleteReturnOrder(returnID string) error
}

func (repo repositoryDB) AllGetReturnOrder() ([]response.ReturnOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ดึงข้อมูล ReturnOrder ทั้งหมด
	var orders []response.ReturnOrder
	queryOrder := `
		SELECT 
			ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Remark, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		ORDER BY OrderNo
	`
	err := repo.db.SelectContext(ctx, &orders, queryOrder)
	if err != nil {
		log.Println("Error querying ReturnOrder:", err)
		return nil, err
	}

	// ดึงข้อมูล ReturnOrderLine ทั้งหมด
	queryOrderLines := `
		SELECT 
			ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, CreateBy, 
			CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
	`
	var orderLines []response.ReturnOrderLine
	err = repo.db.SelectContext(ctx, &orderLines, queryOrderLines)
	if err != nil {
		log.Println("Error querying ReturnOrderLine:", err)
		return nil, err
	}

	// Map ReturnOrderLine ไปยัง ReturnOrder ตาม ReturnID
	orderLinesMap := make(map[string][]response.ReturnOrderLine)
	for _, line := range orderLines {
		orderLinesMap[line.ReturnID] = append(orderLinesMap[line.ReturnID], line)
	}
	for i := range orders {
		orders[i].ReturnOrderLine = orderLinesMap[orders[i].ReturnID]
	}

	return orders, nil
}

func (repo repositoryDB) GetReturnOrderByID(returnID string) (*response.ReturnOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ดึงข้อมูล ReturnOrder
	var order response.ReturnOrder
	queryOrder := `
		SELECT 
			ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Remark, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		WHERE ReturnID = @returnID
	`
	err := repo.db.GetContext(ctx, &order, queryOrder, sql.Named("returnID", returnID))
	if err != nil {
		log.Println("Error querying ReturnOrder:", err)
		return nil, err
	}

	// ดึงข้อมูล ReturnOrderLine ที่เกี่ยวข้องกับ ReturnID
	var orderLines []response.ReturnOrderLine
	queryOrderLines := `
		SELECT 
			ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, CreateBy, 
			CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
		WHERE ReturnID = @returnID
	`
	err = repo.db.SelectContext(ctx, &orderLines, queryOrderLines, sql.Named("returnID", returnID))
	if err != nil {
		log.Println("Error querying ReturnOrderLine:", err)
		return nil, err
	}

	order.ReturnOrderLine = orderLines
	return &order, nil
}

func (repo repositoryDB) CreateReturnOrder(req request.CreateReturnOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer handleTransaction(tx)() //หากเกิดข้อผิดพลาดอะไรสักจุดที่วางไว้ให้ rollback to start point ทันที

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Prepare data for `ReturnOrder`
	returnOrderData := map[string]interface{}{
		"ReturnID":      req.ReturnID,
		"OrderNo":       req.OrderNo,
		"SaleOrder":     req.SaleOrder,
		"SaleReturn":    req.SaleReturn,
		"TrackingNo":    req.TrackingNo,
		"PlatfID":       req.PlatfID,
		"ChannelID":     req.ChannelID,
		"OptStatusID":   req.OptStatusID,
		"AxStatusID":    req.AxStatusID,
		"PlatfStatusID": req.PlatfStatusID,
		"Remark":        req.Remark,
		"CancelID":      req.CancelID,
		"StatusCheckID": req.StatusCheckID,
		"CheckBy":       req.CheckBy,
		"Description":   req.Description,
		"CreateBy":      hostname,
		"CreateDate":    time.Now(),
	}

	// Insert `ReturnOrder` with NamedExec
	insertReturnOrderQuery := `
		INSERT INTO ReturnOrder (
			ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Remark, CancelID, StatusCheckID, 
			CheckBy, Description, CreateBy, CreateDate
		) VALUES (
			:ReturnID, :OrderNo, :SaleOrder, :SaleReturn, :TrackingNo, :PlatfID, :ChannelID, 
			:OptStatusID, :AxStatusID, :PlatfStatusID, :Remark, :CancelID, :StatusCheckID, 
			:CheckBy, :Description, :CreateBy, :CreateDate
		)
	`
	_, err = tx.NamedExecContext(ctx, insertReturnOrderQuery, returnOrderData)
	if err != nil {
		return fmt.Errorf("failed to insert into ReturnOrder: %w", err)
	}

	// Insert `ReturnOrderLine` for each line item
	insertReturnOrderLineQuery := `
		INSERT INTO ReturnOrderLine (
			ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
			AlterSKU, CreateBy, CreateDate, UpdateBy, UpdateDate
		) VALUES (
			:ReturnID, :OrderNo, :TrackingNo, :SKU, :ReturnQTY, :CheckQTY, :Price, 
			:AlterSKU, :CreateBy, :CreateDate, NULL, NULL
		)
	`
	for _, line := range req.ReturnOrderLine {
		lineData := map[string]interface{}{
			"ReturnID":    req.ReturnID,
			"OrderNo":     req.OrderNo,
			"TrackingNo":  req.TrackingNo,
			"SKU":         line.SKU,
			"ReturnQTY":   line.ReturnQTY,
			"CheckQTY":    line.CheckQTY,
			"Price":       line.Price,
			"AlterSKU":    line.AlterSKU,
			"CreateBy":    hostname,
			"CreateDate":  time.Now(),
		}

		_, err := tx.NamedExecContext(ctx, insertReturnOrderLineQuery, lineData)
		if err != nil {
			return fmt.Errorf("failed to insert into ReturnOrderLine: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) UpdateReturnOrder(req request.UpdateReturnOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// เริ่ม Transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer handleTransaction(tx)()

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// ดึงข้อมูลเก่า
	var currentData response.ReturnOrder
	selectQuery := `
        SELECT SaleReturn, TrackingNo, PlatfID, ChannelID, Remark, UpdateBy, UpdateDate
        FROM ReturnOrder
        WHERE ReturnID = @ReturnID
    `
	err = tx.GetContext(ctx, &currentData, selectQuery, sql.Named("ReturnID", req.ReturnID))
	if err != nil {
		return fmt.Errorf("failed to fetch current data: %w", err)
	}

	// เตรียมข้อมูลสำหรับการอัปเดต
	updateData := map[string]interface{}{
		"SaleReturn":  req.SaleReturn,
		"TrackingNo":  req.TrackingNo,
		"PlatfID":     req.PlatfID,
		"ChannelID":   req.ChannelID,
		"Remark":      req.Remark,
		"UpdateBy":    hostname,
		"UpdateDate":  time.Now(),
		"ReturnID":    req.ReturnID,
	}

	// ใช้ข้อมูลเก่าในกรณีที่ไม่มีการอัปเดต
	updateQuery := `
		UPDATE ReturnOrder
		SET 
			SaleReturn = COALESCE(:SaleReturn, SaleReturn),
			TrackingNo = COALESCE(:TrackingNo, TrackingNo),
			PlatfID = COALESCE(:PlatfID, PlatfID),
			ChannelID = COALESCE(:ChannelID, ChannelID),
			Remark = COALESCE(:Remark, Remark),
			UpdateBy = CASE
                  WHEN :SaleReturn IS NOT NULL OR :TrackingNo IS NOT NULL OR :PlatfID IS NOT NULL 
                       OR :ChannelID IS NOT NULL OR :Remark IS NOT NULL 
                  THEN :UpdateBy 
                  ELSE UpdateBy
               END,
			UpdateDate = CASE
                    WHEN :SaleReturn IS NOT NULL OR :TrackingNo IS NOT NULL OR :PlatfID IS NOT NULL 
                         OR :ChannelID IS NOT NULL OR :Remark IS NOT NULL 
                    THEN :UpdateDate 
                    ELSE UpdateDate
                 END
		WHERE ReturnID = :ReturnID
	`
	//ใช้เงื่อนไข CASE ใน SQL เพื่อกำหนดว่า UpdateDate จะอัปเดตเมื่อมีการเปลี่ยนแปลงข้อมูลเท่านั้น
	_, err = tx.NamedExecContext(ctx, updateQuery, updateData)
	if err != nil {
		log.Printf("Error updating ReturnOrder: %v", err)
		return fmt.Errorf("failed to update ReturnOrder: %w", err)
	}

	// อัปเดต ReturnOrderLine
	updateReturnOrderLineQuery := `
		UPDATE ReturnOrderLine
		SET 
			TrackingNo = COALESCE(:TrackingNo, TrackingNo)
			UpdateBy = :UpdateBy,
    		UpdateDate = :UpdateDate
		WHERE ReturnID = :ReturnID
		AND (TrackingNo IS DISTINCT FROM :TrackingNo);
		`
	_, err = tx.NamedExecContext(ctx, updateReturnOrderLineQuery, req)
	if err != nil {
		return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) DeleteReturnOrder(returnID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer handleTransaction(tx)()

	// ลบข้อมูลจากตาราง ReturnOrderLine ก่อน
	deleteReturnOrderLineQuery := `
	DELETE FROM ReturnOrderLine
	WHERE ReturnID = :ReturnID
	`
	_, err = tx.NamedExecContext(ctx, deleteReturnOrderLineQuery, map[string]interface{}{
	"ReturnID": returnID,
	})
	if err != nil {
	log.Printf("Error deleting ReturnOrderLine: %v", err)
	return fmt.Errorf("failed to delete from ReturnOrderLine: %w", err)
	}

	// ลบข้อมูลจากตาราง ReturnOrder
	deleteReturnOrderQuery := `
	DELETE FROM ReturnOrder
	WHERE ReturnID = :ReturnID
	`
	_, err = tx.NamedExecContext(ctx, deleteReturnOrderQuery, map[string]interface{}{
	"ReturnID": returnID,
	})
	if err != nil {
	log.Printf("Error deleting ReturnOrder: %v", err)
	return fmt.Errorf("failed to delete from ReturnOrder: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
	log.Printf("Error committing transaction: %v", err)
	return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}