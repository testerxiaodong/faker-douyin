package dao

import (
	"faker-douyin/model/entity"
	"log"
)

// GetTableUserList 获取所有用户
func GetTableUserList() ([]entity.TableUser, error) {
	var tableUsers []entity.TableUser
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

// GetTableUserByName 根据用户名获取第一个用户
func GetTableUserByName(name string) (entity.TableUser, error) {
	var user entity.TableUser
	if err := Db.Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// GetTableUserById 根据id获取第一个用户
func GetTableUserById(id uint64) (entity.TableUser, error) {
	var user entity.TableUser
	if err := Db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(tableUser *entity.TableUser) bool {
	if err := Db.Create(tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
