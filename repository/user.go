package repository

import (
	response "boilerplate-backend-go/dto/response"
	"context"
	"fmt"
)

type UserRepository interface {
	GetUser(ctx context.Context, username, password string) (response.Login, error)
	GetUserFromLark(ctx context.Context, userID, username string) (response.Login, error)
	GetUserWithPermission(ctx context.Context, userid, username string) (response.UserPermission, error)
}

func (repo repositoryDB) GetUser(ctx context.Context, username, password string) (response.Login, error) {
	var user response.Login
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, 'web' as Platform
        FROM ROM_V_UserPermission
        WHERE UserName = :username AND Password = :password
    `
	params := map[string]interface{}{
		"username": username,
		"password": password,
	}

	nstmt, err := repo.db.PrepareNamed(query)
	if (err != nil) {
		return response.Login{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, &user, params)
	if (err != nil) {
		return response.Login{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (repo repositoryDB) GetUserFromLark(ctx context.Context, userid, username string) (response.Login, error) {
	var user response.Login
	user.UserName = username
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, 'lark' as Platform
        FROM ROM_V_UserPermission
        WHERE UserID = :userID AND UserName = :userName
    `
	params := map[string]interface{}{
		"userID":   userid,
		"userName": username,
	}

	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return response.Login{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, &user, params)
	if err != nil {
		return response.Login{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (repo repositoryDB) GetUserWithPermission(ctx context.Context, userid, username string) (response.UserPermission, error) {
	var user response.UserPermission
	query := `
        SELECT UserID, UserName, RoleID, FullNameTH, NickName, DepartmentNo, RoleName, Description, Permission 
        FROM ROM_V_UserPermission
        WHERE UserID = :userid AND UserName = :username
    `
	params := map[string]interface{}{
		"userid":   userid,
		"username": username,
	}

	nstmt, err := repo.db.PrepareNamed(query)
	if err != nil {
		return response.UserPermission{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.GetContext(ctx, &user, params)
	if err != nil {
		return response.UserPermission{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
