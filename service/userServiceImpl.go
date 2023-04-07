package service

import "C"
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"faker-douyin/global"
	"faker-douyin/model/dao"
	"faker-douyin/model/entity"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type UserServiceImpl struct {
}

func (u UserServiceImpl) GetTableUserList() []entity.TableUser {
	tableUsers, err := dao.GetTableUserList()
	if err != nil {
		log.Println(err.Error())
		return tableUsers
	}
	return tableUsers
}

func (u UserServiceImpl) GetTableUserByUsername(name string) entity.TableUser {
	tableUser, err := dao.GetTableUserByName(name)
	if err != nil {
		log.Println(err.Error())
		return tableUser
	}
	return tableUser
}

func (u UserServiceImpl) GetTableUserById(id uint64) entity.TableUser {
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println(err.Error())
		return tableUser
	}
	return tableUser
}

func (u UserServiceImpl) InsertTableUser(tableUser *entity.TableUser) bool {
	result := dao.InsertTableUser(tableUser)
	if result == false {
		log.Println("插入失败")
		return false
	}
	return true
}

//func (u UserServiceImpl) GetUserById(id uint64) (response.GetUserByIdRes, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u UserServiceImpl) GetUserByIdWithCurId(id uint64, curId uint64) (response.GetUserByIdRes, error) {
//	//TODO implement me
//	panic("implement me")
//}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := UserService.GetTableUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	token := NewToken(u)
	println(token)
	return token
}

// NewToken 根据信息创建token
func NewToken(u entity.TableUser) string {
	expiresTime := time.Now().Unix() + int64(global.OneDayOfHours)
	fmt.Printf("expiresTime: %v\n", expiresTime)
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

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Result: " + sha)
	return sha
}
