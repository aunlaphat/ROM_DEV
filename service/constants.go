package service

import (
	entity "boilerplate-backend-go/Entity"
	"database/sql"
	"fmt"
)

type Constants interface {
	GetThaiProvince() ([]entity.Province, error)
	GetThaiDistrict() ([]entity.District, error)
	GetThaiSubDistrict() ([]entity.SubDistrict, error)
	GetProductAll() ([]entity.ROM_V_ProductAll, error)
	//GetCustomer() ([]entity.ROM_V_Customer, error)

}

func (srv service) GetThaiProvince() ([]entity.Province, error) {
	getProvince, err := srv.constant.GetThaiProvince()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no province data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get province error: %w", err)
		}
	}

	return getProvince, nil
}

func (srv service) GetThaiDistrict() ([]entity.District, error) {
	getDistrict, err := srv.constant.GetThaiDistrict()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no district data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get district error: %w", err)
		}
	}

	return getDistrict, nil
}

func (srv service) GetThaiSubDistrict() ([]entity.SubDistrict, error) {
	getSubDistrict, err := srv.constant.GetThaiSubDistrict()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no sub district data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get sub district error: %w", err)
		}
	}

	return getSubDistrict, nil
}

func (srv service) GetProductAll() ([]entity.ROM_V_ProductAll, error) {
	getProductAll, err := srv.constant.GetProductAll()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no product data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get product error: %w", err)
		}
	}

	return getProductAll, nil
}

// func (srv service) GetCustomer() ([]entity.ROM_V_Customer, error) {
// 	getCustomer, err := srv.constant.GetCustomer()
// 	if err != nil {
// 		switch err {
// 		case sql.ErrNoRows:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("no customer data: %w", err)
// 		default:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("get customer error: %w", err)
// 		}
// 	}

// 	return getCustomer, nil
// }
