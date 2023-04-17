package response

type UserRegisterSuccessRes struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

type UserLoginSuccessRes struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

type UserInfoRes struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
