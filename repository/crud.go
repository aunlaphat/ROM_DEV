package repository

import (
	entity "boilerplate-backend-go/Entity"
	"context"
	"database/sql"
	"time"
)

func (repo repositoryDB) SelectData() ([]entity.Province, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provinces := []entity.Province{}

	sqlQuery := `SELECT Code
					,NameTH
					,NameEN
				FROM V_ThaiAddressProvince
				ORDER BY Code`

	err := repo.db.SelectContext(ctx, &provinces, sqlQuery)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func (repo repositoryDB) SelectOneData() (entity.Province, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provinces := entity.Province{}

	sqlQuery := `SELECT Code
					,NameTH
					,NameEN
				FROM V_ThaiAddressProvince
				Where Code = @Code
				ORDER BY Code`

	err := repo.db.GetContext(ctx, &provinces, sqlQuery,
		sql.Named("Code", "10"),
	)
	if err != nil {
		return provinces, err
	}

	return provinces, nil
}
