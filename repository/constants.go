package repository

import (
	entity "boilerplate-backend-go/Entity"
	"context"
	"fmt"
)

type Constants interface {
	GetThaiProvince(ctx context.Context) ([]entity.Province, error)       // จังหวัด
	GetThaiDistrict(ctx context.Context) ([]entity.District, error)       // เขต
	GetThaiSubDistrict(ctx context.Context) ([]entity.SubDistrict, error) // ตำบล
	// GetPostCode(ctx context.Context) ([]entity.PostCode, error) // เลขไปรษณีย์
	GetProduct(ctx context.Context, offset, limit int) ([]entity.ROM_V_ProductAll, error) // รายการสินค้าแบบแบ่งรายการ
	GetWarehouse(ctx context.Context) ([]entity.Warehouse, error)                         // ชื่อคลังสินค้า
	// GetCustomer(ctx context.Context) ([]entity.ROM_V_Customer, error) // ข้อมูลลูกค้า
	SearchProduct(ctx context.Context, keyword string) ([]entity.ROM_V_ProductAll, error)
}

func (repo repositoryDB) GetThaiProvince(ctx context.Context) ([]entity.Province, error) {
	provinces := []entity.Province{}

	query := `  SELECT Code, NameTH, NameEN
				FROM V_ThaiAddressProvince
				ORDER BY Code
			 `
	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var province entity.Province
		if err := rows.StructScan(&province); err != nil {
			return nil, err
		}
		provinces = append(provinces, province)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return provinces, nil
}

func (repo repositoryDB) GetThaiDistrict(ctx context.Context) ([]entity.District, error) {
	districts := []entity.District{}

	query := `  SELECT ProvinceCode, Code, NameTH, NameEN
				FROM V_ThaiAddressDistrict
				ORDER BY Code
			 `

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var district entity.District
		if err := rows.StructScan(&district); err != nil {
			return nil, err
		}
		districts = append(districts, district)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}

func (repo repositoryDB) GetThaiSubDistrict(ctx context.Context) ([]entity.SubDistrict, error) {
	subDistricts := []entity.SubDistrict{}

	query := `	SELECT Code, DistrictCode, ZipCode, NameTH, NameEN
				FROM V_ThaiAddressSubDistrict
				ORDER BY Code
			 `
	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subDistrict entity.SubDistrict
		if err := rows.StructScan(&subDistrict); err != nil {
			return nil, err
		}
		subDistricts = append(subDistricts, subDistrict)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subDistricts, nil
}

// func (repo repositoryDB) GetPostCode(ctx context.Context) ([]entity.PostCode, error) {

// 	postCodes := []entity.PostCode{}

// 	query := `

//           `

// 	rows, err := repo.db.QueryxContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var postCode entity.PostCode
// 		if err := rows.StructScan(&postCode); err != nil {
// 			return nil, err
// 		}
// 		postCodes = append(postCodes, postCode)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return postCodes, nil
// }

func (repo repositoryDB) GetWarehouse(ctx context.Context) ([]entity.Warehouse, error) {
	warehouses := []entity.Warehouse{}

	query := `  SELECT WarehouseID, WarehouseName, Location
				FROM Warehouse
				ORDER BY WarehouseName
			 `

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse entity.Warehouse
		if err := rows.StructScan(&warehouse); err != nil {
			return nil, fmt.Errorf("failed to scan warehouse: %w", err)
		}
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (repo repositoryDB) GetProduct(ctx context.Context, offset, limit int) ([]entity.ROM_V_ProductAll, error) {
	query := `  SELECT SKU, NAMEALIAS, Size, SizeID, Barcode, Type
				FROM Data_WebReturn.dbo.ROM_V_ProductAll
				ORDER BY SKU
				OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY;
			 `
	countQuery := ` SELECT COUNT(*)
					FROM Data_WebReturn.dbo.ROM_V_ProductAll;
				  `

	var products []entity.ROM_V_ProductAll
	total := 0

	// Fetch total count
	if err := repo.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, fmt.Errorf("failed to fetch total count: %w", err)
	}

	// Fetch paginated data
	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.ROM_V_ProductAll
		if err := rows.StructScan(&product); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// func (repo repositoryDB) GetCustomer(ctx context.Context) ([]entity.SubDistrict, error) {

// }

func (repo repositoryDB) SearchProduct(ctx context.Context, keyword string) ([]entity.ROM_V_ProductAll, error) {
	// *️⃣ จำนวนข้อมูลที่ต้องการดึงในแต่ละ batch
	const chunkSize = 1000
	var products []entity.ROM_V_ProductAll
	offset := 0

	// *️⃣ เพิ่ม `%` เข้าไปใน Go ก่อนใช้กับ SQL
	searchParam := "%" + keyword + "%"

	for {
		query := `	SELECT SKU, NAMEALIAS, Size, SizeID, Barcode, Type
					FROM Data_WebReturn.dbo.ROM_V_ProductAll
					WHERE NAMEALIAS LIKE :Keyword 
						  OR SKU LIKE :Keyword
					ORDER BY SKU
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				 `
		var productBatch []entity.ROM_V_ProductAll
		nstmt, err := repo.db.PrepareNamed(query) // PrepareNamed ลด overhead จากการ compile SQL ซ้ำ ๆ
		if err != nil {
			return nil, fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer nstmt.Close()

		// *️⃣ ดึงข้อมูลจาก View ในแต่ละ batch
		err = nstmt.SelectContext(ctx, &productBatch, map[string]interface{}{
			"Keyword": searchParam,
			"Limit":   chunkSize,
			"Offset":  offset,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product data: %w", err)
		}

		// *️⃣ ถ้าไม่มีข้อมูลให้หยุดการดึงข้อมูล
		if len(productBatch) == 0 {
			break
		}

		// *️⃣ ถ้ามีข้อมูล เพิ่ม batch เข้าไปในผลลัพธ์
		products = append(products, productBatch...)
		// *️⃣ เพิ่ม offset เพื่อดึง batch ถัดไป
		offset += chunkSize
	}

	return products, nil
}
