package v1

import (
	"faker-douyin/model/common"
	"faker-douyin/model/dto/request"
	"faker-douyin/model/dto/response"
	"faker-douyin/model/entity"
	"faker-douyin/service"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
	"log"
)

// Register POST douyin/v1/user/register/ 用户注册
func Register(c *gin.Context) {
	var userRegisterReq request.UserRegisterReq
	if err := c.ShouldBindJSON(&userRegisterReq); err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	usi := service.UserServiceImpl{}
	// 根据用户名查询用户是否存在
	u, err := usi.GetTableUserByUsername(userRegisterReq.Name)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 用户名存在，注册失败
	if u.Name != "" {
		common.FailWithMessage("User already exist", c)
	} else {
		newUser := entity.TableUser{
			Name:     userRegisterReq.Name,
			Password: utils.EnCoder(userRegisterReq.Password),
		}
		// 用户名不存在，插入数据
		user, err := usi.InsertTableUser(newUser)
		if err != nil {
			println("Insert Data Fail")
			common.FailWithMessage(err.Error(), c)
			return
		}
		token := utils.GenerateToken(&user)
		log.Println("注册返回的id: ", user.ID)
		common.OkWithDetailed(response.UserRegisterSuccessRes{Id: uint64(user.ID), Token: token}, "注册成功", c)
	}
}

// Login POST douyin/v1/user/login/ 用户登陆
func Login(c *gin.Context) {
	var userLoginReq request.UserLoginReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&userLoginReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	usi := service.UserServiceImpl{}
	// 根据name查询用户
	u, err := usi.GetTableUserByUsername(userLoginReq.Name)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 查询成功，将密码加密后与数据库密码比较
	if utils.EnCoder(userLoginReq.Password) == u.Password {
		token := utils.GenerateToken(&u)
		common.OkWithDetailed(response.UserLoginSuccessRes{Id: uint64(u.ID), Token: token}, "登陆成功", c)
		return
	} else {
		common.FailWithMessage("username or password error", c)
	}
}

// UserInfo GET douyin/v1/user/ 获取用户信息
func UserInfo(c *gin.Context) {
	var userInfoReq request.UserInfoReq
	// 请求参数绑定和校验
	err := c.ShouldBindJSON(&userInfoReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	usi := service.UserServiceImpl{}
	// 根据id查询用户信息
	userInfo, err := usi.GetTableUserById(userInfoReq.UserId)
	if err != nil {
		common.FailWithMessage("user not exist", c)
		return
	}
	common.OkWithData(userInfo, c)
}
