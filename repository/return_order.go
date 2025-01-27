package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"reflect"
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
	GetReturnOrderByID(ctx context.Context, orderNo string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLinesByReturnID(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error
	DeleteReturnOrder(ctx context.Context, orderNo string) error
	CheckOrderNoExist(ctx context.Context, orderNo string) (bool, error)


}

func (repo repositoryDB) AllGetReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	// Step 1: กำหนดคิวรี่ดึงข้อมูล ReturnOrder ทั้งหมด
	var orders []response.ReturnOrder
	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
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

	// Step 3: จับคู่ข้อมูล ReturnOrderLine กับ ReturnOrder ตาม OrderNo
	linesMap := make(map[string][]response.ReturnOrderLine)
	for _, line := range lines {
		linesMap[line.OrderNo] = append(linesMap[line.OrderNo], line)
	}

	// Step 4: เพิ่มข้อมูล ReturnOrderLine เข้าไปในแต่ละ ReturnOrder
	for i := range orders {
		orders[i].ReturnOrderLine = linesMap[orders[i].OrderNo]
	}

	return orders, nil
}

func (repo repositoryDB) GetReturnOrderByID(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	// Step 1: ดึงข้อมูล ReturnOrder จาก OrderNo
	var order response.ReturnOrder
	queryOrder := `
		SELECT 
			OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
			OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
			UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
		FROM ReturnOrder
		WHERE OrderNo = @orderNo
	`
	err := repo.db.GetContext(ctx, &order, queryOrder, sql.Named("orderNo", orderNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("Database error querying ReturnOrder by OrderNo: %v", err)
		return nil, fmt.Errorf("unexpected database error: %w", err)
	}

	// Step 2: ดึงข้อมูล ReturnOrderLine ที่เกี่ยวข้องกับ OrderNo
	lines, err := repo.GetReturnOrderLinesByReturnID(ctx, orderNo)
	if err != nil {
		log.Printf("Error fetching ReturnOrderLines for OrderNo %s: %v", orderNo, err)
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
		SELECT OrderNo, TrackingNo, SKU, ReturnQTY, QTY, Price, 
		       CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
		FROM ReturnOrderLine
		ORDER BY OrderNo
	`

	err := repo.db.SelectContext(ctx, &lines, query)
	if err != nil {
		log.Printf("Error querying ReturnOrderLine: %v", err)
		return nil, fmt.Errorf("get all return order lines error: %w", err)
	}

	return lines, nil
}

// Get ReturnOrderLines by OrderNo
func (repo repositoryDB) GetReturnOrderLinesByReturnID(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	// ตรวจสอบว่า OrderNo มีอยู่จริงในฐานข้อมูล
	var exists bool
	checkQuery := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM ReturnOrder WHERE OrderNo = @OrderNo
        ) THEN 1 ELSE 0 END
    `
	err := repo.db.GetContext(ctx, &exists, checkQuery, sql.Named("OrderNo", orderNo))
	if err != nil {
		return nil, fmt.Errorf("failed to check OrderNo existence: %w", err)
	}
	if !exists {
		return nil, sql.ErrNoRows // OrderNo ไม่พบ
	}

	// ดึงข้อมูล ReturnOrderLines
	var lines []response.ReturnOrderLine
	query := `
        SELECT OrderNo, TrackingNo, SKU, ReturnQTY, QTY, Price, 
               CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
        FROM ReturnOrderLine
        WHERE OrderNo = @OrderNo
    `
	err = repo.db.SelectContext(ctx, &lines, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows // ไม่มีข้อมูล ReturnOrderLine
		}
		return nil, fmt.Errorf("error querying ReturnOrderLine by OrderNo: %w", err)
	}

	return lines, nil
}

func (repo repositoryDB) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	// Step 1: เริ่มต้น Transaction
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 2: เพิ่มข้อมูล ReturnOrder ลงฐานข้อมูล
		insertReturnOrderQuery := `
            INSERT INTO ReturnOrder (
                OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
                OptStatusID, AxStatusID, PlatfStatusID, Reason, CancelID, StatusCheckID, 
                CheckBy, Description, CreateBy, CreateDate
            ) VALUES (
                :OrderNo, :SoNo, :SrNo, :TrackingNo, :PlatfID, :ChannelID, 
                :OptStatusID, :AxStatusID, :PlatfStatusID, :Reason, :CancelID, :StatusCheckID, 
                :CheckBy, :Description, 'USER', SYSDATETIME()
            )
        `
		// Step 2.1: ตรวจสอบว่ามีค่าสำหรับการใส่ใน Query
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
                OrderNo, TrackingNo, SKU, ReturnQTY, QTY, Price, CreateBy, CreateDate
            ) VALUES (
                :OrderNo, :TrackingNo, :SKU, :ReturnQTY, :QTY, :Price, 
                 'USER', SYSDATETIME()
            )
        `
		// Step 3.1: Loop ผ่าน `ReturnOrderLine` และใส่ข้อมูลลงใน Query
		for _, line := range req.ReturnOrderLine {
			line.OrderNo = req.OrderNo
			line.TrackingNo = req.TrackingNo

			params := map[string]interface{}{
				"OrderNo":    line.OrderNo,
				"TrackingNo": line.TrackingNo,
				"SKU":        line.SKU,
				"ReturnQTY":  line.ReturnQTY,
				"QTY":        line.QTY,
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
        current, err := repo.getCurrentReturnOrder(ctx, tx, req.OrderNo)
        if err != nil {
            return err
        }

        // Step 2: ตรวจสอบค่าที่เปลี่ยนแปลงและสร้าง SQL Update เฉพาะส่วนที่เปลี่ยนแปลง
        updateFields, params := repo.buildUpdateFields(req, current)

        // หากไม่มีการเปลี่ยนแปลง ออกจากฟังก์ชัน
        if len(updateFields) == 0 {
            return nil
        }

        // Step 3: เพิ่ม UpdateBy และ UpdateDate ใน SQL Query
        updateFields = append(updateFields, "UpdateBy = 'USER'", "UpdateDate = SYSDATETIME()")
        updateQuery := fmt.Sprintf(`
            UPDATE ReturnOrder
            SET %s
            WHERE OrderNo = :OrderNo
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
            if err := repo.updateReturnOrderLineTrackingNo(ctx, tx, req.OrderNo, req.TrackingNo); err != nil {
                return err
            }
        }

        return nil
    })
}

func (repo repositoryDB) getCurrentReturnOrder(ctx context.Context, tx *sqlx.Tx, orderNo string) (response.ReturnOrder, error) {
    var current response.ReturnOrder
    queryCurrent := `
        SELECT SrNo, TrackingNo, PlatfID, ChannelID, OptStatusID, AxStatusID, 
               PlatfStatusID, Reason, CancelID, StatusCheckID, CheckBy, Description
        FROM ReturnOrder
        WHERE OrderNo = :OrderNo
    `
    currentParams := map[string]interface{}{"OrderNo": orderNo}

    // ใช้ NamedQuery และ Rebind
    namedQuery, args, err := sqlx.Named(queryCurrent, currentParams)
    if err != nil {
        return current, fmt.Errorf("failed to prepare query: %w", err)
    }
    query := tx.Rebind(namedQuery)

    if err := tx.GetContext(ctx, &current, query, args...); err != nil {
        if err == sql.ErrNoRows {
            return current, fmt.Errorf("OrderNo not found: %w", err)
        }
        return current, fmt.Errorf("failed to fetch current data: %w", err)
    }

    return current, nil
}

func (repo repositoryDB) updateReturnOrderLineTrackingNo(ctx context.Context, tx *sqlx.Tx, orderNo string, trackingNo *string) error {
    updateLineQuery := `
        UPDATE ReturnOrderLine
        SET 
            TrackingNo = :TrackingNo, 
            UpdateBy = 'USER', 
            UpdateDate = SYSDATETIME()
        WHERE OrderNo = :OrderNo
    `
    params := map[string]interface{}{
        "OrderNo":    orderNo,
        "TrackingNo": trackingNo,
    }

    namedLineQuery, lineArgs, err := sqlx.Named(updateLineQuery, params)
    if err != nil {
        return fmt.Errorf("failed to prepare line update query: %w", err)
    }
    lineQueryRebind := tx.Rebind(namedLineQuery)

    if _, err := tx.ExecContext(ctx, lineQueryRebind, lineArgs...); err != nil {
        log.Printf("Error updating ReturnOrderLine: %v", err)
        return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
    }

    return nil
}

// ฟังก์ชันสำหรับสร้าง SQL Update เฉพาะส่วนที่เปลี่ยนแปลง
func (repo repositoryDB) buildUpdateFields(req request.UpdateReturnOrder, current response.ReturnOrder) ([]string, map[string]interface{}) {
    updateFields := []string{}
    params := map[string]interface{}{
        "OrderNo": req.OrderNo,
    }

    reqValue := reflect.ValueOf(req)
    currentValue := reflect.ValueOf(current)
    reqType := reqValue.Type()

    for i := 0; i < reqValue.NumField(); i++ {
        field := reqType.Field(i)
        fieldName := field.Name
        reqField := reqValue.Field(i)
        currentField := currentValue.FieldByName(fieldName)

        // ตรวจสอบชนิดข้อมูลก่อนเรียกใช้ IsNil
        if reqField.Kind() == reflect.Ptr {
            if !reqField.IsNil() && (currentField.IsNil() || reqField.Elem().Interface() != currentField.Elem().Interface()) {
                updateFields = append(updateFields, fmt.Sprintf("%s = :%s", fieldName, fieldName))
                params[fieldName] = reqField.Interface()
            }
        } else {
            if reqField.Interface() != currentField.Interface() {
                updateFields = append(updateFields, fmt.Sprintf("%s = :%s", fieldName, fieldName))
                params[fieldName] = reqField.Interface()
            }
        }
    }

    return updateFields, params
}

// ลบ ReturnOrder และ ReturnOrderLine
func (repo repositoryDB) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	// Step 1: เริ่ม Transaction เพื่อควบคุมการลบ
	return handleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// Step 2: ลบข้อมูลใน ReturnOrderLine โดยอิงจาก OrderNo
		deleteReturnOrderLineQuery := `
			DELETE FROM ReturnOrderLine
			WHERE OrderNo = :OrderNo
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderLineQuery, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			log.Printf("Error deleting ReturnOrderLine for OrderNo %s: %v", orderNo, err)
			return fmt.Errorf("failed to delete ReturnOrderLine: %w", err)
		}

		// Step 3: ลบข้อมูลใน ReturnOrder
		deleteReturnOrderQuery := `
			DELETE FROM ReturnOrder
			WHERE OrderNo = :OrderNo
		`
		if _, err := tx.NamedExecContext(ctx, deleteReturnOrderQuery, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			log.Printf("Error deleting ReturnOrder for OrderNo %s: %v", orderNo, err)
			return fmt.Errorf("failed to delete ReturnOrder: %w", err)
		}

		// Step 4: ส่ง Commit หากการลบสำเร็จทั้งหมด
		return nil
	})
}

// CheckBefOrderNoExists - เพิ่มการตรวจสอบ OrderNo ว่ามีอยู่ในฐานข้อมูล
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
