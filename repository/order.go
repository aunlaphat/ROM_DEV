package repository

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	"os"
	"log"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// defines methods for CRUD operations on orders
type OrderRepository interface {
	AllGetOrder() ([]response.OrderResponse, error)
	GetOrderID(orderNo string) (response.OrderResponse, error)
	AllGetOrderLinesByOrderNo(orderNo string) ([]response.OrderLineResponse, error)
	CreateOrder(req request.CreateOrderRequest) error
	UpdateOrder(req request.UpdateOrderRequest) error
	DeleteOrder(orderNo string) error
	CheckOrderExists(orderNo string) (bool, error)
}

func (repo repositoryDB) AllGetOrder() ([]response.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders := []response.OrderResponse{}
	
	sqlQuery := ` SELECT *
				  FROM OrderHead
				  ORDER BY OrderNo
			    `

	rows, err := repo.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order response.OrderResponse
		if err := rows.Scan(
			&order.OrderNo, &order.BrandName, &order.CustName, &order.CustAddress,
			&order.CustDistrict, &order.CustSubDistrict, &order.CustProvince,
			&order.CustPostCode, &order.CustPhoneNum, &order.CreateDate,
			&order.UserCreated, &order.UpdateDate, &order.UserUpdated,
		); err != nil {
			return nil, err
		}

		// Fetch OrderLine ออกมาเพื่อเพิ่มข้อมูลรายการสินค้าเข้ามารวมกับ OrderHead ข้อมูลลูกค้า
		orderLines, err := repo.AllGetOrderLinesByOrderNo(order.OrderNo)
		if err != nil {
			return nil, err
		}
		order.OrderLines = orderLines
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (repo repositoryDB) GetOrderID(orderNo string) (response.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var order response.OrderResponse
	sqlQuery := ` SELECT *
				  FROM OrderHead
				  WHERE OrderNo = :OrderNo
    			`

	// ใช้ NamedQueryContext เพื่อดึงข้อมูล
	rows, err := repo.db.NamedQueryContext(ctx, sqlQuery, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return response.OrderResponse{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// ตรวจสอบว่ามีผลลัพธ์หรือไม่
	if rows.Next() {
		// ใช้ StructScan เพื่อ Map ข้อมูลเข้า Struct
		if err := rows.StructScan(&order); err != nil {
			return response.OrderResponse{}, fmt.Errorf("failed to scan order: %w", err)
		}
	} else {

		return response.OrderResponse{}, fmt.Errorf("order with OrderNo %s not found", orderNo)
	}

	// Fetch OrderLines
	orderLines, err := repo.AllGetOrderLinesByOrderNo(order.OrderNo)
	if err != nil {
		return response.OrderResponse{}, fmt.Errorf("failed to fetch order lines: %w", err)
	}
	order.OrderLines = orderLines

	return order, nil
}

// ดึงข้อมูล OrderLine ที่สัมพันธ์กับ OrderNo ของ OrderHead
func (repo repositoryDB) AllGetOrderLinesByOrderNo(orderNo string) ([]response.OrderLineResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderLines := []response.OrderLineResponse{}
	sqlQuery := ` SELECT OrderNo, SKU, ItemName, QTY, Price
				  FROM OrderLine
				  WHERE OrderNo = :OrderNo
			    `

	rows, err := repo.db.NamedQueryContext(ctx, sqlQuery, map[string]interface{}{"OrderNo": orderNo})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderLine response.OrderLineResponse
		if err := rows.StructScan(&orderLine); err != nil {
			return nil, err
		}
		orderLines = append(orderLines, orderLine)
	}

	// ตรวจสอบข้อผิดพลาดหลังวนลูป
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderLines, nil
}

func (repo repositoryDB) CreateOrder(req request.CreateOrderRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ดึงชื่อเครื่อง
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	// ตรวจสอบว่ามี OrderHead อยู่แล้วหรือไม่
	var count int
	checkQuery := ` SELECT COUNT(*) 
					FROM OrderHead 
					WHERE OrderNo = @OrderNo
				  `
	err = repo.db.QueryRowContext(ctx, checkQuery, sql.Named("OrderNo", req.OrderNo)).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check OrderHead existence: %w", err)
	}
	log.Printf("OrderNo %s exists count: %d", req.OrderNo, count)

	// หากมี OrderHead อยู่แล้ว ให้แจ้งเตือน
	if count > 0 {
		return fmt.Errorf("OrderNo %s already exists; Please use a unique OrderNo", req.OrderNo)
	}

	// หากไม่มี OrderHead ให้สร้างใหม่
		insertOrderHeadQuery := `
				INSERT INTO OrderHead (OrderNo, BrandName, CustName, CustAddress, CustDistrict, CustSubDistrict, CustProvince, CustPostCode, CustPhoneNum, CreateDate, UserCreated)
				VALUES (@OrderNo, @BrandName, @CustName, @CustAddress, @CustDistrict, @CustSubDistrict, @CustProvince, @CustPostCode, @CustPhoneNum, @CreateDate, @UserCreated)
		    `
		_, err = repo.db.ExecContext(ctx, insertOrderHeadQuery,
			sql.Named("OrderNo", req.OrderNo),
			sql.Named("BrandName", req.BrandName),
			sql.Named("CustName", req.CustName),
			sql.Named("CustAddress", req.CustAddress),
			sql.Named("CustDistrict", req.CustDistrict),
			sql.Named("CustSubDistrict", req.CustSubDistrict),
			sql.Named("CustProvince", req.CustProvince),
			sql.Named("CustPostCode", req.CustPostCode),
			sql.Named("CustPhoneNum", req.CustPhoneNum),
			sql.Named("CreateDate", time.Now()), // ดึงค่าเวลาปัจจุบันของเครื่อง
			sql.Named("UserCreated", hostname),  // กำหนดค่าผู้สร้าง
		)
		if err != nil {
			return fmt.Errorf("failed to insert OrderHead: %w", err)
		}
	

	// เพิ่มข้อมูลใน OrderLine
	insertOrderLineQuery := ` INSERT INTO OrderLine (OrderNo, SKU, ItemName, QTY, Price)
							  VALUES (@OrderNo, @SKU, @ItemName, @QTY, @Price)
	                        `
	for _, line := range req.OrderLines {
		_, err := repo.db.ExecContext(ctx, insertOrderLineQuery,
			sql.Named("OrderNo", req.OrderNo),
			sql.Named("SKU", line.SKU),
			sql.Named("ItemName", line.ItemName),
			sql.Named("QTY", line.QTY),
			sql.Named("Price", line.Price),
		)
		if err != nil {
			return fmt.Errorf("failed to insert OrderLine for OrderNo %s: %w", req.OrderNo, err)
		}
	}

	return nil
}

func (repo repositoryDB) UpdateOrder(req request.UpdateOrderRequest) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // ดึงชื่อเครื่อง (hostname)
    hostname, err := os.Hostname()
    if err != nil {
        return fmt.Errorf("failed to get hostname: %w", err)
    }

	// ขึ้นกับ frontend ว่าจะดึงตัวไหนไป เมื่อไม่ใส่ฟิลด์บางตัวในการเทสจึงยังเป็น null อยู่
    query := `	UPDATE OrderHead
				SET  CustName = @CustName
					,CustAddress = @CustAddress
					,CustDistrict = @CustDistrict
					,CustSubDistrict = @CustSubDistrict
					,CustProvince = @CustProvince
					,CustPostCode = @CustPostCode
					,CustPhoneNum = @CustPhoneNum
					,UpdateDate = @UpdateDate
					,UserUpdated = @UserUpdated
				WHERE OrderNo = @OrderNo
			`
    // อัปเดตข้อมูลในฐานข้อมูล
    _, err = repo.db.ExecContext(ctx, query,
        sql.Named("CustName", req.CustName),
        sql.Named("CustAddress", req.CustAddress),
        sql.Named("CustDistrict", req.CustDistrict),
        sql.Named("CustSubDistrict", req.CustSubDistrict),
        sql.Named("CustProvince", req.CustProvince),
        sql.Named("CustPostCode", req.CustPostCode),
        sql.Named("CustPhoneNum", req.CustPhoneNum),
        sql.Named("UpdateDate", time.Now()),
        sql.Named("UserUpdated", hostname),
        sql.Named("OrderNo", req.OrderNo),
    )

    // ตรวจสอบข้อผิดพลาดจากการอัปเดต
    if err != nil {
        return fmt.Errorf("failed to update order: %w", err)
    }

    return nil
}


func (repo repositoryDB) DeleteOrder(orderNo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ลบข้อมูลจาก OrderLine
	deleteOrderLineQuery := ` DELETE FROM OrderLine WHERE OrderNo = @OrderNo `
	_, err := repo.db.ExecContext(ctx, deleteOrderLineQuery, sql.Named("OrderNo", orderNo))
	if err != nil {
		return fmt.Errorf("failed to delete from OrderLine: %w", err)
	}

	// ลบข้อมูลจาก OrderHead
	deleteOrderHeadQuery := ` DELETE FROM OrderHead WHERE OrderNo = @OrderNo `
	_, err = repo.db.ExecContext(ctx, deleteOrderHeadQuery, sql.Named("OrderNo", orderNo))
	if err != nil {
		return fmt.Errorf("failed to delete from OrderHead: %w", err)
	}

	return nil
}

func (repo repositoryDB) CheckOrderExists(orderNo string) (bool, error) {
	var exists bool
	query := ` SELECT CASE 
			   WHEN EXISTS (SELECT 1 FROM OrderHead WHERE OrderNo = @OrderNo) 
			   THEN 1 ELSE 0 END
	         `
	err := repo.db.QueryRow(query, sql.Named("OrderNo", orderNo)).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}


// func (repo repositoryDB) GetOrder() ([]entity.Order, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// ดึงข้อมูลจาก OrderHead
// 	orders := []entity.Order{}
// 	sqlQuery := `SELECT *
//                  FROM OrderHead
//                  ORDER BY OrderNo
// 				`
// 				// `SELECT OrderNo
// 				// 	, BrandName
// 				// 	, CustName
// 				// 	, CustAddress
// 				// 	, CustDistrict
// 				// 	, CustSubDistrict
// 				// 	, CustProvince
// 				// 	, CustPostCode
// 				// 	, CustPhoneNum
// 				// 	, CreateDate
// 				// 	, UserCreated
// 				// 	, UpdateDate
// 				// 	, UserUpdated
//                 //  FROM OrderHead
//                 //  ORDER BY OrderNo`

// 	rows, err := repo.db.QueryxContext(ctx, sqlQuery)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// ดึงข้อมูลแต่ละ OrderHead
// 	for rows.Next() {
// 		var order entity.Order
// 		if err := rows.StructScan(&order); err != nil {
// 			return nil, err
// 		}

// 		// เพิ่ม log เพื่อดูค่าของ Order
// 		fmt.Println("Fetched Order:", order)

// 		// ดึงข้อมูล OrderLine ที่สัมพันธ์กับ OrderNo
// 		orderLines, err := repo.GetOrderLinesByOrderNo(order.OrderNo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// เพิ่ม log เพื่อดูค่าของ OrderLine
// 		fmt.Println("Fetched OrderLines:", orderLines)

// 		order.OrderLines = orderLines // เพิ่มข้อมูล OrderLine เข้าไปใน Order
// 		orders = append(orders, order)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

// // ดึงข้อมูล OrderLine ที่สัมพันธ์กับ OrderNo
// func (repo repositoryDB) GetOrderLinesByOrderNo(orderNo string) ([]entity.OrderLine, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	orderLines := []entity.OrderLine{}

// 	// SQL query สำหรับดึง OrderLine โดยใช้ OrderNo
// 	sqlQuery := ` SELECT OrderNo, QTY, SKU, ItemName, Price
//                   FROM OrderLine
//                   WHERE OrderNo = :OrderNo
// 				`

// 	// ส่ง parameter @OrderNo ให้กับ QueryxContext
// 	rows, err := repo.db.NamedQueryContext(ctx, sqlQuery, map[string]interface{}{"OrderNo": orderNo})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var orderLine entity.OrderLine
// 		if err := rows.StructScan(&orderLine); err != nil {
// 			return nil, err
// 		}
// 		orderLines = append(orderLines, orderLine)
// 	}

// 	// เพิ่ม log เพื่อดูค่าของ Order
// 	fmt.Println("Fetched Order:", orderLines)

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return orderLines, nil
// }
