package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	//"boilerplate-backend-go/errors"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type ReturnOrderRepository interface {
	AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error
	DeleteReturnOrder(ctx context.Context, returnID string) error
}

// ดึงข้อมูล ReturnOrder ทั้งหมดจากฐานข้อมูล
func (repo repositoryDB) AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	
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
		log.Printf("Error querying ReturnOrder: %v", err)
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	lines, err := repo.GetAllReturnOrderLines(ctx)
	if err != nil {
		log.Printf("Error fetching ReturnOrderLines: %v", err)
		return nil, err
	}

	linesMap := make(map[string][]response.ReturnOrderLine)
	for _, line := range lines {
		linesMap[line.ReturnID] = append(linesMap[line.ReturnID], line)
	}

	for i := range orders {
		orders[i].ReturnOrderLine = linesMap[orders[i].ReturnID]
	}

	return orders, nil
}

// ดึงข้อมูล ReturnOrder จาก ReturnID
func (repo repositoryDB) GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error) {
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
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows 
		}
		log.Printf("Database error querying ReturnOrder by ID: %v", err)
		return nil, fmt.Errorf("unexpected database error: %w", err)
	}

	lines, err := repo.GetReturnOrderLinesByReturnID(ctx, returnID)
	if err != nil {
		log.Printf("Error fetching ReturnOrderLines for ReturnID %s: %v", returnID, err)
		return nil, err
	}

	order.ReturnOrderLine = lines
	return &order, nil
}

// Get All ReturnOrderLines
func (repo repositoryDB) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	var lines []response.ReturnOrderLine
	query := `
		SELECT RecID, ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
		       CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
		ORDER BY RecID
	`

	err := repo.db.SelectContext(ctx, &lines, query)
	if err != nil {
		log.Printf("Error querying ReturnOrderLine: %v", err)
		return nil, fmt.Errorf("get all return order lines error: %w", err)
	}

	return lines, nil
}

// Get ReturnOrderLines by ReturnID
func (repo repositoryDB) GetReturnOrderLinesByReturnID(ctx context.Context, returnID string) ([]response.ReturnOrderLine, error) {
	var lines []response.ReturnOrderLine
	query := `
		SELECT RecID, ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
		       CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
		WHERE ReturnID = @ReturnID
	`

	err := repo.db.SelectContext(ctx, &lines, query, sql.Named("ReturnID", returnID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("Error querying ReturnOrderLine by ReturnID: %v", err)
		return nil, fmt.Errorf("get return order lines by ReturnID error: %w", err)
	}

	return lines, nil
}

// สร้าง ReturnOrder พร้อม ReturnOrderLine ที่ ReturnID เดียวกัน
func (repo repositoryDB) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	// Start a transaction
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		insertReturnOrderQuery := `
            INSERT INTO ReturnOrder (
                ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
                OptStatusID, AxStatusID, PlatfStatusID, Remark, CancelID, StatusCheckID, 
                CheckBy, Description, CreateBy, CreateDate
            ) VALUES (
                :ReturnID, :OrderNo, :SaleOrder, :SaleReturn, :TrackingNo, :PlatfID, :ChannelID, 
                :OptStatusID, :AxStatusID, :PlatfStatusID, :Remark, :CancelID, :StatusCheckID, 
                :CheckBy, :Description, 'user', SYSDATETIME()
            )
        `
		if _, err := tx.NamedExecContext(ctx, insertReturnOrderQuery, req); err != nil {
			log.Printf("Error inserting ReturnOrder: %v", err)
			return fmt.Errorf("failed to insert ReturnOrder: %w", err)
		}

		// Insert `ReturnOrderLine`
		insertReturnOrderLineQuery := `
            INSERT INTO ReturnOrderLine (
                ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
                AlterSKU, CreateBy, CreateDate
            ) VALUES (
                :ReturnID, :OrderNo, :TrackingNo, :SKU, :ReturnQTY, :CheckQTY, :Price, 
                :AlterSKU, 'user', SYSDATETIME()
            )
        `
		for _, line := range req.ReturnOrderLine {
			line.ReturnID = req.ReturnID
			line.OrderNo = req.OrderNo
			line.TrackingNo = req.TrackingNo

			if _, err := tx.NamedExecContext(ctx, insertReturnOrderLineQuery, line); err != nil {
				log.Printf("Error inserting ReturnOrderLine: %v", err)
				return fmt.Errorf("failed to insert ReturnOrderLine: %w", err)
			}
		}

		return nil
	})
}

// อัปเดต ReturnOrder และ ReturnOrderLine
func (repo repositoryDB) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error {
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		updateReturnOrderQuery := `
			   UPDATE ReturnOrder
				SET 
					SaleReturn = COALESCE(:SaleReturn, SaleReturn),
					TrackingNo = COALESCE(:TrackingNo, TrackingNo),
					PlatfID = COALESCE(:PlatfID, PlatfID),
					ChannelID = COALESCE(:ChannelID, ChannelID),
					OptStatusID = COALESCE(:OptStatusID, OptStatusID),
					AxStatusID = COALESCE(:AxStatusID, AxStatusID),
					PlatfStatusID = COALESCE(:PlatfStatusID, PlatfStatusID),
					Remark = COALESCE(:Remark, Remark),
					CancelID = COALESCE(:CancelID, CancelID),
					StatusCheckID = COALESCE(:StatusCheckID, StatusCheckID),
					CheckBy = COALESCE(:CheckBy, CheckBy),
					Description = COALESCE(:Description, Description),
					UpdateBy = CASE 
						WHEN :SaleReturn IS NOT NULL OR :TrackingNo IS NOT NULL OR :PlatfID IS NOT NULL 
							OR :ChannelID IS NOT NULL OR :Remark IS NOT NULL OR :OptStatusID IS NOT NULL
							OR :AxStatusID IS NOT NULL OR :PlatfStatusID IS NOT NULL OR :CancelID IS NOT NULL 
							OR :StatusCheckID IS NOT NULL OR :CheckBy IS NOT NULL OR :Description IS NOT NULL
						THEN 'user'
						ELSE UpdateBy
					END,
					UpdateDate = CASE 
						WHEN :SaleReturn IS NOT NULL OR :TrackingNo IS NOT NULL OR :PlatfID IS NOT NULL 
							OR :ChannelID IS NOT NULL OR :Remark IS NOT NULL OR :OptStatusID IS NOT NULL
							OR :AxStatusID IS NOT NULL OR :PlatfStatusID IS NOT NULL OR :CancelID IS NOT NULL 
							OR :StatusCheckID IS NOT NULL OR :CheckBy IS NOT NULL OR :Description IS NOT NULL
						THEN SYSDATETIME()
						ELSE UpdateDate
					END
				WHERE ReturnID = :ReturnID
		`
		if _, err := tx.NamedExecContext(ctx, updateReturnOrderQuery, req); err != nil {
			log.Printf("Error updating ReturnOrder: %v", err)
			return fmt.Errorf("failed to update ReturnOrder: %w", err)
		}

		if req.TrackingNo != nil {
			updateReturnOrderLineQuery := `
				UPDATE ReturnOrderLine
				SET 
					TrackingNo = :TrackingNo, 
					UpdateBy = 'user', 
					UpdateDate = SYSDATETIME()
				WHERE ReturnID = :ReturnID
			`
			if _, err := tx.NamedExecContext(ctx, updateReturnOrderLineQuery, req); err != nil {
				log.Printf("Error updating ReturnOrderLine: %v", err)
				return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
			}
		}

		return nil
	})
}

// ลบ ReturnOrder และ ReturnOrderLine
func (repo repositoryDB) DeleteReturnOrder(ctx context.Context, returnID string) error {
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		deleteReturnOrderLineQuery := `
			DELETE FROM ReturnOrderLine
			WHERE ReturnID = :ReturnID
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderLineQuery, map[string]interface{}{
			"ReturnID": returnID,
		}); err != nil {
			log.Printf("Error deleting ReturnOrderLine for ReturnID %s: %v", returnID, err)
			return fmt.Errorf("failed to delete ReturnOrderLine: %w", err)
		}

		deleteReturnOrderQuery := `
			DELETE FROM ReturnOrder
			WHERE ReturnID = :ReturnID
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderQuery, map[string]interface{}{
			"ReturnID": returnID,
		}); err != nil {
			log.Printf("Error deleting ReturnOrder for ReturnID %s: %v", returnID, err)
			return fmt.Errorf("failed to delete ReturnOrder: %w", err)
		}

		return nil
	})
}
