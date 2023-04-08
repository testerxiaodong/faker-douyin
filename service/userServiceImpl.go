package service

import "C"
import (
	"faker-douyin/model/dao"
	"faker-douyin/model/entity"
	"log"
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
