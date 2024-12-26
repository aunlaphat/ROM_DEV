package repository

import (
	entity "boilerplate-backend-go/Entity"
	"context"
	"time"
)

type Constants interface {
	GetThaiProvince() ([]entity.Province, error)
	GetThaiDistrict() ([]entity.District, error)
	GetThaiSubDistrict() ([]entity.SubDistrict, error)
	GetProductAll() ([]entity.ROM_V_ProductAll, error)
	// GetCustomer() ([]entity.SubDistrict, error)

}

func (repo repositoryDB) GetThaiProvince() ([]entity.Province, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provinces := []entity.Province{}

	sqlQuery := `SELECT Code
					,NameTH
					,NameEN
				FROM V_ThaiAddressProvince
				ORDER BY Code`

	rows, err := repo.db.QueryxContext(ctx, sqlQuery)
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

func (repo repositoryDB) GetThaiDistrict() ([]entity.District, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	districts := []entity.District{}

	sqlQuery := `SELECT ProvinceCode
					,Code
					,NameTH
					,NameEN
				FROM V_ThaiAddressDistrict
				ORDER BY Code`

	rows, err := repo.db.QueryxContext(ctx, sqlQuery)
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

func (repo repositoryDB) GetThaiSubDistrict() ([]entity.SubDistrict, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subDistricts := []entity.SubDistrict{}

	sqlQuery := `SELECT Code
					,DistrictCode
					,ZipCode
					,NameTH
					,NameEN
				FROM V_ThaiAddressSubDistrict
				ORDER BY Code`

	rows, err := repo.db.QueryxContext(ctx, sqlQuery)
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

func (repo repositoryDB) GetProductAll() ([]entity.ROM_V_ProductAll, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    products := []entity.ROM_V_ProductAll{}

    query := `
        SELECT SKU, NAMEALIAS, Size, SizeID, Barcode, Type
        FROM Data_WebReturn.dbo.ROM_V_ProductAll
        ORDER BY SKU
    `

    rows, err := repo.db.QueryxContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var product entity.ROM_V_ProductAll
        if err := rows.StructScan(&product); err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil

}

// func (repo repositoryDB) GetCustomer() ([]entity.SubDistrict, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()


// }
