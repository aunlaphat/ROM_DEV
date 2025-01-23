package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"strings"

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
		// Step 1: ดึงค่าปัจจุบันจากฐานข้อมูล
		var current response.ReturnOrder
		queryCurrent := `
            SELECT SaleReturn, TrackingNo, PlatfID, ChannelID, OptStatusID, AxStatusID, 
                   PlatfStatusID, Remark, CancelID, StatusCheckID, CheckBy, Description
            FROM ReturnOrder
            WHERE ReturnID = :ReturnID
        `
		currentParams := map[string]interface{}{"ReturnID": req.ReturnID}

		// ใช้ NamedQuery และ Rebind
		namedQuery, args, err := sqlx.Named(queryCurrent, currentParams)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		query := tx.Rebind(namedQuery)

		if err := tx.GetContext(ctx, &current, query, args...); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("ReturnID not found: %w", err)
			}
			return fmt.Errorf("failed to fetch current data: %w", err)
		}

		// Step 2: ตรวจสอบค่าที่เปลี่ยนแปลงและสร้าง SQL Update เฉพาะส่วนที่เปลี่ยนแปลง
		updateFields := []string{}
		params := map[string]interface{}{
			"ReturnID": req.ReturnID,
		}

		if req.SaleReturn != nil && (current.SaleReturn == nil || *req.SaleReturn != *current.SaleReturn) {
			updateFields = append(updateFields, "SaleReturn = :SaleReturn")
			params["SaleReturn"] = req.SaleReturn
		}
		if req.TrackingNo != nil && (current.TrackingNo == nil || *req.TrackingNo != *current.TrackingNo) {
			updateFields = append(updateFields, "TrackingNo = :TrackingNo")
			params["TrackingNo"] = req.TrackingNo
		}
		if req.PlatfID != nil && (current.PlatfID == nil || *req.PlatfID != *current.PlatfID) {
			updateFields = append(updateFields, "PlatfID = :PlatfID")
			params["PlatfID"] = req.PlatfID
		}
		if req.ChannelID != nil && (current.ChannelID == nil || *req.ChannelID != *current.ChannelID) {
			updateFields = append(updateFields, "ChannelID = :ChannelID")
			params["ChannelID"] = req.ChannelID
		}
		if req.OptStatusID != nil && (current.OptStatusID == nil || *req.OptStatusID != *current.OptStatusID) {
			updateFields = append(updateFields, "OptStatusID = :OptStatusID")
			params["OptStatusID"] = req.OptStatusID
		}
		if req.AxStatusID != nil && (current.AxStatusID == nil || *req.AxStatusID != *current.AxStatusID) {
			updateFields = append(updateFields, "AxStatusID = :AxStatusID")
			params["AxStatusID"] = req.AxStatusID
		}
		if req.PlatfStatusID != nil && (current.PlatfStatusID == nil || *req.PlatfStatusID != *current.PlatfStatusID) {
			updateFields = append(updateFields, "PlatfStatusID = :PlatfStatusID")
			params["PlatfStatusID"] = req.PlatfStatusID
		}
		if req.Remark != nil && (current.Remark == nil || *req.Remark != *current.Remark) {
			updateFields = append(updateFields, "Remark = :Remark")
			params["Remark"] = req.Remark
		}
		if req.CancelID != nil && (current.CancelID == nil || *req.CancelID != *current.CancelID) {
			updateFields = append(updateFields, "CancelID = :CancelID")
			params["CancelID"] = req.CancelID
		}
		if req.StatusCheckID != nil && (current.StatusCheckID == nil || *req.StatusCheckID != *current.StatusCheckID) {
			updateFields = append(updateFields, "StatusCheckID = :StatusCheckID")
			params["StatusCheckID"] = req.StatusCheckID
		}
		if req.CheckBy != nil && (current.CheckBy == nil || *req.CheckBy != *current.CheckBy) {
			updateFields = append(updateFields, "CheckBy = :CheckBy")
			params["CheckBy"] = req.CheckBy
		}
		if req.Description != nil && (current.Description == nil || *req.Description != *current.Description) {
			updateFields = append(updateFields, "Description = :Description")
			params["Description"] = req.Description
		}

		// หากไม่มีการเปลี่ยนแปลง ออกจากฟังก์ชัน
		if len(updateFields) == 0 {
			return nil
		}

		// Step 3: เพิ่ม UpdateBy และ UpdateDate ใน SQL Query
		updateFields = append(updateFields, "UpdateBy = 'USER'", "UpdateDate = SYSDATETIME()")
		updateQuery := fmt.Sprintf(`
            UPDATE ReturnOrder
            SET %s
            WHERE ReturnID = :ReturnID
        `, strings.Join(updateFields, ", "))

		// ใช้ NamedQuery และ Rebind
		namedUpdateQuery, updateArgs, err := sqlx.Named(updateQuery, params)
		if err != nil {
			return fmt.Errorf("failed to prepare update query: %w", err)
		}
		updateQueryRebind := tx.Rebind(namedUpdateQuery)

		// Step 4: ดำเนินการอัปเดตใน ReturnOrder
		if _, err := tx.ExecContext(ctx, updateQueryRebind, updateArgs...); err != nil {
			log.Printf("Error updating ReturnOrder: %v", err)
			return fmt.Errorf("failed to update ReturnOrder: %w", err)
		}

		// Step 5: อัปเดต TrackingNo ใน ReturnOrderLine หากเปลี่ยนแปลง
		if req.TrackingNo != nil && (current.TrackingNo == nil || *req.TrackingNo != *current.TrackingNo) {
			updateLineQuery := `
                UPDATE ReturnOrderLine
                SET 
                    TrackingNo = :TrackingNo, 
                    UpdateBy = 'USER', 
                    UpdateDate = SYSDATETIME()
                WHERE ReturnID = :ReturnID
            `
			namedLineQuery, lineArgs, err := sqlx.Named(updateLineQuery, params)
			if err != nil {
				return fmt.Errorf("failed to prepare line update query: %w", err)
			}
			lineQueryRebind := tx.Rebind(namedLineQuery)

			if _, err := tx.ExecContext(ctx, lineQueryRebind, lineArgs...); err != nil {
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
