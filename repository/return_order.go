package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
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
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Get hostname for the `CreateBy` field
	hostname, err := os.Hostname()
	if err != nil {
		//หากเกิดข้อผิดพลาดอะไรสักจุดที่วางไว้ให้ rollback to start point ทันที
		tx.Rollback()
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// Insert into ReturnOrder
	insertReturnOrderQuery := `
        INSERT INTO ReturnOrder (
            ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
            OptStatusID, AxStatusID, PlatfStatusID, Remark, CancelID, StatusCheckID, 
            CheckBy, Description, CreateBy, CreateDate
        ) VALUES (
            @ReturnID, @OrderNo, @SaleOrder, @SaleReturn, @TrackingNo, @PlatfID, @ChannelID, 
            @OptStatusID, @AxStatusID, @PlatfStatusID, @Remark, @CancelID, @StatusCheckID, 
            @CheckBy, @Description, @CreateBy, @CreateDate
        )
    `
	_, err = tx.ExecContext(ctx, insertReturnOrderQuery,
		sql.Named("ReturnID", req.ReturnID),
		sql.Named("OrderNo", req.OrderNo),
		sql.Named("SaleOrder", req.SaleOrder),
		sql.Named("SaleReturn", req.SaleReturn),
		sql.Named("TrackingNo", req.TrackingNo),
		sql.Named("PlatfID", req.PlatfID),
		sql.Named("ChannelID", req.ChannelID),
		sql.Named("OptStatusID", req.OptStatusID),
		sql.Named("AxStatusID", req.AxStatusID),
		sql.Named("PlatfStatusID", req.PlatfStatusID),
		sql.Named("Remark", req.Remark),
		sql.Named("CancelID", req.CancelID),
		sql.Named("StatusCheckID", req.StatusCheckID),
		sql.Named("CheckBy", req.CheckBy),
		sql.Named("Description", req.Description),
		sql.Named("CreateBy", hostname),
		sql.Named("CreateDate", time.Now()),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert into ReturnOrder: %w", err)
	}

	// Insert into ReturnOrderLine
	insertReturnOrderLineQuery := `
        INSERT INTO ReturnOrderLine (
            ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
            AlterSKU, CreateBy, CreateDate, UpdateBy, UpdateDate
        ) VALUES (
            @ReturnID, @OrderNo, @TrackingNo, @SKU, @ReturnQTY, @CheckQTY, @Price, 
            @AlterSKU, @CreateBy, @CreateDate, NULL, NULL
        )
    `
	for _, line := range req.ReturnOrderLine {
		log.Printf("Inserting ReturnOrderLine: %+v", line)
		_, err := tx.ExecContext(ctx, insertReturnOrderLineQuery,
			sql.Named("ReturnID", req.ReturnID),
			sql.Named("OrderNo", req.OrderNo),
			sql.Named("TrackingNo", req.TrackingNo),
			sql.Named("SKU", line.SKU),
			sql.Named("ReturnQTY", line.ReturnQTY),
			sql.Named("CheckQTY", line.CheckQTY),
			sql.Named("Price", line.Price),
			sql.Named("AlterSKU", line.AlterSKU),
			sql.Named("CreateBy", hostname),
			sql.Named("CreateDate", time.Now()),
		)
		if err != nil {
			log.Printf("Error inserting ReturnOrderLine: %v", err)
			tx.Rollback()
			return fmt.Errorf("failed to insert into ReturnOrderLine: %w", err)
		}
	}

	// Commit the transaction ยืนยันว่ามีการบันทึกข้อมูล
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

	hostname, err := os.Hostname()
	if err != nil {
		//หากเกิดข้อผิดพลาดอะไรสักจุดที่วางไว้ให้ rollback to start point ทันที
		tx.Rollback()
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
		tx.Rollback()
		return fmt.Errorf("failed to fetch current data: %w", err)
	}

	// ใช้ข้อมูลเก่าในกรณีที่ไม่มีการอัปเดต
	updateQuery := `
        UPDATE ReturnOrder
        SET 
            SaleReturn = COALESCE(@SaleReturn, SaleReturn),
            TrackingNo = COALESCE(@TrackingNo, TrackingNo),
            PlatfID = COALESCE(@PlatfID, PlatfID),
            ChannelID = COALESCE(@ChannelID, ChannelID),
            Remark = COALESCE(@Remark, Remark),
            UpdateBy = @UpdateBy,
            UpdateDate = @UpdateDate
        WHERE ReturnID = @ReturnID
    `
	_, err = tx.ExecContext(ctx, updateQuery,
		sql.Named("SaleReturn", req.SaleReturn),
		sql.Named("TrackingNo", req.TrackingNo),
		sql.Named("PlatfID", req.PlatfID),
		sql.Named("ChannelID", req.ChannelID),
		sql.Named("Remark", req.Remark),
		sql.Named("UpdateBy", hostname),
		sql.Named("UpdateDate", time.Now()),
		sql.Named("ReturnID", req.ReturnID),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update return order: %w", err)
	}

	// Commit Transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) DeleteReturnOrder(returnID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start a transaction
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Delete from ReturnOrderLine
	deleteOrderLineQuery := `
		DELETE FROM ReturnOrderLine
		WHERE ReturnID = @ReturnID
	`
	_, err = tx.ExecContext(ctx, deleteOrderLineQuery, sql.Named("ReturnID", returnID))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete from ReturnOrderLine: %w", err)
	}

	// Delete from ReturnOrder
	deleteOrderQuery := `
		DELETE FROM ReturnOrder
		WHERE ReturnID = @ReturnID
	`
	_, err = tx.ExecContext(ctx, deleteOrderQuery, sql.Named("ReturnID", returnID))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete from ReturnOrder: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}