package repository

import (
	entity "boilerplate-backend-go/Entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type Constants interface {
	SearchProvince(ctx context.Context, keyword string) ([]entity.Province, error)
	GetProvinces(ctx context.Context) ([]entity.Province, error)
	GetDistrict(ctx context.Context, provinceCode string) ([]entity.District, error)
	GetSubDistrict(ctx context.Context, districtCode string) ([]entity.SubDistrict, error)
	GetPostalCode(ctx context.Context, subdistrictCode string) ([]entity.PostalCode, error)

	GetProduct(ctx context.Context, offset, limit int) ([]entity.ROM_V_ProductAll, error) // รายการสินค้าแบบแบ่งรายการ
	GetWarehouse(ctx context.Context) ([]entity.Warehouse, error)                         // ชื่อคลัง + location

	SearchInvoiceNameByCustomerID(ctx context.Context, customerID string, keyword string, offset int, limit int) ([]entity.InvoiceInformation, error)
	GetCustomerID(ctx context.Context) ([]entity.InvoiceInformation, error)
	GetCustomerInfoByCustomerID(ctx context.Context, customerID string, limit, offset int) ([]entity.InvoiceInformation, error)

	SearchProduct(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.ROM_V_ProductAll, error) // ข้อมูลสินค้า
	SearchSKUByNameAndSize(ctx context.Context, nameAlias string, size string, offset int, limit int) ([]entity.ROM_V_ProductAll, error)
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

// GetProvinces ดึงข้อมูล Province ทั้งหมด
func (repo repositoryDB) GetProvinces(ctx context.Context) ([]entity.Province, error) {
	var provinces []entity.Province

	query := `SELECT DISTINCT ProvicesTH, ProvinceCode
              FROM Data_WebReturn.dbo.V_ThaiAddress
              ORDER BY ProvicesTH`

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var province entity.Province
		if err := rows.StructScan(&province); err != nil {
			return nil, fmt.Errorf("failed to scan province: %w", err)
		}
		provinces = append(provinces, province)
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

	var products []entity.ROM_V_ProductAll

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

// เมื่อผู้ใช้เลือก CustomerID ระบบจะดึงข้อมูล CustomerName, Address, TaxID ที่เกี่ยวข้องกับ CustomerID
// และค้นหาจาก CustomerName ด้วย LIKE
// เมื่อผู้ใช้เลือก InvoiceName ระบบจะดึง CustomerName, Address, TaxID ที่เกี่ยวข้องกับ InvoiceName
func (repo repositoryDB) SearchInvoiceNameByCustomerID(ctx context.Context, customerID string, keyword string, offset int, limit int) ([]entity.InvoiceInformation, error) {
	var invoices []entity.InvoiceInformation

	// ทำให้ customerID และ keyword เป็น lowercase และลบช่องว่าง
	cleanedCustomerID := strings.ToLower(strings.ReplaceAll(customerID, " ", ""))
	cleanedKeyword := strings.ToLower(strings.ReplaceAll(keyword, " ", ""))

	// เพิ่ม % เพื่อใช้กับ LIKE
	searchParam := cleanedKeyword + "%"

	// คิวรีสำหรับค้นหาโดยใช้ CustomerID และค้นหาจาก CustomerName ด้วย LIKE
	query := `SELECT DISTINCT CustomerName, Address, TaxID
			  FROM Data_WebReturn.dbo.InvoiceInformation
			  WHERE CustomerID = :CustomerID
			  AND REPLACE(LOWER(CustomerName), ' ', '') LIKE :Keyword
			  ORDER BY CustomerName
			  OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY`

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูลโดยใช้ OFFSET และ FETCH
	err = nstmt.SelectContext(ctx, &invoices, map[string]interface{}{
		"CustomerID": cleanedCustomerID,
		"Keyword":    searchParam,
		"Offset":     offset,
		"Limit":      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice data: %w", err)
	}

	if len(invoices) == 0 {
		return nil, sql.ErrNoRows
	}

	return invoices, nil
}

var customerCache = cache.New(10*time.Minute, 15*time.Minute)

func (repo repositoryDB) GetCustomerID(ctx context.Context) ([]entity.InvoiceInformation, error) {
	// ดึงจากแคชก่อน
	if cachedData, found := customerCache.Get("customer_list"); found {
		return cachedData.([]entity.InvoiceInformation), nil
	}

	customerID := []entity.InvoiceInformation{}
	query := `	SELECT DISTINCT CustomerID 
				FROM Data_WebReturn.dbo.InvoiceInformation 
				ORDER BY CustomerID
			 `

	err := repo.db.SelectContext(ctx, &customerID, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customerID: %w", err)
	}

	// บันทึกลงแคช
	customerCache.Set("customer_list", customerID, cache.DefaultExpiration)
	return customerID, nil
}

// ค้นหาข้อมูลจาก CustomerID เพื่อดึง CustomerName, Address, TaxID (รองรับแบ่งหน้า)
func (repo repositoryDB) GetCustomerInfoByCustomerID(ctx context.Context, customerID string, limit, offset int) ([]entity.InvoiceInformation, error) {
	var customers []entity.InvoiceInformation

	query := `
		SELECT DISTINCT CustomerID, CustomerName, Address, TaxID
		FROM Data_WebReturn.dbo.InvoiceInformation
		WHERE CustomerID = :CustomerID
		ORDER BY CustomerID
		OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
	`

	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &customers, map[string]interface{}{
		"CustomerID": customerID,
		"Limit":      limit,
		"Offset":     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer info: %w", err)
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
	searchParam := cleanedKeyword + "%"

	var query string

	// เลือกค้นหาตาม SKU หรือ NAMEALIAS
	if searchType == "SKU" {
		// ค้นหาจาก SKU
		query = `	SELECT SKU, NAMEALIAS, Size
					FROM Data_WebReturn.dbo.ROM_V_ProductAll
					WHERE REPLACE(LOWER(SKU), ' ', '') LIKE :Keyword
					ORDER BY SKU, NAMEALIAS, Size
					OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY
				 `
	} else if searchType == "NAMEALIAS" {
		// ค้นหาจาก NAMEALIAS
		query = `	SELECT DISTINCT NAMEALIAS, Size
					FROM Data_WebReturn.dbo.ROM_V_ProductAll
					WHERE REPLACE(LOWER(NAMEALIAS), ' ', '') LIKE :Keyword
					ORDER BY NAMEALIAS, Size
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

// SearchSKUByNameAndSize ค้นหา SKU ด้วย name และ size
func (repo repositoryDB) SearchSKUByNameAndSize(ctx context.Context, nameAlias string, size string, offset int, limit int) ([]entity.ROM_V_ProductAll, error) {
	var products []entity.ROM_V_ProductAll

	query := `SELECT SKU, NAMEALIAS, Size
              FROM Data_WebReturn.dbo.ROM_V_ProductAll
              WHERE NAMEALIAS = :Name
              AND Size = :Size
              ORDER BY SKU, NAMEALIAS, Size
              OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY`

	// เตรียมการคิวรี
	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	// ดึงข้อมูลจากฐานข้อมูลโดยใช้ OFFSET และ FETCH
	err = nstmt.SelectContext(ctx, &products, map[string]interface{}{
		"Name":   nameAlias,
		"Size":   size,
		"Offset": offset,
		"Limit":  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product data: %w", err)
	}

	if len(products) == 0 {
		return nil, sql.ErrNoRows
	}

	return products, nil
}
