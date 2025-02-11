package repository

import (
	entity "boilerplate-backend-go/Entity"
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (entity.User, error)
	GetUserFromLark(ctx context.Context, userID, username string) (response.User, error)
}

func (repo repositoryDB) GetUser(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	query := `
        SELECT UserID, UserName, Password, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, Permission
        FROM ROM_V_UserPermission
        WHERE UserName = :username
    `

	params := map[string]interface{}{"username": username}

	stmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &user, params)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (repo repositoryDB) GetUserFromLark(ctx context.Context, userID, username string) (response.User, error) {
	var user response.User
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo
        FROM ROM_V_UserPermission
        WHERE UserID = :userID AND UserName = :userName
    `
	params := map[string]interface{}{
		"userID":   userID,
		"userName": username,
	}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return response.User{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return response.User{}, fmt.Errorf("failed to scan user from Lark: %w", err)
		}
		return user, nil
	}

	return response.User{}, fmt.Errorf("user not found in Lark")
}
