package v1

import (
	"errors"
	"faker-douyin/model/common"
	"faker-douyin/model/dto/request"
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
	"faker-douyin/service"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type UserController struct {
}

// Register POST douyin/v1/user/register/ 用户注册
func (u *UserController) Register(c *gin.Context) {
	var userRegisterReq request.UserRegisterReq
	if err := c.ShouldBindJSON(&userRegisterReq); err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 依赖倒转原则，面向抽象层进行开发
	usi := service.NewUserService()
	// 根据用户名查询用户是否存在
	user, err := usi.GetByUsername(userRegisterReq.Username)
	// 如果有错误，并且错误不是没有记录
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 用户名不存在，开始注册
	if user == nil {
		newUser := entity.User{
			Username: userRegisterReq.Username,
			Password: utils.EnCoder(userRegisterReq.Password),
		}
		// 用户名不存在，插入数据
		user, err := usi.CreateUser(newUser)
		if err != nil {
			println("Insert Data Fail")
			common.FailWithMessage(err.Error(), c)
			return
		}
		token := utils.GenerateToken(user)
		log.Println("注册返回的id: ", user.ID)
		common.OkWithDetailed(response.UserRegisterSuccessRes{ID: uint64(user.ID), Token: token}, "注册成功", c)
	} else {
		// 用户名存在，不允许注册
		common.FailWithMessage("User already exist", c)
	}
}

// Login POST douyin/v1/user/login/ 用户登陆
func (u *UserController) Login(c *gin.Context) {
	var userLoginReq request.UserLoginReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&userLoginReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 依赖倒转原则，面向抽象层进行开发
	var usi service.UserService
	usi = new(service.UserServiceImpl)
	// 根据name查询用户
	user, err := usi.GetByUsername(userLoginReq.Username)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 查询成功，将密码加密后与数据库密码比较
	if utils.EnCoder(userLoginReq.Password) == user.Password {
		token := utils.GenerateToken(user)
		common.OkWithDetailed(response.UserLoginSuccessRes{ID: uint64(user.ID), Token: token}, "登陆成功", c)
		return
	} else {
		common.FailWithMessage("username or password error", c)
	}
}

// UserInfo GET douyin/v1/user/ 获取用户信息
func (u *UserController) UserInfo(c *gin.Context) {
	var userInfoReq request.UserInfoReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&userInfoReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 依赖倒转原则，面向抽象层进行开发
	var usi service.UserService
	usi = new(service.UserServiceImpl)
	// 根据id查询用户信息
	userInfo, err := usi.GetByID(userInfoReq.UserId)
	if err != nil {
		common.FailWithMessage("user not exist", c)
		return
	}
	common.OkWithData(userInfo, c)
}
