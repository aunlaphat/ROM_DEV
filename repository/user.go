package repository

import (
	"context"

	"boilerplate-backend-go/entity"
	"boilerplate-backend-go/errors"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID string) (entity.ROM_V_UserDetail, error)
	GetUsers(ctx context.Context, isActive *bool, limit, offset int) ([]entity.ROM_V_UserDetail, error)
	CheckUserExists(ctx context.Context, userID string) (bool, error)
	AddUser(ctx context.Context, userID string, roleID int, warehouseID string, adminID string) error
	EditUser(ctx context.Context, userID string, newRoleID int, adminID string) error
	DeleteUser(ctx context.Context, userID, adminID string) error
	GetCurrentPassword(ctx context.Context, userID string) (string, error)
	UpdateUserPassword(ctx context.Context, userID, hashedPassword, adminID string) error
}

// ✅ 1️⃣ GetUser - ดึงข้อมูลผู้ใช้จาก View `ROM_V_UserDetail`
func (repo repositoryDB) GetUser(ctx context.Context, userID string) (entity.ROM_V_UserDetail, error) {
	var user entity.ROM_V_UserDetail
	query := `
        SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive
        FROM ROM_V_UserDetail
        WHERE UserID = :userID
    `
	params := map[string]interface{}{"userID": userID}

	err := repo.db.GetContext(ctx, &user, query, params)
	if err != nil {
		return entity.ROM_V_UserDetail{}, errors.NotFoundError("user not found")
	}
	return user, nil
}

// ✅ 2️⃣ GetUsers - ดึงรายชื่อผู้ใช้ทั้งหมด พร้อมฟิลเตอร์ `isActive`
func (repo repositoryDB) GetUsers(ctx context.Context, isActive *bool, limit, offset int) ([]entity.ROM_V_UserDetail, error) {
	query := `
		SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo, RoleID, RoleName, Description, IsActive
		FROM ROM_V_UserDetail
		WHERE (:isActive IS NULL OR IsActive = :isActive)
		ORDER BY UserID
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
	`
	params := map[string]interface{}{
		"isActive": isActive,
		"limit":    limit,
		"offset":   offset,
	}

	var users []entity.ROM_V_UserDetail
	err := repo.db.SelectContext(ctx, &users, query, params)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ✅ 3️⃣ CheckUserExists - ตรวจสอบว่าผู้ใช้มีอยู่หรือไม่
func (repo repositoryDB) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) FROM ROM_V_User WHERE UserID = :userID`
	params := map[string]interface{}{"userID": userID}

	err := repo.db.GetContext(ctx, &exists, query, params)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ✅ 4️⃣ AddUser - เพิ่มผู้ใช้ใหม่ (กำหนด Role)
func (repo repositoryDB) AddUser(ctx context.Context, user entity.UserRole, adminID string) error {
	query := `
		INSERT INTO UserRole (UserID, RoleID, CreatedBy, CreatedAt)
		VALUES (:userID, :roleID, :adminID, GETDATE())
	`
	params := map[string]interface{}{
		"userID":  user.UserID,
		"roleID":  user.RoleID,
		"adminID": adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
}

// ✅ 5️⃣ EditUser - อัปเดตข้อมูล Role ของผู้ใช้
func (repo repositoryDB) EditUser(ctx context.Context, userID string, updatedUser entity.UserRole, adminID string) error {
	query := `
		UPDATE UserRole
		SET RoleID = :roleID, UpdatedBy = :adminID, UpdatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":  userID,
		"roleID":  updatedUser.RoleID,
		"adminID": adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
}

// ✅ 6️⃣ DeleteUser - ปรับ `IsActive = 0` (Soft Delete)
func (repo repositoryDB) DeleteUser(ctx context.Context, userID, adminID string) error {
	query := `
		UPDATE UserStatus
		SET IsActive = 0, UpdatedBy = :adminID, UpdatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":  userID,
		"adminID": adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
}

// ✅ 7️⃣ GetCurrentPassword - ดึงรหัสผ่านปัจจุบันของ User
func (repo repositoryDB) GetCurrentPassword(ctx context.Context, userID string) (string, error) {
	var currentPassword string
	query := `SELECT Password FROM UserStatus WHERE UserID = :userID`
	params := map[string]interface{}{"userID": userID}

	err := repo.db.GetContext(ctx, &currentPassword, query, params)
	if err != nil {
		return "", errors.NotFoundError("password not found")
	}
	return currentPassword, nil
}

// ✅ 8️⃣ UpdateUserPassword - อัปเดตรหัสผ่านผู้ใช้
func (repo repositoryDB) UpdateUserPassword(ctx context.Context, userID, hashedPassword, adminID string) error {
	query := `
		UPDATE UserStatus
		SET Password = :hashedPassword, UpdatedBy = :adminID, UpdatedAt = GETDATE()
		WHERE UserID = :userID
	`
	params := map[string]interface{}{
		"userID":         userID,
		"hashedPassword": hashedPassword,
		"adminID":        adminID,
	}

	_, err := repo.db.NamedExecContext(ctx, query, params)
	return err
}
