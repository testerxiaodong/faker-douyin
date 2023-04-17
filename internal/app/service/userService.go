package service

import (
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
)

type UserService interface {
	/*
		个人使用
	*/

	GetAllUser() ([]*entity.User, error)

	GetByUsername(username string) (*entity.User, error)

	GetByID(id int64) (*response.UserInfoRes, error)

	CreateUser(user entity.User) (*entity.User, error)
	/*
		他人使用
	*/
	// GetUserById 未登录情况下,根据user_id获得User对象
	//GetUserById(id uint64) (response.GetUserByIdRes, error)

	// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	//GetUserByIdWithCurId(id uint64, curId uint64) (response.GetUserByIdRes, error)

	// 根据token返回id
	// 接口:auth中间件,解析完token,将userid放入context
	//(调用方法:直接在context内拿参数"userId"的值)	fmt.Printf("userInfo: %v\n", c.GetString("userId"))
}
