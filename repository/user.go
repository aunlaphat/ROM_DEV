package repository

import (
	"boilerplate-backend-go/dto/response"
	"context"
	"fmt"
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (response.User, error)
	GetUserFromLark(ctx context.Context, userID, username string) (response.User, error)
	GetUserWithPermission(ctx context.Context, userID, username string) (response.UserPermission, error)
}

func (repo repositoryDB) GetUser(ctx context.Context, username string) (response.User, error) {
	var user response.User
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, 'web' as Platform
        FROM ROM_V_UserPermission
        WHERE UserName = :username
    `
	params := map[string]interface{}{"username": username}

	err := repo.db.GetContext(ctx, &user, query, params)
	if err != nil {
		return response.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (repo repositoryDB) GetUserFromLark(ctx context.Context, userID, username string) (response.User, error) {
	var user response.User
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, 'lark' as Platform
        FROM ROM_V_UserPermission
        WHERE UserID = :userID AND UserName = :userName
    `
	params := map[string]interface{}{
		"userID":   userID,
		"userName": username,
	}

	err := repo.db.GetContext(ctx, &user, query, params)
	if err != nil {
		return response.User{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}

	return user, nil
}

func (repo repositoryDB) GetUserWithPermission(ctx context.Context, userID, username string) (response.UserPermission, error) {
	var user response.UserPermission
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, RoleName, Description, Permission 
        FROM ROM_V_UserPermission
        WHERE UserID = :userID AND UserName = :userName
    `
	params := map[string]interface{}{
		"userID":   userID,
		"userName": username,
	}

	err := repo.db.GetContext(ctx, &user, query, params)
	if err != nil {
		return response.UserPermission{}, fmt.Errorf("failed to get user permission: %w", err)
	}

	return user, nil
}
