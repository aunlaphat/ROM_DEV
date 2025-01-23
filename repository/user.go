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
	AddUser(ctx context.Context, user response.UserPermission) error
	// EditUser(ctx context.Context, user response.User) error
	DeleteUser(ctx context.Context, userID string) error
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

// AddUser adds a new user to the database
func (repo repositoryDB) AddUser(ctx context.Context, user response.UserPermission) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert into ROM_V_User
	queryUser := `
        INSERT INTO ROM_V_User (UserID, Password, UserName, NickName, FullNameTH, DepartmentNo)
        VALUES (:userID, :password, :userName, :nickName, :fullNameTH, :departmentNo)
    `
	paramsUser := map[string]interface{}{
		"userID":       user.UserID,
		"password":     user.Password,
		"userName":     user.UserName,
		"nickName":     user.NickName,
		"fullNameTH":   user.FullNameTH,
		"departmentNo": user.DepartmentNo,
	}

	_, err = tx.NamedExecContext(ctx, queryUser, paramsUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add user: %w", err)
	}

	// Insert into UserRole
	queryRole := `
        INSERT INTO UserRole (UserID, RoleID)
        VALUES (:userID, :roleID)
    `
	paramsRole := map[string]interface{}{
		"userID": user.UserID,
		"roleID": user.RoleID,
	}

	_, err = tx.NamedExecContext(ctx, queryRole, paramsRole)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add user role: %w", err)
	}

	// Insert into ROM_V_UserPermission
	queryPermission := `
        INSERT INTO ROM_V_UserPermission (UserID, UserName, Password, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, Permission)
        SELECT UserID, UserName, Password, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, Permission
        FROM ROM_V_User
        WHERE UserID = :userID
    `
	paramsPermission := map[string]interface{}{
		"userID": user.UserID,
	}

	_, err = tx.NamedExecContext(ctx, queryPermission, paramsPermission)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add user permission: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/*
// EditUser updates the role and warehouse of an existing user
func (repo repositoryDB) EditUser(ctx context.Context, user response.User) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Update UserRole
	queryRole := `
        UPDATE UserRole
        SET RoleID = :roleID
        WHERE UserID = :userID
    `
	paramsRole := map[string]interface{}{
		"userID": user.UserID,
		"roleID": user.RoleID,
	}

	_, err = tx.NamedExecContext(ctx, queryRole, paramsRole)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to edit user role: %w", err)
	}

	// Update ROM_V_UserPermission
	queryPermission := `
        UPDATE ROM_V_UserPermission
        SET RoleID = :roleID, WarehouseID = :warehouseID
        WHERE UserID = :userID
    `
	paramsPermission := map[string]interface{}{
		"userID":      user.UserID,
		"roleID":      user.RoleID,
		"warehouseID": user.WarehouseID,
	}

	_, err = tx.NamedExecContext(ctx, queryPermission, paramsPermission)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to edit user permission: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
} */

// DeleteUser removes a user from the database
func (repo repositoryDB) DeleteUser(ctx context.Context, userID string) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Delete from UserRole
	queryRole := `
        DELETE FROM UserRole
        WHERE UserID = :userID
    `
	paramsRole := map[string]interface{}{
		"userID": userID,
	}

	_, err = tx.NamedExecContext(ctx, queryRole, paramsRole)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user role: %w", err)
	}

	// Delete from ROM_V_UserPermission
	queryPermission := `
        DELETE FROM ROM_V_UserPermission
        WHERE UserID = :userID
    `
	paramsPermission := map[string]interface{}{
		"userID": userID,
	}

	_, err = tx.NamedExecContext(ctx, queryPermission, paramsPermission)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user permission: %w", err)
	}

	// Delete from ROM_V_User
	queryUser := `
        DELETE FROM ROM_V_User
        WHERE UserID = :userID
    `
	paramsUser := map[string]interface{}{
		"userID": userID,
	}

	_, err = tx.NamedExecContext(ctx, queryUser, paramsUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
