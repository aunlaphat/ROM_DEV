package repository

import (
	"context"
	"fmt"

	"boilerplate-backend-go/dto/request"
	"boilerplate-backend-go/entity"
	"boilerplate-backend-go/errors"
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
	GetCurrentPassword(ctx context.Context, userName string) (string, error)
	//UpdateUserPassword(ctx context.Context, req request.ResetPasswordRequest, adminID string) error
}

// ✅ **1️⃣ Login - ตรวจสอบข้อมูล User**
func (repo repositoryDB) Login(ctx context.Context, userName string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, Password, UserName, NickName, FullNameTH, DepartmentNo, RoleID
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
        SELECT UserID, Password, UserName, NickName, FullNameTH, DepartmentNo, RoleID
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

// ✅ **2️⃣ GetUser - ดึงข้อมูลผู้ใช้จาก View**
func (repo repositoryDB) GetUser(ctx context.Context, userID string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive
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

// ✅ **3️⃣ GetUsers - ดึงรายชื่อผู้ใช้ทั้งหมด**
func (repo repositoryDB) GetUsers(ctx context.Context, isActive bool, limit, offset int) ([]entity.ROM_V_UserDetail, error) {
	var users []entity.ROM_V_UserDetail
	query := `
		SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive
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

// ✅ **4️⃣ CheckUserExists - ตรวจสอบว่าผู้ใช้มีอยู่หรือไม่**
// 🔍 ตรวจสอบว่าผู้ใช้มีอยู่ใน ERP หรือไม่
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

// 🔍 ตรวจสอบว่าผู้ใช้มีอยู่ในระบบเว็บของเราหรือไม่
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

// ✅ **5️⃣ AddUser - เพิ่มผู้ใช้ใหม่**
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

	// 🔹 Step 1: Insert into UserRole
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

	// 🔹 Step 2: Insert into UserStatus (Default: Active)
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

	// 🔹 Step 3: Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ✅ **6️⃣ EditUser - แก้ไข Role และ Warehouse ของผู้ใช้**
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

	// 🟢 **Step 1: ตรวจสอบว่าผู้ใช้มีอยู่ในระบบ**
	exists, err := repo.CheckUserExists(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return errors.NotFoundError("user not found")
	}

	// 🟢 **Step 2: สร้าง Dynamic Query เฉพาะค่าที่ถูกส่งมา**
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

	// 🟢 **Step 3: Execute SQL**
	_, err = tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update UserRole: %w", err)
	}

	// 🟢 **Step 4: Commit transaction**
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ✅ **7️⃣ DeleteUser - ปิดการใช้งานบัญชี (Soft Delete)**
func (repo repositoryDB) DeleteUser(ctx context.Context, userID, adminID string) error {
	query := `
		UPDATE UserStatus
		SET IsActive = 0, UpdatedBy = :updatedBy, UpdatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":    userID,
		"updatedBy": adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
}

// ✅ **8️⃣ GetCurrentPassword - ดึงรหัสผ่านปัจจุบันของ User**
func (repo repositoryDB) GetCurrentPassword(ctx context.Context, userName string) (string, error) {
	var currentPassword string
	query := `SELECT Password FROM ROM_V_User WHERE UserName = :userName`
	params := map[string]interface{}{"userName": userName}

	err := repo.db.GetContext(ctx, &currentPassword, query, params)
	if err != nil {
		return "", errors.NotFoundError("password not found")
	}
	return currentPassword, nil
}

/* // ✅ **9️⃣ UpdateUserPassword - อัปเดตรหัสผ่านของ User**
func (repo repositoryDB) UpdateUserPassword(ctx context.Context, req request.ResetPasswordRequest, adminID string) error {
	query := `
		UPDATE UserStatus
		SET Password = :hashedPassword, UpdatedBy = :adminID, UpdatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":         req.UserID,
		"hashedPassword": req.NewPassword,
		"adminID":        adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
} */
