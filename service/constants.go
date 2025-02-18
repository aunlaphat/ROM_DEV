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
	GetThaiProvince(ctx context.Context) ([]entity.Province, error)
	GetThaiDistrict(ctx context.Context) ([]entity.District, error)
	GetThaiSubDistrict(ctx context.Context) ([]entity.SubDistrict, error)
	// GetPostCode(ctx context.Context) ([]entity.PostCode, error)
	GetWarehouse(ctx context.Context) ([]entity.Warehouse, error)
	GetProduct(ctx context.Context, page, limit int) ([]entity.ROM_V_ProductAll, error)
	//GetCustomer(ctx context.Context) ([]entity.ROM_V_Customer, error)

}

func (srv service) GetThaiProvince(ctx context.Context) ([]entity.Province, error) {
    getProvince, err := srv.constant.GetThaiProvince(ctx)
    if err != nil {
        if err == sql.ErrNoRows {
            srv.logger.Warn("[ data not found ]", zap.Error(err))
            return nil, errors.ValidationError("[ no province data: %v ]", err)
        }
        srv.logger.Error("[ get province error ]", zap.Error(err))
        return nil, errors.InternalError("[ get province error: %v ]", err)
    }
    return getProvince, nil
}

func (srv service) GetThaiDistrict(ctx context.Context) ([]entity.District, error) {
    getDistrict, err := srv.constant.GetThaiDistrict(ctx)
    if err != nil {
        if err == sql.ErrNoRows {
            srv.logger.Warn("[ data not found ]", zap.Error(err))
            return nil, errors.ValidationError("[ no district data: %v ]", err)
        }
        srv.logger.Error("[ get district error ]", zap.Error(err))
        return nil, errors.InternalError("[ get district error: %v ]", err)
    }
    return getDistrict, nil
}

func (srv service) GetThaiSubDistrict(ctx context.Context) ([]entity.SubDistrict, error) {
    getSubDistrict, err := srv.constant.GetThaiSubDistrict(ctx)
    if err != nil {
        if err == sql.ErrNoRows {
            srv.logger.Warn("[ data not found ]", zap.Error(err))
            return nil, errors.ValidationError("[ no sub district data: %v ]", err)
        }
        srv.logger.Error("[ get sub district error ]", zap.Error(err))
        return nil, errors.InternalError("[ get sub district error: %v ]", err)
    }
    return getSubDistrict, nil
}

// func (srv service) GetPostCode(ctx context.Context) ([]entity.PostCode, error) {
// 	getPostCode, err := srv.constant.GetPostCode(ctx)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			srv.logger.Warn("[  data not found ]", zap.Error(err))
// 			return nil, errors.ValidationError("[ no post code data: %v ]", err)
// 		}
// 		srv.logger.Error("[  get post code error ]", zap.Error(err))
// 		return nil, errors.InternalError("[ get post code error: %v ]", err)
// 	}
// 	return getPostCode, nil
// }

func (srv service) GetWarehouse(ctx context.Context) ([]entity.Warehouse, error) {
    getWarehouse, err := srv.constant.GetWarehouse(ctx)
    if err != nil {
        if err == sql.ErrNoRows {
            srv.logger.Warn("[  data not found ]", zap.Error(err))
            return nil, errors.ValidationError("[ no warehouse data: %v ]", err)
        }
        srv.logger.Error("[  get warehouse error ]", zap.Error(err))
        return nil, errors.InternalError("[ get warehouse error: %v ]", err)
    }
    return getWarehouse, nil
}

func (srv service) GetProduct(ctx context.Context, page, limit int) ([]entity.ROM_V_ProductAll, error) {
    offset := (page - 1) * limit

    products, err := srv.constant.GetProduct(ctx, offset, limit)
    if err != nil {
        return nil, err
    }

    return products, nil
}

// func (srv service) GetCustomer(ctx context.Context) ([]entity.ROM_V_Customer, error) {
// 	getCustomer, err := srv.constant.GetCustomer(ctx)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			srv.logger.Warn("[  data not found ]", zap.Error(err))
// 			return nil, errors.ValidationError("[ no customer data: %v ]", err)
// 		}
// 		srv.logger.Error("[  get customer error ]", zap.Error(err))
// 		return nil, errors.InternalError("[ get customer error: %v ]", err)
// 	}
// 	return getCustomer, nil
// }