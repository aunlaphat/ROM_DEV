package repository

import (
	request "boilerplate-back-go-2411/dto/request"
	response "boilerplate-back-go-2411/dto/response"
	"boilerplate-back-go-2411/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ReturnOrderRepository interface {
	GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error)
	GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error)
	GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error)
	GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error)
	CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error
	UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error
	UpdateReturnOrderLine(ctx context.Context, req request.UpdateReturnOrderLine) error
	DeleteReturnOrder(ctx context.Context, orderNo string) error
	CheckOrderNoExist(ctx context.Context, orderNo string) (bool, error)
	CheckOrderNoLineExist(ctx context.Context, orderNo string) (bool, error)

	GetCreateReturnOrder(ctx context.Context, orderNo string) (*response.CreateReturnOrder, error)
	GetUpdateReturnOrder(ctx context.Context, orderNo string) (*response.UpdateReturnOrder, error)
	GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error)
	GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error)
}

func (repo repositoryDB) GetAllReturnOrder(ctx context.Context) ([]response.ReturnOrder, error) {
	var orders []response.ReturnOrder

	queryHead := ` SELECT OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
						   OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
						   UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
					FROM ReturnOrder
					ORDER BY RecID
				  `
	err := repo.db.SelectContext(ctx, &orders, queryHead)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return orders: %w", err)
	}

	// *️⃣ ใช้แมปเพื่อจัดกลุ่ม Order
	orderMap := make(map[string]*response.ReturnOrder)
	orderNos := make([]string, 0, len(orders))
	for i := range orders {
		orderNos = append(orderNos, orders[i].OrderNo)
		orderMap[orders[i].OrderNo] = &orders[i]
	}

	// *️⃣ Batch Processing + sqlx.In เพื่อดึงข้อมูลเป็นกลุ่ม
	const batchSize = 1000
	for i := 0; i < len(orderNos); i += batchSize {
		end := i + batchSize
		if end > len(orderNos) {
			end = len(orderNos)
		}

		var lines []response.ReturnOrderLine
		queryLine := `  
					SELECT OrderNo, SKU, QTY, ReturnQTY, ActualQTY, Price, TrackingNo, 
						   CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
					FROM ReturnOrderLine
					WHERE OrderNo IN (?) 
					ORDER BY RecID
				`
		query, args, err := sqlx.In(queryLine, orderNos[i:end])
		if err != nil {
			return nil, fmt.Errorf("failed to prepare query: %w", err)
		}
		query = repo.db.Rebind(query)

		if err := repo.db.SelectContext(ctx, &lines, query, args...); err != nil {
			return nil, fmt.Errorf("failed to fetch return order lines: %w", err)
		}

		// *️⃣ แมปข้อมูล ReturnOrderLine กลับเข้า Order
		for _, line := range lines {
			if order, found := orderMap[line.OrderNo]; found {
				order.ReturnOrderLine = append(order.ReturnOrderLine, line) // เพิ่มเข้าไปใน Slice ต้นฉบับ
			}
		}
	}

	return orders, nil
}

func (repo repositoryDB) GetReturnOrderByOrderNo(ctx context.Context, orderNo string) (*response.ReturnOrder, error) {
	var order response.ReturnOrder

	query := `  SELECT  OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
						OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
						UpdateBy, UpdateDate, CancelID, StatusCheckID, CheckBy, Description
				FROM ReturnOrder
				WHERE OrderNo = :OrderNo
				ORDER BY RecID
			 `
	query, args, err := sqlx.Named(query, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	query = repo.db.Rebind(query)

	err = repo.db.GetContext(ctx, &order, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo %s: %w", orderNo, err)
	}

	lines, err := repo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("error fetching ReturnOrderLines for OrderNo %s: %w", orderNo, err)
	}

	//  *️⃣ เพิ่มข้อมูล ReturnOrderLine เข้าไปใน ReturnOrder ที่ OrderNo เดียวกัน
	order.ReturnOrderLine = lines

	return &order, nil
}

func (repo repositoryDB) GetAllReturnOrderLines(ctx context.Context) ([]response.ReturnOrderLine, error) {
	var lines []response.ReturnOrderLine

	query := `  SELECT OrderNo, SKU, QTY, ReturnQTY, ActualQTY, Price, TrackingNo, 
					   CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
				FROM ReturnOrderLine
				ORDER BY RecID
			 `

	err := repo.db.SelectContext(ctx, &lines, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch return order lines: %w", err)
	}

	return lines, nil
}

func (repo repositoryDB) GetReturnOrderLineByOrderNo(ctx context.Context, orderNo string) ([]response.ReturnOrderLine, error) {
	var lines []response.ReturnOrderLine

	query := `  SELECT OrderNo, SKU, ItemName, QTY, ReturnQTY, ActualQTY, Price, TrackingNo, 
					   CreateBy, CreateDate, AlterSKU, UpdateBy, UpdateDate
				FROM ReturnOrderLine
				WHERE OrderNo = :OrderNo
			 	ORDER BY RecID
			 `
	params := map[string]interface{}{
		"OrderNo": orderNo,
	}

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

func (repo repositoryDB) GetReturnOrdersByStatus(ctx context.Context, statusCheckID int) ([]response.DraftTradeDetail, error) {
	var orders []response.DraftTradeDetail

	query := ` 	SELECT 
					r.OrderNo, r.SoNo, r.CustomerID, r.SrNo, r.TrackingNo, 
					r.Logistic, c.ChannelName, r.CreateDate, 
					w.WarehouseName, r.StatusCheckID
				FROM ReturnOrder r
				LEFT JOIN Warehouse w ON r.WarehouseID = w.WarehouseID
				LEFT JOIN Channel c ON r.ChannelID = c.ChannelID
				WHERE r.StatusCheckID = :StatusCheckID 
				ORDER BY r.CreateDate ASC
			 `
	params := map[string]interface{}{
		"StatusCheckID": statusCheckID,
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

	return orders, nil
}

func (repo repositoryDB) GetReturnOrdersByStatusAndDateRange(ctx context.Context, statusCheckID int, startDate, endDate string) ([]response.DraftTradeDetail, error) {
	var orders []response.DraftTradeDetail

	query := ` 	SELECT 
					r.OrderNo, r.SoNo, r.CustomerID, r.SrNo, r.TrackingNo, 
					r.Logistic, c.ChannelName, r.CreateDate, 
					w.WarehouseName, r.StatusCheckID
				FROM ReturnOrder r
				LEFT JOIN Warehouse w ON r.WarehouseID = w.WarehouseID
				LEFT JOIN Channel c ON r.ChannelID = c.ChannelID
				WHERE r.StatusCheckID = :StatusCheckID 
				AND CAST(r.CreateDate AS DATE) >= :StartDate
				AND CAST(r.CreateDate AS DATE) <= :EndDate
				ORDER BY r.CreateDate ASC
			 `
	params := map[string]interface{}{
		"StatusCheckID": statusCheckID,
		"StartDate":     startDate,
		"EndDate":       endDate,
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

	return orders, nil
}

func (repo repositoryDB) CreateReturnOrder(ctx context.Context, req request.CreateReturnOrder) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		queryHead := `  INSERT INTO ReturnOrder (
							OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, OptStatusID, AxStatusID, 
							PlatfStatusID, Reason, StatusCheckID, Description, CreateBy, CreateDate
						) VALUES (
							:OrderNo, :SoNo, :SrNo, :TrackingNo, :PlatfID, :ChannelID, :OptStatusID, :AxStatusID, 
							:PlatfStatusID, :Reason, :StatusCheckID, :Description, :CreateBy, GETDATE()
						)
					 `
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

		if _, err := tx.NamedExecContext(ctx, queryHead, params); err != nil {
			return fmt.Errorf("failed to insert ReturnOrder: %w", err)
		}

		// *️⃣ บังคับให้สร้างรายการคืน 1 คำสั่งซื้อขึ้นไป จึงจะทำการสร้างได้
		if len(req.ReturnOrderLine) > 0 {
			queryLines := ` INSERT INTO ReturnOrderLine (
								OrderNo, SKU, QTY, ReturnQTY, Price, TrackingNo, CreateBy, CreateDate
							) VALUES (
								:OrderNo, :SKU, :QTY, :ReturnQTY, :Price, :TrackingNo, :CreateBy, GETDATE()
							) 
						  `
			var lineParams []map[string]interface{}
			for _, line := range req.ReturnOrderLine {
				lineParams = append(lineParams, map[string]interface{}{
					"OrderNo":    req.OrderNo,
					"TrackingNo": req.TrackingNo,
					"SKU":        line.SKU,
					"QTY":        line.QTY,
					"ReturnQTY":  line.ReturnQTY,
					"Price":      line.Price,
					"CreateBy":   req.CreateBy,
				})
			}

			if _, err := tx.NamedExecContext(ctx, queryLines, lineParams); err != nil {
				return fmt.Errorf("failed to insert ReturnOrderLines: %w", err)
			}
		}

		return nil
	})
}

// *️⃣ แสดงข้อมูลออเดอร์ที่พึ่งสร้าง
func (repo repositoryDB) GetCreateReturnOrder(ctx context.Context, orderNo string) (*response.CreateReturnOrder, error) {
	var order response.CreateReturnOrder

	queryHead := `  SELECT 
						OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
						OptStatusID, AxStatusID, PlatfStatusID, Reason, CreateBy, CreateDate, 
						StatusCheckID, Description
					FROM ReturnOrder
					WHERE OrderNo = @OrderNo
				 `
	err := repo.db.GetContext(ctx, &order, queryHead, sql.Named("OrderNo", orderNo))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // เมื่อไม่มีแถวที่ตรงกับเงื่อนไขในคิวรี่
			return nil, nil
		}
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo %s: %w", orderNo, err)
	}

	// *️⃣ ดึงข้อมูล ReturnOrderLine ที่เลข OrderNo ตรงกับใน ReturnOrder
	lines, err := repo.GetReturnOrderLineByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("error fetching ReturnOrderLines for OrderNo %s: %w", orderNo, err)
	}
	// *️⃣ เพิ่มข้อมูล ReturnOrderLine เข้าไปใน ReturnOrder
	order.ReturnOrderLine = lines

	return &order, nil
}

// *️⃣ แสดงข้อมูลออเดอร์ที่พึ่งอัพเดต
func (repo repositoryDB) GetUpdateReturnOrder(ctx context.Context, orderNo string) (*response.UpdateReturnOrder, error) {
	var order response.UpdateReturnOrder

	query := `  SELECT  OrderNo, SoNo, SrNo, TrackingNo, PlatfID, ChannelID, 
						OptStatusID, AxStatusID, PlatfStatusID, Reason, UpdateBy, UpdateDate, 
						CancelID, StatusCheckID, CheckBy, Description
				FROM ReturnOrder
				WHERE OrderNo = @OrderNo
			 `
	err := repo.db.GetContext(ctx, &order, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // เมื่อไม่มีแถวที่ตรงกับเงื่อนไขในคิวรี่
			return nil, nil
		}
		return nil, fmt.Errorf("database error querying ReturnOrder by OrderNo %s: %w", orderNo, err)
	}

	return &order, nil
}

// *️⃣ สามารถอัปเดตข้อมูลออเดอได้ทั้งหมด แต่ข้อมูลรายการออเดอจะอัพเดตแค่ตอนเลข tracking มีการเปลี่ยนแปลง
//
//	หากเผลออัพเดตค่าเดิมทั้งหมดจะทำการตรวจสอบกับข้อมูลปจบ.ก่อน เพื่อให้วันเวลาอัพเดตแสดงตามจริง เฉพาะฟิลด์ที่มีการเปลี่ยนแปลงจริง
func (repo repositoryDB) UpdateReturnOrder(ctx context.Context, req request.UpdateReturnOrder) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// *️⃣ ดึงค่าปัจจุบันจากฐานข้อมูล
		var current response.ReturnOrder

		queryCurrent := `   SELECT SrNo, TrackingNo, PlatfID, ChannelID, OptStatusID, AxStatusID, 
								   PlatfStatusID, Reason, CancelID, StatusCheckID, CheckBy, Description
							FROM ReturnOrder
							WHERE OrderNo = :OrderNo
						`
		currentParams := map[string]interface{}{"OrderNo": req.OrderNo}

		namedQuery, args, err := sqlx.Named(queryCurrent, currentParams)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		query := tx.Rebind(namedQuery)

		if err := tx.GetContext(ctx, &current, query, args...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("OrderNo not found: %w", err)
			}
			return fmt.Errorf("failed to fetch current data: %w", err)
		}

		// *️⃣ ตรวจสอบค่าที่เปลี่ยนแปลงและสร้าง SQL Update เฉพาะส่วนที่เปลี่ยนแปลง
		updateFields := []string{}

		params := map[string]interface{}{
			"OrderNo":  req.OrderNo,
			"UpdateBy": *req.UpdateBy, // เพิ่ม *req.UpdateBy = userID.(string) ที่รับมาจาก API
		}

		// *️⃣ ตรวจสอบทุกฟิลด์ที่อาจมีการเปลี่ยนแปลง
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

		// *️⃣ ตรวจสอบว่ามีการเปลี่ยนแปลงหรือไม่
		if len(updateFields) == 0 {
			return nil // ไม่มีการเปลี่ยนแปลง ไม่ต้องอัปเดต
		}

		// *️⃣ เพิ่ม UpdateBy และ UpdateDate ใน SQL Query
		updateFields = append(updateFields, "UpdateBy = :UpdateBy", "UpdateDate = GETDATE()")

		updateQuery := fmt.Sprintf(` UPDATE ReturnOrder
									  SET %s
									  WHERE OrderNo = :OrderNo `,
			strings.Join(updateFields, ", "))

		namedUpdateQuery, updateArgs, err := sqlx.Named(updateQuery, params)
		if err != nil {
			return fmt.Errorf("failed to prepare update query: %w", err)
		}
		updateQueryRebind := tx.Rebind(namedUpdateQuery)

		// *️⃣ ดำเนินการอัปเดตใน ReturnOrder
		if _, err := tx.ExecContext(ctx, updateQueryRebind, updateArgs...); err != nil {
			return fmt.Errorf("failed to update ReturnOrder: %w", err)
		}

		// *️⃣ อัปเดต TrackingNo, UpdateBy, UpdateDate ใน ReturnOrderLine หากมีการเปลี่ยนแปลงที่ฟิลด์ TrackingNo
		if req.TrackingNo != nil && (current.TrackingNo == nil || *req.TrackingNo != *current.TrackingNo) {

			updateLineQuery := `   UPDATE ReturnOrderLine
									SET TrackingNo = :TrackingNo, 
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

func (repo repositoryDB) UpdateReturnOrderLine(ctx context.Context, req request.UpdateReturnOrderLine) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {
		// ดึงข้อมูลปัจจุบันจาก ReturnOrderLine ตาม OrderNo และ SKU
		var current response.ReturnOrderLine
		queryCurrent := ` SELECT ActualQTY, Price 
						  FROM ReturnOrderLine 
						  WHERE OrderNo = :OrderNo AND SKU = :SKU 
						`
		currentParams := map[string]interface{}{
			"OrderNo": req.OrderNo,
			"SKU":     req.SKU,
		}

		namedQuery, args, err := sqlx.Named(queryCurrent, currentParams)
		if err != nil {
			return fmt.Errorf("failed to prepare query: %w", err)
		}
		query := tx.Rebind(namedQuery)

		if err := tx.GetContext(ctx, &current, query, args...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("OrderNo and SKU not found: %w", err)
			}
			return fmt.Errorf("failed to fetch current data: %w", err)
		}

		// ตรวจสอบค่าที่เปลี่ยนแปลง
		updateFields := []string{}
		params := map[string]interface{}{
			"OrderNo":  req.OrderNo,
			"SKU":      req.SKU,
			"UpdateBy": *req.UpdateBy, // รับค่า UserID จาก API
		}

		if req.ActualQTY != nil && *req.ActualQTY != *current.ActualQTY {
			updateFields = append(updateFields, "ActualQTY = :ActualQTY")
			params["ActualQTY"] = req.ActualQTY
		}
		if req.Price != nil && *req.Price != current.Price {
			updateFields = append(updateFields, "Price = :Price")
			params["Price"] = req.Price
		}

		// ถ้าไม่มีการเปลี่ยนแปลง ไม่ต้องอัปเดต
		if len(updateFields) == 0 {
			return nil
		}

		// เพิ่มฟิลด์ UpdateBy และ UpdateDate
		updateFields = append(updateFields, "UpdateBy = :UpdateBy", "UpdateDate = GETDATE()")
		updateQuery := fmt.Sprintf(` 
						UPDATE ReturnOrderLine SET %s 
						WHERE OrderNo = :OrderNo AND SKU = :SKU`,
			strings.Join(updateFields, ", "))

		namedUpdateQuery, updateArgs, err := sqlx.Named(updateQuery, params)
		if err != nil {
			return fmt.Errorf("failed to prepare update query: %w", err)
		}
		updateQueryRebind := tx.Rebind(namedUpdateQuery)

		// ทำการอัปเดต ReturnOrderLine
		if _, err := tx.ExecContext(ctx, updateQueryRebind, updateArgs...); err != nil {
			return fmt.Errorf("failed to update ReturnOrderLine: %w", err)
		}

		return nil
	})
}

// *️⃣ ลบออเดอร์ head+line ที่สินค้าเข้าคลังมาเรียบร้อยแล้วออก
func (repo repositoryDB) DeleteReturnOrder(ctx context.Context, orderNo string) error {
	return utils.HandleTransaction(repo.db, func(tx *sqlx.Tx) error {

		queryLine := ` DELETE FROM ReturnOrderLine
					   WHERE OrderNo = :OrderNo
					 `
		if _, err := tx.NamedExecContext(ctx, queryLine, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			return fmt.Errorf("failed to delete ReturnOrderLine: %w", err)
		}

		queryHead := ` DELETE FROM ReturnOrder
			           WHERE OrderNo = :OrderNo
		             `
		if _, err := tx.NamedExecContext(ctx, queryHead, map[string]interface{}{
			"OrderNo": orderNo,
		}); err != nil {
			return fmt.Errorf("error deleting ReturnOrder for OrderNo %s: %w", orderNo, err)
		}

		return nil
	})
}

// *️⃣ Check ว่ามี OrderNo ในออเดอร์นั้นจริง
func (repo repositoryDB) CheckOrderNoExist(ctx context.Context, orderNo string) (bool, error) {
	var exists bool

	query := `	SELECT CASE WHEN EXISTS (
				SELECT 1 FROM ReturnOrder WHERE OrderNo = @OrderNo
				) THEN 1 ELSE 0 END
    		 `
	err := repo.db.GetContext(ctx, &exists, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		return false, fmt.Errorf("failed to check OrderNo existence: %w", err)
	}

	return exists, nil
}

// *️⃣ Check ว่ามี OrderNo ในรายการออเดอร์นั้นจริง
func (repo repositoryDB) CheckOrderNoLineExist(ctx context.Context, orderNo string) (bool, error) {
	var exists bool

	query := `	SELECT CASE WHEN EXISTS (
				SELECT 1 FROM ReturnOrderLine WHERE OrderNo = @OrderNo
				) THEN 1 ELSE 0 END
   			 `
	err := repo.db.GetContext(ctx, &exists, query, sql.Named("OrderNo", orderNo))
	if err != nil {
		return false, fmt.Errorf("failed to check OrderNo existence: %w", err)
	}

	return exists, nil
}
