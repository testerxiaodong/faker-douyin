package jwt

import (
	"errors"
	"faker-douyin/model/common"
	"faker-douyin/utils"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			// 未携带token
			common.FailWithMessage("Unauthorized", context)
			context.Abort()
			return
		}
		Claims, err := utils.ParseToken(auth)
		if err != nil {
			if errors.Is(utils.TokenExpired, err) {
				// token过期
				common.FailWithMessage("token expired", context)
				context.Abort()
				return
			}
			// token无效
			common.FailWithMessage("Token Error", context)
			context.Abort()
			return
		}
		// 解析成功，将userId放入上下文（其实应该将Claims放入上下文，之后Claims会存放更多东西，角色，rbac权限等）
		context.Set("userId", Claims.Id)
		context.Next()
	}
}
