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
	u := usi.GetTableUserByUsername(userRegisterReq.Name)
	if u.Name != "" {
		common.FailWithMessage("User already exist", c)
	} else {
		newUser := entity.TableUser{
			Name:     userRegisterReq.Name,
			Password: utils.EnCoder(userRegisterReq.Password),
		}
		if usi.InsertTableUser(&newUser) != true {
			println("Insert Data Fail")
		}
		u := usi.GetTableUserByUsername(userRegisterReq.Name)
		token := utils.GenerateToken(&u)
		log.Println("注册返回的id: ", u.ID)
		common.OkWithDetailed(response.UserRegisterSuccessRes{Id: uint64(u.ID), Token: token}, "注册成功", c)
	}
}

// Login POST douyin/v1/user/login/ 用户登陆
func Login(c *gin.Context) {
	var userLoginReq request.UserLoginReq
	err := c.ShouldBindJSON(&userLoginReq)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(userLoginReq.Name)
	if utils.EnCoder(userLoginReq.Password) == u.Password {
		token := utils.GenerateToken(&u)
		common.OkWithDetailed(response.UserLoginSuccessRes{Id: uint64(u.ID), Token: token}, "登陆成功", c)
		return
	} else {
		common.FailWithMessage("username or password error", c)
	}
}
