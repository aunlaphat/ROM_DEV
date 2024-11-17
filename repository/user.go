package repository

import (
	response "boilerplate-backend-go/dto/response"
	"context"
	"database/sql"
	"time"
)

type UserRepository interface {
	GetUser(username, password string) (response.Login, error)
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

func (repo repositoryDB) GetUser(username, password string) (response.Login, error) {
	return response.Login{}, nil
}
