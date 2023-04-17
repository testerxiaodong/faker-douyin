package request

type UserRegisterReq struct {
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserLoginReq struct {
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserInfoReq struct {
	UserId int64 `json:"user_id"`
}
