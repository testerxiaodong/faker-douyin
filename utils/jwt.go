package utils

import (
	"errors"
	"faker-douyin/global"
	"faker-douyin/model/entity"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

// GenerateToken 根据username生成一个token
func GenerateToken(user *entity.User) string {
	//fmt.Printf("generate token: %v\n", u)
	expiresTime := time.Now().Unix() + int64(global.OneDayOfHours)
	id64 := user.ID
	fmt.Printf("id: %v\n", strconv.FormatInt(int64(id64), 10))
	claims := jwt.StandardClaims{
		Audience:  user.Username,
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
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// ParseToken 从token中解析出StandardClaims
func ParseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(global.Config.Jwt.Secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if jwtToken != nil {
		if claims, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
