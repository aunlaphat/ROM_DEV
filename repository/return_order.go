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
	CheckReturnIDExists(ctx context.Context, returnID string) (bool, error)
}

func (repo repositoryDB) AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	// Step 1: กำหนดคิวรี่ดึงข้อมูล ReturnOrder ทั้งหมด
	var orders []response.ReturnOrder
	queryOrder := `
		SELECT 
			ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Remark, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		ORDER BY OrderNo
	`
	// Step 1.2: ใช้ `SelectContext` ดึงข้อมูล ReturnOrder ทั้งหมดจากฐานข้อมูล
	err := repo.db.SelectContext(ctx, &orders, queryOrder)
	if err != nil {
		log.Printf("Error querying ReturnOrder: %v", err)
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	// Step 2: ดึงข้อมูล ReturnOrderLine ทั้งหมด
	lines, err := repo.GetAllReturnOrderLines(ctx)
	if err != nil {
		log.Printf("Error fetching ReturnOrderLines: %v", err)
		return nil, err
	}

	// Step 3: จับคู่ข้อมูล ReturnOrderLine กับ ReturnOrder ตาม ReturnID
	linesMap := make(map[string][]response.ReturnOrderLine)
	for _, line := range lines {
		linesMap[line.ReturnID] = append(linesMap[line.ReturnID], line)
	}

	// Step 4: เพิ่มข้อมูล ReturnOrderLine เข้าไปในแต่ละ ReturnOrder
	for i := range orders {
		orders[i].ReturnOrderLine = linesMap[orders[i].ReturnID]
	}

	return orders, nil
}

func (repo repositoryDB) GetReturnOrderByID(ctx context.Context, returnID string) (*response.ReturnOrder, error) {
	// Step 1: ดึงข้อมูล ReturnOrder จาก ReturnID
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

	// Step 2: ดึงข้อมูล ReturnOrderLine ที่เกี่ยวข้องกับ ReturnID
	lines, err := repo.GetReturnOrderLinesByReturnID(ctx, returnID)
	if err != nil {
		log.Printf("Error fetching ReturnOrderLines for ReturnID %s: %v", returnID, err)
		return nil, err
	}
	// Step 3: เพิ่มข้อมูล ReturnOrderLine เข้าไปใน ReturnOrder
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
    // ตรวจสอบว่า ReturnID มีอยู่จริงในฐานข้อมูล
    var exists bool
    checkQuery := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM ReturnOrder WHERE ReturnID = @ReturnID
        ) THEN 1 ELSE 0 END
    `
    err := repo.db.GetContext(ctx, &exists, checkQuery, sql.Named("ReturnID", returnID))
    if err != nil {
        return nil, fmt.Errorf("failed to check ReturnID existence: %w", err)
    }
    if !exists {
        return nil, sql.ErrNoRows // ReturnID ไม่พบ
    }

    // ดึงข้อมูล ReturnOrderLines
    var lines []response.ReturnOrderLine
    query := `
        SELECT RecID, ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, 
               CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
        FROM ReturnOrderLine
        WHERE ReturnID = @ReturnID
    `
    err = repo.db.SelectContext(ctx, &lines, query, sql.Named("ReturnID", returnID))
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, sql.ErrNoRows // ไม่มีข้อมูล ReturnOrderLine
        }
        return nil, fmt.Errorf("error querying ReturnOrderLine by ReturnID: %w", err)
    }

    return lines, nil
}


func (repo repositoryDB) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	// Step 1: เริ่มต้น Transaction
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 2: เพิ่มข้อมูล ReturnOrder ลงฐานข้อมูล
		insertReturnOrderQuery := `
            INSERT INTO ReturnOrder (
                ReturnID, OrderNo, SaleOrder, SaleReturn, TrackingNo, PlatfID, ChannelID, 
                OptStatusID, AxStatusID, PlatfStatusID, Remark, CancelID, StatusCheckID, 
                CheckBy, Description, CreateBy, CreateDate
            ) VALUES (
                :ReturnID, :OrderNo, :SaleOrder, :SaleReturn, :TrackingNo, :PlatfID, :ChannelID, 
                :OptStatusID, :AxStatusID, :PlatfStatusID, :Remark, :CancelID, :StatusCheckID, 
                :CheckBy, :Description, 'USER', SYSDATETIME()
            )
        `
		// Step 2.1: ตรวจสอบว่ามีค่าสำหรับการใส่ใน Query
		params := map[string]interface{}{
			"ReturnID":       req.ReturnID,
			"OrderNo":        req.OrderNo,
			"SaleOrder":      req.SaleOrder,
			"SaleReturn":     req.SaleReturn,
			"TrackingNo":     req.TrackingNo,
			"PlatfID":        req.PlatfID,
			"ChannelID":      req.ChannelID,
			"OptStatusID":    req.OptStatusID,
			"AxStatusID":     req.AxStatusID,
			"PlatfStatusID":  req.PlatfStatusID,
			"Remark":         req.Remark,
			"CancelID":       req.CancelID,
			"StatusCheckID":  req.StatusCheckID,
			"CheckBy":        req.CheckBy,
			"Description":    req.Description,
		}
		if _, err := tx.NamedExecContext(ctx, insertReturnOrderQuery, params); err != nil {
			log.Printf("Error inserting ReturnOrder: %v", err)
			return fmt.Errorf("failed to insert ReturnOrder: %w", err)
		}

		// Step 3: เพิ่มข้อมูล ReturnOrderLine
		insertReturnOrderLineQuery := `
            INSERT INTO ReturnOrderLine (
                ReturnID, OrderNo, TrackingNo, SKU, ReturnQTY, CheckQTY, Price, CreateBy, CreateDate
            ) VALUES (
                :ReturnID, :OrderNo, :TrackingNo, :SKU, :ReturnQTY, :CheckQTY, :Price, 
                 'USER', SYSDATETIME()
            )
        `
		// Step 3.1: Loop ผ่าน `ReturnOrderLine` และใส่ข้อมูลลงใน Query
		for _, line := range req.ReturnOrderLine {
			line.ReturnID = req.ReturnID
			line.OrderNo = req.OrderNo
			line.TrackingNo = req.TrackingNo

			params := map[string]interface{}{
				"ReturnID":   line.ReturnID,
				"OrderNo":    line.OrderNo,
				"TrackingNo": line.TrackingNo,
				"SKU":        line.SKU,
				"ReturnQTY":  line.ReturnQTY,
				"CheckQTY":   line.CheckQTY,
				"Price":      line.Price,
				//"AlterSKU":   line.AlterSKU,
			}
			if _, err := tx.NamedExecContext(ctx, insertReturnOrderLineQuery, params); err != nil {
				log.Printf("Error inserting ReturnOrderLine: %v", err)
				return fmt.Errorf("failed to insert ReturnOrderLine: %w", err)
			}
		}

		// Step 4: Commit Transaction
		return nil
	})
}

// อัปเดต ReturnOrder และ ReturnOrderLine
func (repo repositoryDB) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error {
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		if req.TrackingNo != nil {
			
		}
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
						THEN 'USER'
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
		params := map[string]interface{}{
			"ReturnID":       req.ReturnID,
			"SaleReturn":     req.SaleReturn,
			"TrackingNo":     req.TrackingNo,
			"PlatfID":        req.PlatfID,
			"ChannelID":      req.ChannelID,
			"OptStatusID":    req.OptStatusID,
			"AxStatusID":     req.AxStatusID,
			"PlatfStatusID":  req.PlatfStatusID,
			"Remark":         req.Remark,
			"CancelID":       req.CancelID,
			"StatusCheckID":  req.StatusCheckID,
			"CheckBy":        req.CheckBy,
			"Description":    req.Description,
		}
		if _, err := tx.NamedExecContext(ctx, updateReturnOrderQuery, params); err != nil {
			log.Printf("Error updating ReturnOrder: %v", err)
			return fmt.Errorf("failed to update ReturnOrder: %w", err)
		}

		if req.TrackingNo != nil {
			updateReturnOrderLineQuery := `
				UPDATE ReturnOrderLine
				SET 
					TrackingNo = :TrackingNo, 
					UpdateBy = 'USER', 
					UpdateDate = SYSDATETIME()
				WHERE ReturnID = :ReturnID
			`
			params := map[string]interface{}{
				"ReturnID":    req.ReturnID,
				"TrackingNo":  req.TrackingNo,
			}
			if _, err := tx.NamedExecContext(ctx, updateReturnOrderLineQuery, params); err != nil {
				log.Printf("Error updating ReturnOrderLine: %v", err)
				return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
			}
		}

		return nil
	})
}

// ลบ ReturnOrder และ ReturnOrderLine
func (repo repositoryDB) DeleteReturnOrder(ctx context.Context, returnID string) error {
	// Step 1: เริ่ม Transaction เพื่อควบคุมการลบ
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 2: ลบข้อมูลใน ReturnOrderLine โดยอิงจาก ReturnID
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

		// Step 3: ลบข้อมูลใน ReturnOrder
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

		// Step 4: ส่ง Commit หากการลบสำเร็จทั้งหมด
		return nil
	})
}

// CheckReturnIDExists - เพิ่มการตรวจสอบ ReturnID ว่ามีอยู่ในฐานข้อมูล
func (repo repositoryDB) CheckReturnIDExists(ctx context.Context, returnID string) (bool, error) {
    var exists bool
    query := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM ReturnOrder WHERE ReturnID = @ReturnID
        ) THEN 1 ELSE 0 END
    `
    err := repo.db.GetContext(ctx, &exists, query, sql.Named("ReturnID", returnID))
    if err != nil {
        return false, fmt.Errorf("failed to check ReturnID existence: %w", err)
    }

    return exists, nil
}