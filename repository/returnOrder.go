package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"boilerplate-backend-go/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReturnOrderRepository interface {
	GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder, updateBy string) error
	DeleteReturnOrder(ctx context.Context, orderNo string) error
	CheckOrderNoExist(ctx context.Context, orderNo string) (bool, error)
	CheckOrderNoLineExist(ctx context.Context, orderNo string) (bool, error)


	GetCreateReturnOrder(ctx context.Context, orderNo string) (*response.CreateReturnOrder, error)
	GetUpdateReturnOrder(ctx context.Context, orderNo string) (*response.UpdateReturnOrder, error)
	GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error)
	GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error)
}
// review
func (repo repositoryDB) GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	var orders []response.ReturnOrder

	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		ORDER BY RecID
	`
	// ใช้ `SelectContext` ดึงข้อมูล ReturnOrder ทั้งหมดจากฐานข้อมูล
	err := repo.db.SelectContext(ctx, &orders, queryOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	// ดึงข้อมูล ReturnOrderLine จับคู่กับ OrderNo ที่ตรงกัน
	for i, order := range orders {
		lines, err := repo.GetReturnOrderLineByOrderNo(ctx, order.OrderNo)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch return order lines for order %s: %w", order.OrderNo, err)
		}
		orders[i].ReturnOrderLine = lines
	}

	return orders, nil
}
// review
func (repo repositoryDB) GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	var order response.ReturnOrder

	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		WHERE OrderNo = @orderNo
        ORDER BY RecID
	`
	err := repo.db.GetContext(ctx, &order, queryOrder, sql.Named("orderNo", orderNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo %s: %w", orderNo, err)
	}

	// ดึงข้อมูล ReturnOrderLine โดย OrderNo
	lines, err := repo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("error fetching ReturnOrderLines for OrderNo %s: %w", orderNo, err)
	}
	// เพิ่มข้อมูล ReturnOrderLine เข้าไปใน ReturnOrder
	order.ReturnOrderLine = lines

	return &order, nil
}
// review
// Get All ReturnOrderLines
func (repo repositoryDB) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	var lines []response.ReturnOrderLine

	query := `
		SELECT OrderNo, SKU, QTY, ReturnQTY, ActualQTY, Price, TrackingNo, 
               CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
        ORDER BY RecID
	`

	err := repo.db.SelectContext(ctx, &lines, query)
	if err != nil {
		return nil, fmt.Errorf("get all return order lines error: %w", err)
	}

	return lines, nil
}
// review
// Get ReturnOrderLines by OrderNo
func (repo repositoryDB) GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	// ดึงข้อมูล ReturnOrderLines
	var lines []response.ReturnOrderLine

	query := `
        SELECT OrderNo, SKU, QTY, ReturnQTY, ActualQTY, Price, TrackingNo, 
               CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
        FROM ReturnOrderLine
        WHERE OrderNo = :OrderNo
        ORDER BY RecID
    `
	params := map[string]interface{}{"OrderNo": orderNo}

	namedQuery, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	query = repo.db.Rebind(namedQuery)

	if err := repo.db.SelectContext(ctx, &lines, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch return order lines: %w", err)
	}

	return lines, nil
}
// review
func (repo repositoryDB) GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error) {
	var orders []response.DraftTradeDetail

	query := `
        SELECT 
            OrderNo, SoNo, SrNo, TrackingNo, ChannelID, Reason, StatusCheckID, CreateBy, CreateDate
        FROM ReturnOrder
        WHERE StatusCheckID = :statusCheckID
        ORDER BY CreateDate ASC
    `
	params := map[string]interface{}{"statusCheckID": statusCheckID}

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to bind named parameters: %w", err)
	}
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	return orders, nil
}
// review
func (repo repositoryDB) GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error) {
	var orders []response.DraftTradeDetail

	// แปลง startDate และ endDate ให้เป็น time.Time ส่งไป params
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startDate: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endDate: %w", err)
	}

	// ใช้เวลา 23:59:59 ของ endDate
	end = end.Add(24 * time.Hour) 

	query := `
        SELECT 
            OrderNo, SoNo, SrNo, TrackingNo, ChannelID, Reason, StatusCheckID, CreateBy, CreateDate
        FROM ReturnOrder
        WHERE StatusCheckID = :StatusCheckID 
        AND CreateDate >= :StartDate
        AND CreateDate < :EndDate
        ORDER BY CreateDate ASC
    `

	params := map[string]interface{}{
		"StatusCheckID": statusCheckID,
		"StartDate":     start.Format("2006-01-02"), // ส่งแค่วันที่เริ่ม
		"EndDate":       end.Format("2006-01-02"),   // ส่งแค่วันที่สิ้นสุด
	}

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to bind named parameters: %w", err)
	}
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	log.Printf("Fetched %d return orders", len(orders))
	return orders, nil
}
// review
func (repo repositoryDB) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	// เริ่มต้น Transaction
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// เพิ่มข้อมูล ReturnOrder ลงฐานข้อมูล
		insertReturnOrderQuery := `
            INSERT INTO ReturnOrder (
                OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
                OptStatusID, AxStatusID, PlatfStatusID, Reason, StatusCheckID, 
                Description, CreateBy, CreateDate
            ) VALUES (
                :OrderNo, :SoNo, :SrNo, :TrackingNo, :PlatfID, :ChannelID, 
                :OptStatusID, :AxStatusID, :PlatfStatusID, :Reason, :StatusCheckID, 
                :Description, :CreateBy, GETDATE()
            )
        `
		// ตรวจสอบว่ามีค่าตรงกันกับที่จะนำไป Query
		params := map[string]interface{}{
			"OrderNo":       req.OrderNo,
			"SoNo":          req.SoNo,
			"SrNo":          req.SrNo,
			"TrackingNo":    req.TrackingNo,
			"PlatfID":       req.PlatfID,
			"ChannelID":     req.ChannelID,
			"OptStatusID":   req.OptStatusID,
			"AxStatusID":    req.AxStatusID,
			"PlatfStatusID": req.PlatfStatusID,
			"Reason":        req.Reason,
			"StatusCheckID": req.StatusCheckID,
			"Description":   req.Description,
			"CreateBy":      req.CreateBy,
		}
		if _, err := tx.NamedExecContext(ctx, insertReturnOrderQuery, params); err != nil {
			return fmt.Errorf("failed to insert ReturnOrder: %w", err)
		}

		// เพิ่มข้อมูล ReturnOrderLine
		insertReturnOrderLineQuery := `
            INSERT INTO ReturnOrderLine (
                OrderNo, SKU, QTY, ReturnQTY, Price, TrackingNo, CreateBy, CreateDate
            ) VALUES (
                :OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :TrackingNo, :CreateBy, GETDATE()
            )
        `
		// Loop ผ่าน `ReturnOrderLine` และใส่ข้อมูลลงใน Query
		for _, line := range req.ReturnOrderLine {
			line.OrderNo = req.OrderNo
			line.TrackingNo = req.TrackingNo

			params := map[string]interface{}{
				"OrderNo":    line.OrderNo,
				"TrackingNo": line.TrackingNo,
				"SKU":        line.SKU,
				"QTY":        line.QTY,
				"ReturnQTY":  line.ReturnQTY,
				"Price":      line.Price,
				"CreateBy":   req.CreateBy,
			}
			if _, err := tx.NamedExecContext(ctx, insertReturnOrderLineQuery, params); err != nil {
				return fmt.Errorf("failed to insert ReturnOrderLine: %w", err)
			}
		}

		// Commit Transaction
		return nil
	})
}

// แสดงข้อมูลออเดอร์ที่เพิ่งสร้างไป
func (repo repositoryDB) GetCreateReturnOrder(ctx context.Context, orderNo string) (*response.CreateReturnOrder, error) {
	// ดึงข้อมูล ReturnOrder จาก OrderNo
	var order response.CreateReturnOrder
	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
			StatusCheckID, Description
		FROM ReturnOrder
		WHERE OrderNo = @orderNo
	`
	err := repo.db.GetContext(ctx, &order, queryOrder, sql.Named("orderNo", orderNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo %s: %w", orderNo, err)
	}

	// ดึงข้อมูล ReturnOrderLine ที่เกี่ยวข้องกับ OrderNo
	lines, err := repo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("error fetching ReturnOrderLines for OrderNo %s: %w", orderNo, err)
	}
	// เพิ่มข้อมูล ReturnOrderLine เข้าไปใน ReturnOrder
	order.ReturnOrderLine = lines

	return &order, nil
}

// แสดงข้อมูลออเดอร์นั้นที่เพิ่งอัพเดตไป
func (repo repositoryDB) GetUpdateReturnOrder(ctx context.Context, orderNo string) (*response.UpdateReturnOrder, error) {
	// ดึงข้อมูล ReturnOrder จาก OrderNo
	var order response.UpdateReturnOrder
	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, UpdateBy, UpdateDate, 
			CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		WHERE OrderNo = @orderNo
	`
	err := repo.db.GetContext(ctx, &order, queryOrder, sql.Named("orderNo", orderNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo: %w", err)
	}

	return &order, nil
}
// review
// สามารถอัปเดตข้อมูลออเดอได้ทั้งหมด แต่ข้อมูลรายการออเดอจะอัพเดตแค่ตอนเลข tracking มีการเปลี่ยนแปลง
// หากเผลออัพเดตค่าเดิมทั้งหมดจะทำการตรวจสอบกับข้อมูลปจบ.ก่อน เพื่อให้วันเวลาอัพเดตแสดงตามจริง เฉพาะฟิลด์ที่มีการเปลี่ยนแปลงจริง
func (repo repositoryDB) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder, updateBy string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 1: ดึงค่าปัจจุบันจากฐานข้อมูล
		var current response.ReturnOrder
		queryCurrent := `
            SELECT SrNo, TrackingNo, PlatfID, ChannelID, OptStatusID, AxStatusID, 
                   PlatfStatusID, Reason, CancelID, StatusCheckID, CheckBy, Description
            FROM ReturnOrder
            WHERE OrderNo = :OrderNo
        `
		currentParams := map[string]interface{}{"OrderNo": req.OrderNo}

		// ใช้ NamedQuery และ Rebind
		namedQuery, args, err := sqlx.Named(queryCurrent, currentParams)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		query := tx.Rebind(namedQuery)

		if err := tx.GetContext(ctx, &current, query, args...); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("OrderNo not found: %w", err)
			}
			return fmt.Errorf("failed to fetch current data: %w", err)
		}

		// Step 2: ตรวจสอบค่าที่เปลี่ยนแปลงและสร้าง SQL Update เฉพาะส่วนที่เปลี่ยนแปลง
		updateFields := []string{}
		params := map[string]interface{}{
			"OrderNo":  req.OrderNo,
			"UpdateBy": updateBy, // เพิ่ม updateBy ที่รับมาจาก API
		}

		if req.SrNo != nil && (current.SrNo == nil || *req.SrNo != *current.SrNo) {
			updateFields = append(updateFields, "SrNo = :SrNo")
			params["SrNo"] = req.SrNo
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
		if req.Reason != nil && (current.Reason == nil || *req.Reason != *current.Reason) {
			updateFields = append(updateFields, "Reason = :Reason")
			params["Reason"] = req.Reason
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
		updateFields = append(updateFields, "UpdateBy = :UpdateBy", "UpdateDate = GETDATE()")
		updateQuery := fmt.Sprintf(`
            UPDATE ReturnOrder
            SET %s
            WHERE OrderNo = :OrderNo
        `, strings.Join(updateFields, ", "))

		namedUpdateQuery, updateArgs, err := sqlx.Named(updateQuery, params)
		if err != nil {
			return fmt.Errorf("failed to prepare update query: %w", err)
		}
		updateQueryRebind := tx.Rebind(namedUpdateQuery)

		// ดำเนินการอัปเดตใน ReturnOrder
		if _, err := tx.ExecContext(ctx, updateQueryRebind, updateArgs...); err != nil {
			return fmt.Errorf("failed to update ReturnOrder: %w", err)
		}

		// อัปเดต TrackingNo, UpdateBy, UpdateDate ใน ReturnOrderLine หากมีการเปลี่ยนแปลงที่ฟิลด์ TrackingNo
		if req.TrackingNo != nil && (current.TrackingNo == nil || *req.TrackingNo != *current.TrackingNo) {
			updateLineQuery := `
                UPDATE ReturnOrderLine
                SET 
                    TrackingNo = :TrackingNo, 
                    UpdateBy = :UpdateBy, 
                    UpdateDate = GETDATE()
                WHERE OrderNo = :OrderNo
            `
			namedLineQuery, lineArgs, err := sqlx.Named(updateLineQuery, params)
			if err != nil {
				return fmt.Errorf("failed to prepare line update query: %w", err)
			}
			lineQueryRebind := tx.Rebind(namedLineQuery)

			if _, err := tx.ExecContext(ctx, lineQueryRebind, lineArgs...); err != nil {
				return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
			}
		}

		return nil
	})
}
// review
// ลบออเดอร์ head+line ที่สินค้าเข้าคลังมาเรียบร้อยแล้วออก
func (repo repositoryDB) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	// Step 1: เริ่ม Transaction
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		deleteReturnOrderLineQuery := `
			DELETE FROM ReturnOrderLine
			WHERE OrderNo = :OrderNo
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderLineQuery, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			return fmt.Errorf("failed to delete ReturnOrderLine: %w", err)
		}

		deleteReturnOrderQuery := `
			DELETE FROM ReturnOrder
			WHERE OrderNo = :OrderNo
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderQuery, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			return fmt.Errorf("error deleting ReturnOrder for OrderNo %s: %w", orderNo, err)
		}

		// Step 4: ส่ง Commit หากการลบสำเร็จทั้งหมด
		return nil
	})
}
// review
// Check ว่ามี OrderNo ในออเดอร์นั้นจริง
func (repo repositoryDB) CheckOrderNoExist(ctx context.Context, orderNo string) (bool, error) {
	var exists bool
	query := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM ReturnOrder WHERE OrderNo = @OrderNo
        ) THEN 1 ELSE 0 END
    `
	err := repo.db.GetContext(ctx, &exists, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		return false, fmt.Errorf("failed to check OrderNo existence: %w", err)
	}

	return exists, nil
}
// review
// Check ว่ามี OrderNo ในรายการออเดอร์นั้นจริง 
func (repo repositoryDB) CheckOrderNoLineExist(ctx context.Context, orderNo string) (bool, error) {
	var exists bool
	query := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM ReturnOrderLine WHERE OrderNo = @OrderNo
        ) THEN 1 ELSE 0 END
    `
	err := repo.db.GetContext(ctx, &exists, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		return false, fmt.Errorf("failed to check OrderNo existence: %w", err)
	}

	return exists, nil
}