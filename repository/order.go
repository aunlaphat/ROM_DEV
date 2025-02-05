package repository

import (
	"boilerplate-backend-go/dto/response"
	"context"
)

type OrderRepository interface {
	SearchOrder(ctx context.Context, orderNo, soNo string) (*response.SearchOrderResponse, error)
}

func (repo repositoryDB) SearchOrder(ctx context.Context, orderNo, soNo string) (*response.SearchOrderResponse, error) {
	var order response.SearchOrderResponse

	// ✅ คำสั่ง SQL ปรับปรุงใหม่ รองรับหลาย SoNo ต่อ OrderNo
	query := `
		SELECT 
			h.OrderNo, h.SoNo, h.StatusMKP, h.SalesStatus, h.CreateDate,
			l.SKU, l.ItemName, l.QTY, l.Price
		FROM ROM_V_OrderHeadDetail h
		JOIN ROM_V_OrderLineDetail l
			ON h.OrderNo = l.OrderNo 
			AND (h.SoNo = l.SoNo OR h.SoNo IS NULL)
		WHERE (:soNo IS NULL OR h.SoNo = :soNo)
		  AND (:orderNo IS NULL OR h.OrderNo = :orderNo)
	`

	// ✅ กำหนด Named Parameters
	params := map[string]interface{}{
		"soNo":    soNo,
		"orderNo": orderNo,
	}

	// ✅ ดึงข้อมูลคำสั่งซื้อและรายการสินค้า
	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ✅ Map ข้อมูลที่ได้จาก Query
	var items []response.SearchOrderItem
	for rows.Next() {
		var item response.SearchOrderItem
		if err := rows.Scan(&order.OrderNo, &order.SoNo, &order.StatusMKP, &order.SalesStatus, &order.CreateDate,
			&item.SKU, &item.ItemName, &item.QTY, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.Items = items

	// ✅ ตรวจสอบว่าพบข้อมูลหรือไม่
	if order.OrderNo == "" {
		return nil, nil
	}

	return &order, nil
}
