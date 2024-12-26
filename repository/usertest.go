package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	response "boilerplate-backend-go/dto/response"
)

type UserTestRepository interface {
	GetUserTest(userID, password string) (response.UserInform, error)
}

func (repo repositoryDB) GetUserTest(userID, password string) (response.UserInform, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// โครงสร้างข้อมูลผู้ใช้งาน
	var user response.UserInform

	// คำสั่ง SQL ดึงข้อมูลจาก View
	query := `
		SELECT UserID, Username, NickName, FullNameTH, DepartmentNo
		FROM Data_WebReturn.dbo.ROM_V_User
		WHERE UserID = @userID AND Password = @password
	`
	fmt.Printf("Executing Query: %s\nParams: UserID=%s, Password=%s\n", query, userID, password)

	// ดึงข้อมูลและตรวจสอบรหัสผ่าน
	err := repo.db.GetContext(ctx, &user, query,
		sql.Named("userID", userID),
		sql.Named("password", password), // ค่าแฮช MD5 ที่ส่งเข้ามา
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No User Found for UserID=%s\n", userID)
			return user, fmt.Errorf("invalid credentials")
		}
		fmt.Printf("Error Querying Database: %s\n", err.Error())
		return user, fmt.Errorf("error querying user: %w", err)
	}

	fmt.Printf("User Retrieved: %+v\n", user)
	return user, nil
}
