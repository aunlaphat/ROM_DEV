package service

import (
	request "boilerplate-backend-go/dto/request"
	response "boilerplate-backend-go/dto/response"
	// "crypto/md5"
	// "encoding/hex"
	"fmt"
)

type UserTestService interface {
	LoginTest(req request.Login) (response.UserInform, error)
}

func (srv service) LoginTest(req request.Login) (response.UserInform, error) {
	// ตรวจสอบว่า UserID และ Password ถูกส่งมาหรือไม่
	if req.UserID == "" || req.Password == "" {
		return response.UserInform{}, fmt.Errorf("userID and password must not be empty")
	}

	// แปลงรหัสผ่านที่รับมาจากผู้ใช้ให้เป็นค่าแฮช MD5
	// hasher := md5.New()
	// hasher.Write([]byte(req.Password))
	// hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	// fmt.Printf("Hashed Password: %s\n", hashedPassword)

	// ดึงข้อมูลผู้ใช้งานจากฐานข้อมูล โดยส่งค่าแฮชเพื่อเปรียบเทียบ
	user, err := srv.usertestRepo.GetUserTest(req.UserID, req.Password)
	if err != nil {
		return response.UserInform{}, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
