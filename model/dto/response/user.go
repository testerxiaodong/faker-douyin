package response

type UserRegisterSuccessRes struct {
	Id    uint64 `json:"id"`
	Token string `json:"token"`
}

type UserLoginSuccessRes struct {
	Id    uint64 `json:"id"`
	Token string `json:"token"`
}

type UserInfoRes struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
