package service

import (
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
)

type UserService interface {
	/*
		个人使用
	*/
	// GetTableUserList 获得全部TableUser对象
	GetTableUserList() []entity.TableUser

	// GetTableUserByUsername 根据username获得TableUser对象
	GetTableUserByUsername(name string) (entity.TableUser, error)

	// GetTableUserById 根据user_id获得TableUser对象
	GetTableUserById(id uint64) (response.UserInfoRes, error)

	// InsertTableUser 将tableUser插入表内
	InsertTableUser(tableUser entity.TableUser) (entity.TableUser, error)
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
