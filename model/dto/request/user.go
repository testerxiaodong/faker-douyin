package request

type UserRegisterReq struct {
	Name     string `json:"name,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserLoginReq struct {
	Name     string `json:"name,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserInfoReq struct {
	UserId uint64 `json:"user_id"`
}
