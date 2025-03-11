package repository

import (
	"context"
	"fmt"

	"boilerplate-back-go-2411/dto/request"
	"boilerplate-back-go-2411/entity"
	"boilerplate-back-go-2411/errors"
)

type UserRepository interface {
	Login(ctx context.Context, userName string) (entity.ROM_V_UserDetail, error)
	LoginLark(ctx context.Context, userID, userName string) (entity.ROM_V_UserDetail, error)
	GetUser(ctx context.Context, userID string) (entity.ROM_V_UserDetail, error)
	GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]entity.ROM_V_UserDetail, error)
	CheckUserExistsInERP(ctx context.Context, userID string) (bool, error)
	CheckUserExists(ctx context.Context, userID string) (bool, error)
	AddUser(ctx context.Context, req request.AddUserRequest, adminID string) error
	EditUser(ctx context.Context, req request.EditUserRequest, adminID string) error
	DeleteUser(ctx context.Context, userID, adminID string) error
}

func (repo repositoryDB) Login(ctx context.Context, userName string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, Password, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName
        FROM ROM_V_UserDetail
        WHERE UserName = :userName
    `
	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"userName": userName,
	})
	if err != nil {
		return entity.ROM_V_UserDetail{}, fmt.Errorf("failed to execute login query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.ROM_V_UserDetail{}, fmt.Errorf("failed to scan user: %w", err)
		}
		return user, nil
	}
	return entity.ROM_V_UserDetail{}, errors.NotFoundError("user not found")
}

func (repo repositoryDB) LoginLark(ctx context.Context, userID, userName string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, Password, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName
        FROM ROM_V_UserDetail
        WHERE UserID = :userID AND UserName = :userName
    `
	params := map[string]interface{}{
		"userID":   userID,
		"userName": userName,
	}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return entity.ROM_V_UserDetail{}, fmt.Errorf("failed to get user from Lark: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.ROM_V_UserDetail{}, fmt.Errorf("failed to scan user from Lark: %w", err)
		}
		return user, nil
	}

	return entity.ROM_V_UserDetail{}, fmt.Errorf("user not found in Lark")
}

func (repo repositoryDB) GetUser(ctx context.Context, userID string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, WarehouseID, WarehouseName, Description, IsActive
        FROM ROM_V_UserDetail
        WHERE UserID = :userID
    `

	params := map[string]interface{}{"userID": userID}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return entity.ROM_V_UserDetail{}, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return entity.ROM_V_UserDetail{}, fmt.Errorf("failed to scan user: %w", err)
		}
	} else {
		return entity.ROM_V_UserDetail{}, errors.NotFoundError("user not found in database")
	}

	return user, nil
}

func (repo repositoryDB) GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]entity.ROM_V_UserDetail, error) {
	var users []entity.ROM_V_UserDetail
	query := `
		SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, WarehouseID, WarehouseName, Description, IsActive
		FROM ROM_V_UserDetail
		WHERE IsActive = :isActive
		ORDER BY UserID
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
	`
	params := map[string]interface{}{
		"isActive": isActive,
		"limit":    limit,
		"offset":   offset,
	}

	rows, err := repo.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.ROM_V_UserDetail
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo repositoryDB) CheckUserExistsInERP(ctx context.Context, userID string) (bool, error) {
	query := `SELECT COUNT(1) AS Count FROM ROM_V_User WHERE UserID = :userID`
	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{"userID": userID})
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	return count > 0, nil
}

func (repo repositoryDB) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT COUNT(1) AS Count FROM UserRole WHERE UserID = :userID`

	rows, err := repo.db.NamedQueryContext(ctx, query, map[string]interface{}{"userID": userID})
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	return count > 0, nil
}

func (repo repositoryDB) AddUser(ctx context.Context, req request.AddUserRequest, adminID string) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	queryUserRole := `
		INSERT INTO UserRole (UserID, RoleID, WarehouseID, CreatedBy, CreatedAt)
		VALUES (:userID, :roleID, :warehouseID, :createdBy, GETDATE())
	`
	_, err = tx.NamedExecContext(ctx, queryUserRole, map[string]interface{}{
		"userID":      req.UserID,
		"roleID":      req.RoleID,
		"warehouseID": req.WarehouseID,
		"createdBy":   adminID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert into UserRole: %w", err)
	}

	queryUserStatus := `
		INSERT INTO UserStatus (UserID, IsActive, CreatedBy, CreatedAt)
		VALUES (:userID, 1, :createdBy, GETDATE())
	`
	_, err = tx.NamedExecContext(ctx, queryUserStatus, map[string]interface{}{
		"userID":    req.UserID,
		"createdBy": adminID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert into UserStatus: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) EditUser(ctx context.Context, req request.EditUserRequest, adminID string) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	exists, err := repo.CheckUserExists(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errors.NotFoundError("user not found")
	}

	query := `UPDATE UserRole SET UpdatedBy = :updatedBy, UpdatedAt = GETDATE()`
	params := map[string]interface{}{
		"userID":    req.UserID,
		"updatedBy": adminID,
	}

	if req.RoleID != nil {
		query += `, RoleID = :roleID`
		params["roleID"] = *req.RoleID
	}

	if req.WarehouseID != nil {
		query += `, WarehouseID = :warehouseID`
		params["warehouseID"] = *req.WarehouseID
	}

	query += ` WHERE UserID = :userID`

	_, err = tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update UserRole: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo repositoryDB) DeleteUser(ctx context.Context, userID, adminID string) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	exists, err := repo.CheckUserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errors.NotFoundError("user not found")
	}

	query := `
		UPDATE UserStatus
		SET IsActive = 0, DeactivatedBy = :deactivatedBy, DeactivatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":        userID,
		"deactivatedBy": adminID,
	}

	_, err = tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update UserStatus: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
