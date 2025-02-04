package service

type Constants interface {
	/* GetThaiProvince() ([]entity.Province, error)
	GetThaiDistrict() ([]entity.District, error)
	GetThaiSubDistrict() ([]entity.SubDistrict, error)
	// GetPostCode() ([]entity.PostCode, error)
	GetWarehouse() ([]entity.Warehouse, error)
	GetProductAll() ([]entity.ROM_V_ProductAll, error)
	GetProductAllWithPagination(page, limit int) ([]entity.ROM_V_ProductAll, error) */
	//GetCustomer() ([]entity.ROM_V_Customer, error)

}

/* func (srv service) GetThaiProvince() ([]entity.Province, error) {

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
} */

/* func (srv service) GetThaiDistrict() ([]entity.District, error) {

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
} */

/* func (srv service) GetThaiSubDistrict() ([]entity.SubDistrict, error) {
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
} */

// func (srv service) GetPostCode() ([]entity.PostCode, error) {

// 	getPostCode, err := srv.constant.GetPostCode()
// 	if err != nil {
// 		switch err {
// 		case sql.ErrNoRows:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("no post code data: %w", err)
// 		default:
// 			srv.logger.Error(err)
// 			return nil, fmt.Errorf("get post code error: %w", err)
// 		}
// 	}

// 	return getPostCode, nil
// }

/* func (srv service) GetWarehouse() ([]entity.Warehouse, error) {

	getWarehouse, err := srv.constant.GetWarehouse()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			srv.logger.Error(err)
			return nil, fmt.Errorf("no warehouse data: %w", err)
		default:
			srv.logger.Error(err)
			return nil, fmt.Errorf("get warehouse error: %w", err)
		}
	}

	return getWarehouse, nil
} */

/* func (srv service) GetProductAll() ([]entity.ROM_V_ProductAll, error) {

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
} */

/* func (srv service) GetProductAllWithPagination(page, limit int) ([]entity.ROM_V_ProductAll, error) {

	offset := (page - 1) * limit

	products, err := srv.constant.GetProductAllWithPagination(context.Background(), offset, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
} */

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
