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

// Register POST douyin/user/register/ 用户注册
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
		token := utils.GenerateToken(userRegisterReq.Name)
		log.Println("注册返回的id: ", u.ID)
		common.OkWithData(response.UserLoginSuccessRes{Id: uint64(u.ID), Token: token}, c)
	}
}
