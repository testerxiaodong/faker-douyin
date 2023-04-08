package utils

import (
	"faker-douyin/global"
	"faker-douyin/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := service.UserService.GetTableUserByUsername(new(service.UserServiceImpl), username)
	//fmt.Printf("generate token: %v\n", u)
	expiresTime := time.Now().Unix() + int64(global.OneDayOfHours)
	id64 := u.ID
	fmt.Printf("id: %v\n", strconv.FormatInt(int64(id64), 10))
	claims := jwt.StandardClaims{
		Audience:  u.Name,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(int64(id64), 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(global.Config.Jwt.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}
