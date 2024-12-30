package repository

import (
	response "boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// UserRepository interface กำหนด method สำหรับการทำงานกับฐานข้อมูลผู้ใช้
type UserRepository interface {
	GetUser(ctx context.Context, username, password string) (response.Login, error)
	GetUserFromLark(username, userID string) (response.Login, error)
}

func (repo repositoryDB) GetUserFromLark(username, userID string) (response.Login, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user response.Login
	user.UserName = username
	query := `
        SELECT UserName, UserID, RoleID, NickName, FullNameTH , 'lark' as  Platform
        FROM V_User_Login
        WHERE UserName = @userName AND UserID = @userID
    `
	err := repo.db.GetContext(ctx, &user, query,
		sql.Named("userName", username),
		sql.Named("userID", userID),
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Implementation สำหรับ GetUser
func (repo repositoryDB) GetUser(ctx context.Context, username, password string) (response.Login, error) {
	var user response.Login
	query := `
        SELECT UserID, UserName, NickName, FullNameTH, DepartmentNo
        FROM ROM_V_User
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
