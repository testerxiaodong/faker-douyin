package service

import (
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/model/dto/response"
	"faker-douyin/internal/app/model/entity"
)

type UserServiceImpl struct {
	DataRepo *dao.DataRepo
}

func (u *UserServiceImpl) GetAllUser() ([]*entity.User, error) {
	users, err := u.DataRepo.Db.User.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserServiceImpl) GetByUsername(username string) (*entity.User, error) {
	user, err := u.DataRepo.Db.User.Where(u.DataRepo.Db.User.Username.Eq(username)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserServiceImpl) GetByID(id int64) (*response.UserInfoRes, error) {
	user, err := u.DataRepo.Db.User.Where(u.DataRepo.Db.User.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	var userInfo response.UserInfoRes
	userInfo.ID = user.ID
	userInfo.Username = user.Username
	return &userInfo, nil
}

func (u *UserServiceImpl) CreateUser(user entity.User) (*entity.User, error) {
	newUser := &entity.User{
		Username: user.Username,
		Password: user.Password,
	}
	err := u.DataRepo.Db.User.Create(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
