package service

import (
	entity "boilerplate-backend-go/Entity"
	"boilerplate-backend-go/errors"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

// ส่วนรับ dropdown
type Constants interface {
	SearchProvince(ctx context.Context, keyword string) ([]entity.Province, error)
	GetDistrict(ctx context.Context, provinceCode string) ([]entity.District, error)
	GetSubDistrict(ctx context.Context, districtCode string) ([]entity.SubDistrict, error)
	GetPostalCode(ctx context.Context, subdistrictCode string) ([]entity.PostalCode, error)

	GetWarehouse(ctx context.Context) ([]entity.Warehouse, error)
	GetProduct(ctx context.Context, page, limit int) ([]entity.ROM_V_ProductAll, error)
	SearchCustomer(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.InvoiceInformation, error)
	SearchProduct(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.ROM_V_ProductAll, error)
}

// Service Method ที่ค้นหาจังหวัด (Province)
func (srv service) SearchProvince(ctx context.Context, keyword string) ([]entity.Province, error) {
	provinces, err := srv.constant.SearchProvince(ctx, keyword)
	if err != nil {
		if err == sql.ErrNoRows {
			srv.logger.Warn("[  Data not found ]", zap.Error(err))
			return nil, sql.ErrNoRows
		}
		srv.logger.Error("[ Failed to fetch province data ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch province data: %v ]", err)
	}
	return provinces, nil
}

// Service Method ที่ค้นหาจังหวัดตาม ProvinceCode
func (srv service) GetDistrict(ctx context.Context, provinceCode string) ([]entity.District, error) {
	districts, err := srv.constant.GetDistrict(ctx, provinceCode)
	if err != nil {
		srv.logger.Error("[ Failed to fetch district data ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch district data: %v ]", err)
	}
	return districts, nil
}

// Service Method ที่ค้นหาตำบล (Subdistrict) ตาม DistrictCode
func (srv service) GetSubDistrict(ctx context.Context, districtCode string) ([]entity.SubDistrict, error) {
	subdistricts, err := srv.constant.GetSubDistrict(ctx, districtCode)
	if err != nil {
		srv.logger.Error("[ Failed to fetch subdistrict data ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch subdistrict data: %v ]", err)
	}
	return subdistricts, nil
}

// Service Method ที่ค้นหาข้อมูลรหัสไปรษณีย์ (PostalCode) ตาม SubdistrictCode
func (srv service) GetPostalCode(ctx context.Context, subdistrictCode string) ([]entity.PostalCode, error) {
	postalCodes, err := srv.constant.GetPostalCode(ctx, subdistrictCode)
	if err != nil {
		srv.logger.Error("[ Failed to fetch postal code data ]", zap.Error(err))
		return nil, errors.InternalError("[ Failed to fetch postal code data: %v ]", err)
	}
	return postalCodes, nil
}

func (srv service) GetWarehouse(ctx context.Context) ([]entity.Warehouse, error) {
	getWarehouse, err := srv.constant.GetWarehouse(ctx)
	if err != nil {
		srv.logger.Error("[  get warehouse error ]", zap.Error(err))
		return nil, errors.InternalError("[ get warehouse error: %v ]", err)
	}
	return getWarehouse, nil
}

func (srv service) GetProduct(ctx context.Context, page, limit int) ([]entity.ROM_V_ProductAll, error) {
	offset := (page - 1) * limit

	getProducts, err := srv.constant.GetProduct(ctx, offset, limit)
	if err != nil {
		srv.logger.Error("[  get product error ]", zap.Error(err))
		return nil, errors.InternalError("[ get product error: %v ]", err)
	}

	return getProducts, nil
}

func (srv service) SearchCustomer(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.InvoiceInformation, error) {
	getCustomer, err := srv.constant.SearchCustomer(ctx, keyword, searchType, offset, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			srv.logger.Warn("[  Data not found ]", zap.Error(err))
			return nil, sql.ErrNoRows
		}
		srv.logger.Error("[  get customer error ]", zap.Error(err))
		return nil, errors.InternalError("[ get customer error: %v ]", err)
	}

	return getCustomer, nil
}

func (srv service) SearchProduct(ctx context.Context, keyword string, searchType string, offset int, limit int) ([]entity.ROM_V_ProductAll, error) {
	getProducts, err := srv.constant.SearchProduct(ctx, keyword, searchType, offset, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			srv.logger.Warn("[  Data not found ]", zap.Error(err))
			return nil, sql.ErrNoRows
		}
		srv.logger.Error("[  search product error ]", zap.Error(err))
		return nil, errors.InternalError("[ search product error: %v ]", err)
	}

	return getProducts, nil
}
