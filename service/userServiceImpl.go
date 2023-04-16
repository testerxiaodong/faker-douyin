package service

import "C"
import (
	"faker-douyin/model/dao"
	"faker-douyin/model/dto/response"
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

func (u UserServiceImpl) GetTableUserByUsername(name string) (entity.TableUser, error) {
	tableUser, err := dao.GetTableUserByName(name)
	if err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

func (u UserServiceImpl) GetTableUserById(id uint64) (response.UserInfoRes, error) {
	var userInfo response.UserInfoRes
	tableUser, err := dao.GetTableUserById(id)
	userInfo.Id = uint64(tableUser.ID)
	userInfo.Name = tableUser.Name
	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	}
	return userInfo, nil
}

func (u UserServiceImpl) InsertTableUser(tableUser entity.TableUser) (entity.TableUser, error) {
	var user entity.TableUser
	user, err := dao.InsertTableUser(tableUser)
	if err != nil {
		log.Println("插入失败")
		return user, err
	}
	return user, nil
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
