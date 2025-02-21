package repository

import (
	entity "boilerplate-backend-go/Entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Constants interface {
	SearchProvince(ctx context.Context, keyword string) ([]entity.Province, error)
	GetDistrict(ctx context.Context, provinceCode string) ([]entity.District, error)
	GetSubDistrict(ctx context.Context, districtCode string) ([]entity.SubDistrict, error)
	GetPostalCode(ctx context.Context, subdistrictCode string) ([]entity.PostalCode, error)
	
	GetProduct(ctx context.Context, offset, limit int) ([]entity.ROM_V_ProductAll, error)                                              // รายการสินค้าแบบแบ่งรายการ
	GetWarehouse(ctx context.Context) ([]entity.Warehouse, error)                                                                      // ชื่อคลัง + location
	SearchCustomer(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.InvoiceInformation, error) // ข้อมูลลูกค้า + invoice
	SearchProduct(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.ROM_V_ProductAll, error)    // ข้อมูลสินค้า
}

// ค้นหาในฐานข้อมูล Province ที่ตรงกับคำค้นหาของผู้ใช้ และคืนค่าผลลัพธ์เป็นจังหวัดที่มีอยู่ทั้งหมด
func (repo repositoryDB) SearchProvince(ctx context.Context, keyword string) ([]entity.Province, error) {
	var provinces []entity.Province
	// ทำให้ keyword เป็น lowercase และลบช่องว่าง
	cleanedKeyword := strings.ToLower(strings.ReplaceAll(keyword, " ", ""))
	// เพิ่ม % เพื่อใช้กับ LIKE
	searchParam := cleanedKeyword + "%"

	query := `	SELECT DISTINCT ProvinceCode, ProvicesTH
				FROM Data_WebReturn.dbo.V_ThaiAddress
				WHERE LOWER(ProvicesTH) LIKE :Keyword
			  `

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูล
	err = nstmt.SelectContext(ctx, &provinces, map[string]interface{}{"Keyword": searchParam})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch province data: %w", err)
	}

	if len(provinces) == 0 {
		return nil, sql.ErrNoRows
	}

	return provinces, nil
}

// หลังจากเลือกจังหวัดแล้ว ระบบจะใช้ ProvinceCode เพื่อดึงข้อมูล District ที่สัมพันธ์กับจังหวัดนั้น
func (repo repositoryDB) GetDistrict(ctx context.Context, provinceCode string) ([]entity.District, error) {
	var districts []entity.District

	query := `	SELECT DISTINCT ProvinceCode, DistrictCode, DistrictTH
				FROM Data_WebReturn.dbo.V_ThaiAddress
				WHERE ProvinceCode = :ProvinceCode
				ORDER BY DistrictTH
			  `

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูล
	err = nstmt.SelectContext(ctx, &districts, map[string]interface{}{"ProvinceCode": provinceCode})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district data: %w", err)
	}

	return districts, nil
}

// หลังจากเลือก District ระบบจะจำกัดตัวเลือก Subdistrict ตาม DistrictCode
func (repo repositoryDB) GetSubDistrict(ctx context.Context, districtCode string) ([]entity.SubDistrict, error) {
	var subdistricts []entity.SubDistrict

	query := `	SELECT DISTINCT DistrictCode, SubdistrictCode, SubdistrictTH
				FROM Data_WebReturn.dbo.V_ThaiAddress
				WHERE DistrictCode = :DistrictCode
			  `

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูล
	err = nstmt.SelectContext(ctx, &subdistricts, map[string]interface{}{"DistrictCode": districtCode})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subdistrict data: %w", err)
	}

	return subdistricts, nil
}

// เมื่อเลือก Subdistrict ระบบจะแสดง PostalCode ที่ตรงกับ SubdistrictCode
func (repo repositoryDB) GetPostalCode(ctx context.Context, subdistrictCode string) ([]entity.PostalCode, error) {
	var postalCodes []entity.PostalCode

	query := `	SELECT DISTINCT SubdistrictCode, ZipCode
				FROM Data_WebReturn.dbo.V_ThaiAddress
				WHERE SubdistrictCode = :SubdistrictCode
			  `

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูล
	err = nstmt.SelectContext(ctx, &postalCodes, map[string]interface{}{"SubdistrictCode": subdistrictCode})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch postal code data: %w", err)
	}

	return postalCodes, nil
}

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

// เมื่อผู้ใช้เลือก CustomerID จาก dropdown ระบบจะดึงข้อมูล CustomerName, Address, TaxID ที่เกี่ยวข้องกับ CustomerID
// เมื่อผู้ใช้เลือก InvoiceName ระบบจะดึง CustomerName, Address, TaxID ที่เกี่ยวข้องกับ InvoiceName
func (repo repositoryDB) SearchCustomer(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.InvoiceInformation, error) {
	var customers []entity.InvoiceInformation

	// ทำให้ keyword เป็น lowercase และลบช่องว่าง
	cleanedKeyword := strings.ToLower(strings.ReplaceAll(keyword, " ", ""))

	// เพิ่ม % เพื่อใช้กับ LIKE
	searchParam := cleanedKeyword + "%"

	var query string

	// เลือกค้นหาตาม searchType
	if searchType == "CustomerID" {
		// ค้นหาจาก CustomerID
		query = `	SELECT CustomerID, CustomerName, Address, TaxID
					FROM Data_WebReturn.dbo.InvoiceInformation
					WHERE REPLACE(LOWER(CustomerID), ' ', '') LIKE :Keyword
					ORDER BY CustomerID
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				`
	} else if searchType == "InvoiceName" {
		// ค้นหาจาก InvoiceName
		query = `	SELECT CustomerName, Address, TaxID
					FROM Data_WebReturn.dbo.InvoiceInformation
					WHERE REPLACE(LOWER(CustomerName), ' ', '') LIKE :Keyword
					ORDER BY CustomerName
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				`
	} else {
		return nil, fmt.Errorf("invalid search type")
	}

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูลโดยใช้ OFFSET และ FETCH
	err = nstmt.SelectContext(ctx, &customers, map[string]interface{}{
		"Keyword": searchParam,
		"Offset":  offset,
		"Limit":   limit,
	})
	if err != nil {

		return nil, fmt.Errorf("failed to fetch customer data: %w", err)
	}

	if len(customers) == 0 {
		return nil, sql.ErrNoRows
	}

	return customers, nil
}

// เมื่อเลือก SKU ระบบจะแสดง NAMEALIAS ที่ตรงกับ SKU
// เมื่อเลือก NAMEALIAS ระบบจะแสดง SKU ที่ตรงกับ NAMEALIAS
func (repo repositoryDB) SearchProduct(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.ROM_V_ProductAll, error) {
	var products []entity.ROM_V_ProductAll

	// ทำให้ keyword เป็น lowercase และลบช่องว่าง
	cleanedKeyword := strings.ToLower(strings.ReplaceAll(keyword, " ", ""))

	// เพิ่ม % เพื่อใช้กับ LIKE
	searchParam := "%" + cleanedKeyword + "%"

	var query string

	// เลือกค้นหาตาม SKU หรือ NAMEALIAS
	if searchType == "SKU" {
		// ค้นหาจาก SKU
		query = `	SELECT SKU, NAMEALIAS, Size, SizeID, Barcode, Type
					FROM Data_WebReturn.dbo.ROM_V_ProductAll
					WHERE REPLACE(LOWER(SKU), ' ', '') LIKE :Keyword
					ORDER BY SKU
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				 `
	} else if searchType == "NAMEALIAS" {
		// ค้นหาจาก NAMEALIAS
		query = `	SELECT SKU, NAMEALIAS, Size, SizeID, Barcode, Type
					FROM Data_WebReturn.dbo.ROM_V_ProductAll
					WHERE REPLACE(LOWER(NAMEALIAS), ' ', '') LIKE :Keyword
					ORDER BY NAMEALIAS
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				 `
	} else {
		return nil, fmt.Errorf("invalid search type")
	}

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูลโดยใช้ OFFSET และ FETCH
	err = nstmt.SelectContext(ctx, &products, map[string]interface{}{
		"Keyword": searchParam,
		"Offset":  offset,
		"Limit":   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product data: %w", err)
	}

	if len(products) == 0 {
		return nil, sql.ErrNoRows
	}

	return products, nil
}
